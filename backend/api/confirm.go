package api

import (
	"bilibili-analyzer/backend/sse"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfirmRequest 确认请求
// 用户确认AI解析结果后，前端会发送这个请求来启动实际的抓取和分析任务
type ConfirmRequest struct {
	Category   string   `json:"category"` // 商品类目
	Brands     []string `json:"brands"`   // 用户确认的品牌列表
	Dimensions []struct {
		Name        string `json:"name"`        // 维度名称
		Description string `json:"description"` // 维度描述
	} `json:"dimensions"` // 用户确认的评价维度
	Keywords []string `json:"keywords"` // 用户确认的搜索关键词
}

// HandleConfirm 处理确认请求
// 用户确认AI解析结果后，创建后台任务开始抓取和分析
func HandleConfirm(c *gin.Context) {
	// 1. 解析请求参数
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 2. 验证必填字段
	if req.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品类目不能为空"})
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

	// 3. 生成唯一的任务ID
	taskID := uuid.New().String()

	// 4. 创建任务状态通道（用于SSE推送进度）
	sse.CreateTaskChannel(taskID)

	// 5. 启动后台任务
	go func() {
		// 确保任务结束时关闭通道
		defer sse.CloseTaskChannel(taskID)

		// 推送初始状态：开始搜索视频
		sse.PushStatus(taskID, sse.TaskStatus{
			TaskID:  taskID,
			Status:  "searching",
			Message: "正在搜索视频...",
			Progress: &sse.Progress{
				Current: 0,
				Total:   100,
			},
		})

		// TODO: 实际的抓取和分析逻辑（在Task 8和Task 9实现）
		// 这里只是占位代码，后续会实现：
		// 1. 使用keywords在B站搜索视频
		// 2. 抓取视频评论
		// 3. 使用AI分析评论（按brands和dimensions分类）
		// 4. 生成分析报告
		// 5. 推送完成状态

		// 暂时推送一个示例状态
		sse.PushStatus(taskID, sse.TaskStatus{
			TaskID:  taskID,
			Status:  "pending",
			Message: "任务已创建，等待后续功能实现",
			Progress: &sse.Progress{
				Current: 10,
				Total:   100,
			},
		})
	}()

	// 6. 立即返回任务ID给前端
	// 前端可以使用这个task_id连接SSE接口来接收实时进度
	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "任务已创建，请通过SSE接口获取实时进度",
	})
}
