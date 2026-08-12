package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"

	"github.com/gin-gonic/gin"
)

func setupHistoryTestDB(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	if err := database.InitDB(filepath.Join(t.TempDir(), "history.db")); err != nil {
		t.Fatalf("init db: %v", err)
	}
}

func TestHandleGetHistoryDetail_LooksUpByTaskID(t *testing.T) {
	setupHistoryTestDB(t)
	row := models.AnalysisHistory{
		TaskID:      "task-uuid-1",
		Category:    "吸尘器",
		Status:      models.StatusFailed,
		ProgressMsg: "请先配置AI API Key",
		ReportID:    0,
	}
	if err := database.DB.Create(&row).Error; err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "task-uuid-1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/history/task-uuid-1", nil)
	HandleGetHistoryDetail(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["status"] != models.StatusFailed {
		t.Fatalf("status=%v", payload["status"])
	}
	if payload["taskId"] != "task-uuid-1" {
		t.Fatalf("taskId=%v", payload["taskId"])
	}
	if payload["progressMsg"] != "请先配置AI API Key" {
		t.Fatalf("progressMsg=%v", payload["progressMsg"])
	}
}

func TestHandleGetHistoryDetail_CompletedIncludesReportID(t *testing.T) {
	setupHistoryTestDB(t)
	row := models.AnalysisHistory{
		TaskID:   "task-uuid-done",
		Category: "吸尘器",
		Status:   models.StatusCompleted,
		ReportID: 42,
	}
	if err := database.DB.Create(&row).Error; err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "task-uuid-done"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/history/task-uuid-done", nil)
	HandleGetHistoryDetail(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["reportId"] != float64(42) {
		t.Fatalf("reportId=%v", payload["reportId"])
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", row.ID)}}
	c2.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/history/%d", row.ID), nil)
	HandleGetHistoryDetail(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("numeric id status %d", w2.Code)
	}
}
