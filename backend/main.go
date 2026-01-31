package main

import (
	"bilibili-analyzer/backend/database"
	"log"
)

func main() {
	// 初始化数据库
	// 数据库文件将创建在项目根目录的 data/ 文件夹下
	if err := database.InitDB("data/bilibili-analyzer.db"); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	log.Println("🚀 Bilibili Analyzer - Backend Server Starting...")
	// TODO: 启动Gin服务器
}
