package ai

import (
	"fmt"
	"strings"
)

func BuildRecommendationUserPrompt(input RecommendationInput) string {
	var rankingText string
	for _, r := range input.Rankings {
		analysis := input.BrandAnalysis[r.Brand]
		rankingText += fmt.Sprintf("第%d名：%s（%.1f分）", r.Rank, r.Brand, r.OverallScore)
		if len(analysis.Strengths) > 0 {
			rankingText += fmt.Sprintf("，优势：%v", analysis.Strengths)
		}
		if len(analysis.Weaknesses) > 0 {
			rankingText += fmt.Sprintf("，劣势：%v", analysis.Weaknesses)
		}
		rankingText += "\n"
	}

	var modelText string
	if len(input.ModelRankings) > 0 {
		modelText = "\n\n型号排名：\n"
		for _, m := range input.ModelRankings {
			modelText += fmt.Sprintf("第%d名：%s %s（%.1f分，%d条评论）\n",
				m.Rank, m.Brand, m.Model, m.OverallScore, m.CommentCount)
		}
	}

	var contextLines []string
	if strings.TrimSpace(input.Budget) != "" {
		contextLines = append(contextLines, "预算："+input.Budget)
	}
	if strings.TrimSpace(input.Scenario) != "" {
		contextLines = append(contextLines, "使用场景："+input.Scenario)
	}
	if len(input.SpecialNeeds) > 0 {
		contextLines = append(contextLines, "特殊需求："+strings.Join(input.SpecialNeeds, "、"))
	}
	contextBlock := ""
	if len(contextLines) > 0 {
		contextBlock = "\n用户约束：\n" + strings.Join(contextLines, "\n") + "\n"
	}

	return fmt.Sprintf("商品类别：%s\n%s\n品牌排名及分析：\n%s%s\n请生成购买建议：",
		input.Category, contextBlock, rankingText, modelText)
}
