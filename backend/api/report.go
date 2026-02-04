package api

import (
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"bilibili-analyzer/backend/pdf"
	"bilibili-analyzer/backend/report"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "报告ID不能为空"})
		return
	}

	var reportModel models.Report
	if err := database.DB.First(&reportModel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报告不存在"})
		return
	}

	var reportData map[string]interface{}
	if err := json.Unmarshal([]byte(reportModel.ReportData), &reportData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析报告数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         reportModel.ID,
		"history_id": reportModel.HistoryID,
		"category":   reportModel.Category,
		"data":       reportData,
		"created_at": reportModel.CreatedAt,
	})
}

func HandleExportPDF(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "报告ID不能为空"})
		return
	}

	var reportModel models.Report
	if err := database.DB.First(&reportModel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报告不存在"})
		return
	}

	var reportData report.ReportData
	if err := json.Unmarshal([]byte(reportModel.ReportData), &reportData); err != nil {
		log.Printf("[PDF] 解析报告数据失败: %v", err)
		dataPreview := reportModel.ReportData
		if len(dataPreview) > 500 {
			dataPreview = dataPreview[:500]
		}
		log.Printf("[PDF] 报告数据内容前500字符: %s", dataPreview)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析报告数据失败"})
		return
	}

	log.Printf("[PDF] 开始生成PDF，报告ID: %d, 类目: %s", reportModel.ID, reportData.Category)
	pdfBytes, err := pdf.GeneratePDF(&reportData, reportModel.ID)
	if err != nil {
		log.Printf("[PDF] 生成PDF失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成PDF失败: %v", err)})
		return
	}
	log.Printf("[PDF] PDF生成成功，大小: %d bytes", len(pdfBytes))

	filename := fmt.Sprintf("report_%d.pdf", reportModel.ID)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
