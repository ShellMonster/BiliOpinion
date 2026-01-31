package ai

import (
	"context"
	"testing"
)

func TestParseAnalysisResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
		check    func(*AnalyzeCommentResponse) bool
	}{
		{
			name:     "valid JSON",
			response: `{"scores":{"吸力":8.5,"续航":7.0,"噪音":null}}`,
			wantErr:  false,
			check: func(r *AnalyzeCommentResponse) bool {
				if r.Scores["吸力"] == nil || *r.Scores["吸力"] != 8.5 {
					return false
				}
				if r.Scores["续航"] == nil || *r.Scores["续航"] != 7.0 {
					return false
				}
				if r.Scores["噪音"] != nil {
					return false
				}
				return true
			},
		},
		{
			name:     "JSON with extra text",
			response: `根据评论分析结果如下：{"scores":{"吸力":9.0,"续航":8.0}}`,
			wantErr:  false,
			check: func(r *AnalyzeCommentResponse) bool {
				return r.Scores["吸力"] != nil && *r.Scores["吸力"] == 9.0
			},
		},
		{
			name:     "invalid JSON",
			response: `这不是JSON`,
			wantErr:  true,
			check:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAnalysisResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAnalysisResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(result) {
				t.Errorf("parseAnalysisResponse() result check failed")
			}
		})
	}
}

func TestAnalyzeCommentValidation(t *testing.T) {
	client := NewClient(Config{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})

	ctx := context.Background()

	t.Run("empty comment", func(t *testing.T) {
		_, err := client.AnalyzeComment(ctx, AnalyzeCommentRequest{
			Comment:    "",
			Dimensions: []Dimension{{Name: "吸力", Description: "吸尘效果"}},
		})
		if err == nil {
			t.Error("expected error for empty comment")
		}
	})

	t.Run("empty dimensions", func(t *testing.T) {
		_, err := client.AnalyzeComment(ctx, AnalyzeCommentRequest{
			Comment:    "这个吸尘器很好用",
			Dimensions: []Dimension{},
		})
		if err == nil {
			t.Error("expected error for empty dimensions")
		}
	})
}

func TestAnalyzeCommentsBatchValidation(t *testing.T) {
	client := NewClient(Config{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})

	ctx := context.Background()

	t.Run("empty comments", func(t *testing.T) {
		_, err := client.AnalyzeCommentsBatch(ctx, []CommentInput{}, []Dimension{{Name: "吸力", Description: "吸尘效果"}})
		if err == nil {
			t.Error("expected error for empty comments")
		}
	})

	t.Run("empty dimensions", func(t *testing.T) {
		_, err := client.AnalyzeCommentsBatch(ctx, []CommentInput{{ID: "1", Content: "好用"}}, []Dimension{})
		if err == nil {
			t.Error("expected error for empty dimensions")
		}
	})
}

func TestCommentInput(t *testing.T) {
	input := CommentInput{
		ID:      "12345",
		Content: "这个产品非常好用，吸力很强",
	}

	if input.ID != "12345" {
		t.Errorf("expected ID '12345', got '%s'", input.ID)
	}
	if input.Content != "这个产品非常好用，吸力很强" {
		t.Errorf("unexpected content")
	}
}

func TestCommentAnalysisResult(t *testing.T) {
	score := 8.5
	result := CommentAnalysisResult{
		CommentID: "12345",
		Content:   "吸力很强",
		Scores:    map[string]*float64{"吸力": &score},
		Error:     "",
	}

	if result.CommentID != "12345" {
		t.Errorf("expected CommentID '12345', got '%s'", result.CommentID)
	}
	if result.Scores["吸力"] == nil || *result.Scores["吸力"] != 8.5 {
		t.Error("unexpected score")
	}
	if result.Error != "" {
		t.Error("expected no error")
	}
}
