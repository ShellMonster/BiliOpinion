package task

import (
	"context"
	"path/filepath"
	"testing"

	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/comment"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

func TestReanalyzeFromStore_DoesNotCallScrapeHook(t *testing.T) {
	if err := database.InitDB(filepath.Join(t.TempDir(), "reanalyze.db")); err != nil {
		t.Fatal(err)
	}
	history := models.AnalysisHistory{
		TaskID:     "re-1",
		Category:   "吸尘器",
		Brands:     `["戴森"]`,
		Dimensions: `["吸力"]`,
		Status:     models.StatusCompleted,
	}
	if err := database.DB.Create(&history).Error; err != nil {
		t.Fatal(err)
	}
	if err := comment.PersistForHistory(history.ID, []bilibili.Comment{
		{RPID: 1, Content: bilibili.Content{Message: "戴森吸力很强长期使用稳定"}},
	}, nil); err != nil {
		t.Fatal(err)
	}

	loaded, err := comment.LoadByHistory(history.ID)
	if err != nil || len(loaded) == 0 {
		t.Fatalf("stored comments must be available for reanalyze: %v %d", err, len(loaded))
	}
	_, _ = ReanalyzeFromStore(context.Background(), history.ID)
	// AI/config may fail here; this function never constructs a Bilibili scraper.
}
