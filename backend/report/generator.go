package report

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"
	"unicode"
)

// normalizeModelKey 生成归一化的型号key用于比对
// 规则：品牌小写 + "|" + 型号去空格小写
// 例如：("OPPO", "TWS 5") -> "oppo|tws5"
func normalizeModelKey(brand, model string) string {
	brandKey := strings.ToLower(strings.TrimSpace(brand))
	modelKey := strings.ToLower(strings.TrimSpace(model))
	modelKey = strings.ReplaceAll(modelKey, " ", "")
	modelKey = strings.ReplaceAll(modelKey, "-", "")
	modelKey = strings.ReplaceAll(modelKey, "_", "")
	return brandKey + "|" + modelKey
}

// getDisplayModel 从多个型号变体中选择最佳显示名称
// 优先选择：有空格分隔的 > 首字母大写的 > 第一个出现的
func getDisplayModel(variants []string) string {
	if len(variants) == 0 {
		return ""
	}
	if len(variants) == 1 {
		return variants[0]
	}
	for _, v := range variants {
		if strings.Contains(v, " ") {
			return v
		}
	}
	for _, v := range variants {
		if len(v) > 0 && v[0] >= 'A' && v[0] <= 'Z' {
			return v
		}
	}
	return variants[0]
}

// VideoSource 视频来源信息
type VideoSource struct {
	BVID        string `json:"bvid"`         // BV号
	Title       string `json:"title"`        // 视频标题
	Author      string `json:"author"`       // UP主
	Play        int    `json:"play"`         // 播放量
	VideoReview int    `json:"video_review"` // 评论数
}

// KeywordItem 关键词词频项
// Word为关键词文本，Count为出现次数
type KeywordItem struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// ReportData 报告数据结构
// 包含分析的完整结果，包括品牌、维度、得分、排名和购买建议
type ReportData struct {
	Category       string                        `json:"category"`       // 商品类别（如"手机"）
	Brands         []string                      `json:"brands"`         // 参与对比的品牌列表
	Dimensions     []ai.Dimension                `json:"dimensions"`     // 评价维度列表
	Scores         map[string]map[string]float64 `json:"scores"`         // 品牌 -> 维度 -> 得分
	Rankings       []BrandRanking                `json:"rankings"`       // 品牌排名列表
	Recommendation string                        `json:"recommendation"` // 购买建议文本
	// 新增字段
	Stats                 ReportStats                 `json:"stats"`                  // 统计数据
	SentimentDistribution SentimentStats              `json:"sentiment_distribution"` // 情感分布（基于评分阈值统计）
	TopComments           map[string][]TypicalComment `json:"top_comments"`           // 品牌 -> 好评列表
	BadComments           map[string][]TypicalComment `json:"bad_comments"`           // 品牌 -> 差评列表
	BrandAnalysis         map[string]BrandAnalysis    `json:"brand_analysis"`         // 品牌 -> 优劣势分析
	ModelRankings         []ModelRanking              `json:"model_rankings"`         // 型号排名列表
	VideoSources          []VideoSource               `json:"video_sources"`          // 视频来源列表
	KeywordFrequency      []KeywordItem               `json:"keyword_frequency"`      // 关键词词频（用于词云）
}

// BrandRanking 品牌排名信息
// 包含单个品牌的综合得分、排名和各维度得分
type BrandRanking struct {
	Brand        string             `json:"brand"`         // 品牌名称
	OverallScore float64            `json:"overall_score"` // 综合得分（所有维度平均）
	Rank         int                `json:"rank"`          // 排名（1表示第一名）
	Scores       map[string]float64 `json:"scores"`        // 各维度得分
}

// ReportStats 报告统计数据
// 包含视频数、评论数、各品牌评论数等统计信息
type ReportStats struct {
	TotalVideos     int            `json:"total_videos"`      // 搜索到的视频总数
	TotalComments   int            `json:"total_comments"`    // 抓取的评论总数
	CommentsByBrand map[string]int `json:"comments_by_brand"` // 各品牌评论数
}

// SentimentStats 情感分布统计
// 注意：这里不做任何AI情感分析，只按评分阈值划分好评/中性/差评
type SentimentStats struct {
	PositiveCount int     `json:"positive_count"`
	NeutralCount  int     `json:"neutral_count"`
	NegativeCount int     `json:"negative_count"`
	PositivePct   float64 `json:"positive_pct"`
	NeutralPct    float64 `json:"neutral_pct"`
	NegativePct   float64 `json:"negative_pct"`
}

// TypicalComment 典型评论
// 包含评论内容和平均得分
type TypicalComment struct {
	Content string  `json:"content"` // 评论内容
	Score   float64 `json:"score"`   // 平均得分
}

