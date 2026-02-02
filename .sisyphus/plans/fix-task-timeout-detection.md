# 修复任务超时检测逻辑

## TL;DR

> **Quick Summary**: 修复 `recovery.go` 中的超时检测逻辑，处理 `LastHeartbeat` 零值导致任务被错误标记为超时的问题，并添加调试日志帮助追踪抓取进度。
> 
> **Deliverables**:
> - 修复 `backend/task/recovery.go` 的超时检测逻辑
> - 在 `backend/bilibili/scraper.go` 添加调试日志
> - 验证进度显示正常工作
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: NO - sequential
> **Critical Path**: Task 1 → Task 2 → Task 3

---

## Context

### Original Request
用户反馈：抓取视频评论阶段仍然卡住，看不到具体进度（如 "正在抓取 (1/50): 视频标题..."）

### Interview Summary
**Key Discussions**:
- 后端日志显示任务被标记为超时失败
- `LastHeartbeat` 零值导致 `time.Since()` 返回非常大的值
- 进度回调代码看起来正确，但需要添加日志确认

**Research Findings**:
- `recovery.go` 第 31 行：`time.Since(task.LastHeartbeat) > time.Hour` 没有处理零值情况
- `executor.go` 第 268 行已设置 `LastHeartbeat: time.Now()`，但旧数据可能没有这个值
- `scraper.go` 的进度回调代码正确，但缺少调试日志

---

## Work Objectives

### Core Objective
修复任务超时检测逻辑，确保抓取进度能正确显示

### Concrete Deliverables
- `backend/task/recovery.go` - 添加零值检查和详细日志
- `backend/bilibili/scraper.go` - 添加调试日志确认回调被调用

### Definition of Done
- [ ] 新任务不会被错误标记为超时
- [ ] 抓取阶段能看到 "正在抓取 (X/N): 视频标题..." 进度
- [ ] 后端日志显示进度回调被正确调用

### Must Have
- 零值检查逻辑
- 详细的调试日志
- 不破坏现有功能

### Must NOT Have (Guardrails)
- 不要删除超时检测功能（只是修复零值问题）
- 不要修改 SSE 推送逻辑
- 不要修改前端代码

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: NO (Go 项目没有测试框架配置)
- **User wants tests**: Manual-only
- **Framework**: none
- **QA approach**: Manual verification via backend logs and frontend UI

---

## Execution Strategy

### Sequential Execution
```
Task 1: 修复 recovery.go 超时检测
    ↓
Task 2: 添加 scraper.go 调试日志
    ↓
Task 3: 重启后端并验证
```

---

## TODOs

- [ ] 1. 修复 recovery.go 超时检测逻辑

  **What to do**:
  - 在超时检查前添加 `LastHeartbeat.IsZero()` 检查
  - 如果是零值，更新为当前时间并尝试恢复（而不是标记为失败）
  - 添加详细的日志输出，包括 `LastHeartbeat` 的实际值和时间差

  **Must NOT do**:
  - 不要删除超时检测功能
  - 不要修改超时时间（保持 1 小时）

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 单文件修改，逻辑简单
  - **Skills**: []
    - 无需特殊技能

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 2, Task 3
  - **Blocked By**: None

  **References**:
  - `backend/task/recovery.go:29-39` - 当前超时检测逻辑，需要修改
  - `backend/task/executor.go:268` - `LastHeartbeat: time.Now()` 设置位置

  **Acceptance Criteria**:
  
  **修改后的代码应该是**:
  ```go
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
  ```

  **Commit**: YES
  - Message: `fix(recovery): handle zero LastHeartbeat value in timeout detection`
  - Files: `backend/task/recovery.go`

---

