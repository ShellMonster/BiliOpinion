package api

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/comment"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"bilibili-analyzer/backend/report"
	"bilibili-analyzer/backend/sse"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VideoParseRequest 视频解析请求
// 用于解析单个视频URL，获取视频信息
type VideoParseRequest struct {
	VideoURL string `json:"video_url" binding:"required"` // B站视频链接
}

// VideoParseResponse 视频解析响应
// 返回视频的详细信息
type VideoParseResponse struct {
	BVID         string `json:"bvid"`          // 视频BV号
	Title        string `json:"title"`         // 视频标题
	Author       string `json:"author"`        // UP主名称
	PlayCount    int    `json:"play_count"`    // 播放量
	CommentCount int    `json:"comment_count"` // 评论数
	PubDate      string `json:"pub_date"`      // 发布时间
	Cover        string `json:"cover"`         // 封面图URL
}

// VideoAnalyzeRequest 视频分析请求
// 用于启动单个视频的评论分析任务
type VideoAnalyzeRequest struct {
	VideoURL    string `json:"video_url" binding:"required"` // B站视频链接
	MaxComments int    `json:"max_comments,omitempty"`       // 最大分析评论数，默认1000
}

// VideoAnalyzeResponse 视频分析响应
// 返回任务ID，前端通过SSE获取实时进度
type VideoAnalyzeResponse struct {
	TaskID  string `json:"task_id"` // 任务ID
	Message string `json:"message"` // 提示信息
}

// HandleVideoParse 处理视频URL解析
// POST /api/video/parse
// 解析视频URL，提取BV号并获取视频详细信息
//
// 请求示例：
//
//	{"video_url": "https://www.bilibili.com/video/BV1mH4y1u7UA"}
//
// 响应示例：
//
//	{"bvid": "BV1mH4y1u7UA", "title": "...", "author": "...", "play_count": 1000, ...}
func HandleVideoParse(c *gin.Context) {
	var req VideoParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证视频URL不为空
	if req.VideoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "视频链接不能为空"})
		return
	}

	// 解析视频URL，提取BV号
	bvid, err := bilibili.ParseVideoURL(req.VideoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取B站Cookie配置（用于访问需要登录的接口）
	cookie := getBilibiliCookie()

	// 创建B站客户端并获取视频信息
	client := bilibili.NewClient(cookie)
	videoInfo, err := client.GetVideoInfo(bvid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("获取视频信息失败: %v", err)})
		return
	}

	// 构建响应
	response := VideoParseResponse{
		BVID:         videoInfo.BVID,
		Title:        videoInfo.Title,
		Author:       videoInfo.Author,
		PlayCount:    videoInfo.PlayCount,
		CommentCount: videoInfo.CommentCount,
		PubDate:      videoInfo.PubDate,
		Cover:        videoInfo.Cover,
	}

	c.JSON(http.StatusOK, response)
}

// HandleVideoAnalyze 处理视频评论分析
// POST /api/video/analyze
// 启动单个视频的评论分析任务，通过SSE推送实时进度
//
// 请求示例：
//
//	{"video_url": "https://www.bilibili.com/video/BV1mH4y1u7UA", "max_comments": 1000}
//
// 响应示例：
//
//	{"task_id": "xxx-xxx-xxx", "message": "任务已创建"}
func HandleVideoAnalyze(c *gin.Context) {
	var req VideoAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证视频URL不为空
	if req.VideoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "视频链接不能为空"})
		return
	}

	// 设置默认最大评论数
	maxComments := req.MaxComments
	if maxComments <= 0 {
		maxComments = 1000 // 默认分析1000条评论
	}

	// 生成任务ID
	taskID := uuid.New().String()

	// 创建SSE任务通道
	sse.CreateTaskChannel(taskID)

	// 异步启动视频分析任务
	go executeVideoAnalyzeTask(taskID, req.VideoURL, maxComments)

	// 立即返回任务ID
	c.JSON(http.StatusOK, VideoAnalyzeResponse{
		TaskID:  taskID,
		Message: "任务已创建，请通过SSE接口获取实时进度",
	})
}

