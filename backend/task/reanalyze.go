package task

import (
	"context"
	"encoding/json"
	"fmt"

	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/comment"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"bilibili-analyzer/backend/report"
)

func ReanalyzeFromStore(ctx context.Context, historyID uint) (uint, error) {
	var history models.AnalysisHistory
	if err := database.DB.First(&history, historyID).Error; err != nil {
		return 0, fmt.Errorf("history not found")
	}
	stored, err := comment.LoadByHistory(historyID)
	if err != nil {
		return 0, err
	}
	if len(stored) == 0 {
		return 0, fmt.Errorf("没有可重分析的已存评论")
	}

	settings, err := loadAISettingsOnly()
	if err != nil {
		return 0, err
	}
	aiClient := ai.NewClient(ai.Config{
		APIBase: settings.AIBaseURL,
		APIKey:  settings.AIAPIKey,
		Model:   settings.AIModel,
	})

	var dimNames []string
	_ = json.Unmarshal([]byte(history.Dimensions), &dimNames)
	dimensions := make([]ai.Dimension, 0, len(dimNames))
	for _, name := range dimNames {
		dimensions = append(dimensions, ai.Dimension{Name: name})
	}
	if len(dimensions) == 0 {
		return 0, fmt.Errorf("历史记录缺少评价维度")
	}

	inputs := make([]ai.CommentInput, 0, len(stored))
	for i, c := range stored {
		inputs = append(inputs, ai.CommentInput{
			ID:         fmt.Sprintf("stored_%d", i),
			Content:    c.Content,
			VideoBVID:  c.VideoID,
			VideoTitle: "",
		})
	}

	results, err := aiClient.AnalyzeCommentsWithRateLimit(ctx, inputs, dimensions, 5)
	if err != nil {
		return 0, err
	}

	analysisResults := make(map[string][]report.CommentWithScore)
	for _, r := range results {
		if r.Error != "" || r.Scores == nil {
			continue
		}
		brand := r.Brand
		if brand == "" {
			brand = "未知"
		}
		analysisResults[brand] = append(analysisResults[brand], report.CommentWithScore{
			Content: r.Content,
			Scores:  r.Scores,
			Brand:   brand,
			Model:   r.Model,
		})
	}

	var brands []string
	_ = json.Unmarshal([]byte(history.Brands), &brands)
	reportData, err := report.GenerateReportWithInput(report.GenerateReportInput{
		Category:        history.Category,
		Brands:          brands,
		Dimensions:      dimensions,
		AnalysisResults: analysisResults,
		Stats: report.ReportStats{
			TotalComments: len(stored),
		},
		Videos: []bilibili.VideoInfo{},
	})
	if err != nil {
		return 0, err
	}

	exec := NewExecutor(nil)
	reportID, err := exec.saveReport(history.ID, reportData)
	if err != nil {
		return 0, err
	}
	exec.updateHistoryWithReport(history.ID, reportID)
	return reportID, nil
}

func loadAISettingsOnly() (*AppSettings, error) {
	get := func(key string) string {
		var setting models.Settings
		if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			return ""
		}
		return setting.Value
	}
	settings := &AppSettings{
		AIBaseURL: get(models.SettingKeyAIAPIBase),
		AIAPIKey:  get(models.SettingKeyAIAPIKey),
		AIModel:   get(models.SettingKeyAIModel),
	}
	if settings.AIAPIKey == "" {
		return nil, fmt.Errorf("请先配置AI API Key")
	}
	return settings, nil
}
