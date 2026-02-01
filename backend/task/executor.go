package task

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"bilibili-analyzer/backend/report"
	"bilibili-analyzer/backend/sse"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// AppSettings 应用配置（从数据库读取后的结构化配置）
type AppSettings struct {
	AIBaseURL      string
	AIAPIKey       string
	AIModel        string
	BilibiliCookie string
}

// TaskConfig 任务配置
type TaskConfig struct {
	MaxVideosPerKeyword int // 每个关键词最大视频数（默认10）
	MaxCommentsPerVideo int // 每个视频最大评论数（默认100）
	MaxConcurrency      int // 最大并发数（默认3）
	AIBatchSize         int // AI分析批次大小（默认5）
}

// DefaultTaskConfig 默认任务配置
// 性能优化：增加并发数和批次大小以提升分析速度
func DefaultTaskConfig() TaskConfig {
	return TaskConfig{
		MaxVideosPerKeyword: 10,
		MaxCommentsPerVideo: 100,
		MaxConcurrency:      5,  // 从3增加到5，提升抓取速度
		AIBatchSize:         10, // 从5增加到10，减少AI API调用次数
	}
}

// TaskRequest 任务请求
type TaskRequest struct {
	TaskID      string         // 任务ID
	Requirement string         // 用户原始需求
	Brands      []string       // 品牌列表
	Dimensions  []ai.Dimension // 评价维度
	Keywords    []string       // 搜索关键词
}

// Executor 任务执行器
// 整合搜索、抓取、分析、报告生成的完整流程
type Executor struct {
	config TaskConfig
}

// NewExecutor 创建任务执行器
func NewExecutor(config *TaskConfig) *Executor {
	cfg := DefaultTaskConfig()
	if config != nil {
		cfg = *config
	}
	return &Executor{config: cfg}
}

// Execute 执行完整的分析任务
// 流程：搜索视频 -> 抓取评论 -> AI分析 -> 生成报告 -> 保存数据库
func (e *Executor) Execute(ctx context.Context, req TaskRequest) error {
	taskID := req.TaskID
	log.Printf("[Task %s] Starting execution...", taskID)

	// 阶段0：创建历史记录
	history, err := e.createHistory(req)
	if err != nil {
		sse.PushError(taskID, fmt.Sprintf("创建任务记录失败: %v", err))
		return err
	}
	log.Printf("[Task %s] History created: ID=%d", taskID, history.ID)

	// 阶段1：获取配置
	sse.PushProgress(taskID, sse.StatusSearching, 0, 100, "正在加载配置...")

	settings, err := e.loadSettings()
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("加载配置失败: %v", err))
		return err
	}

	// 阶段2：搜索视频
	sse.PushProgress(taskID, sse.StatusSearching, 5, 100, "正在搜索相关视频...")

	biliClient := bilibili.NewClient(settings.BilibiliCookie)
	allVideos, err := e.searchVideos(ctx, taskID, biliClient, req.Keywords)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("搜索视频失败: %v", err))
		return err
	}

	if len(allVideos) == 0 {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, "未找到相关视频，请尝试其他关键词")
		return fmt.Errorf("no videos found")
	}

	log.Printf("[Task %s] Found %d videos", taskID, len(allVideos))

	// 阶段3：抓取评论
	sse.PushProgress(taskID, sse.StatusScraping, 20, 100, fmt.Sprintf("开始抓取%d个视频的评论...", len(allVideos)))

	scraper := bilibili.NewScraper(biliClient, &bilibili.ScraperConfig{
		MaxVideos:           len(allVideos),
		MaxCommentsPerVideo: e.config.MaxCommentsPerVideo,
		MaxConcurrency:      int64(e.config.MaxConcurrency),
		FetchReplies:        true,
		RequestDelay:        200 * time.Millisecond,
	})

	// 设置进度回调
	scraper.SetProgressCallback(func(stage string, current, total int, message string) {
		progress := 20 + (current * 30 / max(total, 1)) // 20-50%
		sse.PushProgress(taskID, sse.StatusScraping, progress, 100, message)
	})

	scrapeResult, err := scraper.ScrapeByVideos(ctx, allVideos)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("抓取评论失败: %v", err))
		return err
	}

	log.Printf("[Task %s] Scraped %d comments from %d videos",
		taskID, scrapeResult.Stats.TotalComments, scrapeResult.Stats.TotalVideos)

	// 更新历史记录的统计信息
	e.updateHistoryStats(history.ID, scrapeResult.Stats.TotalVideos, scrapeResult.Stats.TotalComments)

	// 阶段4：AI分析评论
	sse.PushProgress(taskID, sse.StatusAnalyzing, 50, 100, "正在使用AI分析评论...")

	aiClient := ai.NewClient(ai.Config{
		APIBase: settings.AIBaseURL,
		APIKey:  settings.AIAPIKey,
		Model:   settings.AIModel,
	})
	analysisResults, err := e.analyzeComments(ctx, taskID, aiClient, scrapeResult, req.Brands, req.Dimensions)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("AI分析失败: %v", err))
		return err
	}

	log.Printf("[Task %s] Analysis completed for %d brands", taskID, len(analysisResults))

	// 阶段5：生成报告
	sse.PushProgress(taskID, sse.StatusGenerating, 85, 100, "正在生成分析报告...")

	// 构建统计数据
	commentsByBrand := make(map[string]int)
	for brand, results := range analysisResults {
		commentsByBrand[brand] = len(results)
	}

	reportInput := report.GenerateReportInput{
		Category:        req.Requirement,
		Brands:          req.Brands,
		Dimensions:      req.Dimensions,
		AnalysisResults: analysisResults,
		Stats: report.ReportStats{
			TotalVideos:     scrapeResult.Stats.TotalVideos,
			TotalComments:   scrapeResult.Stats.TotalComments,
			CommentsByBrand: commentsByBrand,
		},
	}

	reportData, err := report.GenerateReportWithInput(reportInput)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("生成报告失败: %v", err))
		return err
	}

	// 使用AI生成更专业的购买建议
	sse.PushProgress(taskID, sse.StatusGenerating, 90, 100, "正在生成AI购买建议...")
	aiRecommendation, err := e.generateAIRecommendation(ctx, aiClient, reportData)
	if err == nil && aiRecommendation != "" {
		reportData.Recommendation = aiRecommendation
	}

	// 阶段6：保存报告到数据库
	sse.PushProgress(taskID, sse.StatusGenerating, 95, 100, "正在保存报告...")

	reportID, err := e.saveReport(history.ID, reportData)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("保存报告失败: %v", err))
		return err
	}

	// 更新历史记录状态为完成
	e.updateHistoryWithReport(history.ID, reportID)

	log.Printf("[Task %s] Completed! Report ID: %d", taskID, reportID)

	// 推送完成状态（包含报告ID）
	sse.PushStatus(taskID, sse.TaskStatus{
		TaskID:  taskID,
		Status:  sse.StatusCompleted,
		Message: fmt.Sprintf("分析完成！共分析%d个视频，%d条评论", scrapeResult.Stats.TotalVideos, scrapeResult.Stats.TotalComments),
		Progress: &sse.Progress{
			Current: 100,
			Total:   100,
			Stage:   fmt.Sprintf("%d", reportID), // 使用Stage字段传递报告ID
		},
	})

	return nil
}