// executeVideoAnalyzeTask 执行视频分析任务
// 在后台goroutine中运行，完成以下步骤：
// 1. 解析视频URL获取BV号
// 2. 获取视频详细信息
// 3. 抓取视频评论
// 4. AI分析评论
// 5. 生成分析报告
func executeVideoAnalyzeTask(taskID, videoURL string, maxComments int) {
	// 确保任务结束时关闭SSE通道
	defer sse.CloseTaskChannel(taskID)

	// 推送初始状态
	sse.PushProgress(taskID, sse.StatusParsing, 0, 100, "正在解析视频链接...")

	// 步骤1：解析视频URL
	bvid, err := bilibili.ParseVideoURL(videoURL)
	if err != nil {
		sse.PushError(taskID, fmt.Sprintf("解析视频链接失败: %v", err))
		return
	}

	// 推送进度：正在获取视频信息
	sse.PushProgress(taskID, sse.StatusSearching, 5, 100, "正在获取视频信息...")

	// 步骤2：获取配置和创建客户端
	settings, err := loadTaskSettings()
	if err != nil {
		sse.PushError(taskID, err.Error())
		return
	}

	// 创建B站客户端
	biliClient := bilibili.NewClient(settings.BilibiliCookie)

	// 获取视频详细信息
	videoInfo, err := biliClient.GetVideoInfo(bvid)
	if err != nil {
		sse.PushError(taskID, fmt.Sprintf("获取视频信息失败: %v", err))
		return
	}

	log.Printf("[Task %s] Video info: %s, comments: %d", taskID, videoInfo.Title, videoInfo.CommentCount)

	// 步骤3：创建历史记录
	history, err := createVideoAnalyzeHistory(taskID, videoInfo, maxComments)
	if err != nil {
		sse.PushError(taskID, fmt.Sprintf("创建任务记录失败: %v", err))
		return
	}

	// 推送进度：正在抓取评论
	sse.PushProgress(taskID, sse.StatusScraping, 10, 100,
		fmt.Sprintf("正在抓取评论（预计 %d 条）...", min(maxComments, videoInfo.CommentCount)))

	// 步骤4：抓取评论
	scraper := bilibili.NewScraper(biliClient, &bilibili.ScraperConfig{
		MaxVideos:           1,
		MaxCommentsPerVideo: maxComments,
		MaxConcurrency:      1, // 单个视频，不需要并发
		FetchReplies:        true,
		RequestDelay:        200 * time.Millisecond,
	})

	// 构建视频信息用于抓取
	videos := []bilibili.VideoInfo{{
		BVID:        videoInfo.BVID,
		Title:       videoInfo.Title,
		Author:      videoInfo.Author,
		VideoReview: videoInfo.CommentCount,
		Pic:         videoInfo.Cover,
	}}

	// 设置评论分配
	commentAllocation := map[string]int{
		videoInfo.BVID: min(maxComments, videoInfo.CommentCount),
	}

	scrapeResult, err := scraper.ScrapeByVideos(context.Background(), videos, commentAllocation)
	if err != nil {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("抓取评论失败: %v", err))
		return
	}

	actualCommentCount := scrapeResult.Stats.TotalComments
	log.Printf("[Task %s] Scraped %d comments", taskID, actualCommentCount)

	// 更新历史记录统计
	updateHistoryStats(history.ID, 1, actualCommentCount)

	if actualCommentCount == 0 {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, "该视频没有评论可分析")
		return
	}

	// 推送进度：正在AI分析
	sse.PushProgress(taskID, sse.StatusAnalyzing, 40, 100,
		fmt.Sprintf("正在AI分析 %d 条评论...", actualCommentCount))

	// 步骤5：创建AI客户端
	aiClient := ai.NewClient(ai.Config{
		APIBase: settings.AIBaseURL,
		APIKey:  settings.AIAPIKey,
		Model:   settings.AIModel,
	})

	// 设置AI进度回调
	aiClient.SetProgressCallback(func(stage string, current, total int, message string) {
		progress := 40 + (current * 40 / max(total, 1)) // 40-80%
		sse.PushProgress(taskID, sse.StatusAnalyzing, progress, 100, message)
	})

	// 使用固定的6个默认评价维度
	dimensions := getDefaultDimensions()

	// 准备评论数据用于AI分析
	comments := getAllCommentsWithVideo(scrapeResult)

	// 过滤短评论
	var inputs []ai.CommentInput
	for i, c := range comments {
		if len(strings.TrimSpace(c.Content)) < 5 {
			continue
		}
		inputs = append(inputs, ai.CommentInput{
			ID:         fmt.Sprintf("comment_%d", i),
			Content:    c.Content,
			VideoTitle: c.VideoTitle,
			VideoBVID:  c.VideoBVID,
		})
	}

	if len(inputs) == 0 {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, "没有有效的评论可分析")
		return
	}

	// 截断到最大评论数
	if len(inputs) > maxComments {
		inputs = inputs[:maxComments]
	}

	// 执行AI分析
	analysisResults, err := aiClient.AnalyzeCommentsWithRateLimit(context.Background(), inputs, dimensions, 10)
	if err != nil {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("AI分析失败: %v", err))
		return
	}

	// 处理分析结果，按品牌分组
	resultsByBrand := processAnalysisResults(analysisResults)

	log.Printf("[Task %s] Analysis completed for %d brands", taskID, len(resultsByBrand))

	// 推送进度：正在生成报告
	sse.PushProgress(taskID, sse.StatusGenerating, 85, 100, "正在生成分析报告...")

	// 步骤6：生成报告
	reportInput := report.GenerateReportInput{
		Category:        videoInfo.Title, // 使用视频标题作为类目
		Brands:          getBrandList(resultsByBrand),
		Dimensions:      dimensions,
		AnalysisResults: resultsByBrand,
		Stats: report.ReportStats{
			TotalVideos:     1,
			TotalComments:   actualCommentCount,
			CommentsByBrand: getCommentsByBrand(resultsByBrand),
		},
		Videos: videos,
	}

	reportData, err := report.GenerateReportWithInput(reportInput)
	if err != nil {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("生成报告失败: %v", err))
		return
	}

	// 推送进度：正在生成购买建议
	sse.PushProgress(taskID, sse.StatusGenerating, 90, 100, "正在生成AI购买建议...")

	// 生成AI购买建议
	aiRecommendation, err := generateAIRecommendation(context.Background(), aiClient, reportData)
	if err == nil && aiRecommendation != "" {
		reportData.Recommendation = aiRecommendation
	}

	// 推送进度：正在保存报告
	sse.PushProgress(taskID, sse.StatusGenerating, 95, 100, "正在保存报告...")

	// 步骤7：保存报告到数据库
	reportID, err := saveReport(history.ID, reportData)
	if err != nil {
		updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("保存报告失败: %v", err))
		return
	}

	// 更新历史记录状态为完成
	updateHistoryWithReport(history.ID, reportID)

	log.Printf("[Task %s] Completed! Report ID: %d", taskID, reportID)

	// 推送完成状态（包含报告ID）
	sse.PushStatus(taskID, sse.TaskStatus{
		TaskID:  taskID,
		Status:  sse.StatusCompleted,
		Message: fmt.Sprintf("分析完成！共分析1个视频，%d条评论", actualCommentCount),
		Progress: &sse.Progress{
			Current: 100,
			Total:   100,
			Stage:   fmt.Sprintf("%d", reportID), // 使用Stage字段传递报告ID
		},
	})
}

