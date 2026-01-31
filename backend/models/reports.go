package models

import (
	"time"
)

// Report 报告数据表
// 永久保存生成的分析报告，ReportData字段存储完整的JSON格式报告
type Report struct {
	ID         uint      `gorm:"primaryKey"`         // 主键ID
	HistoryID  uint      `gorm:"index;not null"`     // 关联的分析历史ID（外键引用analysis_history表）
	Category   string    `gorm:"index"`              // 商品类目（冗余字段，便于快速查询）
	ReportData string    `gorm:"type:text;not null"` // 报告JSON数据（完整报告内容）
	CreatedAt  time.Time `gorm:"index"`              // 创建时间（用于时间范围查询）
	UpdatedAt  time.Time // 更新时间
}

// ReportData JSON结构示例：
// {
//   "brands": ["戴森", "小米", "石头"],
//   "dimensions": [
//     {"name": "吸力", "description": "吸尘效果强弱"},
//     {"name": "续航", "description": "单次充电使用时长"}
//   ],
//   "scores": {
//     "戴森": {"吸力": 9.2, "续航": 7.5, "噪音": 6.8},
//     "小米": {"吸力": 8.1, "续航": 8.8, "噪音": 7.2}
//   },
//   "recommendation": "综合评价：戴森吸力最强但续航较弱，小米性价比高..."
// }
