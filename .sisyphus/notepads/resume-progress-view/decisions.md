# Decisions - Resume Progress View

## Architectural Choices
- 使用 UUID 作为 task_id（而不是数据库自增 ID）
- SSE 重连后进度从 0% 开始（可接受的限制）

## Trade-offs
- 不添加进度状态持久化（超出范围）
- 不修改 SSE manager 核心逻辑

## Add "查看进度" Button to History Page
- **Decision**: Added a "查看进度" button to the History page for tasks in `processing` status.
- **Rationale**: Users need a way to re-enter the progress view for ongoing tasks if they navigated away.
- **Implementation**:
  - Updated `HistoryItem` interface to include optional `taskId`.
  - Added button conditional rendering: `status === 'processing' && taskId`.
  - Navigation uses `taskId` (UUID) instead of `id` (DB ID) to match the routing configuration `/progress/:taskId`.
