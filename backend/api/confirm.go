package api

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/sse"
	"bilibili-analyzer/backend/task"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConfirmRequest struct {
	Requirement string   `json:"requirement"`
	Brands      []string `json:"brands"`
	Dimensions  []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"dimensions"`
	Keywords []string `json:"keywords"`
}

func HandleConfirm(c *gin.Context) {
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if req.Requirement == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需求描述不能为空"})
		return
	}
	if len(req.Brands) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "品牌列表不能为空"})
		return
	}
	if len(req.Dimensions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评价维度不能为空"})
		return
	}
	if len(req.Keywords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	taskID := uuid.New().String()
	sse.CreateTaskChannel(taskID)

	dimensions := make([]ai.Dimension, len(req.Dimensions))
	for i, d := range req.Dimensions {
		dimensions[i] = ai.Dimension{
			Name:        d.Name,
			Description: d.Description,
		}
	}

	go func() {
		defer sse.CloseTaskChannel(taskID)

		executor := task.NewExecutor(nil)
		err := executor.Execute(context.Background(), task.TaskRequest{
			TaskID:      taskID,
			Requirement: req.Requirement,
			Brands:      req.Brands,
			Dimensions:  dimensions,
			Keywords:    req.Keywords,
		})

		if err != nil {
			sse.PushError(taskID, err.Error())
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "任务已创建，请通过SSE接口获取实时进度",
	})
}
