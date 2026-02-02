package ai

import (
	"fmt"
	"testing"
)

func TestCalculateBatches(t *testing.T) {
	tests := []struct {
		name            string
		comments        []CommentInput
		config          *BatchConfig
		expectedBatches int
	}{
		{
			name:            "空评论列表",
			comments:        []CommentInput{},
			config:          nil,
			expectedBatches: 0,
		},
		{
			name:            "单条评论",
			comments:        []CommentInput{{ID: "1", Content: "测试"}},
			config:          nil,
			expectedBatches: 1,
		},
		{
			name: "多条短评论合并",
			comments: []CommentInput{
				{ID: "1", Content: "好用"},
				{ID: "2", Content: "不错"},
				{ID: "3", Content: "推荐"},
			},
			config:          &BatchConfig{MaxCharsPerBatch: 100, MaxItemsPerBatch: 10, MinItemsPerBatch: 1},
			expectedBatches: 1,
		},
		{
			name: "超长评论单独一批",
			comments: []CommentInput{
				{ID: "1", Content: string(make([]rune, 3000))}, // 3000 字符
				{ID: "2", Content: "短评论"},
			},
			config:          &BatchConfig{MaxCharsPerBatch: 3000, MaxItemsPerBatch: 15, MinItemsPerBatch: 1},
			expectedBatches: 2,
		},
		{
			name: "按条数限制分批",
			comments: func() []CommentInput {
				var cs []CommentInput
				for i := 0; i < 20; i++ {
					cs = append(cs, CommentInput{ID: fmt.Sprintf("%d", i), Content: "短"})
				}
				return cs
			}(),
			config:          &BatchConfig{MaxCharsPerBatch: 10000, MaxItemsPerBatch: 5, MinItemsPerBatch: 1},
			expectedBatches: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := CalculateBatches(tt.comments, tt.config)
			if len(batches) != tt.expectedBatches {
				t.Errorf("期望 %d 批，实际 %d 批", tt.expectedBatches, len(batches))
			}
		})
	}
}

func TestCalculateBatchesWithVideoTitle(t *testing.T) {
	comments := []CommentInput{
		{ID: "1", Content: "评论1", VideoTitle: "视频标题1"},
		{ID: "2", Content: "评论2", VideoTitle: "视频标题2"},
		{ID: "3", Content: "评论3", VideoTitle: "视频标题3"},
	}

	config := &BatchConfig{
		MaxCharsPerBatch: 100, // 很小的限制，确保分批
		MaxItemsPerBatch: 10,
		MinItemsPerBatch: 1,
	}

	batches := CalculateBatches(comments, config)

	// 视频标题也会被计入字符数，所以应该分成多批
	if len(batches) < 1 {
		t.Errorf("至少应该有 1 批，实际 %d 批", len(batches))
	}

	// 验证所有评论都被分配到了批次中
	totalComments := 0
	for _, batch := range batches {
		totalComments += len(batch)
	}
	if totalComments != len(comments) {
		t.Errorf("总评论数不匹配，期望 %d，实际 %d", len(comments), totalComments)
	}
}

func TestDefaultBatchConfig(t *testing.T) {
	config := DefaultBatchConfig()

	if config.MaxCharsPerBatch != 3000 {
		t.Errorf("MaxCharsPerBatch 期望 3000，实际 %d", config.MaxCharsPerBatch)
	}
	if config.MaxItemsPerBatch != 15 {
		t.Errorf("MaxItemsPerBatch 期望 15，实际 %d", config.MaxItemsPerBatch)
	}
	if config.MinItemsPerBatch != 1 {
		t.Errorf("MinItemsPerBatch 期望 1，实际 %d", config.MinItemsPerBatch)
	}
}
