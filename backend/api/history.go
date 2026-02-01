package api

import (
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HistoryListResponse 历史记录列表响应结构
type HistoryListResponse struct {
	ID           uint   `json:"id"`           // 历史记录ID
	TaskId       string `json:"taskId"`       // 任务ID
	Category     string `json:"category"`     // 商品类目
	VideoCount   int    `json:"videoCount"`   // 视频数量
	CommentCount int    `json:"commentCount"` // 评论数量
	Status       string `json:"status"`       // 任务状态
	ReportID     uint   `json:"reportId"`     // 关联的报告ID
	CreatedAt    string `json:"createdAt"`    // 创建时间
}

// HistoryDetailResponse 历史记录详情响应结构
type HistoryDetailResponse struct {
	ID           uint     `json:"id"`           // 历史记录ID
	Category     string   `json:"category"`     // 商品类目
	Keywords     []string `json:"keywords"`     // 搜索关键词
	Brands       []string `json:"brands"`       // 品牌列表
	Dimensions   []string `json:"dimensions"`   // 评价维度
	VideoCount   int      `json:"videoCount"`   // 视频数量
	CommentCount int      `json:"commentCount"` // 评论数量
	Status       string   `json:"status"`       // 任务状态
	ReportData   string   `json:"reportData"`   // 报告JSON数据
	CreatedAt    string   `json:"createdAt"`    // 创建时间
}

// HandleGetHistory 获取历史记录列表
// GET /api/history
// 返回所有历史记录的简要信息，按创建时间倒序排列
func HandleGetHistory(c *gin.Context) {
	var histories []models.AnalysisHistory

	// 查询所有历史记录，按创建时间倒序排列
	if err := database.DB.Order("created_at DESC").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch history records",
		})
		return
	}

	// 转换为响应格式
	response := make([]HistoryListResponse, len(histories))
	for i, h := range histories {
		response[i] = HistoryListResponse{
			ID:           h.ID,
			TaskId:       h.TaskID,
			Category:     h.Category,
			VideoCount:   h.VideoCount,
			CommentCount: h.CommentCount,
			Status:       h.Status,
			ReportID:     h.ReportID,
			CreatedAt:    h.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// HandleGetHistoryDetail 获取历史记录详情
// GET /api/history/:id
// 返回指定ID的历史记录详细信息，包括关联的报告数据
func HandleGetHistoryDetail(c *gin.Context) {
	// 获取URL参数中的ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid history ID",
		})
		return
	}

	// 查询历史记录
	var history models.AnalysisHistory
	if err := database.DB.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "History record not found",
		})
		return
	}

	// 查询关联的报告数据
	var report models.Report
	reportData := ""
	if history.ReportID > 0 {
		if err := database.DB.First(&report, history.ReportID).Error; err == nil {
			reportData = report.ReportData
		}
	}

	// 解析JSON字段（简化处理，前端会进一步解析）
	// 这里直接返回原始JSON字符串，前端负责解析
	response := HistoryDetailResponse{
		ID:           history.ID,
		Category:     history.Category,
		Keywords:     parseJSONArray(history.Keywords),
		Brands:       parseJSONArray(history.Brands),
		Dimensions:   parseJSONArray(history.Dimensions),
		VideoCount:   history.VideoCount,
		CommentCount: history.CommentCount,
		Status:       history.Status,
		ReportData:   reportData,
		CreatedAt:    history.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// HandleDeleteHistory 删除历史记录
// DELETE /api/history/:id
// 删除指定ID的历史记录及其关联的报告数据
func HandleDeleteHistory(c *gin.Context) {
	// 获取URL参数中的ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid history ID",
		})
		return
	}

	// 查询历史记录
	var history models.AnalysisHistory
	if err := database.DB.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "History record not found",
		})
		return
	}

	// 删除关联的报告数据
	if history.ReportID > 0 {
		database.DB.Delete(&models.Report{}, history.ReportID)
	}

	// 删除历史记录
	if err := database.DB.Delete(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete history record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "History record deleted successfully",
	})
}

// parseJSONArray 解析JSON数组字符串为字符串切片
// 简化处理：如果解析失败，返回空切片
func parseJSONArray(jsonStr string) []string {
	if jsonStr == "" || jsonStr == "[]" {
		return []string{}
	}

	var result []string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return []string{}
	}
	return result
}