// getBilibiliCookie 获取B站Cookie配置
func getBilibiliCookie() string {
	var setting models.Settings
	if err := database.DB.Where("key = ?", models.SettingKeyBilibiliCookie).First(&setting).Error; err != nil {
		return ""
	}
	return setting.Value
}

// loadTaskSettings 加载任务配置
func loadTaskSettings() (*taskSettings, error) {
	getSettingValue := func(key string) string {
		var setting models.Settings
		if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			return ""
		}
		return setting.Value
	}

	settings := &taskSettings{
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

// taskSettings 任务配置
type taskSettings struct {
	AIBaseURL      string
	AIAPIKey       string
	AIModel        string
	BilibiliCookie string
}

// createVideoAnalyzeHistory 创建视频分析历史记录
func createVideoAnalyzeHistory(taskID string, videoInfo *bilibili.VideoDetail, maxComments int) (*models.AnalysisHistory, error) {
	// 构建任务配置
	configJSON, _ := json.Marshal(map[string]interface{}{
		"max_comments": maxComments,
	})

	history := &models.AnalysisHistory{
		TaskID:        taskID,
		Category:      videoInfo.Title,
		Keywords:      "[]",
		Brands:        "[]",
		Dimensions:    "[]",
		Status:        models.StatusProcessing,
		Stage:         "initializing",
		Progress:      0,
		ProgressMsg:   "任务初始化中...",
		TaskConfig:    string(configJSON),
		LastHeartbeat: time.Now(),
	}

	if err := database.DB.Create(history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// updateHistoryStatus 更新历史记录状态
func updateHistoryStatus(historyID uint, status string) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Update("status", status)
}

// updateHistoryStats 更新历史记录统计信息
func updateHistoryStats(historyID uint, videoCount, commentCount int) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"video_count":   videoCount,
		"comment_count": commentCount,
	})
}

