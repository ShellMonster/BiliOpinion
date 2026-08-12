package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"

	"github.com/gin-gonic/gin"
)

func setupConfigTestDB(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	if err := database.InitDB(filepath.Join(dir, "test.db")); err != nil {
		t.Fatalf("init db: %v", err)
	}
}

func seedSetting(t *testing.T, key, value string) {
	t.Helper()
	if err := database.DB.Create(&models.Settings{Key: key, Value: value}).Error; err != nil {
		t.Fatalf("seed %s: %v", key, err)
	}
}

func TestHandleGetConfig_MasksSecrets(t *testing.T) {
	setupConfigTestDB(t)
	fullKey := "sk-abcdefghijklmnopqrstuvwxyz"
	fullCookie := "SESSDATA=very-secret-session-token; bili_jct=abc"
	seedSetting(t, models.SettingKeyAIAPIKey, fullKey)
	seedSetting(t, models.SettingKeyBilibiliCookie, fullCookie)
	seedSetting(t, models.SettingKeyAIAPIBase, "https://api.openai.com/v1")
	seedSetting(t, models.SettingKeyAIModel, "gemini-3-flash-preview")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/config", nil)
	HandleGetConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if bytes.Contains(w.Body.Bytes(), []byte(fullKey)) {
		t.Fatal("response leaked full AI API key")
	}
	if bytes.Contains(w.Body.Bytes(), []byte(fullCookie)) {
		t.Fatal("response leaked full Bilibili cookie")
	}
	if bytes.Contains(w.Body.Bytes(), []byte("SESSDATA=very-secret-session-token")) {
		t.Fatal("response leaked SESSDATA value")
	}
	if bytes.Contains(w.Body.Bytes(), []byte("sk-abcdefghijklmnopqrstuv")) {
		t.Fatal("response leaked AI key prefix")
	}
	if bytes.Contains(w.Body.Bytes(), []byte("very-secret-session-token")) {
		t.Fatal("response leaked cookie token fragment")
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	key, _ := payload["ai_api_key"].(string)
	cookie, _ := payload["bilibili_cookie"].(string)
	if key != MaskSecret(fullKey) {
		t.Fatalf("ai_api_key want %q got %q", MaskSecret(fullKey), key)
	}
	if cookie != MaskSecret(fullCookie) {
		t.Fatalf("bilibili_cookie want %q got %q", MaskSecret(fullCookie), cookie)
	}
}

func TestHandleSaveConfig_EmptySecretDoesNotOverwrite(t *testing.T) {
	setupConfigTestDB(t)
	fullKey := "sk-keep-this-original-key-9999"
	seedSetting(t, models.SettingKeyAIAPIKey, fullKey)
	seedSetting(t, models.SettingKeyBilibiliCookie, "SESSDATA=keep-me-cookie-aaaa")

	reqBody := `{"ai_base_url":"https://api.openai.com/v1","ai_api_key":"","ai_model":"m","bilibili_cookie":"","scrape_max_concurrency":"5","ai_max_concurrency":"10"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	HandleSaveConfig(c)
	if w.Code != http.StatusOK {
		t.Fatalf("save status %d body %s", w.Code, w.Body.String())
	}

	var stored models.Settings
	if err := database.DB.Where("key = ?", models.SettingKeyAIAPIKey).First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Value != fullKey {
		t.Fatalf("empty save overwrote API key: got %q", stored.Value)
	}
	var cookie models.Settings
	if err := database.DB.Where("key = ?", models.SettingKeyBilibiliCookie).First(&cookie).Error; err != nil {
		t.Fatal(err)
	}
	if cookie.Value != "SESSDATA=keep-me-cookie-aaaa" {
		t.Fatalf("empty save overwrote cookie: got %q", cookie.Value)
	}
	if bytes.Contains(w.Body.Bytes(), []byte(fullKey)) {
		t.Fatal("save response leaked API key")
	}
	if bytes.Contains(w.Body.Bytes(), []byte("SESSDATA=keep-me-cookie-aaaa")) {
		t.Fatal("save response leaked cookie")
	}
}

func TestHandleSaveConfig_MaskedSecretDoesNotOverwrite(t *testing.T) {
	setupConfigTestDB(t)
	fullKey := "sk-keep-this-original-key-9999"
	fullCookie := "SESSDATA=keep-me-cookie-aaaa"
	seedSetting(t, models.SettingKeyAIAPIKey, fullKey)
	seedSetting(t, models.SettingKeyBilibiliCookie, fullCookie)

	reqBody := `{"ai_api_key":"` + MaskSecret(fullKey) + `","bilibili_cookie":"` + MaskSecret(fullCookie) + `"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	HandleSaveConfig(c)

	var stored models.Settings
	if err := database.DB.Where("key = ?", models.SettingKeyAIAPIKey).First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Value != fullKey {
		t.Fatalf("masked save overwrote API key: got %q", stored.Value)
	}
	var cookie models.Settings
	if err := database.DB.Where("key = ?", models.SettingKeyBilibiliCookie).First(&cookie).Error; err != nil {
		t.Fatal(err)
	}
	if cookie.Value != fullCookie {
		t.Fatalf("masked save overwrote cookie: got %q", cookie.Value)
	}
	if bytes.Contains(w.Body.Bytes(), []byte(fullKey)) {
		t.Fatal("save response leaked API key")
	}
	if bytes.Contains(w.Body.Bytes(), []byte(fullCookie)) {
		t.Fatal("save response leaked cookie")
	}
}