// BrandAnalysis 品牌优劣势分析
// 包含品牌的优势和劣势维度列表
type BrandAnalysis struct {
	Strengths  []string `json:"strengths"`  // 优势维度（得分>=8.0）
	Weaknesses []string `json:"weaknesses"` // 劣势维度（得分<6.0）
}

// ModelRanking 型号排名信息
type ModelRanking struct {
	Model        string             `json:"model"`         // 型号名称
	Brand        string             `json:"brand"`         // 品牌名称
	OverallScore float64            `json:"overall_score"` // 综合得分
	Rank         int                `json:"rank"`          // 排名
	Scores       map[string]float64 `json:"scores"`        // 各维度得分
	CommentCount int                `json:"comment_count"` // 评论数量
}

// CommentWithScore 带得分的评论
type CommentWithScore struct {
	Content     string
	Scores      map[string]*float64
	Brand       string
	Model       string
	PublishTime time.Time
}

// GenerateReportInput 报告生成输入参数
type GenerateReportInput struct {
	Category        string
	Brands          []string
	Dimensions      []ai.Dimension
	AnalysisResults map[string][]CommentWithScore // brand -> 评论及得分列表
	Stats           ReportStats
	Videos          []bilibili.VideoInfo
}

// GenerateReport 生成分析报告
// 参数：
//   - category: 商品类别
//   - brands: 品牌列表
//   - dimensions: 评价维度列表
//   - analysisResults: 每个品牌的评论分析结果（品牌 -> 分析结果列表）
//
// 返回：
//   - *ReportData: 生成的报告数据
//   - error: 错误信息
func GenerateReport(
	category string,
	brands []string,
	dimensions []ai.Dimension,
	analysisResults map[string][]ai.AnalyzeCommentResponse, // brand -> results
) (*ReportData, error) {
	input := GenerateReportInput{
		Category:        category,
		Brands:          brands,
		Dimensions:      dimensions,
		AnalysisResults: make(map[string][]CommentWithScore),
		Stats:           ReportStats{CommentsByBrand: make(map[string]int)},
	}
	for brand, results := range analysisResults {
		for _, r := range results {
			input.AnalysisResults[brand] = append(input.AnalysisResults[brand], CommentWithScore{
				Content: "",
				Scores:  r.Scores,
			})
		}
		input.Stats.CommentsByBrand[brand] = len(results)
	}
	return GenerateReportWithInput(input)
}

// GenerateReportWithInput 使用完整输入生成报告（支持典型评论筛选）
func GenerateReportWithInput(input GenerateReportInput) (*ReportData, error) {
	scores := make(map[string]map[string]float64)
	for brand, results := range input.AnalysisResults {
		brandScores := make(map[string]float64)
		dimCounts := make(map[string]int)

		for _, result := range results {
			for dimName, score := range result.Scores {
				if score != nil {
					brandScores[dimName] += *score
					dimCounts[dimName]++
				}
			}
		}

		for dimName, total := range brandScores {
			if count := dimCounts[dimName]; count > 0 {
				brandScores[dimName] = math.Round((total/float64(count))*10) / 10
			}
		}

		scores[brand] = brandScores
	}

	rankings := generateRankings(input.Brands, input.Dimensions, scores)
	recommendation := generateRecommendation(rankings, input.Dimensions)

	// 收集所有发现的品牌（用于品牌分析）
	allBrands := make([]string, 0, len(scores))
	for brand := range scores {
		if brand != "" {
			allBrands = append(allBrands, brand)
		}
	}
	brandAnalysis := generateBrandAnalysis(allBrands, input.Dimensions, scores)
	topComments, badComments := selectTypicalComments(input.AnalysisResults)

	// 生成型号排名（使用归一化key合并相似型号）
	modelRankings := generateModelRankings(input.AnalysisResults, input.Dimensions)

	// 收集所有品牌名称用于报告（按排名顺序）
	allBrandNames := make([]string, 0, len(rankings))
	for _, r := range rankings {
		allBrandNames = append(allBrandNames, r.Brand)
	}

	log.Printf("[Report] Input.Videos count: %d", len(input.Videos))

	videoSources := make([]VideoSource, len(input.Videos))
	for i, v := range input.Videos {
		videoSources[i] = VideoSource{
			BVID:        v.BVID,
			Title:       v.Title,
			Author:      v.Author,
			Play:        v.Play,
			VideoReview: v.VideoReview,
		}
	}

	// 计算整体情感分布：仅按评分阈值划分，不做任何AI情感分析
	sentimentDistribution := calculateSentiment(input.AnalysisResults)
	log.Printf("[GenerateReport] SentimentDistribution: %+v", sentimentDistribution)

	// 提取关键词词频：用于词云（简单分词 + 停用词/单字过滤）
	keywordFrequency := extractKeywords(input.AnalysisResults)
	log.Printf("[GenerateReport] KeywordFrequency count: %d", len(keywordFrequency))
	log.Printf("[GenerateReport] VideoSources count: %d", len(videoSources))

	return &ReportData{
		Category:              input.Category,
		Brands:                allBrandNames,
		Dimensions:            input.Dimensions,
		Scores:                scores,
		Rankings:              rankings,
		Recommendation:        recommendation,
		Stats:                 input.Stats,
		SentimentDistribution: sentimentDistribution,
		TopComments:           topComments,
		BadComments:           badComments,
		BrandAnalysis:         brandAnalysis,
		ModelRankings:         modelRankings,
		VideoSources:          videoSources,
		KeywordFrequency:      keywordFrequency,
	}, nil
}

