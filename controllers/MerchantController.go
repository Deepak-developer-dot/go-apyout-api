package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payout-backend/config"
	"payout-backend/models"
	"payout-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AddMerchant(c *gin.Context) {
	db := config.DB
	var req struct {
		OwnerName              string  `json:"owner_name" binding:"required"`
		Email                  string  `json:"email" binding:"required"`
		Password               string  `json:"password" binding:"required"`
		AppName                string  `json:"app_name" binding:"required"`
		Mobile                 uint    `json:"mobile_number" binding:"required"`
		PayoutProvider         uint    `json:"payout_provider" binding:"required"`
		AppUrl                 string  `json:"app_url" binding:"required"`
		WebUrl                 string  `json:"web_url" binding:"required"`
		WebhookUrl             string  `json:"webhook_url" binding:"required"`
		IpAddress              string  `json:"ip_address" binding:"required"`
		WithdrawalCommission   float32 `json:"withdrawal_commission" binding:"required"`
		DirectPayoutEnabled    float32 `json:"direct_payout_enabled" binding:"required"`
		CommissionChargeType   uint    `json:"commission_charge_type" binding:"required"`
		AutomaticCreditEnabled float32 `json:"automatic_credit_enabled" binding:"required"`
		IsWebhookActive        float32 `json:"is_webhook_active" binding:"required"`
		// SortBy                 float32 `json:"sort_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	// loggedInAdminID:= c.GetUint("user_id")
	// loggedInAdminRole:= c.GetUint("role")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": false, "message": "error while hashing password"})
		return
	}

	merchant := models.Merchant{
		OwnerName:              req.OwnerName,
		AppName:                req.AppName,
		Email:                  req.Email,
		Password:               string(hashedPassword),
		Mobile:                 req.Mobile,
		PayoutProvider:         req.PayoutProvider,
		AppUrl:                 req.AppUrl,
		WebUrl:                 req.WebUrl,
		WebhookUrl:             req.WebhookUrl,
		IpAddress:              req.IpAddress,
		IsAdmin:                2,
		Status:                 1,
		IsDeleted:              0,
		WithdrawalCommission:   float64(req.WithdrawalCommission),
		DirectPayoutEnabled:    uint(req.DirectPayoutEnabled),
		CommissionChargeType:   uint(req.CommissionChargeType),
		AutomaticCreditEnabled: uint(req.AutomaticCreditEnabled),
		IsWebhookActive:        uint(req.IsWebhookActive),
		// SortBy:                 uint(req.SortBy),
	}
	// if(loggedInAdminRole ==4){
	// 	merchant.ReferenceID= loggedInAdminID
	// }
	if err := db.Create(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "false", "message": "failed to create merchant"})
		return
	}

	apiKey, secretKey, err := utils.GenerateApiAndSecretKey(merchant.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to generate api key and secret key"})
		return
	}

	merchantAccountDetails := models.MerchantAccountDetail{
		MerchantID:    merchant.ID,
		ApiKey:        apiKey,
		SecretKey:     secretKey,
		AccountNumber: "576510110008967",
		IfscCode:      "BKID0005765",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&merchantAccountDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to generate merchant app key and sceret key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"message":  "Merchant Created Successfully!",
		"merchant": merchant,
	})

}

func GetMerchants(c *gin.Context) {
	db := config.DB
	role := c.GetInt("role")
	userID := c.GetUint("user_id")

	keywordSearch := c.Query("keyword_search")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("merchant:keywod_search=%s:page=%d", keywordSearch, page)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	query := db.Model(&models.Merchant{}).Where("is_admin =?", 2)
	if role == 2 {
		query = query.Where("id", userID)
	}
	if keywordSearch != "" {
		query = query.Where("name ILIKE ?", keywordSearch)
	}

	var merchants []models.Merchant
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&merchants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch merchants"})
		return
	}

	response := gin.H{
		"message":   "Merchant Fetched Successfully",
		"status":    true,
		"merchants": merchants,
		"page":      page,
		"per_page":  limit,
	}
	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Second)
	c.JSON(http.StatusOK, response)
}

