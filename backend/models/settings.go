package models

import (
	"time"
)

// Settings 系统配置表
// 使用Key-Value模式存储各类配置项，灵活扩展
type Settings struct {
	ID        uint      `gorm:"primaryKey"`           // 主键ID
	Key       string    `gorm:"uniqueIndex;not null"` // 配置键（唯一索引）
	Value     string    `gorm:"type:text"`            // 配置值（支持大文本）
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

// 常用配置键常量
const (
	SettingKeyAIAPIKey       = "ai_api_key"      // OpenAI API Key
	SettingKeyAIAPIBase      = "ai_api_base"     // API Base URL
	SettingKeyAIModel        = "ai_model"        // 模型名称
	SettingKeyBilibiliCookie = "bilibili_cookie" // B站完整Cookie字符串
)