// createHistory 创建分析历史记录
func (e *Executor) createHistory(req TaskRequest) (*models.AnalysisHistory, error) {
	keywordsJSON, _ := json.Marshal(req.Keywords)
	brandsJSON, _ := json.Marshal(req.Brands)

	dimNames := make([]string, len(req.Dimensions))
	for i, d := range req.Dimensions {
		dimNames[i] = d.Name
	}
	dimensionsJSON, _ := json.Marshal(dimNames)

	history := &models.AnalysisHistory{
		Category:   req.Requirement,
		Keywords:   string(keywordsJSON),
		Brands:     string(brandsJSON),
		Dimensions: string(dimensionsJSON),
		Status:     models.StatusProcessing,
	}

	if err := database.DB.Create(history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

func (e *Executor) loadSettings() (*AppSettings, error) {
	getSettingValue := func(key string) string {
		var setting models.Settings
		if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			return ""
		}
		return setting.Value
	}

	settings := &AppSettings{
		AIBaseURL:      getSettingValue(models.SettingKeyAIAPIBase),
		AIAPIKey:       getSettingValue(models.SettingKeyAIAPIKey),
		AIModel:        getSettingValue(models.SettingKeyAIModel),
		BilibiliCookie: getSettingValue(models.SettingKeyBilibiliCookie),
	}

	if settings.AIAPIKey == "" {
		return nil, fmt.Errorf("请先配置AI API Key")
	}
	if settings.BilibiliCookie == "" {
		return nil, fmt.Errorf("请先配置B站Cookie")
	}

	return settings, nil
}

// searchVideos 搜索视频
func (e *Executor) searchVideos(ctx context.Context, taskID string, client *bilibili.Client, keywords []string) ([]bilibili.VideoInfo, error) {
	var allVideos []bilibili.VideoInfo
	videoMap := make(map[string]bool) // 用于去重

	for i, keyword := range keywords {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		progress := 5 + (i * 15 / max(len(keywords), 1))
		sse.PushProgress(taskID, sse.StatusSearching, progress, 100,
			fmt.Sprintf("正在搜索: %s (%d/%d)", keyword, i+1, len(keywords)))

		videos, err := client.SearchVideosWithLimit(keyword, e.config.MaxVideosPerKeyword)
		if err != nil {
			log.Printf("[Task %s] Search failed for keyword '%s': %v", taskID, keyword, err)
			continue // 单个关键词失败不影响整体
		}

		// 去重添加
		for _, v := range videos {
			if !videoMap[v.BVID] {
				videoMap[v.BVID] = true
				allVideos = append(allVideos, v)
			}
		}

		// 添加延迟避免请求过快
		time.Sleep(500 * time.Millisecond)
	}

	return allVideos, nil
}

// analyzeComments 分析评论
func (e *Executor) analyzeComments(
	ctx context.Context,
	taskID string,
	aiClient *ai.Client,
	scrapeResult *bilibili.ScrapeResult,
	brands []string,
	dimensions []ai.Dimension,
) (map[string][]report.CommentWithScore, error) {

	allComments := bilibili.GetAllCommentTexts(scrapeResult)
	if len(allComments) == 0 {
		return nil, fmt.Errorf("没有获取到任何评论")
	}

	log.Printf("[Task %s] Analyzing %d comments...", taskID, len(allComments))

	brandComments := make(map[string][]string)
	for _, comment := range allComments {
		commentLower := strings.ToLower(comment)
		for _, brand := range brands {
			brandLower := strings.ToLower(brand)
			if strings.Contains(commentLower, brandLower) {
				brandComments[brand] = append(brandComments[brand], comment)
			}
		}
	}

	results := make(map[string][]report.CommentWithScore)
	totalBrands := len(brands)

	for i, brand := range brands {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		comments := brandComments[brand]
		if len(comments) == 0 {
			log.Printf("[Task %s] No comments found for brand: %s", taskID, brand)
			continue
		}

		maxComments := 100
		if len(comments) > maxComments {
			comments = comments[:maxComments]
		}

		progress := 50 + (i * 35 / max(totalBrands, 1))
		sse.PushProgress(taskID, sse.StatusAnalyzing, progress, 100,
			fmt.Sprintf("正在分析 %s 的评论 (%d/%d品牌, %d条评论)", brand, i+1, totalBrands, len(comments)))

		inputs := make([]ai.CommentInput, len(comments))
		for j, c := range comments {
			inputs[j] = ai.CommentInput{
				ID:      fmt.Sprintf("%s_%d", brand, j),
				Content: c,
			}
		}

		analysisResults, err := aiClient.AnalyzeCommentsWithRateLimit(ctx, inputs, dimensions, e.config.AIBatchSize)
		if err != nil {
			log.Printf("[Task %s] Analysis failed for brand %s: %v", taskID, brand, err)
			continue
		}

		for _, r := range analysisResults {
			if r.Error == "" && r.Scores != nil {
				results[brand] = append(results[brand], report.CommentWithScore{
					Content: r.Content,
					Scores:  r.Scores,
				})
			}
		}

		log.Printf("[Task %s] Brand %s: analyzed %d comments, got %d valid results",
			taskID, brand, len(comments), len(results[brand]))

		time.Sleep(1 * time.Second)
	}

	return results, nil
}

// saveReport 保存报告到数据库
func (e *Executor) saveReport(historyID uint, reportData *report.ReportData) (uint, error) {
	data, err := json.Marshal(reportData)
	if err != nil {
		return 0, err
	}

	reportRecord := &models.Report{
		HistoryID:  historyID,
		Category:   reportData.Category,
		ReportData: string(data),
	}

	if err := database.DB.Create(reportRecord).Error; err != nil {
		return 0, err
	}

	return reportRecord.ID, nil
}

// updateHistoryStatus 更新历史记录状态
func (e *Executor) updateHistoryStatus(historyID uint, status string) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Update("status", status)
}

