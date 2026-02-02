# 修复评论抓取进度显示问题

## TL;DR

> **Quick Summary**: 修复"抓取视频评论"阶段进度不更新的问题，在开始抓取每个视频时就推送进度，而不是等完成后
> 
> **Deliverables**:
> - 修改 scraper.go，在开始抓取每个视频时推送进度
> - 显示当前正在抓取的视频标题
> - 提供更细粒度的进度反馈
> 
> **Estimated Effort**: Quick (10分钟)
> **Parallel Execution**: NO - sequential
> **Critical Path**: Task 1 → Task 2 → Task 3

---

## Context

### Original Request
用户反馈：抓取视频评论阶段只有转圈，看不到进度更新

### 问题分析
1. **当前行为**：进度回调只在每个视频的评论抓取**完成后**才触发
2. **问题**：单个视频抓取可能需要30秒-2分钟，期间前端看不到任何进度
3. **用户体验**：前端一直显示"抓取视频评论"转圈，不知道是否在工作

### 根本原因
`scraper.go:350-351` 的进度回调在 goroutine 内部，只有完成后才触发：
```go
s.reportProgress("scraping", completedCount, len(videos),
    fmt.Sprintf("已完成%d/%d", completedCount, len(videos)))
```

---

## Work Objectives

### Core Objective
让用户在抓取评论阶段能看到实时进度，知道当前正在抓取哪个视频

### Concrete Deliverables
- 在开始抓取每个视频时推送进度
- 显示当前正在抓取的视频标题（截断显示）
- 完成后也推送进度

### Definition of Done
- [ ] 开始抓取每个视频时前端能看到进度更新
- [ ] 进度消息显示当前视频标题
- [ ] 完成后显示统计信息

---

## TODOs

- [ ] 1. 修改 ScrapeByVideos 函数，在开始抓取时推送进度

  **What to do**:
  - 在 goroutine 开始时（抓取评论之前）推送进度
  - 显示当前正在抓取的视频标题
  - 使用 startedCount 计数器跟踪已开始的视频数

  **File**: `backend/bilibili/scraper.go`
  
  **Current Code** (Lines 307-352):
  ```go
  var wg sync.WaitGroup
  var mu sync.Mutex
  var completedCount int

  for _, video := range videos {
      // ...
      wg.Add(1)
      go func(v VideoInfo) {
          defer wg.Done()
          defer sem.Release(1)

          comments, err := s.scrapeVideoComments(ctx, v.BVID)

          mu.Lock()
          defer mu.Unlock()

          completedCount++
          // ... 只在完成后推送进度
          s.reportProgress("scraping", completedCount, len(videos),
              fmt.Sprintf("已完成%d/%d", completedCount, len(videos)))
      }(video)
  }
  ```
  
  **New Code**:
  ```go
  var wg sync.WaitGroup
  var mu sync.Mutex
  var completedCount int
  var startedCount int

  for _, video := range videos {
      // ...
      wg.Add(1)
      go func(v VideoInfo) {
          defer wg.Done()
          defer sem.Release(1)

          // 开始抓取时推送进度
          mu.Lock()
          startedCount++
          currentStarted := startedCount
          mu.Unlock()
          
          // 截断视频标题（最多20个字符）
          title := v.Title
          if len([]rune(title)) > 20 {
              title = string([]rune(title)[:20]) + "..."
          }
          s.reportProgress("scraping", currentStarted, len(videos),
              fmt.Sprintf("正在抓取 (%d/%d): %s", currentStarted, len(videos), title))

          comments, err := s.scrapeVideoComments(ctx, v.BVID)

          mu.Lock()
          defer mu.Unlock()

          completedCount++

          if err != nil {
              result.Stats.Errors = append(result.Stats.Errors,
                  fmt.Sprintf("视频%s抓取失败: %v", v.BVID, err))
              s.reportProgress("scraping", completedCount, len(videos),
                  fmt.Sprintf("已完成 %d/%d (失败: %s)", completedCount, len(videos), v.BVID))
              return
          }

          result.Comments[v.BVID] = comments

          commentCount := len(comments)
          replyCount := 0
          for _, c := range comments {
              replyCount += len(c.Replies)
          }
          result.Stats.TotalComments += commentCount
          result.Stats.TotalReplies += replyCount

          s.reportProgress("scraping", completedCount, len(videos),
              fmt.Sprintf("已完成 %d/%d，共%d条评论", completedCount, len(videos), result.Stats.TotalComments))
      }(video)
  }
  ```

  **Verification**: `cd /Users/daozhang/Trae_AI/B站商品评论解析 && go build ./backend/...`

  **Acceptance Criteria**:
  - [ ] 编译通过
  - [ ] 开始抓取每个视频时推送进度
  - [ ] 进度消息包含视频标题
  - [ ] 完成后显示评论总数

---

- [ ] 2. 重启后端服务

  **What to do**:
  ```bash
  pkill -f "go run backend/main.go"
  cd /Users/daozhang/Trae_AI/B站商品评论解析 && nohup go run backend/main.go > /tmp/backend.log 2>&1 &
  sleep 3 && curl -s http://localhost:8080/api/config | head -1
  ```

  **Verification**: 后端正常响应

---

- [ ] 3. 手动测试验证

  **What to do**:
  1. 打开 http://localhost:5173
  2. 输入一个商品需求（如"无线耳机"）
  3. 确认后观察进度页面
  4. 验证"抓取视频评论"阶段能看到实时进度更新
  5. 验证进度消息显示当前视频标题

  **Acceptance Criteria**:
  - [ ] 进度页面显示"正在抓取 (1/N): 视频标题..."
  - [ ] 进度随着抓取进行而更新
  - [ ] 完成后显示评论总数

---

## Success Criteria

### Verification Commands
```bash
# 后端编译检查
cd /Users/daozhang/Trae_AI/B站商品评论解析 && go build ./backend/...

# 查看后端日志
tail -f /tmp/backend.log | grep -E "正在抓取|已完成"
```

### Final Checklist
- [ ] 开始抓取每个视频时前端能看到进度更新
- [ ] 进度消息显示当前视频标题（截断到20字符）
- [ ] 完成后显示统计信息（评论总数）
- [ ] 抓取失败时也显示进度更新

---

## Technical Notes

### 为什么需要两个计数器？
- `startedCount`: 跟踪已开始抓取的视频数（用于"正在抓取"消息）
- `completedCount`: 跟踪已完成抓取的视频数（用于"已完成"消息）

### 为什么要截断标题？
- B站视频标题可能很长（50-100字符）
- 截断到20字符保持进度消息简洁
- 使用 `[]rune` 正确处理中文字符

### 并发安全
- 使用 `mu.Lock()` 保护计数器的读写
- 在推送进度前先获取当前值，避免长时间持有锁
