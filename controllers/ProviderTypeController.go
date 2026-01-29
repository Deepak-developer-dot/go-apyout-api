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

func AddProviderType(c *gin.Context) {
	db := config.DB

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	ProviderType := models.ProviderTypes{
		Name: req.Name,
	}

	if err := db.Create(&ProviderType).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Provider Type Not Create"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        true,
		"message":       "Provider Type Created Successfully!",
		"provider_type": ProviderType,
	})
}

func GetProviderType(c *gin.Context) {
	db := config.DB
	keywordSearch := c.Query("keyword_search")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	limit := 10
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("provider_type:keyword_search=%s:page=%d", keywordSearch, page)
	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	query := db.Model(&models.ProviderTypes{})

	if keywordSearch != "" {
		query = query.Where("name ILIKE ?", keywordSearch)
	}
	var providerType []models.ProviderTypes
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&providerType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch provider type"})
		return
	}

	response := gin.H{
		"message":      "Provider Type Fetch Successfully",
		"status":       true,
		"providerType": providerType,
		"page":         page,
		"limit":        limit,
	}

	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Second)
	c.JSON(http.StatusOK, response)
}

func GetProviderTypeByID(c *gin.Context) {
	id := c.Param("id")
	db := config.DB

	cacheKey := fmt.Sprintf("get_provider_type_by_id:%s", id)
	if val, found := config.Cache.Get(cacheKey); found {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	var providerType models.ProviderTypes
	if err := db.Where("id=?", id).First(&providerType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "providerType Not Found"})
		return
	}

	response := gin.H{
		"message":      "Provider Type Fetch Successfully",
		"status":       true,
		"providerType": providerType,
	}
	data, _ := json.Marshal(response)
	config.Cache.Set(cacheKey, string(data), 10*time.Minute)
	c.JSON(http.StatusOK, response)
}

func UpdateProviderType(c *gin.Context) {

	id := c.Param("id")
	db := config.DB

	providerTypeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	var providerType models.ProviderTypes
	if err := db.First(&providerType, providerTypeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "providertype not found"})
		return
	}

	providerType.Name = req.Name
	if err := db.Save(&providerType).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"message":      "Provider Type Update",
		"providerType": providerType,
	})
}
