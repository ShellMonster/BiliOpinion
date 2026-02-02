# 评论条数限制与智能排序

## TL;DR

> **Quick Summary**: 添加前端评论条数限制设置，实现智能评论筛选算法（综合热度、长度、关键词相关性），替换后端硬编码的500条限制。
> 
> **Deliverables**:
> - 前端 Confirm 页面新增评论条数下拉框
> - 后端 API 扩展支持 `max_comments` 参数
> - 新建 `backend/comment/filter.go` 评论筛选模块
> - 单元测试覆盖筛选逻辑
> 
> **Estimated Effort**: Medium (4-6小时)
> **Parallel Execution**: YES - 2 waves
> **Critical Path**: Task 2 (后端API) → Task 3 (筛选模块) → Task 4 (集成)

---

## Context

### Original Request
用户希望：
1. 在前端添加评论条数限制设置（与视频时长筛选并列）
2. 实现智能评论排序/筛选算法，综合考虑热度(A)、内容长度(B)、关键词相关性(C)
3. 移除后端硬编码的500条限制

### Interview Summary
**Key Discussions**:
- 发现关键问题：`executor.go:392` 硬编码截断所有评论到500条，导致90%数据丢失
- 前端UI位置：`Confirm.tsx:162-176`（视频时长筛选器旁边）
- B站评论可用字段：`Like`(点赞)、`Count`(回复)、`Content.Message`(内容)

**User Decisions**:
- 默认评论数：500条
- 最短长度阈值：10个字
- 剔除纯表情评论：是
- 评分权重：热度40 + 长度30 + 关键词30

### Metis Review
**Identified Gaps** (addressed):
- 关键词来源：使用 `brands` + `keywords` 作为相关性评分的关键词
- 回复处理：回复独立评分，不继承父评论分数
- 平局处理：按点赞数降序作为次要排序
- Unicode长度：使用 rune count（字符数）而非字节数
- 大小写匹配：关键词匹配不区分大小写

---

## Work Objectives

### Core Objective
实现可配置的评论条数限制，并通过智能排序算法优先保留高质量评论。

### Concrete Deliverables
- `frontend/src/pages/Confirm.tsx` - 新增评论条数下拉框
- `backend/api/confirm.go` - 扩展 `ConfirmRequest` 结构体
- `backend/comment/filter.go` - 新建评论筛选模块
- `backend/comment/filter_test.go` - 单元测试
- `backend/task/executor.go` - 集成筛选逻辑

### Definition of Done
- [x] 前端下拉框显示并可选择评论条数限制
- [x] 后端正确接收并使用 `max_comments` 参数
- [x] 筛选模块单元测试全部通过
- [x] 分析任务使用智能排序而非简单截断

### Must Have
- 评论条数限制下拉框（100/200/500/1000/2000）
- 智能评分算法（热度40 + 长度30 + 关键词30）
- 无效评论过滤（<10字、纯表情）
- 单元测试覆盖核心逻辑

### Must NOT Have (Guardrails)
- ❌ 不实现 UI 可配置权重（权重固定为40/30/30）
- ❌ 不实现 每品牌单独限制（使用全局限制）
- ❌ 不实现 筛选预览功能
- ❌ 不实现 去重功能（后续优化）
- ❌ 不修改 抓取逻辑（`scraper.go`）
- ❌ 不修改 AI分析逻辑（`ai/analysis.go`）
- ❌ 不添加 数据库迁移

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: YES (Go testing)
- **User wants tests**: YES (单元测试)
- **Framework**: `go test`

### Test Coverage Requirements
筛选模块必须包含以下测试用例：
1. 空输入返回空输出
2. 纯表情评论被过滤
3. 短评论（<10字）被过滤
4. 评分排序正确（高分在前）
5. 限制条数生效
6. 关键词匹配不区分大小写

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: 前端UI改动 [no dependencies]
└── Task 2: 后端API扩展 [no dependencies]

Wave 2 (After Wave 1):
├── Task 3: 筛选模块实现 [depends: 2]
└── Task 4: 集成到执行器 [depends: 2, 3]

Critical Path: Task 2 → Task 3 → Task 4
Parallel Speedup: ~30% faster than sequential
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | None | 2 |
| 2 | None | 3, 4 | 1 |
| 3 | 2 | 4 | None |
| 4 | 2, 3 | None | None (final) |

---

## TODOs

