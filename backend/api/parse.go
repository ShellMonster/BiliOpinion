package api

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ParseRequest 解析请求
// 前端发送用户输入的需求描述（自由格式）
type ParseRequest struct {
	Requirement string `json:"requirement" binding:"required"` // 用户需求描述，必填字段
	// 为了向后兼容，同时支持旧的 category 字段
	Category string `json:"category"` // 兼容旧版本，如果 requirement 为空则使用 category
}

// HandleParse 处理解析请求
// 这个API接收用户输入的需求描述（自由格式），调用AI进行解析，返回品牌、维度和关键词
func HandleParse(c *gin.Context) {
	// 1. 解析请求参数
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 向后兼容：如果 requirement 为空但 category 不为空，使用 category
	if req.Requirement == "" && req.Category != "" {
		req.Requirement = req.Category
	}

	// 验证至少有一个字段不为空
	if req.Requirement == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: requirement 或 category 字段不能为空"})
		return
	}

	// 2. 从数据库获取AI配置
	// 需要读取用户在设置页面配置的API Key、API Base和模型名称
	var apiKey, apiBase, model string
	database.DB.Model(&models.Settings{}).Where("key = ?", "ai_api_key").Pluck("value", &apiKey)
	database.DB.Model(&models.Settings{}).Where("key = ?", "ai_api_base").Pluck("value", &apiBase)
	database.DB.Model(&models.Settings{}).Where("key = ?", "ai_model").Pluck("value", &model)

	// 3. 验证AI配置是否完整
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AI API密钥未配置，请先在设置页面配置"})
		return
	}

	// 如果没有配置API Base，使用默认值
	if apiBase == "" {
		apiBase = "https://api.openai.com/v1"
	}

	// 如果没有配置模型，使用默认值
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	// 4. 创建AI客户端
	aiClient := ai.NewClient(ai.Config{
		APIBase: apiBase,
		APIKey:  apiKey,
		Model:   model,
	})

	// 5. 调用AI解析关键词
	result, err := aiClient.ParseKeyword(c.Request.Context(), ai.ParseKeywordRequest{
		Requirement: req.Requirement,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI解析失败: " + err.Error()})
		return
	}

	// 6. 返回解析结果
	c.JSON(http.StatusOK, result)
}
