# Learnings - Resume Progress View

## Conventions
- TaskID 字段使用 UUID 格式（36 字符）
- 导航使用 `taskId`（UUID），不是 `id`（数据库 ID）
- 按钮仅在 `status === 'processing' && taskId` 时显示

## Patterns
- GORM 自动迁移：添加字段后重启服务即可
- 前端条件渲染：使用 `&&` 短路运算符

## Gotchas
- 必须使用 `history.taskId` 而不是 `history.id`
- 旧数据（无 taskId）不显示按钮

## TaskID 字段添加完成
- 位置：`backend/models/analysis_history.go` 第11行
- 字段定义：`TaskID string \`gorm:"index;size:36"\``
- 注释：任务ID（UUID格式，用于SSE连接和进度查看）
- GORM标签：index索引 + size:36（UUID标准长度）
- 重启后端服务即可自动迁移数据库表结构
## TaskId Field Added to History API

### Changes Made
1. Added TaskId field to HistoryListResponse struct (line 16)
   - Field type: string
   - JSON tag: taskId (camelCase)
   - Comment: 任务ID

2. Added TaskId mapping in response handler (line 57)
   - Maps h.TaskID from database model to response struct
   - Maintains field order consistency

### Verification
- Go compilation: ✅ PASSED (go build ./backend/...)
- No syntax errors
- Field properly mapped from models.AnalysisHistory.TaskID

### Pattern Consistency
- Follows existing struct field pattern (json tag + comment)
- Maintains alphabetical-ish field ordering (ID, TaskId, Category...)
- Uses camelCase for JSON field name (taskId)
- Uses PascalCase for Go field name (TaskId)

## Frontend-Backend Contract
- **Pattern**: Verified JSON response keys from backend (`json:"taskId"`) match frontend interface (`taskId`).
- **Learning**: Always check the API response structure when adding new fields to frontend interfaces, even if the task seems frontend-only.

## Task Completion Summary

### Task 1: 后端模型添加 TaskID 字段
- 在 AnalysisHistory 结构体中添加 TaskID 字段
- 使用 `gorm:"index;size:36"` 标签
- GORM 会在重启服务时自动迁移数据库

### Task 2: 后端 executor 保存 task_id
- 修改 createHistory 函数签名，添加 taskID 参数
- 在创建 history 时保存 TaskID 字段
- 修改调用处传入 req.TaskID

### Task 3: 后端 API 返回 taskId
- 在 HistoryListResponse 结构体中添加 TaskId 字段
- JSON 标签使用 camelCase: `json:"taskId"`
- 在响应映射中添加 TaskId 字段赋值

### Task 4: 前端添加"查看进度"按钮
- 在 HistoryItem 接口中添加 taskId?: string 字段（可选）
- 添加条件渲染：`status === 'processing' && taskId`
- 使用 navigate(`/progress/${history.taskId}`) 跳转
- 按钮位置：在"查看报告"按钮之前

## Key Learnings
- UUID 格式的 task_id 用于 SSE 连接和进度查看
- 旧数据（无 taskId）通过可选字段和条件渲染优雅处理
- 前端使用 taskId（UUID）而不是 id（数据库 ID）进行导航