// generateModelRankings 生成型号排名
// 按"品牌+型号"聚合，使用归一化key合并相似型号（如TWS5、TWS 5、Tws5）
func generateModelRankings(analysisResults map[string][]CommentWithScore, dimensions []ai.Dimension) []ModelRanking {
	modelScores := make(map[string]map[string][]float64) // normalizedKey -> 维度 -> 分数列表
	modelCommentCounts := make(map[string]int)           // normalizedKey -> 评论数
	modelVariants := make(map[string][]string)           // normalizedKey -> 原始型号变体列表
	modelBrands := make(map[string]string)               // normalizedKey -> 品牌

	for brandKey, results := range analysisResults {
		for _, result := range results {
			model := strings.TrimSpace(result.Model)
			if model == "" || model == "未知" || model == "通用" {
				continue
			}
			brand := strings.TrimSpace(result.Brand)
			if brand == "" {
				brand = strings.TrimSpace(brandKey)
			}
			if brand == "" || brand == "未知" {
				continue
			}

			normalizedKey := normalizeModelKey(brand, model)
			if modelScores[normalizedKey] == nil {
				modelScores[normalizedKey] = make(map[string][]float64)
				modelBrands[normalizedKey] = brand
			}

			modelVariants[normalizedKey] = append(modelVariants[normalizedKey], model)

			for dimName, score := range result.Scores {
				if score != nil {
					modelScores[normalizedKey][dimName] = append(modelScores[normalizedKey][dimName], *score)
				}
			}
			modelCommentCounts[normalizedKey]++
		}
	}

	modelRankings := make([]ModelRanking, 0, len(modelScores))
	for normalizedKey, dimScores := range modelScores {
		commentCount := modelCommentCounts[normalizedKey]
		if commentCount < 1 {
			continue
		}

		brand := modelBrands[normalizedKey]
		displayModel := getDisplayModel(modelVariants[normalizedKey])

		avgScores := make(map[string]float64)
		var total float64
		var dimCount int
		for dimName, scores := range dimScores {
			if len(scores) == 0 {
				continue
			}
			var sum float64
			for _, s := range scores {
				sum += s
			}
			avg := sum / float64(len(scores))
			avgScores[dimName] = math.Round(avg*10) / 10
			total += avg
			dimCount++
		}
		overallScore := 0.0
		if dimCount > 0 {
			overallScore = total / float64(dimCount)
		}

		modelRankings = append(modelRankings, ModelRanking{
			Model:        displayModel,
			Brand:        brand,
			OverallScore: overallScore,
			Scores:       avgScores,
			CommentCount: commentCount,
		})
	}

	// 按综合得分排序
	sort.Slice(modelRankings, func(i, j int) bool {
		if modelRankings[i].OverallScore == modelRankings[j].OverallScore {
			if modelRankings[i].Brand == modelRankings[j].Brand {
				return modelRankings[i].Model < modelRankings[j].Model
			}
			return modelRankings[i].Brand < modelRankings[j].Brand
		}
		return modelRankings[i].OverallScore > modelRankings[j].OverallScore
	})

	// 设置排名
	for i := range modelRankings {
		modelRankings[i].Rank = i + 1
	}

	return modelRankings
}

