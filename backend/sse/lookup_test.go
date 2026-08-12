package sse

import (
	"path/filepath"
	"testing"

	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

func TestLookupTerminalTask_UnknownIsNotFound(t *testing.T) {
	if err := database.InitDB(filepath.Join(t.TempDir(), "sse.db")); err != nil {
		t.Fatal(err)
	}
	if _, ok := lookupTerminalTask("missing"); ok {
		t.Fatal("unknown task must not look terminal")
	}
}

func TestLookupTerminalTask_FailedHistory(t *testing.T) {
	if err := database.InitDB(filepath.Join(t.TempDir(), "sse2.db")); err != nil {
		t.Fatal(err)
	}
	row := models.AnalysisHistory{
		TaskID:      "dead-task",
		Category:    "x",
		Status:      models.StatusFailed,
		ProgressMsg: "请先配置AI API Key",
	}
	if err := database.DB.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	status, ok := lookupTerminalTask("dead-task")
	if !ok {
		t.Fatal("expected terminal status")
	}
	if status.Status != StatusError {
		t.Fatalf("status=%s", status.Status)
	}
	if status.Error != "请先配置AI API Key" {
		t.Fatalf("error=%s", status.Error)
	}
}
