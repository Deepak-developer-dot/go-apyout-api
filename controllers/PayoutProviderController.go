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
)

func AddPayoutProvider(c *gin.Context) {
	db := config.DB
	var req struct {
		ProviderType  uint    `json:"provider_type" binding:"required"`
		Name          string  `json:"name" binding:"required"`
		MerchantName  string  `json:"merchant_name" binding:"required"`
		Email         string  `json:"email" binding:"required"`
		AccountNumber string  `json:"account_number" binding:"required"`
		AppKey        string  `json:"app_key" binding:"required"`
		SecretKey     string  `json:"secret_key" binding:"required"`
		Balance       float64 `json:"balance" binding:"required"`

		Status uint `json:"status" binding:"required"`
		// SortBy                 float32 `json:"sort_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	payoutProvider := models.PayoutProviders{
		Name:          req.Name,
		MerchantName:  req.MerchantName,
		Email:         req.Email,
		AccountNumber: req.AccountNumber,
		AppKey:        req.AppKey,
		SecretKey:     req.SecretKey,
		Balance:       req.Balance,
		Status:        req.Status,
	}
	// if(loggedInAdminRole ==4){
	// 	merchant.ReferenceID= loggedInAdminID
	// }
	if err := db.Create(&payoutProvider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "false", "message": "failed to create payoutProvider"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          true,
		"message":         "Payout Provider Created Successfully!",
		"payout provider": payoutProvider,
	})

}

func GetPayoutProvider(c *gin.Context) {
	db := config.DB
	keywordSearch := c.Query("keyword_search")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("payout_provider:keywod_search=%s:page=%d", keywordSearch, page)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	query := db.Model(&models.PayoutProviders{})

	if keywordSearch != "" {
		query = query.Where("name ILIKE ?", keywordSearch)
	}

	var payoutProvider []models.PayoutProviders
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&payoutProvider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch payoutProvider"})
		return
	}

	response := gin.H{
		"message":        "Payout Provider Fetched Successfully",
		"status":         true,
		"payoutProvider": payoutProvider,
		"page":           page,
		"per_page":       limit,
	}
	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Second)
	c.JSON(http.StatusOK, response)
}

func GetPayoutProviderByID(c *gin.Context) {
	id := c.Param("id")

	db := config.DB

	cacheKey := fmt.Sprintf("get_payoutproviderbyId=%s", id)

	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err != nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	var payoutProvider models.PayoutProviders
	if err := db.Where("id=?", id).
		First(&payoutProvider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "payoutProvider Not Found",
		})
	}
	response := gin.H{
		"message":        "PayoutProvider Fetch Sucessfully",
		"status":         true,
		"payoutProvider": payoutProvider,
	}
	data, _ := json.Marshal(payoutProvider)
	config.Cache.Set(cacheKey, string(data), 10*time.Minute)
	c.JSON(http.StatusOK, response)
}

func UpdatePayoutProvider(c *gin.Context) {
	db := config.DB
	id := c.Param("id")

	payoutProviderID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var req struct {
		ProviderType  uint    `json:"provider_type" binding:"required"`
		Name          string  `json:"name" binding:"required"`
		MerchantName  string  `json:"merchant_name" binding:"required"`
		Email         string  `json:"email" binding:"required"`
		AccountNumber string  `json:"account_number" binding:"required"`
		AppKey        string  `json:"app_key" binding:"required"`
		SecretKey     string  `json:"secret_key" binding:"required"`
		Balance       float64 `json:"balance" binding:"required"`

		Status uint `json:"status" binding:"required"`
		// SortBy                 float32 `json:"sort_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	var payoutProvider models.PayoutProviders
	if err := db.First(&payoutProvider, payoutProviderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "payoutProvider Not Found"})
		return
	}

	payoutProvider = models.PayoutProviders{
		Name:          req.Name,
		MerchantName:  req.MerchantName,
		Email:         req.Email,
		AccountNumber: req.AccountNumber,
		AppKey:        req.AppKey,
		SecretKey:     req.SecretKey,
		Balance:       req.Balance,
		Status:        req.Status,
	}

	if err := db.Save(&payoutProvider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "failed to update PayoutProvider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         true,
		"message":        "PayoutProvider Updated Successfully",
		"payoutProvider": payoutProvider,
	})

}
