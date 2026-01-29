package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payout-backend/config"
	"payout-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MakePayout(c *gin.Context) {
	db := config.DB
	var req struct {
		TransactionID string `json:"transaction_id" binding:"required"`
		Mobile        uint   `json:"mobile" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Amount        uint   `json:"amount" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
		IfscCode      string `json:"ifsc_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var count int64
	db.Model(&models.WalletTransaction{}).
		Where(" client_transaction_id = ?", req.TransactionID).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Transaction Id already processed",
		})
		return
	}
	merchant, _ := c.MustGet("merchant").(models.Merchant)
	serviceCharge, gst, totalAmount := calculateCharges(float64(req.Amount), merchant.WithdrawalCommission)

	var walletBalance float64
	db.Model(&models.Wallet{}).
		Where("merchant_id", merchant.ID).
		Select("COALESCE(SUM(amount),0)").
		Scan(&walletBalance)

	if walletBalance < totalAmount {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Insufficient Fund!",
		})
		return
	}

	tx := db.Begin()
	transactionID := "TXN" + strconv.FormatInt(time.Now().Unix(), 10)

	walletTransaction := models.WalletTransaction{
		MerchantID:          merchant.ID,
		Amount:              c.GetFloat64(req.Amount),
		TransactionID:       transactionID,
		ClientTransactionID: req.TransactionID,
		ServiceCharge:       serviceCharge,
		GstCharge:           gst,
		PaymentType:         2,
		PaymentTypeString:   "Debit",
		CurWalletBal:        walletBalance,
		PayoutProvider:      merchant.PayoutProvider,
		ContactMobileNumber: req.AccountNumber,
		Status:              0,
	}
	// if(loggedInAdminRole ==4){
	// 	merchant.ReferenceID= loggedInAdminID
	// }
	if err := tx.Create(&walletTransaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "false", "message": "wallet Transaction failed"})
		return
	}

	customer := models.CustomerDetails{
		TransactionID: transactionID,
		Name:          req.Name,
		Email:         req.Email,
		AccountNumber: string(req.AccountNumber),
		IfscCode:      req.IfscCode,
	}
	if err := tx.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Customer save failed"})
		return
	}

	tx.Model(&models.Wallet{}).
		Where("merchant_id", merchant.ID).
		Update("amount", gorm.Expr("amount -?", totalAmount))
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status":              true,
		"message":             "Payment is in processing",
		"exch_transaction_id": transactionID,
	})

}

func calculateCharges(amount float64, commission float64) (serviceCharge, gst, total float64) {
	gstRate := 18.0

	platformCommission := amount * commission / 100
	gst = platformCommission * gstRate / 100
	serviceCharge = platformCommission + gst
	total = amount + serviceCharge
	return
}

func GetTransaction(c *gin.Context) {
	db := config.DB
	role := c.GetInt("role")
	userID := c.GetUint("user_id")

	keywordSearch := c.Query("keyword_search")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("tranactions:keywod_search=%s:page=%d", keywordSearch, page)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	query := db.Model(&models.WalletTransaction{})
	if role == 2 {
		query = query.Where("merchant_id", userID)
	}
	// if keywordSearch != "" {
	// 	query = query.Where("name ILIKE ?", keywordSearch)
	// }

	var tranactions []models.WalletTransaction
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&tranactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tranactions"})
		return
	}

	response := gin.H{
		"message":     "Wallet Transaction Fetched Successfully",
		"status":      true,
		"tranactions": tranactions,
		"page":        page,
		"per_page":    limit,
	}
	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Second)
	c.JSON(http.StatusOK, response)
}
func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	db := config.DB

	cacheKey := fmt.Sprintf("get_transactionbyId=%s", id)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err != nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	var tranactions models.WalletTransaction
	if err := db.Preload("CustomerDetail").
		Where("id = ?", id).
		First(&tranactions).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Transaction Not Found",
		})
		return
	}
	response := gin.H{
		"message":     "Tranactions Fetch Sucessfully",
		"status":      true,
		"tranactions": tranactions,
	}
	data, _ := json.Marshal(tranactions)
	config.Cache.Set(cacheKey, string(data), 10*time.Minute)
	c.JSON(http.StatusOK, response)
}