- [ ] 2. 添加 scraper.go 调试日志

  **What to do**:
  - 在 `ScrapeByVideos` 函数的进度回调处添加 `log.Printf` 调试日志
  - 确认回调函数被正确调用
  - 日志应包含：当前视频索引、总数、视频标题

  **Must NOT do**:
  - 不要修改进度回调的逻辑
  - 不要修改 SSE 推送代码

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 单文件修改，只添加日志
  - **Skills**: []
    - 无需特殊技能

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

  **References**:
  - `backend/bilibili/scraper.go:339-340` - 开始抓取时的进度推送
  - `backend/bilibili/scraper.go:353-354` - 失败时的进度推送
  - `backend/bilibili/scraper.go:369-370` - 完成时的进度推送

  **Acceptance Criteria**:
  
  **在第 339-340 行之前添加日志**:
  ```go
  log.Printf("[Scraper] Starting video %d/%d: %s (BVID: %s)", currentStarted, len(videos), title, v.BVID)
  s.reportProgress("scraping", currentStarted, len(videos),
      fmt.Sprintf("正在抓取 (%d/%d): %s", currentStarted, len(videos), title))
  ```

  **在第 369-370 行之前添加日志**:
  ```go
  log.Printf("[Scraper] Completed video %d/%d: %s, got %d comments", completedCount, len(videos), v.BVID, commentCount)
  s.reportProgress("scraping", completedCount, len(videos),
      fmt.Sprintf("已完成 %d/%d，共%d条评论", completedCount, len(videos), result.Stats.TotalComments))
  ```

  **Commit**: YES
  - Message: `feat(scraper): add debug logging for progress tracking`
  - Files: `backend/bilibili/scraper.go`

---

- [ ] 3. 重启后端并验证

  **What to do**:
  - 重启后端服务
  - 创建新的分析任务
  - 观察后端日志，确认：
    1. 没有 "timed out" 错误日志
    2. 看到 "[Scraper] Starting video..." 日志
    3. 看到 "[Scraper] Completed video..." 日志
  - 观察前端进度页面，确认：
    1. 能看到 "正在抓取 (X/N): 视频标题..." 进度
    2. 进度条正常更新

  **Must NOT do**:
  - 不要修改任何代码
  - 不要删除数据库

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 简单的验证任务
  - **Skills**: [`playwright`]
    - `playwright`: 用于验证前端进度显示

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: None
  - **Blocked By**: Task 1, Task 2

  **References**:
  - `backend/main.go` - 后端入口文件

  **Acceptance Criteria**:
  
  **后端验证** (using Bash):
  ```bash
  # 重启后端
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  # 如果有正在运行的后端进程，先停止
  pkill -f "go run backend/main.go" || true
  # 启动后端（后台运行）
  go run backend/main.go &
  # 等待启动
  sleep 3
  # 检查是否启动成功
  curl -s http://localhost:8080/api/config | head -c 100
  ```

  **前端验证** (using playwright skill):
  ```
  1. Navigate to: http://localhost:5173
  2. 输入商品类目（如 "无线吸尘器"）
  3. 点击解析按钮
  4. 确认解析结果后开始分析
  5. 观察进度页面，确认能看到 "正在抓取 (X/N): 视频标题..."
  6. Screenshot: .sisyphus/evidence/task-3-progress-display.png
  ```

  **后端日志验证**:
  - 应该看到: `[Scraper] Starting video 1/N: ...`
  - 应该看到: `[Scraper] Completed video 1/N: ...`
  - 不应该看到: `[Recovery] Task xxx timed out`

  **Commit**: NO (验证任务，无代码修改)

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `fix(recovery): handle zero LastHeartbeat value in timeout detection` | backend/task/recovery.go | 编译通过 |
| 2 | `feat(scraper): add debug logging for progress tracking` | backend/bilibili/scraper.go | 编译通过 |

---

## Success Criteria

### Verification Commands
```bash
# 编译检查
cd /Users/daozhang/Trae_AI/B站商品评论解析
go build ./backend/...  # Expected: 无错误

# 启动后端
go run backend/main.go  # Expected: 服务启动成功
```

### Final Checklist
- [ ] `recovery.go` 添加了零值检查
- [ ] `scraper.go` 添加了调试日志
- [ ] 后端编译通过
- [ ] 新任务不会被错误标记为超时
- [ ] 前端能看到抓取进度
