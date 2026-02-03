package bilibili

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// ScraperConfig 抓取器配置
type ScraperConfig struct {
	MaxVideos           int           // 最大视频数量（默认50）
	MaxCommentsPerVideo int           // 每个视频最大评论数（默认500）
	MaxConcurrency      int64         // 最大并发数（默认5）
	FetchReplies        bool          // 是否获取楼中楼（默认true）
	RequestDelay        time.Duration // 请求间隔（默认200ms）
}

// DefaultScraperConfig 默认抓取器配置
func DefaultScraperConfig() ScraperConfig {
	return ScraperConfig{
		MaxVideos:           50,
		MaxCommentsPerVideo: 500,
		MaxConcurrency:      5,
		FetchReplies:        true,
		RequestDelay:        200 * time.Millisecond,
	}
}

// ScrapeResult 抓取结果
type ScrapeResult struct {
	Videos   []VideoInfo          // 视频列表
	Comments map[string][]Comment // 评论映射（key: BVID）
	Stats    ScrapeStats          // 统计信息
}

// ScrapeStats 抓取统计
type ScrapeStats struct {
	TotalVideos   int           // 总视频数
	TotalComments int           // 总评论数
	TotalReplies  int           // 总楼中楼数
	Duration      time.Duration // 耗时
	Errors        []string      // 错误列表
}

// ProgressCallback 进度回调函数类型
// 参数：
//   - stage: 当前阶段（searching/scraping）
//   - current: 当前进度
//   - total: 总数
//   - message: 状态消息
type ProgressCallback func(stage string, current, total int, message string)

// Scraper 评论抓取器
// 支持并发抓取、进度回调、数量限制
type Scraper struct {
	client   *Client          // B站API客户端
	config   ScraperConfig    // 抓取配置
	callback ProgressCallback // 进度回调
}

// NewScraper 创建新的抓取器
//
// 参数：
//   - client: B站API客户端
//   - config: 抓取配置（可选，传nil使用默认配置）
//
// 返回：
//   - *Scraper: 抓取器实例
//
// 示例：
//
//	client := NewClient("")
//	scraper := NewScraper(client, nil)
func NewScraper(client *Client, config *ScraperConfig) *Scraper {
	cfg := DefaultScraperConfig()
	if config != nil {
		cfg = *config
	}
	return &Scraper{
		client: client,
		config: cfg,
	}
}

// SetProgressCallback 设置进度回调
//
// 参数：
//   - callback: 进度回调函数
//
// 示例：
//
//	scraper.SetProgressCallback(func(stage string, current, total int, message string) {
//	    fmt.Printf("[%s] %d/%d: %s\n", stage, current, total, message)
//	})
func (s *Scraper) SetProgressCallback(callback ProgressCallback) {
	s.callback = callback
}

// reportProgress 报告进度
func (s *Scraper) reportProgress(stage string, current, total int, message string) {
	// 添加调试日志，确认进度回调被调用
	log.Printf("[Scraper] Progress: stage=%s, current=%d, total=%d, message=%s", stage, current, total, message)
	if s.callback != nil {
		s.callback(stage, current, total, message)
	} else {
		log.Printf("[Scraper] WARNING: No callback set, progress not pushed to SSE")
	}
}

