package task

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
	"regexp"
	"strings"
	"time"
)

var modelPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(iPhone|Galaxy|Pixel|Mate|Mi|Redmi|V|G|X|S)\s*(\d+)\s*(Pro|Max|Plus|Ultra|Detect|Slim)(\s+(Pro|Max|Plus|Ultra))?`),
	regexp.MustCompile(`(?i)\b(Pura\s+(X|Max)|T[1-4]S?|SCOOPER(\s+SE)?|Young|M1(\s+Pro)?|T\s+Air)\b`),
	regexp.MustCompile(`(?i)([A-Z]+)(\d+)\s*(Pro|Max|Plus|Ultra|Detect|Slim)`),
	regexp.MustCompile(`(?i)\b([A-Z]+)(\d+)\b`),
	regexp.MustCompile(`(?i)\s(Pro|Max|Plus|Ultra)\s`),
}

// AppSettings 应用配置（从数据库读取后的结构化配置）
type AppSettings struct {
	AIBaseURL      string
	AIAPIKey       string
	AIModel        string
	BilibiliCookie string
}

// TaskConfig 任务配置
type TaskConfig struct {
	MaxVideosPerKeyword   int // 每个关键词最大视频数（默认10）
	MaxCommentsPerVideo   int // 每个视频最大评论数（默认100）
	MaxConcurrency        int // 最大并发数（默认3）
	AIBatchSize           int // AI分析批次大小（默认5）
	VideoDateRangeMonths  int // 视频时间范围（月），0表示不限制，默认24（2年）
	MinVideoDuration      int // 最小视频时长（秒），0表示不过滤
	MaxComments           int // 最大分析评论数（默认500）
	MinVideoComments      int // 最小视频评论数过滤（默认0，表示不限制）
	MinCommentsPerVideo   int // 每视频最少抓取数（默认10）
	MaxCommentsPerVideoV2 int // 每视频最多抓取数（默认200）
}

// DefaultTaskConfig 默认任务配置
// 性能优化：增加并发数和批次大小以提升分析速度
func DefaultTaskConfig() TaskConfig {
	return TaskConfig{
		MaxVideosPerKeyword:   20,
		MaxCommentsPerVideo:   200,
		MaxConcurrency:        5,   // 从3增加到5，提升抓取速度
		AIBatchSize:           10,  // 从5增加到10，减少AI API调用次数
		VideoDateRangeMonths:  0,   // 默认不限时间
		MinVideoDuration:      30,  // 默认过滤30秒以下短视频
		MaxComments:           500, // 默认分析500条评论
		MinVideoComments:      0,   // 默认不限制视频评论数
		MinCommentsPerVideo:   10,  // 默认每视频最少抓取10条
		MaxCommentsPerVideoV2: 200, // 默认每视频最多抓取200条
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

// CommentWithVideo 带视频信息的评论
type CommentWithVideo struct {
	Content    string // 评论内容
	VideoTitle string // 视频标题
	VideoBVID  string // 视频BVID
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
		// 只覆盖非零值的字段，保留默认值
		if config.MaxVideosPerKeyword > 0 {
			cfg.MaxVideosPerKeyword = config.MaxVideosPerKeyword
		}
		if config.MaxCommentsPerVideo > 0 {
			cfg.MaxCommentsPerVideo = config.MaxCommentsPerVideo
		}
		if config.MaxConcurrency > 0 {
			cfg.MaxConcurrency = config.MaxConcurrency
		}
		if config.AIBatchSize > 0 {
			cfg.AIBatchSize = config.AIBatchSize
		}
		// VideoDateRangeMonths 和 MinVideoDuration 可以是 0（表示不限制）
		cfg.VideoDateRangeMonths = config.VideoDateRangeMonths
		cfg.MinVideoDuration = config.MinVideoDuration
		if config.MaxComments > 0 {
			cfg.MaxComments = config.MaxComments
		}
		// MinVideoComments 可以是 0（表示不限制）
		cfg.MinVideoComments = config.MinVideoComments
		if config.MinCommentsPerVideo > 0 {
			cfg.MinCommentsPerVideo = config.MinCommentsPerVideo
		}
		if config.MaxCommentsPerVideoV2 > 0 {
			cfg.MaxCommentsPerVideoV2 = config.MaxCommentsPerVideoV2
		}
	}
	return &Executor{config: cfg}
}

// Execute 执行完整的分析任务
// 流程：搜索视频 -> 抓取评论 -> AI分析 -> 生成报告 -> 保存数据库
func (e *Executor) Execute(ctx context.Context, req TaskRequest) error {
	taskID := req.TaskID
	log.Printf("[Task %s] Starting execution...", taskID)

	// 阶段0：创建历史记录
	history, err := e.createHistory(req, req.TaskID)
	if err != nil {
		sse.PushError(taskID, fmt.Sprintf("创建任务记录失败: %v", err))
		return err
	}
	log.Printf("[Task %s] History created: ID=%d", taskID, history.ID)

	// 阶段1：获取配置
	sse.PushProgress(taskID, sse.StatusSearching, 0, 100, "正在加载配置...")
	e.updateTaskProgress(history.ID, sse.StatusSearching, 0, "正在加载配置...")

	settings, err := e.loadSettings()
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("加载配置失败: %v", err))
		return err
	}

	// 阶段2：搜索视频
	sse.PushProgress(taskID, sse.StatusSearching, 5, 100, "正在搜索相关视频...")
	e.updateTaskProgress(history.ID, sse.StatusSearching, 5, "正在搜索相关视频...")

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

	// 阶段3：计算按比例分配
	commentAllocation := e.calculateProportionalAllocation(
		allVideos,
		e.config.MaxComments,
		e.config.MinCommentsPerVideo,
		e.config.MaxCommentsPerVideoV2,
	)

	log.Printf("[Task %s] Comment allocation calculated for %d videos", taskID, len(commentAllocation))

	// 阶段4：抓取评论
	sse.PushProgress(taskID, sse.StatusScraping, 20, 100, fmt.Sprintf("开始抓取%d个视频的评论...", len(allVideos)))
	e.updateTaskProgress(history.ID, sse.StatusScraping, 20, fmt.Sprintf("开始抓取%d个视频的评论...", len(allVideos)))

	scraper := bilibili.NewScraper(biliClient, &bilibili.ScraperConfig{
		MaxVideos:           len(allVideos),
		MaxCommentsPerVideo: e.config.MaxCommentsPerVideo,
		MaxConcurrency:      int64(e.config.MaxConcurrency),
		FetchReplies:        true,
		RequestDelay:        200 * time.Millisecond,
	})

	scraper.SetProgressCallback(func(stage string, current, total int, message string) {
		progress := 20 + (current * 30 / max(total, 1))
		sse.PushProgress(taskID, sse.StatusScraping, progress, 100, message)
	})

	scrapeResult, err := scraper.ScrapeByVideos(ctx, allVideos, commentAllocation)
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
	e.updateTaskProgress(history.ID, sse.StatusAnalyzing, 50, "正在使用AI分析评论...")

	aiClient := ai.NewClient(ai.Config{
		APIBase: settings.AIBaseURL,
		APIKey:  settings.AIAPIKey,
		Model:   settings.AIModel,
	})

	// 设置AI分析进度回调
	aiClient.SetProgressCallback(func(stage string, current, total int, message string) {
		progress := 50 + (current * 35 / max(total, 1)) // 50-85%
		sse.PushProgress(taskID, sse.StatusAnalyzing, progress, 100, message)
	})

	analysisResults, err := e.analyzeComments(ctx, taskID, aiClient, scrapeResult, req.Brands, req.Keywords, req.Dimensions, req.Requirement)
	if err != nil {
		e.updateHistoryStatus(history.ID, models.StatusFailed)
		sse.PushError(taskID, fmt.Sprintf("AI分析失败: %v", err))
		return err
	}

	log.Printf("[Task %s] Analysis completed for %d brands", taskID, len(analysisResults))

	// 阶段5：生成报告
	sse.PushProgress(taskID, sse.StatusGenerating, 85, 100, "正在生成分析报告...")
	e.updateTaskProgress(history.ID, sse.StatusGenerating, 85, "正在生成分析报告...")

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
		Videos: scrapeResult.Videos,
	}

	log.Printf("[Executor] scrapeResult.Videos count: %d", len(scrapeResult.Videos))

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
	e.updateTaskProgress(history.ID, sse.StatusGenerating, 95, "正在保存报告...")

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
func (e *Executor) createHistory(req TaskRequest, taskID string) (*models.AnalysisHistory, error) {
	keywordsJSON, err := json.Marshal(req.Keywords)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal keywords: %w", err)
	}
	brandsJSON, err := json.Marshal(req.Brands)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal brands: %w", err)
	}

	dimNames := make([]string, len(req.Dimensions))
	for i, d := range req.Dimensions {
		dimNames[i] = d.Name
	}
	dimensionsJSON, err := json.Marshal(dimNames)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dimensions: %w", err)
	}

	// 序列化任务配置
	configJSON, err := json.Marshal(e.config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	history := &models.AnalysisHistory{
		TaskID:        taskID,
		Category:      req.Requirement,
		Keywords:      string(keywordsJSON),
		Brands:        string(brandsJSON),
		Dimensions:    string(dimensionsJSON),
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

		videos, err := client.SearchVideosWithLimit(keyword, e.config.MaxVideosPerKeyword, e.config.MinVideoDuration)
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

	// 视频时间过滤：过滤掉发布时间超过指定月数的旧视频
	if e.config.VideoDateRangeMonths > 0 {
		cutoffTime := time.Now().AddDate(0, -e.config.VideoDateRangeMonths, 0)
		var filteredVideos []bilibili.VideoInfo
		filteredCount := 0
		for _, v := range allVideos {
			if time.Unix(v.Pubdate, 0).After(cutoffTime) {
				filteredVideos = append(filteredVideos, v)
			} else {
				filteredCount++
			}
		}
		if filteredCount > 0 {
			log.Printf("[Task %s] 过滤了 %d 个超过 %d 个月的旧视频", taskID, filteredCount, e.config.VideoDateRangeMonths)
		}
		allVideos = filteredVideos
	}

	// 评论数过滤：过滤掉评论数低于最小值的视频
	if e.config.MinVideoComments > 0 {
		var filteredVideos []bilibili.VideoInfo
		filteredCount := 0
		for _, v := range allVideos {
			if v.VideoReview >= e.config.MinVideoComments {
				filteredVideos = append(filteredVideos, v)
			} else {
				filteredCount++
			}
		}
		if filteredCount > 0 {
			log.Printf("[Task %s] 过滤了 %d 个评论数低于 %d 的视频", taskID, filteredCount, e.config.MinVideoComments)
		}
		allVideos = filteredVideos
	}

	return allVideos, nil
}

// calculateProportionalAllocation 计算按比例分配评论数
// 根据视频评论数按比例分配抓取数量，确保总数不超过 maxComments
func (e *Executor) calculateProportionalAllocation(
	videos []bilibili.VideoInfo,
	maxComments int,
	minPerVideo int,
	maxPerVideo int,
) map[string]int {
	result := make(map[string]int)

	if len(videos) == 0 {
		return result
	}

	totalComments := 0
	for _, v := range videos {
		totalComments += v.VideoReview
	}

	if totalComments == 0 {
		perVideo := maxComments / len(videos)
		if perVideo < minPerVideo {
			perVideo = minPerVideo
		}
		if perVideo > maxPerVideo {
			perVideo = maxPerVideo
		}
		for _, v := range videos {
			result[v.BVID] = perVideo
		}
		return result
	}

	for _, v := range videos {
		ratio := float64(v.VideoReview) / float64(totalComments)
		allocated := int(ratio * float64(maxComments))

		if allocated < minPerVideo {
			allocated = minPerVideo
		}
		if allocated > maxPerVideo {
			allocated = maxPerVideo
		}
		if v.VideoReview > 0 && allocated > v.VideoReview {
			allocated = v.VideoReview
		}

		result[v.BVID] = allocated
	}

	return result
}

// analyzeComments 分析评论
func (e *Executor) analyzeComments(
	ctx context.Context,
	taskID string,
	aiClient *ai.Client,
	scrapeResult *bilibili.ScrapeResult,
	brands []string,
	keywords []string,
	dimensions []ai.Dimension,
	category string,
) (map[string][]report.CommentWithScore, error) {

	// 1. 使用 GetAllCommentsWithVideo 获取评论
	allComments := GetAllCommentsWithVideo(scrapeResult)
	if len(allComments) == 0 {
		return nil, fmt.Errorf("没有获取到任何评论")
	}

	// 2. 构建 AI 输入
	// 简单过滤：保留长度大于等于5的评论
	var inputs []ai.CommentInput
	for i, c := range allComments {
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

	// 如果过滤后没有评论，返回错误
	if len(inputs) == 0 {
		return nil, fmt.Errorf("过滤后没有有效评论")
	}

	// 截断到最大评论数
	if len(inputs) > e.config.MaxComments {
		inputs = inputs[:e.config.MaxComments]
	}

	log.Printf("[Task %s] Prepared %d comments for analysis", taskID, len(inputs))

	sse.PushProgress(taskID, sse.StatusAnalyzing, 55, 100,
		fmt.Sprintf("正在AI分析 %d 条评论...", len(inputs)))

	// 3. AI 分析
	analysisResults, err := aiClient.AnalyzeCommentsWithRateLimit(ctx, inputs, dimensions, e.config.AIBatchSize)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// === 批量识别未知品牌 ===
	// 收集品牌为"未知"但有型号的评论
	unknownBrandModels := make(map[string]bool)
	for _, r := range analysisResults {
		if r.Error != "" || r.Scores == nil {
			continue
		}
		brand := strings.TrimSpace(r.Brand)
		model := strings.TrimSpace(r.Model)
		if (brand == "" || brand == "未知") && model != "" && model != "未知" && model != "通用" {
			unknownBrandModels[model] = true
		}
	}

	// 批量调用AI识别品牌
	var modelToBrand map[string]string
	if len(unknownBrandModels) > 0 {
		models := make([]string, 0, len(unknownBrandModels))
		for m := range unknownBrandModels {
			models = append(models, m)
		}
		log.Printf("[Task %s] 🔍 AI识别未知品牌: %v", taskID, models)

		var err error
		// 收集已发现的品牌（从analysisResults中提取）
		discoveredBrands := collectDiscoveredBrands(analysisResults)

		identifyCtx := ai.BrandIdentifyContext{
			Category:         category,
			KnownBrands:      brands,
			DiscoveredBrands: discoveredBrands,
		}
		modelToBrand, err = aiClient.IdentifyBrandsForModels(ctx, models, identifyCtx)
		if err != nil {
			log.Printf("[Task %s] ⚠️ 品牌识别失败: %v", taskID, err)
			modelToBrand = make(map[string]string)
		}
	}

	// 更新分析结果中的品牌
	for i := range analysisResults {
		r := &analysisResults[i]
		brand := strings.TrimSpace(r.Brand)
		model := strings.TrimSpace(r.Model)

		// 如果品牌未知，尝试从AI识别结果获取
		if (brand == "" || brand == "未知") && model != "" {
			if identifiedBrand, ok := modelToBrand[model]; ok && identifiedBrand != "" && identifiedBrand != "未知" {
				r.Brand = identifiedBrand
				brand = identifiedBrand
			}
		}

		// 格式化品牌名称（纯字母转大写）
		if brand != "" {
			r.Brand = formatBrandName(brand)
		}

		// AI分析返回后，清洗品牌和型号
		analysisResults[i].Brand = comment.CleanBrandName(r.Brand, brands)
		analysisResults[i].Model = comment.CleanModelName(r.Model)
	}

	// === DISCOVERY MODE: 收集所有AI识别的品牌，不仅仅是用户指定的 ===

	// 创建用户指定品牌的映射（用于分类，不是过滤）
	specifiedBrands := make(map[string]string) // lowercase -> original
	for _, brand := range brands {
		specifiedBrands[strings.ToLower(brand)] = brand
	}

	// 分类收集结果：指定品牌 vs 发现的新品牌
	specifiedResults := make(map[string][]report.CommentWithScore)
	discoveredResults := make(map[string][]report.CommentWithScore)

	for _, r := range analysisResults {
		if r.Error != "" || r.Scores == nil {
			continue
		}

		// 从AI结果获取品牌
		brand := strings.TrimSpace(r.Brand)
		brand = normalizeBrand(brand)
		if brand == "" || brand == "未知" {
			// 尝试从评论内容中匹配用户指定的品牌
			contentLower := strings.ToLower(r.Content)
			for specBrandLower, origBrand := range specifiedBrands {
				if strings.Contains(contentLower, specBrandLower) {
					brand = origBrand
					break
				}
			}
		}

		if brand == "" {
			continue // 仍然没有品牌则跳过
		}

		// 型号后备提取：仅在AI未提取到型号时使用正则匹配
		model := r.Model
		if model == "" || model == "未知" {
			model = extractModelFromContent(r.Content)
		}

		// 由于 CommentWithVideo 不包含时间信息，这里使用当前时间或默认值
		// 注意：如果需要准确时间，需要在 CommentWithVideo 中添加 Time 字段
		publishTime := time.Time{}

		commentItem := report.CommentWithScore{
			Content:     r.Content,
			Scores:      r.Scores,
			Brand:       r.Brand,
			Model:       model,
			PublishTime: publishTime,
		}

		// 分类：指定品牌还是发现的新品牌
		brandLower := strings.ToLower(brand)
		brand = normalizeBrand(brand)
		brandLower = strings.ToLower(brand)
		isSpecified := false
		for specBrandLower, origBrand := range specifiedBrands {
			if strings.Contains(brandLower, specBrandLower) || strings.Contains(specBrandLower, brandLower) {
				specifiedResults[origBrand] = append(specifiedResults[origBrand], commentItem)
				isSpecified = true
				break
			}
		}

		if !isSpecified {
			// 这是新发现的品牌 - 保留它！
			discoveredResults[brand] = append(discoveredResults[brand], commentItem)
		}
	}

	// 记录发现的新品牌
	if len(discoveredResults) > 0 {
		var discoveredBrandNames []string
		for brand := range discoveredResults {
			discoveredBrandNames = append(discoveredBrandNames, brand)
		}
		log.Printf("[Task %s] 🔍 发现新品牌: %v", taskID, discoveredBrandNames)
	}

	// 合并结果：先指定品牌，再发现的品牌
	results := make(map[string][]report.CommentWithScore)
	for brand, comments := range specifiedResults {
		results[brand] = comments
		log.Printf("[Task %s] 指定品牌 %s: %d 条评论", taskID, brand, len(comments))
	}
	for brand, comments := range discoveredResults {
		results[brand] = comments
		log.Printf("[Task %s] 发现品牌 %s: %d 条评论", taskID, brand, len(comments))
	}

	return results, nil
}

// saveReport 保存报告到数据库
func (e *Executor) saveReport(historyID uint, reportData *report.ReportData) (uint, error) {
	// 添加调试日志：检查字段是否存在
	log.Printf("[saveReport] VideoSources count: %d", len(reportData.VideoSources))
	log.Printf("[saveReport] SentimentDistribution: %+v", reportData.SentimentDistribution)
	log.Printf("[saveReport] KeywordFrequency count: %d", len(reportData.KeywordFrequency))

	data, err := json.Marshal(reportData)
	if err != nil {
		return 0, err
	}

	// 添加调试日志：检查序列化后的JSON
	log.Printf("[saveReport] Serialized JSON length: %d", len(data))
	log.Printf("[saveReport] Contains video_sources: %v", strings.Contains(string(data), "video_sources"))
	log.Printf("[saveReport] Contains sentiment_distribution: %v", strings.Contains(string(data), "sentiment_distribution"))
	log.Printf("[saveReport] Contains keyword_frequency: %v", strings.Contains(string(data), "keyword_frequency"))

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
	if err := database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Update("status", status).Error; err != nil {
		log.Printf("[Executor] Failed to update history status for ID %d: %v", historyID, err)
	}
}

// updateHistoryStats 更新历史记录统计信息
func (e *Executor) updateHistoryStats(historyID uint, videoCount, commentCount int) {
	if err := database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"video_count":   videoCount,
		"comment_count": commentCount,
	}).Error; err != nil {
		log.Printf("[Executor] Failed to update history stats for ID %d: %v", historyID, err)
	}
}

// updateHistoryWithReport 更新历史记录的报告ID和状态
func (e *Executor) updateHistoryWithReport(historyID uint, reportID uint) {
	if err := database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"report_id": reportID,
		"status":    models.StatusCompleted,
	}).Error; err != nil {
		log.Printf("[Executor] Failed to update history with report for ID %d: %v", historyID, err)
	}
}

// updateTaskProgress 更新任务进度到数据库
func (e *Executor) updateTaskProgress(historyID uint, stage string, progress int, message string) {
	if err := database.DB.Model(&models.AnalysisHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"stage":          stage,
		"progress":       progress,
		"progress_msg":   message,
		"last_heartbeat": time.Now(),
	}).Error; err != nil {
		log.Printf("[Executor] Failed to update task progress for ID %d: %v", historyID, err)
	}
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

// normalizeBrand 品牌名称归一化
// 处理常见的品牌别名，返回统一的品牌名称
func normalizeBrand(brand string) string {
	brandLower := strings.ToLower(strings.TrimSpace(brand))

	// 品牌别名映射
	aliases := map[string][]string{
		"苹果":      {"apple", "iphone", "ipad", "mac", "airpods"},
		"戴森":      {"dyson"},
		"小米":      {"xiaomi", "mi", "redmi"},
		"华为":      {"huawei", "honor", "荣耀"},
		"三星":      {"samsung", "galaxy"},
		"索尼":      {"sony", "playstation", "ps5"},
		"小佩":      {"petkit"},
		"CATLINK": {"catlink", "猫猫狗狗"},
	}

	for normalized, aliasList := range aliases {
		for _, alias := range aliasList {
			if strings.Contains(brandLower, alias) || strings.Contains(alias, brandLower) {
				return normalized
			}
		}
		// 也检查归一化名称本身
		if strings.Contains(brandLower, strings.ToLower(normalized)) {
			return normalized
		}
	}

	return brand // 返回原始名称
}

// formatBrandName 格式化品牌名称
// 纯字母品牌转全大写，中文品牌保持原样
func formatBrandName(brand string) string {
	brand = strings.TrimSpace(brand)
	if brand == "" {
		return brand
	}

	// 检查是否为纯字母（ASCII字母）
	isPureAlpha := true
	for _, r := range brand {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			isPureAlpha = false
			break
		}
	}

	if isPureAlpha {
		return strings.ToUpper(brand)
	}

	return brand
}

// extractModelFromContent 从评论内容中提取型号（正则匹配后备方案）
// 仅在AI未能提取型号时使用
// 参数：
//   - content: 评论内容
//
// 返回：
//   - 提取到的型号，如果未找到则返回空字符串
func extractModelFromContent(content string) string {
	for _, re := range modelPatterns {
		if match := re.FindString(content); match != "" {
			return strings.TrimSpace(match)
		}
	}
	return ""
}

// collectDiscoveredBrands 从分析结果中收集已发现的品牌
func collectDiscoveredBrands(results []ai.CommentAnalysisResult) []string {
	brandSet := make(map[string]bool)
	for _, r := range results {
		if r.Error != "" || r.Brand == "" || r.Brand == "未知" {
			continue
		}
		brand := strings.TrimSpace(r.Brand)
		if brand != "" {
			brandSet[brand] = true
		}
	}

	brands := make([]string, 0, len(brandSet))
	for brand := range brandSet {
		brands = append(brands, brand)
	}
	return brands
}

// GetAllCommentsWithVideo 获取所有评论（带视频信息）
func GetAllCommentsWithVideo(result *bilibili.ScrapeResult) []CommentWithVideo {
	var comments []CommentWithVideo

	videoTitleMap := make(map[string]string)
	for _, video := range result.Videos {
		videoTitleMap[video.BVID] = video.Title
	}

	for bvid, videoComments := range result.Comments {
		videoTitle := videoTitleMap[bvid]
		for _, c := range videoComments {
			comments = append(comments, CommentWithVideo{
				Content:    c.Content.Message,
				VideoTitle: videoTitle,
				VideoBVID:  bvid,
			})
			for _, r := range c.Replies {
				comments = append(comments, CommentWithVideo{
					Content:    r.Content.Message,
					VideoTitle: videoTitle,
					VideoBVID:  bvid,
				})
			}
		}
	}

	return comments
}
