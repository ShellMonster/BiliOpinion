package report

import (
	"encoding/json"
	"strings"
	"testing"

	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/bilibili"
)

// TestReportDataJSONSerialization 测试ReportData的JSON序列化
// 验证video_sources、sentiment_distribution、keyword_frequency字段是否正确序列化
func TestReportDataJSONSerialization(t *testing.T) {
	// 创建包含所有字段的ReportData
	input := GenerateReportInput{
		Category:   "测试商品",
		Brands:     []string{"品牌A", "品牌B"},
		Dimensions: []ai.Dimension{{Name: "测试维度", Description: "测试"}},
		AnalysisResults: map[string][]CommentWithScore{
			"品牌A": {
				{Scores: map[string]*float64{"测试维度": floatPtr(8.5)}},
			},
		},
		Stats: ReportStats{
			TotalVideos:   2,
			TotalComments: 10,
		},
		Videos: []bilibili.VideoInfo{
			{BVID: "BV1", Title: "视频1", Author: "UP1", Play: 1000, VideoReview: 50},
			{BVID: "BV2", Title: "视频2", Author: "UP2", Play: 2000, VideoReview: 100},
		},
	}

	reportData, err := GenerateReportWithInput(input)
	if err != nil {
		t.Fatalf("GenerateReportWithInput failed: %v", err)
	}

	// 序列化为JSON
	data, err := json.Marshal(reportData)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)

	// 验证必需字段存在
	requiredFields := []string{
		"video_sources",
		"sentiment_distribution",
		"keyword_frequency",
	}

	for _, field := range requiredFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing required field: %s", field)
			t.Errorf("JSON content: %s", jsonStr[:min(len(jsonStr), 500)])
		}
	}

	// 验证video_sources内容
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// 检查video_sources
	if vs, ok := result["video_sources"].([]interface{}); !ok {
		t.Errorf("video_sources is not an array")
	} else if len(vs) != 2 {
		t.Errorf("Expected 2 video_sources, got %d", len(vs))
	}

	// 检查sentiment_distribution
	if sd, ok := result["sentiment_distribution"].(map[string]interface{}); !ok {
		t.Errorf("sentiment_distribution is not an object")
	} else {
		if _, ok := sd["positive_count"]; !ok {
			t.Errorf("sentiment_distribution missing positive_count")
		}
	}

	// 检查keyword_frequency
	if _, ok := result["keyword_frequency"].([]interface{}); !ok {
		t.Errorf("keyword_frequency is not an array")
	}

	t.Logf("JSON serialization test passed!")
	t.Logf("JSON length: %d bytes", len(data))
}

func floatPtr(f float64) *float64 {
	return &f
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
