package ai

import (
	"strings"
	"testing"
)

func TestBuildRecommendationUserPrompt_IncludesBudgetScenarioNeeds(t *testing.T) {
	prompt := BuildRecommendationUserPrompt(RecommendationInput{
		Category:      "无线吸尘器",
		Budget:        "2000元左右",
		Scenario:      "家庭使用",
		SpecialNeeds:  []string{"宠物毛发清理"},
		Rankings:      []BrandRankingInfo{{Brand: "戴森", OverallScore: 8.5, Rank: 1}},
		BrandAnalysis: map[string]BrandStrengthWeakness{},
	})
	for _, part := range []string{"2000元左右", "家庭使用", "宠物毛发清理"} {
		if !strings.Contains(prompt, part) {
			t.Fatalf("prompt missing %q:\n%s", part, prompt)
		}
	}
}
