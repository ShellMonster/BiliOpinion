package models

import (
	"time"
)

// RawComment 原始评论数据表
// 临时存储抓取的原始评论，3天后自动清理以节省存储空间
type RawComment struct {
	ID          uint      `gorm:"primaryKey"`     // 主键ID
	HistoryID   uint      `gorm:"index;not null"` // 关联的分析历史ID（外键引用analysis_history表）
	VideoID     string    `gorm:"index"`          // B站视频BV号（如：BV1xx411c7mD）
	CommentID   string    `gorm:"uniqueIndex"`    // 评论唯一ID（B站的rpid）
	Content     string    `gorm:"type:text"`      // 评论内容（完整文本）
	Author      string    `gorm:"index"`          // 评论作者昵称
	Likes       int       `gorm:"default:0"`      // 点赞数
	ReplyCount  int       `gorm:"default:0"`      // 回复数
	PublishTime time.Time `gorm:"index"`          // 评论发布时间（B站原始时间）
	CreatedAt   time.Time `gorm:"index"`          // 数据抓取时间（用于3天清理判断）
}

// 注意：此表数据会在程序启动时自动清理超过3天的记录
// 清理条件：created_at < NOW() - 3 days