// ScrapeByKeyword 根据关键词抓取评论
// 搜索视频并抓取评论
//
// 参数：
//   - ctx: 上下文（用于取消）
//   - keyword: 搜索关键词
//
// 返回：
//   - *ScrapeResult: 抓取结果
//   - error: 错误信息
//
// 示例：
//
//	result, err := scraper.ScrapeByKeyword(ctx, "iPhone 15 评测")
func (s *Scraper) ScrapeByKeyword(ctx context.Context, keyword string) (*ScrapeResult, error) {
	startTime := time.Now()
	result := &ScrapeResult{
		Comments: make(map[string][]Comment),
		Stats:    ScrapeStats{},
	}

	// 阶段1：搜索视频
	s.reportProgress("searching", 0, s.config.MaxVideos, "正在搜索视频...")

	videos, err := s.client.SearchVideosWithLimit(keyword, s.config.MaxVideos, 0)
	if err != nil {
		return nil, fmt.Errorf("搜索视频失败: %w", err)
	}

	result.Videos = videos
	result.Stats.TotalVideos = len(videos)

	s.reportProgress("searching", len(videos), s.config.MaxVideos,
		fmt.Sprintf("搜索完成，找到%d个视频", len(videos)))

	// 阶段2：并发抓取评论
	if len(videos) == 0 {
		result.Stats.Duration = time.Since(startTime)
		return result, nil
	}

	s.reportProgress("scraping", 0, len(videos), "开始抓取评论...")

	// 使用semaphore控制并发
	log.Printf("[Scraper] MaxConcurrency: %d", s.config.MaxConcurrency)
	sem := semaphore.NewWeighted(s.config.MaxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var completedCount int

	for i, video := range videos {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 获取信号量
		log.Printf("[Scraper] Acquiring semaphore for video: %s", video.BVID)
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(idx int, v VideoInfo) {
			defer wg.Done()
			defer sem.Release(1)

			// 抓取评论
			comments, err := s.scrapeVideoComments(ctx, v.BVID)

			mu.Lock()
			defer mu.Unlock()

			completedCount++

			if err != nil {
				result.Stats.Errors = append(result.Stats.Errors,
					fmt.Sprintf("视频%s抓取失败: %v", v.BVID, err))
				s.reportProgress("scraping", completedCount, len(videos),
					fmt.Sprintf("视频%s抓取失败", v.BVID))
				return
			}

			result.Comments[v.BVID] = comments

			// 统计评论数
			commentCount := len(comments)
			replyCount := 0
			for _, c := range comments {
				replyCount += len(c.Replies)
			}
			result.Stats.TotalComments += commentCount
			result.Stats.TotalReplies += replyCount

			s.reportProgress("scraping", completedCount, len(videos),
				fmt.Sprintf("已完成%d/%d，当前视频%d条评论", completedCount, len(videos), commentCount))
		}(i, video)

		// 添加请求间隔
		time.Sleep(s.config.RequestDelay)
	}

	// 等待所有goroutine完成
	wg.Wait()

	result.Stats.Duration = time.Since(startTime)
	s.reportProgress("scraping", len(videos), len(videos),
		fmt.Sprintf("抓取完成，共%d条评论，%d条楼中楼", result.Stats.TotalComments, result.Stats.TotalReplies))

	return result, nil
}

// scrapeVideoComments 抓取单个视频的评论
func (s *Scraper) scrapeVideoComments(ctx context.Context, bvid string) ([]Comment, error) {
	var allComments []Comment
	page := 1
	pageSize := 20

	for len(allComments) < s.config.MaxCommentsPerVideo {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 获取评论
		comments, _, err := s.client.GetCommentsWithReplies(GetCommentsRequest{
			BVID:     bvid,
			Page:     page,
			PageSize: pageSize,
			Sort:     1, // 按点赞排序
		}, s.config.FetchReplies, 10)

		if err != nil {
			return nil, err
		}

		// 没有更多评论
		if len(comments) == 0 {
			break
		}

		allComments = append(allComments, comments...)
		page++

		// 防止无限循环
		if page > 50 {
			break
		}

		// 添加请求间隔
		time.Sleep(s.config.RequestDelay)
	}

	// 截取到指定数量
	if len(allComments) > s.config.MaxCommentsPerVideo {
		allComments = allComments[:s.config.MaxCommentsPerVideo]
	}

	return allComments, nil
}

// ScrapeByVideos 根据视频列表抓取评论
// 直接抓取指定视频的评论
//
// 参数：
//   - ctx: 上下文（用于取消）
//   - videos: 视频列表
//
// 返回：
//   - *ScrapeResult: 抓取结果
//   - error: 错误信息
func (s *Scraper) ScrapeByVideos(ctx context.Context, videos []VideoInfo) (*ScrapeResult, error) {
	startTime := time.Now()
	result := &ScrapeResult{
		Videos:   videos,
		Comments: make(map[string][]Comment),
		Stats: ScrapeStats{
			TotalVideos: len(videos),
		},
	}

	if len(videos) == 0 {
		return result, nil
	}

	// 限制视频数量
	if len(videos) > s.config.MaxVideos {
		videos = videos[:s.config.MaxVideos]
		result.Videos = videos
		result.Stats.TotalVideos = len(videos)
	}

	log.Printf("[Scraper] Starting ScrapeByVideos with %d videos", len(videos))
	s.reportProgress("scraping", 0, len(videos), "开始抓取评论...")

	// 使用semaphore控制并发
	log.Printf("[Scraper] MaxConcurrency: %d", s.config.MaxConcurrency)
	sem := semaphore.NewWeighted(s.config.MaxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var completedCount int
	var startedCount int

	for i, video := range videos {
		log.Printf("[Scraper] Loop iteration %d, video: %s", i, video.BVID)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		log.Printf("[Scraper] Acquiring semaphore for video: %s", video.BVID)
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(v VideoInfo) {
			log.Printf("[Scraper] Goroutine started for video: %s", v.BVID)
			defer wg.Done()
			defer sem.Release(1)

			// 开始抓取时推送进度，让用户知道当前正在处理哪个视频
			mu.Lock()
			startedCount++
			currentStarted := startedCount
			mu.Unlock()

			// 截断视频标题（最多20个字符），避免进度消息过长
			title := v.Title
			if len([]rune(title)) > 20 {
				title = string([]rune(title)[:20]) + "..."
			}

			log.Printf("[Scraper] Starting video %d/%d: %s (BVID: %s)", currentStarted, len(videos), title, v.BVID)
			s.reportProgress("scraping", currentStarted, len(videos),
				fmt.Sprintf("正在抓取 (%d/%d): %s", currentStarted, len(videos), title))

			comments, err := s.scrapeVideoComments(ctx, v.BVID)

			mu.Lock()
			defer mu.Unlock()

			completedCount++

			if err != nil {
				result.Stats.Errors = append(result.Stats.Errors,
					fmt.Sprintf("视频%s抓取失败: %v", v.BVID, err))
				// 失败时也推送进度，让用户知道进度在继续
				log.Printf("[Scraper] Failed video %d/%d: %s, error: %v", completedCount, len(videos), v.BVID, err)
				s.reportProgress("scraping", completedCount, len(videos),
					fmt.Sprintf("已完成 %d/%d (失败: %s)", completedCount, len(videos), v.BVID))
				return
			}

			result.Comments[v.BVID] = comments

			commentCount := len(comments)
			replyCount := 0
			for _, c := range comments {
				replyCount += len(c.Replies)
			}
			result.Stats.TotalComments += commentCount
			result.Stats.TotalReplies += replyCount

			// 完成时推送进度，显示累计评论数
			log.Printf("[Scraper] Completed video %d/%d: %s, got %d comments", completedCount, len(videos), v.BVID, commentCount)
			s.reportProgress("scraping", completedCount, len(videos),
				fmt.Sprintf("已完成 %d/%d，共%d条评论", completedCount, len(videos), result.Stats.TotalComments))
		}(video)

		time.Sleep(s.config.RequestDelay)
	}

	wg.Wait()

	log.Printf("[Scraper] Returning result with %d videos", len(result.Videos))

	result.Stats.Duration = time.Since(startTime)
	log.Printf("[Scraper] ScrapeByVideos completed: %d videos, %d comments, duration: %v",
		result.Stats.TotalVideos, result.Stats.TotalComments, result.Stats.Duration)
	return result, nil
}

// GetAllCommentsFlat 获取所有评论的扁平列表
// 将所有视频的评论合并为一个列表
//
// 参数：
//   - result: 抓取结果
//
// 返回：
//   - []Comment: 所有评论列表
func GetAllCommentsFlat(result *ScrapeResult) []Comment {
	var allComments []Comment
	for _, comments := range result.Comments {
		allComments = append(allComments, comments...)
	}
	return allComments
}

// GetAllCommentTexts 获取所有评论文本
// 提取所有评论的文本内容（包括楼中楼）
//
// 参数：
//   - result: 抓取结果
//
// 返回：
//   - []string: 所有评论文本列表
func GetAllCommentTexts(result *ScrapeResult) []string {
	var texts []string
	for _, comments := range result.Comments {
		for _, c := range comments {
			texts = append(texts, c.Content.Message)
			for _, r := range c.Replies {
				texts = append(texts, r.Content.Message)
			}
		}
	}
	return texts
}
