package report

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/models"
	"encoding/json"
	"fmt"
	"sort"
)

// ReportData 报告数据结构
// 包含分析的完整结果，包括品牌、维度、得分、排名和购买建议
type ReportData struct {
	Category       string                        `json:"category"`       // 商品类别（如"手机"）
	Brands         []string                      `json:"brands"`         // 参与对比的品牌列表
	Dimensions     []ai.Dimension                `json:"dimensions"`     // 评价维度列表
	Scores         map[string]map[string]float64 `json:"scores"`         // 品牌 -> 维度 -> 得分
	Rankings       []BrandRanking                `json:"rankings"`       // 品牌排名列表
	Recommendation string                        `json:"recommendation"` // 购买建议文本
}

// BrandRanking 品牌排名信息
// 包含单个品牌的综合得分、排名和各维度得分
type BrandRanking struct {
	Brand        string             `json:"brand"`         // 品牌名称
	OverallScore float64            `json:"overall_score"` // 综合得分（所有维度平均）
	Rank         int                `json:"rank"`          // 排名（1表示第一名）
	Scores       map[string]float64 `json:"scores"`        // 各维度得分
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

	// 第一步：计算各品牌各维度的平均得分
	scores := make(map[string]map[string]float64)
	for brand, results := range analysisResults {
		brandScores := make(map[string]float64)
		dimCounts := make(map[string]int) // 记录每个维度有多少条有效评分

		// 遍历该品牌的所有评论分析结果
		for _, result := range results {
			for dimName, score := range result.Scores {
				// 只统计非空的评分
				if score != nil {
					brandScores[dimName] += *score
					dimCounts[dimName]++
				}
			}
		}

		// 计算平均值
		for dimName, total := range brandScores {
			if count := dimCounts[dimName]; count > 0 {
				brandScores[dimName] = total / float64(count)
			}
		}

		scores[brand] = brandScores
	}

	// 第二步：生成品牌排名
	rankings := generateRankings(brands, dimensions, scores)

	// 第三步：生成购买建议
	recommendation := generateRecommendation(rankings, dimensions)

	return &ReportData{
		Category:       category,
		Brands:         brands,
		Dimensions:     dimensions,
		Scores:         scores,
		Rankings:       rankings,
		Recommendation: recommendation,
	}, nil
}

// generateRankings 生成品牌排名
// 根据各维度得分计算综合得分，并按综合得分排序
func generateRankings(brands []string, dimensions []ai.Dimension, scores map[string]map[string]float64) []BrandRanking {
	rankings := make([]BrandRanking, 0, len(brands))

	// 为每个品牌计算综合得分
	for _, brand := range brands {
		brandScores := scores[brand]

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
			overallScore = total / float64(count)
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
