package api

import (
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetConfig(c *gin.Context) {
	getSettingValue := func(key string) string {
		var setting models.Settings
		if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			return ""
		}
		return setting.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"ai_base_url":            getSettingValue(models.SettingKeyAIAPIBase),
		"ai_api_key":             getSettingValue(models.SettingKeyAIAPIKey),
		"ai_model":               getSettingValue(models.SettingKeyAIModel),
		"bilibili_cookie":        getSettingValue(models.SettingKeyBilibiliCookie),
		"scrape_max_concurrency": getSettingValue(models.SettingKeyScrapeMaxConcurrency),
		"ai_max_concurrency":     getSettingValue(models.SettingKeyAIMaxConcurrency),
	})
}

func HandleSaveConfig(c *gin.Context) {
	var req struct {
		AIBaseURL            string `json:"ai_base_url"`
		AIAPIKey             string `json:"ai_api_key"`
		AIModel              string `json:"ai_model"`
		BilibiliCookie       string `json:"bilibili_cookie"`
		ScrapeMaxConcurrency string `json:"scrape_max_concurrency"`
		AIMaxConcurrency     string `json:"ai_max_concurrency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	saveOrUpdate := func(key, value string) error {
		var setting models.Settings
		result := database.DB.Where("key = ?", key).First(&setting)
		if result.Error != nil {
			setting = models.Settings{Key: key, Value: value}
			return database.DB.Create(&setting).Error
		}
		setting.Value = value
		return database.DB.Save(&setting).Error
	}

	if err := saveOrUpdate(models.SettingKeyAIAPIBase, req.AIBaseURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	if err := saveOrUpdate(models.SettingKeyAIAPIKey, req.AIAPIKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	if err := saveOrUpdate(models.SettingKeyAIModel, req.AIModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	if err := saveOrUpdate(models.SettingKeyBilibiliCookie, req.BilibiliCookie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	if err := saveOrUpdate(models.SettingKeyScrapeMaxConcurrency, req.ScrapeMaxConcurrency); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	if err := saveOrUpdate(models.SettingKeyAIMaxConcurrency, req.AIMaxConcurrency); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config saved successfully"})
}
