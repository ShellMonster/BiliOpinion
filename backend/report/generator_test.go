package report

import (
	"bilibili-analyzer/backend/ai"
	"testing"
)

// TestGenerateReport 测试报告生成功能
func TestGenerateReport(t *testing.T) {
	// 准备测试数据
	category := "手机"
	brands := []string{"小米", "华为", "苹果"}
	dimensions := []ai.Dimension{
		{Name: "性能", Description: "处理器性能和运行速度"},
		{Name: "拍照", Description: "相机质量和拍照效果"},
		{Name: "续航", Description: "电池续航能力"},
	}

	// 模拟分析结果
	score1 := 8.5
	score2 := 7.0
	score3 := 9.0
	score4 := 8.0
	score5 := 7.5
	score6 := 9.5
	score7 := 6.5
	score8 := 9.0
	score9 := 8.5

	analysisResults := map[string][]ai.AnalyzeCommentResponse{
		"小米": {
			{Scores: map[string]*float64{"性能": &score1, "拍照": &score2, "续航": &score3}},
			{Scores: map[string]*float64{"性能": &score4, "拍照": &score5, "续航": &score6}},
		},
		"华为": {
			{Scores: map[string]*float64{"性能": &score7, "拍照": &score8, "续航": &score9}},
		},
		"苹果": {
			{Scores: map[string]*float64{"性能": &score6, "拍照": &score8, "续航": &score4}},
		},
	}

	// 生成报告
	report, err := GenerateReport(category, brands, dimensions, analysisResults)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	// 验证基本信息
	if report.Category != category {
		t.Errorf("Expected category %s, got %s", category, report.Category)
	}

	if len(report.Brands) != len(brands) {
		t.Errorf("Expected %d brands, got %d", len(brands), len(report.Brands))
	}

	// 验证得分计算
	if len(report.Scores) != len(brands) {
		t.Errorf("Expected scores for %d brands, got %d", len(brands), len(report.Scores))
	}

	// 验证小米的性能得分（应该是8.5和8.0的平均值）
	xiaomiScores := report.Scores["小米"]
	expectedPerf := (8.5 + 8.0) / 2
	if xiaomiScores["性能"] != expectedPerf {
		t.Errorf("Expected 小米 性能 score %.2f, got %.2f", expectedPerf, xiaomiScores["性能"])
	}

	// 验证排名
	if len(report.Rankings) != len(brands) {
		t.Errorf("Expected %d rankings, got %d", len(brands), len(report.Rankings))
	}

	// 验证排名顺序（第一名应该有最高的综合得分）
	if report.Rankings[0].Rank != 1 {
		t.Errorf("First ranking should have rank 1, got %d", report.Rankings[0].Rank)
	}

	// 验证综合得分递减
	for i := 1; i < len(report.Rankings); i++ {
		if report.Rankings[i].OverallScore > report.Rankings[i-1].OverallScore {
			t.Errorf("Rankings not sorted correctly: rank %d score %.2f > rank %d score %.2f",
				i+1, report.Rankings[i].OverallScore,
				i, report.Rankings[i-1].OverallScore)
		}
	}

	// 验证购买建议不为空
	if report.Recommendation == "" {
		t.Error("Recommendation should not be empty")
	}

	t.Logf("Report generated successfully:")
	t.Logf("  Category: %s", report.Category)
	t.Logf("  Rankings:")
	for _, r := range report.Rankings {
		t.Logf("    %d. %s (%.2f分)", r.Rank, r.Brand, r.OverallScore)
	}
	t.Logf("  Recommendation: %s", report.Recommendation)
}

// TestGenerateRankings 测试排名生成
func TestGenerateRankings(t *testing.T) {
	brands := []string{"品牌A", "品牌B", "品牌C"}
	dimensions := []ai.Dimension{
		{Name: "维度1", Description: "测试维度1"},
		{Name: "维度2", Description: "测试维度2"},
	}

	scores := map[string]map[string]float64{
		"品牌A": {"维度1": 8.0, "维度2": 7.0}, // 平均7.5
		"品牌B": {"维度1": 9.0, "维度2": 9.0}, // 平均9.0
		"品牌C": {"维度1": 6.0, "维度2": 8.0}, // 平均7.0
	}

	rankings := generateRankings(brands, dimensions, scores)

	// 验证排名数量
	if len(rankings) != len(brands) {
		t.Errorf("Expected %d rankings, got %d", len(brands), len(rankings))
	}

	// 验证第一名是品牌B
	if rankings[0].Brand != "品牌B" {
		t.Errorf("Expected 品牌B to be first, got %s", rankings[0].Brand)
	}

	// 验证综合得分
	if rankings[0].OverallScore != 9.0 {
		t.Errorf("Expected 品牌B overall score 9.0, got %.2f", rankings[0].OverallScore)
	}

	// 验证排名号
	for i, r := range rankings {
		if r.Rank != i+1 {
			t.Errorf("Expected rank %d, got %d for %s", i+1, r.Rank, r.Brand)
		}
	}
}

// TestGenerateRecommendation 测试购买建议生成
func TestGenerateRecommendation(t *testing.T) {
	dimensions := []ai.Dimension{
		{Name: "性能", Description: "性能测试"},
		{Name: "价格", Description: "价格测试"},
	}

	rankings := []BrandRanking{
		{
			Brand:        "品牌A",
			OverallScore: 8.5,
			Rank:         1,
			Scores:       map[string]float64{"性能": 9.0, "价格": 8.0},
		},
		{
			Brand:        "品牌B",
			OverallScore: 7.5,
			Rank:         2,
			Scores:       map[string]float64{"性能": 7.0, "价格": 8.0},
		},
	}

	recommendation := generateRecommendation(rankings, dimensions)

	// 验证推荐文本包含第一名品牌
	if recommendation == "" {
		t.Error("Recommendation should not be empty")
	}

	// 应该包含品牌A
	if !contains(recommendation, "品牌A") {
		t.Error("Recommendation should mention the top brand")
	}

	// 应该包含得分
	if !contains(recommendation, "8.5") {
		t.Error("Recommendation should mention the overall score")
	}

	t.Logf("Generated recommendation: %s", recommendation)
}

// TestGenerateRecommendationEmpty 测试空排名的购买建议
func TestGenerateRecommendationEmpty(t *testing.T) {
	dimensions := []ai.Dimension{}
	rankings := []BrandRanking{}

	recommendation := generateRecommendation(rankings, dimensions)

	if recommendation != "暂无足够数据生成购买建议" {
		t.Errorf("Expected default message for empty rankings, got: %s", recommendation)
	}
}

// contains 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