func GetMerchantsByID(c *gin.Context) {
	id := c.Param("id")

	db := config.DB

	cacheKey := fmt.Sprintf("get_merchantbyId=%s", id)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err != nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	var merchant models.Merchant
	if err := db.Preload("MerchantAccountDetail").Where("id=? AND is_admin=?", id, 2).
		First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Merchant Not Found",
		})
	}
	response := gin.H{
		"message":  "Merchant Fetch Sucessfully",
		"status":   true,
		"merchant": merchant,
	}
	data, _ := json.Marshal(merchant)
	config.Cache.Set(cacheKey, string(data), 10*time.Minute)
	c.JSON(http.StatusOK, response)
}

func UpdateMerchant(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	merchantID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var req struct {
		OwnerName              string  `json:"owner_name" binding:"required"`
		Email                  string  `json:"email" binding:"required"`
		Password               string  `json:"password" binding:"required"`
		AppName                string  `json:"app_name" binding:"required"`
		Mobile                 uint    `json:"mobile_number" binding:"required"`
		PayoutProvider         uint    `json:"payout_provider" binding:"required"`
		AppUrl                 string  `json:"app_url" binding:"required"`
		WebUrl                 string  `json:"web_url" binding:"required"`
		WebhookUrl             string  `json:"webhook_url" binding:"required"`
		IpAddress              string  `json:"ip_address" binding:"required"`
		WithdrawalCommission   float32 `json:"withdrawal_commission" binding:"required"`
		DirectPayoutEnabled    float32 `json:"direct_payout_enabled" binding:"required"`
		CommissionChargeType   uint    `json:"commission_charge_type" binding:"required"`
		AutomaticCreditEnabled float32 `json:"automatic_credit_enabled" binding:"required"`
		IsWebhookActive        float32 `json:"is_webhook_active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	var merchant models.Merchant
	if err := db.First(&merchant, merchantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Merchant Not Found"})
		return
	}

	merchant.OwnerName = req.OwnerName
	merchant.AppName = req.AppName
	merchant.Email = req.Email
	merchant.Mobile = req.Mobile
	merchant.PayoutProvider = req.PayoutProvider
	merchant.AppUrl = req.AppUrl
	merchant.WebUrl = req.WebUrl
	merchant.WebhookUrl = req.WebhookUrl
	merchant.IpAddress = req.IpAddress
	merchant.IsAdmin = 2
	merchant.Status = 1
	merchant.IsDeleted = 0
	merchant.WithdrawalCommission = float64(req.WithdrawalCommission)
	merchant.DirectPayoutEnabled = uint(req.DirectPayoutEnabled)
	merchant.CommissionChargeType = uint(req.CommissionChargeType)
	merchant.AutomaticCreditEnabled = uint(req.AutomaticCreditEnabled)
	merchant.IsWebhookActive = uint(req.IsWebhookActive)
	// SortBy:                 uint(req.SortBy),

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error while hasing the password"})
			return
		}
		merchant.Password = string(hashedPassword)

	}

	if err := db.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to update merchant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"message":  "Merchant Updated Successfully",
		"merchant": merchant,
	})

}

func GetMerchantAccountDetails(c *gin.Context) {
	db := config.DB
	role := c.GetInt("role")
	userID := c.GetUint("user_id")

	keywordSearch := c.Query("keyword_search")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10
	offset := (page - 1) * limit

	fmt.Println("role", role)
	fmt.Println("userID", userID)

	cacheKey := fmt.Sprintf("merchant_account_details:keyword_search=%s:page=%d", keywordSearch, offset)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	query := db.Model(&models.MerchantAccountDetail{})
	if role == 2 {
		query = query.Where("merchant_id", userID)
	}
	if keywordSearch != "" {
		query = query.Where("app_key ILIKE ?", keywordSearch)
	}

	var merchantAccountDetails []models.MerchantAccountDetail
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&merchantAccountDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch merchantAccountDetails"})
		return
	}

	response := gin.H{
		"message":                "Merchant Account Details Fetched Successfully",
		"status":                 true,
		"merchantAccountDetails": merchantAccountDetails,
		"page":                   page,
		"per_page":               limit,
	}
	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Second)
	c.JSON(http.StatusOK, response)

}
