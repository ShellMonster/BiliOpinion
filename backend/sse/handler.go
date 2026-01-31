package sse

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleSSE SSE连接处理器
// 客户端通过此接口建立SSE连接，接收任务状态更新
//
// 请求参数：
//   - task_id: 任务ID（必填）
//
// 响应：
//   - Content-Type: text/event-stream
//   - 持续推送TaskStatus JSON数据
//
// 示例：
//
//	GET /api/sse?task_id=abc123
func HandleSSE(c *gin.Context) {
	// 获取任务ID
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id参数必填"})
		return
	}

	// 设置SSE响应头
	// Content-Type: text/event-stream 表示SSE流
	// Cache-Control: no-cache 禁止缓存
	// Connection: keep-alive 保持长连接
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲

	// 获取任务状态通道
	statusChan, exists := GetTaskChannel(taskID)
	if !exists {
		// 任务不存在，创建一个新通道
		statusChan = CreateTaskChannel(taskID)
	}

	// 获取Flusher接口，用于实时推送数据
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式响应"})
		return
	}

	// 发送初始连接成功消息
	sendSSEMessage(c.Writer, flusher, TaskStatus{
		TaskID:  taskID,
		Status:  "connected",
		Message: "SSE连接已建立",
	})

	// 监听状态更新
	for {
		select {
		case status, ok := <-statusChan:
			if !ok {
				// 通道关闭，任务结束
				return
			}

			// 发送状态更新
			sendSSEMessage(c.Writer, flusher, status)

			// 如果任务完成或出错，关闭连接
			if status.Status == StatusCompleted || status.Status == StatusError || status.Status == StatusCancelled {
				return
			}

		case <-c.Request.Context().Done():
			// 客户端断开连接
			// 清理任务通道
			CloseTaskChannel(taskID)
			return
		}
	}
}

// sendSSEMessage 发送SSE消息
// 将TaskStatus序列化为JSON并按SSE格式发送
//
// 参数：
//   - w: HTTP响应写入器
//   - flusher: Flusher接口
//   - status: 任务状态
func sendSSEMessage(w http.ResponseWriter, flusher http.Flusher, status TaskStatus) {
	// 序列化状态为JSON
	data, err := json.Marshal(status)
	if err != nil {
		return
	}

	// 按SSE格式发送
	// 格式：data: {json}\n\n
	fmt.Fprintf(w, "data: %s\n\n", data)

	// 立即刷新缓冲区，确保数据发送到客户端
	flusher.Flush()
}

// HandleSSEWithHeartbeat SSE连接处理器（带心跳）
// 定期发送心跳消息，防止连接超时
//
// 请求参数：
//   - task_id: 任务ID（必填）
//
// 响应：
//   - Content-Type: text/event-stream
//   - 持续推送TaskStatus JSON数据
//   - 每30秒发送心跳消息
func HandleSSEWithHeartbeat(c *gin.Context) {
	// 获取任务ID
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id参数必填"})
		return
	}

	// 设置SSE响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// 获取任务状态通道
	statusChan, exists := GetTaskChannel(taskID)
	if !exists {
		statusChan = CreateTaskChannel(taskID)
	}

	// 获取Flusher接口
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式响应"})
		return
	}

	// 发送初始连接成功消息
	sendSSEMessage(c.Writer, flusher, TaskStatus{
		TaskID:  taskID,
		Status:  "connected",
		Message: "SSE连接已建立",
	})

	// 创建心跳定时器（30秒）
	// heartbeatTicker := time.NewTicker(30 * time.Second)
	// defer heartbeatTicker.Stop()

	// 监听状态更新
	for {
		select {
		case status, ok := <-statusChan:
			if !ok {
				return
			}

			sendSSEMessage(c.Writer, flusher, status)

			if status.Status == StatusCompleted || status.Status == StatusError || status.Status == StatusCancelled {
				return
			}

		// case <-heartbeatTicker.C:
		// 	// 发送心跳消息
		// 	sendSSEMessage(c.Writer, flusher, TaskStatus{
		// 		TaskID:  taskID,
		// 		Status:  "heartbeat",
		// 		Message: "保持连接",
		// 	})

		case <-c.Request.Context().Done():
			CloseTaskChannel(taskID)
			return
		}
	}
}
