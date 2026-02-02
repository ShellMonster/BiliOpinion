package sse

import (
	"sync"
)

// taskChannels 存储所有任务的状态通道
// key: 任务ID, value: 状态通道
var taskChannels = make(map[string]chan TaskStatus)

// taskLastStatus 缓存每个任务的最后状态，用于新客户端连接时发送当前进度
var taskLastStatus = make(map[string]TaskStatus)

// mu 保护taskChannels和taskLastStatus的读写锁
var mu sync.RWMutex

// TaskStatus 任务状态结构
// 用于SSE推送任务执行进度
type TaskStatus struct {
	TaskID   string    `json:"task_id"`            // 任务ID（唯一标识）
	Status   string    `json:"status"`             // 任务状态
	Progress *Progress `json:"progress,omitempty"` // 进度信息（可选）
	Message  string    `json:"message,omitempty"`  // 状态消息（可选）
	Error    string    `json:"error,omitempty"`    // 错误信息（可选）
}

// 任务状态常量
const (
	StatusParsing        = "parsing"         // 正在解析用户输入
	StatusWaitingConfirm = "waiting_confirm" // 等待用户确认
	StatusSearching      = "searching"       // 正在搜索视频
	StatusScraping       = "scraping"        // 正在抓取评论
	StatusAnalyzing      = "analyzing"       // 正在AI分析
	StatusGenerating     = "generating"      // 正在生成报告
	StatusCompleted      = "completed"       // 任务完成
	StatusError          = "error"           // 任务出错
	StatusCancelled      = "cancelled"       // 任务取消
)

// Progress 进度信息结构
type Progress struct {
	Current int    `json:"current"` // 当前进度
	Total   int    `json:"total"`   // 总数
	Stage   string `json:"stage"`   // 当前阶段描述
}

// CreateTaskChannel 创建任务状态通道
// 为新任务创建一个带缓冲的通道，用于推送状态更新
//
// 参数：
//   - taskID: 任务ID
//
// 返回：
//   - chan TaskStatus: 任务状态通道
//
// 示例：
//
//	ch := CreateTaskChannel("task_123")
//	ch <- TaskStatus{TaskID: "task_123", Status: "processing"}
func CreateTaskChannel(taskID string) chan TaskStatus {
	mu.Lock()
	defer mu.Unlock()

	// 如果通道已存在，先关闭旧通道
	if ch, exists := taskChannels[taskID]; exists {
		close(ch)
	}

	// 创建带缓冲的通道（缓冲100条消息，防止阻塞）
	ch := make(chan TaskStatus, 100)
	taskChannels[taskID] = ch
	return ch
}

// GetTaskChannel 获取任务状态通道
// 根据任务ID获取对应的状态通道
//
// 参数：
//   - taskID: 任务ID
//
// 返回：
//   - chan TaskStatus: 任务状态通道（如果不存在返回nil）
//   - bool: 通道是否存在
//
// 示例：
//
//	ch, exists := GetTaskChannel("task_123")
//	if exists {
//	    status := <-ch
//	}
func GetTaskChannel(taskID string) (chan TaskStatus, bool) {
	mu.RLock()
	defer mu.RUnlock()
	ch, exists := taskChannels[taskID]
	return ch, exists
}

// CloseTaskChannel 关闭任务状态通道
// 任务完成或取消时调用，关闭通道并从map中删除
//
// 参数：
//   - taskID: 任务ID
//
// 示例：
//
//	CloseTaskChannel("task_123")
func CloseTaskChannel(taskID string) {
	mu.Lock()
	defer mu.Unlock()

	if ch, exists := taskChannels[taskID]; exists {
		close(ch)
		delete(taskChannels, taskID)
	}
	delete(taskLastStatus, taskID)
}

// GetLastStatus 获取任务的最后状态，用于新客户端连接时发送当前进度
func GetLastStatus(taskID string) (TaskStatus, bool) {
	mu.RLock()
	defer mu.RUnlock()
	status, exists := taskLastStatus[taskID]
	return status, exists
}

// PushStatus 推送任务状态
// 向指定任务的通道推送状态更新
// 如果通道已满或不存在，消息会被丢弃（非阻塞）
//
// 参数：
//   - taskID: 任务ID
//   - status: 任务状态
//
// 示例：
//
//	PushStatus("task_123", TaskStatus{
//	    TaskID: "task_123",
//	    Status: StatusSearching,
//	    Progress: &Progress{Current: 5, Total: 50},
//	    Message: "正在搜索视频...",
//	})
func PushStatus(taskID string, status TaskStatus) {
	mu.Lock()
	taskLastStatus[taskID] = status
	ch := taskChannels[taskID]
	mu.Unlock()

	if ch != nil {
		select {
		case ch <- status:
		default:
		}
	}
}

// PushProgress 推送进度更新（便捷方法）
// 快速推送进度信息
//
// 参数：
//   - taskID: 任务ID
//   - status: 状态字符串
//   - current: 当前进度
//   - total: 总数
//   - message: 状态消息
//
// 示例：
//
//	PushProgress("task_123", StatusScraping, 10, 50, "正在抓取第10个视频的评论...")
func PushProgress(taskID, status string, current, total int, message string) {
	PushStatus(taskID, TaskStatus{
		TaskID: taskID,
		Status: status,
		Progress: &Progress{
			Current: current,
			Total:   total,
		},
		Message: message,
	})
}

// PushError 推送错误状态（便捷方法）
// 快速推送错误信息
//
// 参数：
//   - taskID: 任务ID
//   - errMsg: 错误消息
//
// 示例：
//
//	PushError("task_123", "搜索失败：网络超时")
func PushError(taskID, errMsg string) {
	PushStatus(taskID, TaskStatus{
		TaskID:  taskID,
		Status:  StatusError,
		Message: errMsg,
		Error:   errMsg,
	})
}

// PushCompleted 推送完成状态（便捷方法）
// 快速推送任务完成信息
//
// 参数：
//   - taskID: 任务ID
//   - message: 完成消息
//
// 示例：
//
//	PushCompleted("task_123", "分析完成，共处理50个视频，500条评论")
func PushCompleted(taskID, message string) {
	PushStatus(taskID, TaskStatus{
		TaskID:  taskID,
		Status:  StatusCompleted,
		Message: message,
	})
}

// GetActiveTaskCount 获取活跃任务数量
// 返回当前正在进行的任务数量
//
// 返回：
//   - int: 活跃任务数量
func GetActiveTaskCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(taskChannels)
}

// GetActiveTaskIDs 获取所有活跃任务ID
// 返回当前所有正在进行的任务ID列表
//
// 返回：
//   - []string: 任务ID列表
func GetActiveTaskIDs() []string {
	mu.RLock()
	defer mu.RUnlock()

	ids := make([]string, 0, len(taskChannels))
	for id := range taskChannels {
		ids = append(ids, id)
	}
	return ids
}
