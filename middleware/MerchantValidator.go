package middleware

import (
	"net/http"
	"payout-backend/config"
	"payout-backend/models"

	"github.com/gin-gonic/gin"
)

func MerchantValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		ipAddress := c.GetHeader("X-Forwarded-For")
		apiKey := c.GetHeader("X-API-KEY")
		secretKey := c.GetHeader("X-SECRET-KEY")
		ownerName := c.GetHeader("X-OWNER-NAME")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statue":     false,
				"code":       401,
				"message":    "Ip address required!",
				"ip-address": ipAddress,
			})
			c.Abort()
			return
		}

		var merchant models.Merchant
		err := config.DB.
			// Where("ip_address", ipAddress).
			Where("owner_name", ownerName).
			Where("status", 1).
			Where("is_deleted", 0).
			First(&merchant).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     false,
				"code":       401,
				"message":    "Merchant not exists or Ip not whitelisted!",
				"ip-address": ipAddress,
			})
			c.Abort()
			return
		}

		var account models.MerchantAccountDetail
		err = config.DB.
			Where("api_key", apiKey).
			Where("secret_key", secretKey).
			Where("merchant_id", merchant.ID).
			First(&account).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     false,
				"code":       403,
				"message":    "Invalid Header Value",
				"ip-address": ipAddress,
			})
			c.Abort()
			return
		}

		c.Set("merchant", merchant)
		c.Set("merchant_account", account)
	}

}
