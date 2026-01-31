package models

import (
	"time"
)

// AnalysisHistory 分析历史记录表
// 永久保存所有分析任务的历史记录，用于历史查询和追溯
type AnalysisHistory struct {
	ID           uint      `gorm:"primaryKey"`              // 主键ID
	Category     string    `gorm:"index;not null"`          // 商品类目（如：吸尘器、手机）
	Keywords     string    `gorm:"type:text"`               // 搜索关键词JSON数组（如：["吸尘器评测","吸尘器推荐"]）
	Brands       string    `gorm:"type:text"`               // 品牌列表JSON数组（如：["戴森","小米","石头"]）
	Dimensions   string    `gorm:"type:text"`               // 评价维度JSON数组（如：["吸力","续航","噪音"]）
	VideoCount   int       `gorm:"default:0"`               // 抓取的视频数量
	CommentCount int       `gorm:"default:0"`               // 抓取的评论数量
	Status       string    `gorm:"index;default:'pending'"` // 任务状态：pending/processing/completed/failed
	ReportID     uint      `gorm:"index"`                   // 关联的报告ID（外键引用reports表）
	CreatedAt    time.Time `gorm:"index"`                   // 创建时间（用于时间范围查询）
	UpdatedAt    time.Time // 更新时间
}

// 任务状态常量
const (
	StatusPending    = "pending"    // 待处理
	StatusProcessing = "processing" // 处理中
	StatusCompleted  = "completed"  // 已完成
	StatusFailed     = "failed"     // 失败
)