- [x] 1. 前端：添加评论条数限制下拉框

  **What to do**:
  - 在 `Confirm.tsx` 中添加 `maxComments` state，默认值 500
  - 在视频时长筛选器旁边添加评论条数下拉框
  - 更新 `handleConfirm` 函数，传递 `max_comments` 参数
  - 下拉选项：100 / 200 / 500(默认) / 1000 / 2000

  **Must NOT do**:
  - 不添加新的 UI 区块（只在现有筛选区域添加）
  - 不添加客户端筛选逻辑

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 前端UI改动，需要保持样式一致性
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: 确保UI风格与现有下拉框一致

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 2)
  - **Blocks**: None
  - **Blocked By**: None (can start immediately)

  **References**:
  - `frontend/src/pages/Confirm.tsx:25-27` - 现有 state 定义模式（`videoDateRangeMonths`, `minVideoDuration`）
  - `frontend/src/pages/Confirm.tsx:162-176` - 视频时长筛选器UI模式（复制此结构）
  - `frontend/src/pages/Confirm.tsx:63-70` - `handleConfirm` 中的参数传递模式

  **Acceptance Criteria**:

  ```bash
  # AC1: 前端编译通过
  cd frontend && npm run build
  # Expected: 编译成功，无错误

  # AC2: 下拉框渲染正确（Playwright验证）
  # 1. Navigate to: http://localhost:5173/confirm?requirement=test
  # 2. Wait for: 页面加载完成
  # 3. Assert: 存在评论条数下拉框
  # 4. Assert: 默认值为 500
  # 5. Screenshot: .sisyphus/evidence/task-1-dropdown.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): add comment limit dropdown in Confirm page`
  - Files: `frontend/src/pages/Confirm.tsx`

---

- [x] 2. 后端：扩展 API 支持 max_comments 参数

  **What to do**:
  - 在 `ConfirmRequest` 结构体中添加 `MaxComments int` 字段
  - 在 `TaskConfig` 结构体中添加 `MaxComments int` 字段
  - 在 `HandleConfirm` 中传递 `MaxComments` 到任务配置
  - 默认值处理：如果未提供或为0，默认500

  **Must NOT do**:
  - 不修改现有字段的行为
  - 不添加新的 API 端点

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 简单的结构体字段添加，改动量小
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocks**: Task 3, Task 4
  - **Blocked By**: None (can start immediately)

  **References**:
  - `backend/api/confirm.go:14-24` - `ConfirmRequest` 结构体定义
  - `backend/api/confirm.go:26-70` - `HandleConfirm` 函数，参数传递逻辑
  - `backend/task/executor.go:24-45` - `TaskConfig` 结构体定义

  **Acceptance Criteria**:

  ```bash
  # AC1: 后端编译通过
  cd backend && go build ./...
  # Expected: 编译成功，无错误

  # AC2: API 接受 max_comments 参数
  curl -s -X POST http://localhost:8080/api/confirm \
    -H "Content-Type: application/json" \
    -d '{"requirement":"测试","brands":["A"],"dimensions":[{"name":"X","description":"Y"}],"keywords":["test"],"max_comments":100}' \
    | jq '.task_id'
  # Expected: 返回 task_id（非空字符串）

  # AC3: 默认值生效（不传 max_comments 时）
  curl -s -X POST http://localhost:8080/api/confirm \
    -H "Content-Type: application/json" \
    -d '{"requirement":"测试","brands":["A"],"dimensions":[{"name":"X","description":"Y"}],"keywords":["test"]}' \
    | jq '.task_id'
  # Expected: 返回 task_id，后端日志显示使用默认值500
  ```

  **Commit**: YES
  - Message: `feat(api): add max_comments parameter to confirm endpoint`
  - Files: `backend/api/confirm.go`, `backend/task/executor.go`

---

