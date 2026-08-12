package comment

import (
	"path/filepath"
	"testing"

	"bilibili-analyzer/backend/bilibili"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

func TestPersistAndLoadByHistory(t *testing.T) {
	if err := database.InitDB(filepath.Join(t.TempDir(), "comments.db")); err != nil {
		t.Fatal(err)
	}
	history := models.AnalysisHistory{TaskID: "t1", Category: "吸尘器", Status: models.StatusCompleted}
	if err := database.DB.Create(&history).Error; err != nil {
		t.Fatal(err)
	}

	comments := []bilibili.Comment{
		{RPID: 11, Like: 3, Content: bilibili.Content{Message: "戴森吸力很强，用了半年还行"}, Member: bilibili.Member{Uname: "u1"}},
		{RPID: 12, Like: 1, Content: bilibili.Content{Message: "小米续航一般但是便宜"}, Member: bilibili.Member{Uname: "u2"}},
	}
	if err := PersistForHistory(history.ID, comments, map[string]string{"11": "BV1xx", "12": "BV1yy"}); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadByHistory(history.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("want 2 comments, got %d", len(loaded))
	}
	if loaded[0].Content == "" || loaded[1].Content == "" {
		t.Fatal("stored comments lost content")
	}
}
