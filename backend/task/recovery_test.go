package task

import (
	"path/filepath"
	"testing"
	"time"

	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

func setupRecoveryDB(t *testing.T) {
	t.Helper()
	if err := database.InitDB(filepath.Join(t.TempDir(), "recovery.db")); err != nil {
		t.Fatalf("init db: %v", err)
	}
}

func TestRecoverIncompleteTasks_MarksProcessingFailedWithoutNewRow(t *testing.T) {
	setupRecoveryDB(t)
	row := models.AnalysisHistory{
		TaskID:        "task-same-id",
		Category:      "吸尘器",
		Keywords:      `["吸尘器评测"]`,
		Brands:        `["戴森"]`,
		Dimensions:    `["吸力"]`,
		Status:        models.StatusProcessing,
		Stage:         "scraping",
		Progress:      30,
		LastHeartbeat: time.Now().Add(-2 * time.Hour),
	}
	if err := database.DB.Create(&row).Error; err != nil {
		t.Fatalf("create history: %v", err)
	}

	RecoverIncompleteTasks()

	var count int64
	if err := database.DB.Model(&models.AnalysisHistory{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("recovery must not insert a second history row, got %d", count)
	}

	var stored models.AnalysisHistory
	if err := database.DB.Where("task_id = ?", "task-same-id").First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Status != models.StatusFailed {
		t.Fatalf("timed-out processing must become failed, got %q", stored.Status)
	}
}

func TestRecoverIncompleteTasks_RecentProcessingAlsoFailsNoRerun(t *testing.T) {
	setupRecoveryDB(t)
	row := models.AnalysisHistory{
		TaskID:        "task-recent",
		Category:      "键盘",
		Keywords:      `["机械键盘"]`,
		Brands:        `["HHKB"]`,
		Dimensions:    `["手感"]`,
		Status:        models.StatusProcessing,
		Stage:         "analyzing",
		Progress:      60,
		LastHeartbeat: time.Now(),
	}
	if err := database.DB.Create(&row).Error; err != nil {
		t.Fatalf("create history: %v", err)
	}

	RecoverIncompleteTasks()

	var stored models.AnalysisHistory
	if err := database.DB.Where("task_id = ?", "task-recent").First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Status != models.StatusFailed {
		t.Fatalf("leftover processing must be marked failed (no re-run), got %q", stored.Status)
	}

	var count int64
	database.DB.Model(&models.AnalysisHistory{}).Count(&count)
	if count != 1 {
		t.Fatalf("must keep a single row, got %d", count)
	}
}
