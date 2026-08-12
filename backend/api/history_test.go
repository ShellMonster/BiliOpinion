package api

import (
	"encoding/json"
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
