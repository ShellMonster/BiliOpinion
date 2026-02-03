package task

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/database"
	"bilibili-analyzer/backend/models"
	"bilibili-analyzer/backend/sse"
	"context"
	"encoding/json"
	"log"
	"time"
)

// RecoverIncompleteTasks 恢复未完成的任务
// 在后端启动时调用，检查 processing 状态的任务并尝试恢复
func RecoverIncompleteTasks() {
	var tasks []models.AnalysisHistory

	// 查找所有 processing 状态的任务
	database.DB.Where("status = ?", models.StatusProcessing).Find(&tasks)

	if len(tasks) == 0 {
		log.Println("[Recovery] No incomplete tasks found")
		return
	}

	log.Printf("[Recovery] Found %d incomplete tasks", len(tasks))

	for _, task := range tasks {
		// 检查 LastHeartbeat 是否为零值（可能是旧数据或数据库迁移问题）
		if task.LastHeartbeat.IsZero() {
			log.Printf("[Recovery] Task %s has zero LastHeartbeat (old data?), skipping timeout check and attempting recovery", task.TaskID)
			// 更新 LastHeartbeat 为当前时间，然后尝试恢复
			database.DB.Model(&task).Update("last_heartbeat", time.Now())
			go recoverTask(task)
			continue
		}

		// 检查是否超时（超过1小时）
		timeSinceHeartbeat := time.Since(task.LastHeartbeat)
		if timeSinceHeartbeat > time.Hour {
			log.Printf("[Recovery] Task %s timed out (last heartbeat: %v, %v ago), marking as failed",
				task.TaskID, task.LastHeartbeat.Format("2006-01-02 15:04:05"), timeSinceHeartbeat)
			database.DB.Model(&task).Update("status", models.StatusFailed)
			continue
		}

		log.Printf("[Recovery] Task %s is recent (last heartbeat: %v ago), attempting recovery",
			task.TaskID, timeSinceHeartbeat)
		// 尝试恢复任务
		go recoverTask(task)
	}
}

func recoverTask(history models.AnalysisHistory) {
	taskID := history.TaskID
	log.Printf("[Recovery] Recovering task %s from stage %s", taskID, history.Stage)

	// 创建 SSE 通道
	sse.CreateTaskChannel(taskID)
	defer sse.CloseTaskChannel(taskID)

	// 解析任务配置
	var config TaskConfig
	if history.TaskConfig != "" {
		if err := json.Unmarshal([]byte(history.TaskConfig), &config); err != nil {
			log.Printf("[Recovery] Failed to unmarshal task config for task %s: %v, using defaults", taskID, err)
			config = DefaultTaskConfig()
		}
	} else {
		config = DefaultTaskConfig()
	}

	// 解析任务请求参数
	var keywords, brands []string
	if err := json.Unmarshal([]byte(history.Keywords), &keywords); err != nil {
		log.Printf("[Recovery] Failed to unmarshal keywords for task %s: %v", taskID, err)
		database.DB.Model(&history).Update("status", models.StatusFailed)
		return
	}
	if err := json.Unmarshal([]byte(history.Brands), &brands); err != nil {
		log.Printf("[Recovery] Failed to unmarshal brands for task %s: %v", taskID, err)
		database.DB.Model(&history).Update("status", models.StatusFailed)
		return
	}

	// 解析评价维度
	var dimNames []string
	if err := json.Unmarshal([]byte(history.Dimensions), &dimNames); err != nil {
		log.Printf("[Recovery] Failed to unmarshal dimensions for task %s: %v", taskID, err)
		database.DB.Model(&history).Update("status", models.StatusFailed)
		return
	}
	dimensions := make([]ai.Dimension, len(dimNames))
	for i, name := range dimNames {
		dimensions[i] = ai.Dimension{
			Name:        name,
			Description: "", // 恢复时不需要描述
		}
	}

	// 推送恢复状态
	sse.PushProgress(taskID, sse.StatusSearching, history.Progress, 100,
		"任务恢复中: "+history.ProgressMsg)

	// 重新执行任务
	executor := NewExecutor(&config)
	err := executor.Execute(context.Background(), TaskRequest{
		TaskID:      taskID,
		Requirement: history.Category,
		Brands:      brands,
		Dimensions:  dimensions,
		Keywords:    keywords,
	})

	if err != nil {
		log.Printf("[Recovery] Task %s recovery failed: %v", taskID, err)
		database.DB.Model(&history).Update("status", models.StatusFailed)
	}
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
