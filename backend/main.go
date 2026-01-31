package main

import (
	"bilibili-analyzer/backend/api"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/sse"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	// 数据库文件将创建在项目根目录的 data/ 文件夹下
	if err := database.InitDB("data/bilibili-analyzer.db"); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	log.Println("🚀 Bilibili Analyzer - Backend Server Starting...")

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

		// SSE接口 - 前端通过此接口接收任务实时进度
		apiGroup.GET("/sse", sse.HandleSSE)
	}

	// 启动服务器
	log.Println("✅ Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