// generateRankings 生成品牌排名
// 根据各维度得分计算综合得分，并按综合得分排序
// 注意：遍历 scores 中的所有品牌（包括AI发现的新品牌），而不仅仅是用户指定的品牌
func generateRankings(brands []string, dimensions []ai.Dimension, scores map[string]map[string]float64) []BrandRanking {
	rankings := make([]BrandRanking, 0, len(scores))

	// 为每个品牌计算综合得分（遍历所有发现的品牌）
	for brand, brandScores := range scores {
		if brand == "" {
			continue
		}

		// 计算综合得分（所有维度的平均值）
		var total float64
		var count int
		for _, dim := range dimensions {
			if score, ok := brandScores[dim.Name]; ok {
				total += score
				count++
			}
		}

		overallScore := 0.0
		if count > 0 {
			overallScore = math.Round((total/float64(count))*10) / 10
		}

		rankings = append(rankings, BrandRanking{
			Brand:        brand,
			OverallScore: overallScore,
			Scores:       brandScores,
		})
	}

	// 按综合得分从高到低排序
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].OverallScore > rankings[j].OverallScore
	})

	// 设置排名（1表示第一名）
	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings
}

// generateRecommendation 生成购买建议
// 基于排名和得分生成人性化的购买建议文本
func generateRecommendation(rankings []BrandRanking, dimensions []ai.Dimension) string {
	if len(rankings) == 0 {
		return "暂无足够数据生成购买建议"
	}

	topBrand := rankings[0]

	// 找出该品牌的优势维度（得分>=8.0的维度）
	var strengths []string
	for _, dim := range dimensions {
		if score, ok := topBrand.Scores[dim.Name]; ok && score >= 8.0 {
			strengths = append(strengths, dim.Name)
		}
	}

	// 构建推荐文本
	recommendation := fmt.Sprintf("综合评价最高的是 %s（综合得分：%.1f分）", topBrand.Brand, topBrand.OverallScore)

	if len(strengths) > 0 {
		recommendation += fmt.Sprintf("，在 %v 方面表现突出", strengths)
	}

	// 如果有第二名，也提及
	if len(rankings) > 1 {
		secondBrand := rankings[1]
		recommendation += fmt.Sprintf("。%s（%.1f分）紧随其后", secondBrand.Brand, secondBrand.OverallScore)
	}

	recommendation += "。建议根据个人需求和预算选择合适的产品。"

	return recommendation
}

// SaveReport 保存报告到数据库
// 将报告数据序列化为JSON并保存到数据库
func SaveReport(historyID uint, reportData *ReportData, db interface{}) error {
	// 序列化报告数据为JSON
	data, err := json.Marshal(reportData)
	if err != nil {
		return fmt.Errorf("marshal report failed: %w", err)
	}

	// 创建报告记录
	report := models.Report{
		HistoryID:  historyID,
		Category:   reportData.Category,
		ReportData: string(data),
	}

	// TODO: 实际保存到数据库（需要GORM实例）
	_ = report

	return nil
}

// generateBrandAnalysis 生成品牌优劣势分析
// 遍历每个品牌的各维度得分，得分>=8.0归为优势，<6.0归为劣势
func generateBrandAnalysis(brands []string, dimensions []ai.Dimension, scores map[string]map[string]float64) map[string]BrandAnalysis {
	result := make(map[string]BrandAnalysis)

	for _, brand := range brands {
		brandScores := scores[brand]
		analysis := BrandAnalysis{
			Strengths:  []string{},
			Weaknesses: []string{},
		}

		for _, dim := range dimensions {
			if score, ok := brandScores[dim.Name]; ok {
				if score >= 8.0 {
					analysis.Strengths = append(analysis.Strengths, dim.Name)
				} else if score < 6.0 {
					analysis.Weaknesses = append(analysis.Weaknesses, dim.Name)
				}
			}
		}

		result[brand] = analysis
	}

	return result
}

// selectTypicalComments 筛选典型好评和差评
// 好评：平均得分>=8.0，取前3条；差评：平均得分<5.0，取前3条
func selectTypicalComments(analysisResults map[string][]CommentWithScore) (map[string][]TypicalComment, map[string][]TypicalComment) {
	topComments := make(map[string][]TypicalComment)
	badComments := make(map[string][]TypicalComment)

	for brand, results := range analysisResults {
		var goodList, badList []TypicalComment

		for _, r := range results {
			if r.Content == "" {
				continue
			}
			avgScore := calculateAverageScore(r.Scores)
			if avgScore >= 8.0 {
				goodList = append(goodList, TypicalComment{Content: r.Content, Score: avgScore})
			} else if avgScore < 5.0 && avgScore > 0 {
				badList = append(badList, TypicalComment{Content: r.Content, Score: avgScore})
			}
		}

		sort.Slice(goodList, func(i, j int) bool { return goodList[i].Score > goodList[j].Score })
		sort.Slice(badList, func(i, j int) bool { return badList[i].Score < badList[j].Score })

		if len(goodList) > 3 {
			goodList = goodList[:3]
		}
		if len(badList) > 3 {
			badList = badList[:3]
		}

		topComments[brand] = goodList
		badComments[brand] = badList
	}

	return topComments, badComments
}

