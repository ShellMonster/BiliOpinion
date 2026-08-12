package sse

import (
	"fmt"

	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
)

// lookupTerminalTask returns a terminal SSE status for a finished history row.
func lookupTerminalTask(taskID string) (TaskStatus, bool) {
	if database.DB == nil {
		return TaskStatus{}, false
	}
	var history models.AnalysisHistory
	if err := database.DB.Where("task_id = ?", taskID).First(&history).Error; err != nil {
		return TaskStatus{}, false
	}
	switch history.Status {
	case models.StatusCompleted:
		status := TaskStatus{
			TaskID:  taskID,
			Status:  StatusCompleted,
			Message: history.ProgressMsg,
		}
		if history.ReportID > 0 {
			status.Progress = &Progress{
				Current: 100,
				Total:   100,
				Stage:   fmt.Sprintf("%d", history.ReportID),
			}
		}
		return status, true
	case models.StatusFailed:
		msg := history.ProgressMsg
		if msg == "" {
			msg = "任务执行失败"
		}
		return TaskStatus{
			TaskID:  taskID,
			Status:  StatusError,
			Message: msg,
			Error:   msg,
		}, true
	default:
		return TaskStatus{}, false
	}
}
