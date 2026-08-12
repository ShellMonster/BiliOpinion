package task

import (
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"log"
	"time"
)

// RecoverIncompleteTasks marks leftover processing tasks as failed.
// It never re-runs search/scrape/AI and never inserts a new history row.
func RecoverIncompleteTasks() {
	result := database.DB.Model(&models.AnalysisHistory{}).
		Where("status = ?", models.StatusProcessing).
		Update("status", models.StatusFailed)

	if result.Error != nil {
		log.Printf("[Recovery] Failed to mark incomplete tasks: %v", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		log.Println("[Recovery] No incomplete tasks found")
		return
	}
	log.Printf("[Recovery] Marked %d leftover processing task(s) as failed (no re-run)", result.RowsAffected)
}

// CleanupTimedOutTasks 清理超时任务
func CleanupTimedOutTasks() {
	cutoff := time.Now().Add(-time.Hour)
	result := database.DB.Model(&models.AnalysisHistory{}).
		Where("status = ? AND last_heartbeat < ?", models.StatusProcessing, cutoff).
		Update("status", models.StatusFailed)

	if result.RowsAffected > 0 {
		log.Printf("[Cleanup] Marked %d timed out tasks as failed", result.RowsAffected)
	}
}
