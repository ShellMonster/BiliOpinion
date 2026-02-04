package api

import (
	"bilibili-analyzer/backend/ai"
	"bilibili-analyzer/backend/sse"
	"bilibili-analyzer/backend/task"
	"context"
	"log"
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
	Keywords              []string `json:"keywords"`
	VideoDateRangeMonths  int      `json:"video_date_range_months,omitempty"`   // 视频时间范围（月），0表示不限制，默认24
	MinVideoDuration      int      `json:"min_video_duration,omitempty"`        // 最小视频时长（秒），0表示不过滤
	MaxComments           int      `json:"max_comments,omitempty"`              // 最大分析评论数，默认500
	MinVideoComments      int      `json:"min_video_comments,omitempty"`        // 最小视频评论数过滤（默认0，表示不限制）
	MinCommentsPerVideo   int      `json:"min_comments_per_video,omitempty"`    // 每视频最少抓取数（默认10）
	MaxCommentsPerVideoV2 int      `json:"max_comments_per_video_v2,omitempty"` // 每视频最多抓取数（默认200）
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

		config := &task.TaskConfig{
			VideoDateRangeMonths:  req.VideoDateRangeMonths,
			MinVideoDuration:      req.MinVideoDuration,
			MaxComments:           req.MaxComments,
			MinVideoComments:      req.MinVideoComments,
			MinCommentsPerVideo:   req.MinCommentsPerVideo,
			MaxCommentsPerVideoV2: req.MaxCommentsPerVideoV2,
		}

		executor := task.NewExecutor(config)
		err := executor.Execute(context.Background(), task.TaskRequest{
			TaskID:      taskID,
			Requirement: req.Requirement,
			Brands:      req.Brands,
			Dimensions:  dimensions,
			Keywords:    req.Keywords,
		})

		if err != nil {
			log.Printf("[Task %s] Execution failed: %v", taskID, err)
			sse.PushError(taskID, err.Error())
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "任务已创建，请通过SSE接口获取实时进度",
	})
}
