package controllers

import (
	"net/http"
	"payout-backend/config"
	"payout-backend/models"

	"github.com/gin-gonic/gin"
)

func DashboardData(c *gin.Context) {
	db := config.DB
	role := c.GetUint("role")
	userID := c.GetUint("user_id")

	var (
		totalMerchant                 int64
		totalTransaction              int64
		totalSuccessTransaction       int64
		totalSuccessTransactionAmount int64
		totalFailedTransaction        int64
		totalGstCharge                float64
		totalFlatCharge               float64
		successRate                   float64
	)

	if role == 2 {
		totalMerchant = 0
		db.Model(&models.WalletTransaction{}).
			Where("merchant_id", userID).Count(&totalTransaction)

		db.Model(&models.WalletTransaction{}).
			Where("status=? AND merchant_id", 1, userID).Count(&totalSuccessTransaction)
		db.Model(&models.WalletTransaction{}).
			Where("status=? AND merchant_id", 3, userID).Count(&totalFailedTransaction)
		db.Model(&models.WalletTransaction{}).
			Where("status=? AND merchant_id", 1, userID).Select("COALESCE(SUM(gst_charge), 0)").Scan(&totalGstCharge)
		db.Model(&models.WalletTransaction{}).
			Where("status=? AND merchant_id", 1, userID).Select("COALESCE(SUM(amount), 0)").Scan(&totalGstCharge)
		db.Model(&models.WalletTransaction{}).
			Where("status=? AND merchant_id", 1, userID).Select("COALESCE(SUM(service_charge), 0)").Scan(&totalFlatCharge)
	} else {
		db.Model(&models.Merchant{}).
			Where("is_admin=? AND status", 2, 1).Count(&totalTransaction)
		db.Model(&models.WalletTransaction{}).
			Count(&totalTransaction)

		db.Model(&models.WalletTransaction{}).
			Where("status=", 1).Count(&totalSuccessTransaction)
		db.Model(&models.WalletTransaction{}).
			Where("status=", 3).Count(&totalFailedTransaction)
		db.Model(&models.WalletTransaction{}).
			Where("status=", 1).Select("COALESCE(SUM(gst_charge), 0)").Scan(&totalGstCharge)
		db.Model(&models.WalletTransaction{}).
			Where("status=", 1).Select("COALESCE(SUM(amount), 0)").Scan(&totalSuccessTransactionAmount)
		db.Model(&models.WalletTransaction{}).
			Where("status=", 1).Select("COALESCE(SUM(service_charge), 0)").Scan(&totalFlatCharge)

	}

	if totalTransaction > 0 {
		successRate = (float64(totalSuccessTransaction) / float64(totalTransaction)) * 100
	}
	c.JSON(http.StatusOK, gin.H{
		"totalMerchant":           totalMerchant,
		"totalTransaction":        totalTransaction,
		"totalSuccessTransaction": totalSuccessTransaction,
		"totalFailedTransaction":  totalFailedTransaction,
		"totalSuccessAmount":      totalSuccessTransactionAmount,
		"totalGstCharge":          totalGstCharge,
		"totalFlatCharge":         totalFlatCharge,
		"successRate":             successRate,
	})
}
