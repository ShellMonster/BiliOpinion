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

	scrapeCalled := false
	prev := scrapeHook
	scrapeHook = func() { scrapeCalled = true }
	defer func() { scrapeHook = prev }()

	_, err := ReanalyzeFromStore(context.Background(), history.ID)
	if scrapeCalled {
		t.Fatal("re-analyze invoked the scrape hook")
	}
	if err == nil {
		return
	}
	// AI/config may fail in this environment; the contract is no scrape.
}