- [x] 3. 后端：实现评论筛选模块

  **What to do**:
  - 创建 `backend/comment/filter.go`
  - 实现 `FilterConfig` 结构体（MinLength, FilterEmoji, Keywords, MaxComments）
  - 实现 `FilterAndRank(comments, config)` 函数
  - 实现 `scoreComment(comment, keywords)` 评分函数
  - 实现 `isValidComment(comment, minLength, filterEmoji)` 验证函数
  - 实现 `removeEmojiAndSymbols(text)` 辅助函数
  - 创建 `backend/comment/filter_test.go` 单元测试

  **评分算法**:
  ```
  总分 = 热度分(0-40) + 长度分(0-30) + 关键词分(0-30)
  
  热度分：
  - 点赞分 = min(likes/100, 20)
  - 回复分 = min(replies/10, 20)
  
  长度分：
  - 长度分 = min(字符数/10, 30)
  
  关键词分：
  - 每命中一个关键词 +10分，上限30分
  - 不区分大小写
  ```

  **过滤规则**:
  - 字符数 < 10 → 过滤
  - 去除表情/符号后字符数 < 10 → 过滤

  **Must NOT do**:
  - 不访问数据库（纯函数）
  - 不修改输入切片（返回新切片）
  - 不实现去重逻辑

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
    - Reason: 核心业务逻辑实现，需要算法设计
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (sequential)
  - **Blocks**: Task 4
  - **Blocked By**: Task 2 (需要知道 TaskConfig 结构)

  **References**:
  - `backend/bilibili/comment.go:32-50` - `Comment` 结构体定义（Like, Count, Content.Message）
  - `backend/task/executor.go:392-395` - 当前的简单截断逻辑（需要替换）

  **Acceptance Criteria**:

  ```bash
  # AC1: 单元测试全部通过
  cd backend && go test ./comment/... -v
  # Expected: 所有测试通过

  # AC2: 测试覆盖率 >= 80%
  cd backend && go test ./comment/... -cover
  # Expected: coverage >= 80%
  ```

  **必须包含的测试用例**:
  ```go
  // filter_test.go
  
  func TestFilterAndRank_EmptyInput(t *testing.T)
  // Input: [], Output: []
  
  func TestFilterAndRank_PureEmojiFiltered(t *testing.T)
  // Input: [{Content: "😀😀😀"}], Output: []
  
  func TestFilterAndRank_ShortCommentFiltered(t *testing.T)
  // Input: [{Content: "好"}], Output: []
  
  func TestFilterAndRank_ValidCommentKept(t *testing.T)
  // Input: [{Content: "这个产品非常好用，推荐购买"}], Output: [same]
  
  func TestFilterAndRank_SortByScore(t *testing.T)
  // Input: [low_score, high_score], Output: [high_score, low_score]
  
  func TestFilterAndRank_LimitApplied(t *testing.T)
  // Input: 10 comments, MaxComments: 5, Output: 5 comments
  
  func TestScoreComment_KeywordCaseInsensitive(t *testing.T)
  // Input: "IPHONE很好", keywords: ["iphone"], Expected: keyword score > 0
  
  func TestIsValidComment_WhitespaceOnly(t *testing.T)
  // Input: "   ", Expected: false
  ```

  **Commit**: YES
  - Message: `feat(comment): add smart filtering and ranking module`
  - Files: `backend/comment/filter.go`, `backend/comment/filter_test.go`

---

- [x] 4. 后端：集成筛选模块到执行器

  **What to do**:
  - 在 `executor.go` 中导入 `comment` 包
  - 替换 `analyzeComments` 函数中的硬编码截断逻辑（第392-395行）
  - 使用 `comment.FilterAndRank()` 进行智能筛选
  - 传递关键词（brands + keywords）给筛选函数
  - 添加日志：筛选前后的评论数量

  **Must NOT do**:
  - 不修改抓取逻辑
  - 不修改 AI 分析逻辑
  - 不修改数据库操作

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 简单的集成工作，替换几行代码
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (final)
  - **Blocks**: None
  - **Blocked By**: Task 2, Task 3

  **References**:
  - `backend/task/executor.go:392-395` - 需要替换的硬编码截断逻辑
  - `backend/task/executor.go:380-416` - `analyzeComments` 函数完整上下文
  - `backend/comment/filter.go` - 新建的筛选模块（Task 3 产出）

  **Acceptance Criteria**:

  ```bash
  # AC1: 后端编译通过
  cd backend && go build ./...
  # Expected: 编译成功

  # AC2: 集成测试 - 运行完整分析任务
  # 1. 启动后端: go run backend/main.go
  # 2. 发起分析请求（max_comments=100）
  # 3. 检查日志输出
  # Expected: 日志显示 "Filtered X→Y comments" 且 Y <= 100

  # AC3: 验证筛选生效
  # 检查后端日志，确认：
  # - 筛选前评论数 > 筛选后评论数（如果原始数据 > max_comments）
  # - 筛选后评论数 <= max_comments
  ```

  **Commit**: YES
  - Message: `feat(executor): integrate smart comment filtering`
  - Files: `backend/task/executor.go`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `feat(frontend): add comment limit dropdown` | Confirm.tsx | `npm run build` |
| 2 | `feat(api): add max_comments parameter` | confirm.go, executor.go | `go build ./...` |
| 3 | `feat(comment): add filtering module` | filter.go, filter_test.go | `go test ./comment/...` |
| 4 | `feat(executor): integrate filtering` | executor.go | `go build ./...` |

---

## Success Criteria

### Verification Commands
```bash
# 前端编译
cd frontend && npm run build
# Expected: 成功

# 后端编译
cd backend && go build ./...
# Expected: 成功

# 单元测试
cd backend && go test ./comment/... -v
# Expected: 全部通过

# 集成测试（手动）
# 1. 启动服务
# 2. 访问 /confirm 页面
# 3. 选择评论条数限制
# 4. 运行分析
# 5. 检查日志确认筛选生效
```

### Final Checklist
- [x] 前端下拉框显示正确，默认值500
- [x] 后端 API 接受 max_comments 参数
- [x] 筛选模块单元测试全部通过
- [x] 执行器使用智能筛选替代简单截断
- [x] 日志显示筛选统计信息
