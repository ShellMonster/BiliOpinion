package database

import (
	"bilibili-analyzer/backend/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接和表结构
// 参数：
//   - dbPath: 数据库文件路径（如：data/bilibili-analyzer.db）
//
// 返回：
//   - error: 初始化失败时返回错误信息
func InitDB(dbPath string) error {
	var err error

	// 打开SQLite数据库连接
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移表结构（如果表不存在则创建，如果字段有变化则更新）
	// 迁移顺序：先迁移基础表，再迁移有外键关联的表
	err = DB.AutoMigrate(
		&models.Settings{},        // 系统配置表
		&models.AnalysisHistory{}, // 分析历史记录表
		&models.Report{},          // 报告数据表
		&models.RawComment{},      // 原始评论临时表
	)
	if err != nil {
		return err
	}

	log.Println("✅ Database initialized successfully")

	// 启动时清理3天前的临时数据
	// 注意：清理失败不影响程序启动，只记录警告日志
	if err := CleanOldComments(); err != nil {
		log.Printf("⚠️  Warning: Failed to clean old comments: %v", err)
	}

	return nil
}

// CleanOldComments 清理3天前的原始评论数据
// 此函数在程序启动时自动调用，用于节省存储空间
// 清理规则：删除 created_at < NOW() - 3 days 的所有 raw_comments 记录
// 返回：
//   - error: 清理失败时返回错误信息
func CleanOldComments() error {
	// 计算3天前的时间点
	threeDaysAgo := time.Now().AddDate(0, 0, -3)

	// 执行删除操作
	result := DB.Where("created_at < ?", threeDaysAgo).Delete(&models.RawComment{})
	if result.Error != nil {
		return result.Error
	}

	// 记录清理结果
	if result.RowsAffected > 0 {
		log.Printf("🗑️  Cleaned %d old comments (older than 3 days)", result.RowsAffected)
	} else {
		log.Println("✅ No old comments to clean")
	}

	return nil
}