// updateHistoryStats 更新历史记录统计信息
func (e *Executor) updateHistoryStats(historyID uint, videoCount, commentCount int) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"video_count":   videoCount,
		"comment_count": commentCount,
	})
}

// updateHistoryWithReport 更新历史记录的报告ID和状态
func (e *Executor) updateHistoryWithReport(historyID uint, reportID uint) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"report_id": reportID,
		"status":    models.StatusCompleted,
	})
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (e *Executor) generateAIRecommendation(ctx context.Context, aiClient *ai.Client, reportData *report.ReportData) (string, error) {
	rankings := make([]ai.BrandRankingInfo, len(reportData.Rankings))
	for i, r := range reportData.Rankings {
		rankings[i] = ai.BrandRankingInfo{
			Brand:        r.Brand,
			OverallScore: r.OverallScore,
			Rank:         r.Rank,
		}
	}

	brandAnalysis := make(map[string]ai.BrandStrengthWeakness)
	for brand, analysis := range reportData.BrandAnalysis {
		brandAnalysis[brand] = ai.BrandStrengthWeakness{
			Strengths:  analysis.Strengths,
			Weaknesses: analysis.Weaknesses,
		}
	}

	return aiClient.GenerateRecommendation(ctx, ai.RecommendationInput{
		Category:      reportData.Category,
		Rankings:      rankings,
		BrandAnalysis: brandAnalysis,
	})
}
