package main

import (
	"bilibili-analyzer/backend/api"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/sse"
	"bilibili-analyzer/backend/task"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	// 数据库文件将创建在项目根目录的 data/ 文件夹下
	if err := database.InitDB("data/bilibili-analyzer.db"); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	log.Println("🚀 Bilibili Analyzer - Backend Server Starting...")

	// 恢复未完成的任务（后端重启后）
	go task.RecoverIncompleteTasks()

	// 启动定时清理任务（每5分钟检查一次超时任务）
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			task.CleanupTimedOutTasks()
		}
	}()

	// 创建Gin路由器
	r := gin.Default()

	// 配置CORS（允许前端跨域访问）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 注册API路由
	apiGroup := r.Group("/api")
	{
		// 解析API - 用户输入商品类目，AI解析返回品牌、维度、关键词
		apiGroup.POST("/parse", api.HandleParse)

		// 确认API - 用户确认AI解析结果，启动抓取和分析任务
		apiGroup.POST("/confirm", api.HandleConfirm)

		// 视频API - 解析视频URL和分析单个视频评论
		apiGroup.POST("/video/parse", api.HandleVideoParse)
		apiGroup.POST("/video/analyze", api.HandleVideoAnalyze)

		// SSE接口 - 前端通过此接口接收任务实时进度
		apiGroup.GET("/sse", sse.HandleSSE)

		// 历史记录API
		apiGroup.GET("/history", api.HandleGetHistory)           // 获取历史记录列表
		apiGroup.GET("/history/:id", api.HandleGetHistoryDetail) // 获取历史记录详情
		apiGroup.DELETE("/history/:id", api.HandleDeleteHistory) // 删除历史记录

		// 报告API
		apiGroup.GET("/report/:id", api.HandleGetReport)     // 获取报告详情
		apiGroup.GET("/report/:id/pdf", api.HandleExportPDF) // 导出PDF

		// 配置API
		apiGroup.GET("/config", api.HandleGetConfig)   // 获取配置
		apiGroup.POST("/config", api.HandleSaveConfig) // 保存配置
	}

	// 启动服务器
	log.Println("✅ Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
