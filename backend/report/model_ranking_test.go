package report

import (
	"bilibili-analyzer/backend/ai"
	"testing"
)

func TestModelLevelAggregation(t *testing.T) {
	// 准备测试数据
	score8 := 8.0
	score9 := 9.0
	score6 := 6.0

	analysisResults := map[string][]CommentWithScore{
		"BrandA": {
			{
				Brand: "BrandA",
				Model: "ModelX",
				Scores: map[string]*float64{
					"Dim1": &score8,
				},
			},
			{
				Brand: "BrandA",
				Model: "ModelX",
				Scores: map[string]*float64{
					"Dim1": &score9,
				},
			},
			{
				Brand: "BrandA",
				Model: "ModelY",
				Scores: map[string]*float64{
					"Dim1": &score6,
				},
			},
		},
		"BrandB": {
			{
				Brand: "BrandB",
				Model: "ModelZ",
				Scores: map[string]*float64{
					"Dim1": &score9,
				},
			},
		},
	}

	dimensions := []ai.Dimension{
		{Name: "Dim1", Description: "Dimension 1"},
	}

	// 调用内部函数 generateModelRankings
	rankings := generateModelRankings(analysisResults, dimensions)

	// 验证结果
	if len(rankings) != 3 {
		t.Fatalf("Expected 3 model rankings, got %d", len(rankings))
	}

	// 验证 ModelX (BrandA)
	// 平均分 should be (8+9)/2 = 8.5
	var modelX *ModelRanking
	for i := range rankings {
		if rankings[i].Model == "ModelX" && rankings[i].Brand == "BrandA" {
			modelX = &rankings[i]
			break
		}
	}

	if modelX == nil {
		t.Fatal("ModelX not found")
	}

	if modelX.CommentCount != 2 {
		t.Errorf("Expected ModelX count 2, got %d", modelX.CommentCount)
	}

	if modelX.Scores["Dim1"] != 8.5 {
		t.Errorf("Expected ModelX Dim1 score 8.5, got %f", modelX.Scores["Dim1"])
	}

	// 验证排序 (ModelZ 9.0 > ModelX 8.5 > ModelY 6.0)
	if rankings[0].Model != "ModelZ" {
		t.Errorf("Expected ModelZ to be first (9.0), got %s (%f)", rankings[0].Model, rankings[0].OverallScore)
	}
	if rankings[1].Model != "ModelX" {
		t.Errorf("Expected ModelX to be second (8.5), got %s (%f)", rankings[1].Model, rankings[1].OverallScore)
	}
	if rankings[2].Model != "ModelY" {
		t.Errorf("Expected ModelY to be third (6.0), got %s (%f)", rankings[2].Model, rankings[2].OverallScore)
	}
}