func calculateAverageScore(scores map[string]*float64) float64 {
	var total float64
	var count int
	for _, score := range scores {
		if score != nil {
			total += *score
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// calculateSentiment 计算情感分布（仅基于评分阈值）
// 规则：评分 >= 8 为好评；评分 5-8 为中性；评分 < 5 为差评
// 注意：不做任何AI情感分析；当评论没有有效评分（平均分<=0）时跳过
// 百分比保留1位小数
func calculateSentiment(analysisResults map[string][]CommentWithScore) SentimentStats {
	var positiveCount, neutralCount, negativeCount int
	var total int

	for _, results := range analysisResults {
		for _, r := range results {
			avgScore := calculateAverageScore(r.Scores)
			if avgScore <= 0 {
				continue
			}
			total++
			if avgScore >= 8.0 {
				positiveCount++
			} else if avgScore >= 5.0 {
				neutralCount++
			} else {
				negativeCount++
			}
		}
	}

	stats := SentimentStats{
		PositiveCount: positiveCount,
		NeutralCount:  neutralCount,
		NegativeCount: negativeCount,
	}
	if total == 0 {
		return stats
	}

	stats.PositivePct = math.Round((float64(positiveCount)/float64(total)*100)*10) / 10
	stats.NeutralPct = math.Round((float64(neutralCount)/float64(total)*100)*10) / 10
	stats.NegativePct = math.Round((float64(negativeCount)/float64(total)*100)*10) / 10

	return stats
}

// extractKeywords 从所有评论内容中提取关键词并做简单词频统计
// 规则：
//  1. 仅做“按空白/标点切分”的粗粒度分词（适配中英文混合）
//  2. 过滤常见停用词（如“的、了、是、等”）
//  3. 过滤单字词（避免词云被无意义的单字占满）
//  4. 返回按出现次数降序的Top50
func extractKeywords(analysisResults map[string][]CommentWithScore) []KeywordItem {
	// 常见停用词：这里保持简单可控，避免误伤业务关键词
	stopWords := map[string]struct{}{
		"的": {}, "了": {}, "是": {}, "等": {}, "也": {}, "就": {}, "都": {}, "还": {}, "很": {},
		"我": {}, "你": {}, "他": {}, "她": {}, "它": {}, "我们": {}, "你们": {}, "他们": {},
		"这": {}, "那": {}, "这个": {}, "那个": {}, "一个": {}, "一些": {}, "不是": {}, "没有": {},
		"在": {}, "和": {}, "与": {}, "及": {}, "而": {}, "且": {}, "或": {}, "或者": {},
		"因为": {}, "所以": {}, "如果": {}, "但是": {}, "而且": {}, "以及": {},
		"啊": {}, "呢": {}, "吗": {}, "吧": {}, "哦": {}, "呀": {}, "哈": {},
	}

	counts := make(map[string]int)

	flushToken := func(buf *strings.Builder, tokens *[]string) {
		if buf.Len() == 0 {
			return
		}
		w := strings.TrimSpace(buf.String())
		buf.Reset()
		if w == "" {
			return
		}
		*tokens = append(*tokens, w)
	}

	for _, results := range analysisResults {
		for _, r := range results {
			text := strings.TrimSpace(r.Content)
			if text == "" {
				continue
			}

			// 以“字母/数字/汉字”为token字符，其它（空格、标点等）视为分隔符
			var buf strings.Builder
			tokens := make([]string, 0, 32)
			for _, ch := range text {
				if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
					buf.WriteRune(ch)
					continue
				}
				flushToken(&buf, &tokens)
			}
			flushToken(&buf, &tokens)

			for _, raw := range tokens {
				word := strings.ToLower(strings.TrimSpace(raw))
				if word == "" {
					continue
				}
				// 过滤单字词
				if len([]rune(word)) <= 1 {
					continue
				}
				// 过滤停用词
				if _, ok := stopWords[word]; ok {
					continue
				}
				counts[word]++
			}
		}
	}

	items := make([]KeywordItem, 0, len(counts))
	for w, c := range counts {
		items = append(items, KeywordItem{Word: w, Count: c})
	}

	// 稳定、可复现：次数降序；次数相同按词典序升序
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Word < items[j].Word
		}
		return items[i].Count > items[j].Count
	})

	if len(items) > 50 {
		items = items[:50]
	}
	return items
}
