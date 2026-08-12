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
	now := time.Now()
	rows := []models.AnalysisHistory{
		{
			TaskID:        "task-recent",
			Category:      "键盘",
			Keywords:      `["机械键盘"]`,
			Brands:        `["HHKB"]`,
			Dimensions:    `["手感"]`,
			Status:        models.StatusProcessing,
			Stage:         "analyzing",
			Progress:      60,
			LastHeartbeat: now,
		},
		{
			TaskID:        "task-done",
			Category:      "耳机",
			Keywords:      `[]`,
			Brands:        `[]`,
			Dimensions:    `[]`,
			Status:        models.StatusCompleted,
			Stage:         "generating",
			Progress:      100,
			LastHeartbeat: now,
		},
		{
			TaskID:        "task-pending",
			Category:      "投影",
			Keywords:      `[]`,
			Brands:        `[]`,
			Dimensions:    `[]`,
			Status:        models.StatusPending,
			Stage:         "initializing",
			Progress:      0,
			LastHeartbeat: now,
		},
	}
	for i := range rows {
		if err := database.DB.Create(&rows[i]).Error; err != nil {
			t.Fatalf("create history: %v", err)
		}
	}

	RecoverIncompleteTasks()

	var stored models.AnalysisHistory
	if err := database.DB.Where("task_id = ?", "task-recent").First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Status != models.StatusFailed {
		t.Fatalf("leftover processing must be marked failed (no re-run), got %q", stored.Status)
	}
	if stored.Stage != "analyzing" || stored.Progress != 60 {
		t.Fatalf("recovery must not mutate stage/progress, got %q %d", stored.Stage, stored.Progress)
	}

	assertStatus(t, "task-done", models.StatusCompleted)
	assertStatus(t, "task-pending", models.StatusPending)

	var count int64
	database.DB.Model(&models.AnalysisHistory{}).Count(&count)
	if count != 3 {
		t.Fatalf("must keep existing rows only, got %d", count)
	}
}

func assertStatus(t *testing.T, taskID, want string) {
	t.Helper()
	var row models.AnalysisHistory
	if err := database.DB.Where("task_id = ?", taskID).First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.Status != want {
		t.Fatalf("%s status want %q got %q", taskID, want, row.Status)
	}
}