// updateHistoryWithReport 更新历史记录的报告ID和状态
func updateHistoryWithReport(historyID uint, reportID uint) {
	database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"report_id": reportID,
		"status":    models.StatusCompleted,
	})
}

// saveReport 保存报告到数据库
func saveReport(historyID uint, reportData *report.ReportData) (uint, error) {
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

// getDefaultDimensions 获取默认的6个评价维度
func getDefaultDimensions() []ai.Dimension {
	return []ai.Dimension{
		{Name: "性能表现", Description: "产品的核心性能、功能实现程度"},
		{Name: "质量做工", Description: "产品质量、材料、做工精细度"},
		{Name: "性价比", Description: "价格与性能的匹配程度，是否物有所值"},
		{Name: "使用体验", Description: "日常使用感受、操作便捷性"},
		{Name: "外观设计", Description: "产品外观、颜值、设计感"},
		{Name: "售后服务", Description: "售后服务质量、客服响应、保修政策"},
	}
}

// getAllCommentsWithVideo 获取所有评论（带视频信息）
type commentWithVideo struct {
	Content    string
	VideoTitle string
	VideoBVID  string
}

func getAllCommentsWithVideo(result *bilibili.ScrapeResult) []commentWithVideo {
	var comments []commentWithVideo

	videoTitleMap := make(map[string]string)
	for _, video := range result.Videos {
		videoTitleMap[video.BVID] = video.Title
	}

	for bvid, videoComments := range result.Comments {
		videoTitle := videoTitleMap[bvid]
		for _, c := range videoComments {
			comments = append(comments, commentWithVideo{
				Content:    c.Content.Message,
				VideoTitle: videoTitle,
				VideoBVID:  bvid,
			})
			// 添加回复
			for _, r := range c.Replies {
				comments = append(comments, commentWithVideo{
					Content:    r.Content.Message,
					VideoTitle: videoTitle,
					VideoBVID:  bvid,
				})
			}
		}
	}

	return comments
}

// processAnalysisResults 处理AI分析结果，按品牌分组
func processAnalysisResults(results []ai.CommentAnalysisResult) map[string][]report.CommentWithScore {
	brandResults := make(map[string][]report.CommentWithScore)

	for _, r := range results {
		if r.Error != "" || r.Scores == nil {
			continue
		}

		brand := strings.TrimSpace(r.Brand)
		if brand == "" || brand == "未知" {
			brand = "未知品牌"
		}

		// 清洗品牌和型号
		brand = comment.CleanBrandName(brand, nil)
		model := comment.CleanModelName(r.Model)

		commentItem := report.CommentWithScore{
			Content:     r.Content,
			Scores:      r.Scores,
			Brand:       brand,
			Model:       model,
			PublishTime: time.Time{},
		}

		brandResults[brand] = append(brandResults[brand], commentItem)
	}

	return brandResults
}

// getBrandList 获取品牌列表
func getBrandList(results map[string][]report.CommentWithScore) []string {
	brands := make([]string, 0, len(results))
	for brand := range results {
		brands = append(brands, brand)
	}
	return brands
}

// getCommentsByBrand 获取每个品牌的评论数
func getCommentsByBrand(results map[string][]report.CommentWithScore) map[string]int {
	counts := make(map[string]int)
	for brand, comments := range results {
		counts[brand] = len(comments)
	}
	return counts
}

// generateAIRecommendation 生成AI购买建议
func generateAIRecommendation(ctx context.Context, aiClient *ai.Client, reportData *report.ReportData) (string, error) {
	rankings := make([]ai.BrandRankingInfo, len(reportData.Rankings))
	for i, r := range reportData.Rankings {
		rankings[i] = ai.BrandRankingInfo{
			Brand:        r.Brand,
			OverallScore: r.OverallScore,
			Rank:         r.Rank,
		}
	}

	modelRankings := make([]ai.ModelRankingInfo, len(reportData.ModelRankings))
	for i, m := range reportData.ModelRankings {
		modelRankings[i] = ai.ModelRankingInfo{
			Model:        m.Model,
			Brand:        m.Brand,
			OverallScore: m.OverallScore,
			Rank:         m.Rank,
			CommentCount: m.CommentCount,
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
		ModelRankings: modelRankings,
	})
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// parseIntSafely 安全解析整数
func parseIntSafely(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}
