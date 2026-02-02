# AI 评论分析批次优化

## TL;DR

> **Quick Summary**: 将 AI 评论分析从"每条评论一个请求"优化为"多条评论合并到一个请求"，通过动态批次策略控制字符数，大幅减少 API 调用次数和重复的 system_prompt 发送。
> 
> **Deliverables**:
> - 新增 `AnalyzeCommentsBatchMerged` 函数，支持真正的批量分析
> - 动态批次计算函数，按字符数而非固定条数分批
> - 更新 `executor.go` 使用新的批量分析方法
> - 保留原有单条分析作为降级方案
> 
> **Estimated Effort**: Medium
> **Parallel Execution**: NO - sequential (有依赖关系)
> **Critical Path**: Task 1 → Task 2 → Task 3 → Task 4

---

## Context

### Original Request
用户希望优化 AI 评论分析的效率，通过批次提交多条评论到一个请求中，避免每次都重复发送 system_prompt。

### Interview Summary
**Key Discussions**:
- 当前每条评论一个 API 请求，500 条评论 = 500 次请求
- 每次请求都发送约 1000 字符的 system_prompt
- 需要考虑单条评论字数波动，按字符数动态分批

**Research Findings**:
- B 站评论通常 20-200 字，平均约 50 字
- 大多数 AI 模型支持 4K-128K tokens
- 中文约 1.5 token/字
- 保守估计每批可处理 10-15 条评论（总字符数 3000 以内）

### 技术设计

#### 批次大小计算
```
可用 tokens = 4000 (保守估计)
System Prompt = 1000 tokens
预留输出 = 1500 tokens (每条约 100 tokens × 15 条)
可用输入 = 1500 tokens ≈ 1000 中文字符

考虑视频标题等额外信息：
- 默认最大字符数：3000 字符
- 默认最大条数：15 条（硬上限）
- 实际按字符数动态调整
```

#### 新 Prompt 设计
```
System: 你是商品评论分析助手。分析以下多条评论，为每条评论提取品牌、型号并打分。

User:
评论列表：
[1] 视频：xxx | 内容：这个耳机音质很好
[2] 视频：yyy | 内容：苹果降噪太强了
...

返回 JSON：
{"results":[{"id":"1","brand":"xxx","model":"xxx","scores":{"维度1":8.5}}, ...]}
```

---

## Work Objectives

### Core Objective
将 AI 评论分析从 N 次请求优化为 N/15 次请求，提升 10-15 倍效率

### Concrete Deliverables
- `backend/ai/analysis.go` - 新增批量合并分析函数
- `backend/ai/batch.go` - 新文件，批次计算逻辑
- `backend/task/executor.go` - 使用新的批量分析方法

### Definition of Done
- [x] 500 条评论的 API 调用次数从 500 次降到 ~35 次
- [x] 分析耗时从 5-10 分钟降到 1-2 分钟
- [x] 单条评论过长时自动降级到单独处理
- [x] 保持分析结果准确性

### Must Have
- 动态批次计算（按字符数）
- JSON 响应解析和错误处理
- 降级机制（批量失败时回退到单条）
- 进度回调支持

### Must NOT Have (Guardrails)
- 不要删除原有的单条分析函数（保留作为降级方案）
- 不要硬编码批次大小（要可配置）
- 不要忽略解析错误（要有重试机制）

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: YES (Go test)
- **User wants tests**: YES (Tests-after)
- **Framework**: Go testing
- **QA approach**: 单元测试 + 集成测试

---

## Execution Strategy

### Sequential Execution
```
Task 1: 创建批次计算模块 (batch.go)
    ↓
Task 2: 实现批量合并分析函数 (analysis.go)
    ↓
Task 3: 更新 executor.go 使用新方法
    ↓
Task 4: 添加测试和验证
```

---

## TODOs

- [x] 1. 创建批次计算模块

  **What to do**:
  - 创建 `backend/ai/batch.go` 文件
  - 实现 `CalculateBatches` 函数，按字符数动态分批
  - 实现 `BatchConfig` 配置结构体

  **Must NOT do**:
  - 不要硬编码字符数限制
  - 不要忽略空评论

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 2, Task 3
  - **Blocked By**: None

  **References**:
  - `backend/ai/analysis.go:229-236` - CommentInput 结构体定义

  **Acceptance Criteria**:
  
  **新文件 `backend/ai/batch.go` 内容**:
  ```go
  package ai

  // BatchConfig 批次配置
  type BatchConfig struct {
      MaxCharsPerBatch int // 每批最大字符数（默认 3000）
      MaxItemsPerBatch int // 每批最大条数（默认 15）
      MinItemsPerBatch int // 每批最小条数（默认 1）
  }

  // DefaultBatchConfig 默认批次配置
  func DefaultBatchConfig() BatchConfig {
      return BatchConfig{
          MaxCharsPerBatch: 3000,
          MaxItemsPerBatch: 15,
          MinItemsPerBatch: 1,
      }
  }

  // CalculateBatches 按字符数动态计算批次
  // 返回分好批的评论列表
  func CalculateBatches(comments []CommentInput, config *BatchConfig) [][]CommentInput {
      if config == nil {
          cfg := DefaultBatchConfig()
          config = &cfg
      }

      var batches [][]CommentInput
      var currentBatch []CommentInput
      currentChars := 0

      for _, c := range comments {
          // 计算当前评论的字符数（内容 + 视频标题）
          commentLen := len([]rune(c.Content)) + len([]rune(c.VideoTitle))
          
          // 如果当前批次加上这条评论会超限，且当前批次不为空，则开始新批次
          shouldStartNewBatch := (currentChars+commentLen > config.MaxCharsPerBatch || 
                                  len(currentBatch) >= config.MaxItemsPerBatch) && 
                                 len(currentBatch) >= config.MinItemsPerBatch
          
          if shouldStartNewBatch {
              batches = append(batches, currentBatch)
              currentBatch = nil
              currentChars = 0
          }
          
          currentBatch = append(currentBatch, c)
          currentChars += commentLen
      }

      // 添加最后一批
      if len(currentBatch) > 0 {
          batches = append(batches, currentBatch)
      }

      return batches
  }
  ```

  **验证**:
  ```bash
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  go build ./backend/...
  ```

  **Commit**: YES
  - Message: `feat(ai): add batch calculation module for comment analysis`
  - Files: `backend/ai/batch.go`

---

- [x] 2. 实现批量合并分析函数

  **What to do**:
  - 在 `backend/ai/analysis.go` 中新增 `AnalyzeCommentsBatchMerged` 函数
  - 设计新的 prompt，支持多条评论一次性分析
  - 实现 JSON 数组响应解析
  - 添加错误处理和降级逻辑

  **Must NOT do**:
  - 不要删除原有的 `AnalyzeComment` 函数
  - 不要忽略 JSON 解析错误

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

  **References**:
  - `backend/ai/analysis.go:50-130` - 现有 AnalyzeComment 函数
  - `backend/ai/analysis.go:30-39` - CommentAnalysisResult 结构体
  - `backend/ai/keyword.go:40-65` - cleanJSONResponse 函数（可复用）

  **Acceptance Criteria**:
  
  **在 `backend/ai/analysis.go` 中添加以下函数**:
  ```go
  // BatchAnalysisResult 批量分析结果（用于 JSON 解析）
  type BatchAnalysisResult struct {
      Results []struct {
          ID     string              `json:"id"`
          Brand  string              `json:"brand"`
          Model  string              `json:"model"`
          Scores map[string]*float64 `json:"scores"`
      } `json:"results"`
  }

  // AnalyzeCommentsBatchMerged 真正的批量分析（多条评论合并到一个请求）
  // 将多条评论合并到一个 API 请求中，大幅减少请求次数
  func (c *Client) AnalyzeCommentsBatchMerged(ctx context.Context, comments []CommentInput, dimensions []Dimension) ([]CommentAnalysisResult, error) {
      if len(comments) == 0 {
          return nil, fmt.Errorf("评论列表不能为空")
      }
      if len(dimensions) == 0 {
          return nil, fmt.Errorf("评价维度不能为空")
      }

      // 构建维度列表
      var dimList []string
      for _, dim := range dimensions {
          dimList = append(dimList, fmt.Sprintf("- %s：%s", dim.Name, dim.Description))
      }

      // 构建评论列表文本
      var commentList []string
      for i, c := range comments {
          if c.VideoTitle != "" {
              commentList = append(commentList, fmt.Sprintf("[%d] 视频：%s | 内容：%s", i+1, c.VideoTitle, c.Content))
          } else {
              commentList = append(commentList, fmt.Sprintf("[%d] 内容：%s", i+1, c.Content))
          }
      }

      systemPrompt := fmt.Sprintf(`你是商品评论分析助手。分析以下多条评论，为每条评论：
  1. 提取品牌名称和具体型号
  2. 对以下维度打分（1-10分，未提及则为null）：
  %s

  评分标准：1-3差评，4-5一般，6-7较好，8-10优秀

  重要规则：
  - 每条评论独立分析，用评论编号[1][2]等标识
  - 无法确定品牌填"未知"，无法确定型号填"通用"
  - 必须返回JSON格式，不要添加任何其他文字
  - results数组的顺序必须与输入评论顺序一致

  返回格式：
  {"results":[{"id":"1","brand":"品牌","model":"型号","scores":{"维度1":8.5,"维度2":null}},{"id":"2",...}]}`, strings.Join(dimList, "\n"))

      userPrompt := fmt.Sprintf("评论列表（共%d条）：\n%s", len(comments), strings.Join(commentList, "\n"))

      messages := []Message{
          {Role: "system", Content: systemPrompt},
          {Role: "user", Content: userPrompt},
      }

      response, err := c.ChatCompletion(ctx, messages)
      if err != nil {
          return nil, fmt.Errorf("AI请求失败: %w", err)
      }

      // 清理响应（移除可能的 markdown 代码块）
      cleanedResponse := cleanJSONResponse(response)

      // 解析批量结果
      var batchResult BatchAnalysisResult
      if err := json.Unmarshal([]byte(cleanedResponse), &batchResult); err != nil {
          // 尝试提取 JSON
          jsonPattern := regexp.MustCompile(`\{[\s\S]*"results"[\s\S]*\}`)
          match := jsonPattern.FindString(response)
          if match != "" {
              if err := json.Unmarshal([]byte(match), &batchResult); err != nil {
                  return nil, fmt.Errorf("解析批量响应失败: %w, 原始响应: %s", err, response[:min(len(response), 500)])
              }
          } else {
              return nil, fmt.Errorf("无法从响应中提取JSON: %s", response[:min(len(response), 500)])
          }
      }

      // 转换为 CommentAnalysisResult 格式
      results := make([]CommentAnalysisResult, len(comments))
      for i, c := range comments {
          results[i] = CommentAnalysisResult{
              CommentID: c.ID,
              Content:   c.Content,
          }
          
          // 查找对应的分析结果
          for _, r := range batchResult.Results {
              // 支持 "1" 或 "comment_0" 格式的 ID
              if r.ID == fmt.Sprintf("%d", i+1) || r.ID == c.ID {
                  results[i].Brand = r.Brand
                  results[i].Model = r.Model
                  results[i].Scores = r.Scores
                  break
              }
          }
          
          // 如果没找到对应结果，标记错误
          if results[i].Scores == nil {
              results[i].Error = "未在批量响应中找到对应结果"
          }
      }

      return results, nil
  }

  // cleanJSONResponse 清理 AI 返回的 JSON 响应（移除 markdown 代码块）
  // 注意：如果 keyword.go 中已有此函数，可以复用
  func cleanJSONResponse(response string) string {
      response = strings.TrimSpace(response)
      
      // 移除 ```json 或 ``` 开头
      if strings.HasPrefix(response, "```json") {
          response = strings.TrimPrefix(response, "```json")
      } else if strings.HasPrefix(response, "```") {
          response = strings.TrimPrefix(response, "```")
      }
      
      // 移除 ``` 结尾
      if strings.HasSuffix(response, "```") {
          response = strings.TrimSuffix(response, "```")
      }
      
      return strings.TrimSpace(response)
  }
  ```

  **验证**:
  ```bash
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  go build ./backend/...
  ```

  **Commit**: YES
  - Message: `feat(ai): implement batch merged comment analysis`
  - Files: `backend/ai/analysis.go`

---

- [x] 3. 更新 AnalyzeCommentsWithRateLimit 使用批量合并

  **What to do**:
  - 修改 `AnalyzeCommentsWithRateLimit` 函数
  - 使用 `CalculateBatches` 动态分批
  - 使用 `AnalyzeCommentsBatchMerged` 进行批量分析
  - 添加降级逻辑：批量失败时回退到单条分析
  - 添加进度日志

  **Must NOT do**:
  - 不要删除原有的并发分析逻辑（保留作为降级方案）
  - 不要忽略部分失败的情况

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 4
  - **Blocked By**: Task 2

  **References**:
  - `backend/ai/analysis.go:249-288` - 现有 AnalyzeCommentsWithRateLimit 函数
  - `backend/ai/batch.go` - 批次计算函数

  **Acceptance Criteria**:
  
  **修改 `backend/ai/analysis.go` 中的 `AnalyzeCommentsWithRateLimit` 函数**:
  ```go
  // AnalyzeCommentsWithRateLimit 带速率限制的批量分析（优化版）
  // 使用批量合并策略，大幅减少 API 调用次数
  func (c *Client) AnalyzeCommentsWithRateLimit(ctx context.Context, comments []CommentInput, dimensions []Dimension, batchSize int) ([]CommentAnalysisResult, error) {
      // 使用动态批次计算
      config := DefaultBatchConfig()
      if batchSize > 0 && batchSize < config.MaxItemsPerBatch {
          config.MaxItemsPerBatch = batchSize
      }
      
      batches := CalculateBatches(comments, &config)
      log.Printf("[AI] 评论分析：共 %d 条评论，分为 %d 批处理", len(comments), len(batches))

      var allResults []CommentAnalysisResult
      
      for i, batch := range batches {
          log.Printf("[AI] 正在分析第 %d/%d 批（%d 条评论）...", i+1, len(batches), len(batch))
          
          // 尝试批量合并分析
          results, err := c.AnalyzeCommentsBatchMerged(ctx, batch, dimensions)
          
          if err != nil {
              log.Printf("[AI] 批量分析失败，降级到单条分析: %v", err)
              // 降级：使用原有的并发单条分析
              results, err = c.AnalyzeCommentsBatch(ctx, batch, dimensions)
              if err != nil {
                  // 如果单条分析也失败，记录错误但继续
                  for _, comment := range batch {
                      allResults = append(allResults, CommentAnalysisResult{
                          CommentID: comment.ID,
                          Content:   comment.Content,
                          Error:     err.Error(),
                      })
                  }
                  continue
              }
          }
          
          allResults = append(allResults, results...)
          
          // 统计当前批次成功率
          successCount := 0
          for _, r := range results {
              if r.Error == "" {
                  successCount++
              }
          }
          log.Printf("[AI] 第 %d/%d 批完成，成功 %d/%d 条", i+1, len(batches), successCount, len(batch))
      }

      return allResults, nil
  }
  ```

  **验证**:
  ```bash
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  go build ./backend/...
  ```

  **Commit**: YES
  - Message: `refactor(ai): use batch merged analysis in rate limit function`
  - Files: `backend/ai/analysis.go`

---

- [x] 4. 添加测试和端到端验证

  **What to do**:
  - 添加 `CalculateBatches` 单元测试
  - 添加 `AnalyzeCommentsBatchMerged` 单元测试
  - 运行完整的分析任务验证效果
  - 记录性能对比数据

  **Must NOT do**:
  - 不要跳过边界条件测试
  - 不要忽略错误处理测试

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: None
  - **Blocked By**: Task 3

  **References**:
  - `backend/ai/analysis_test.go` - 现有测试文件
  - `backend/ai/batch.go` - 批次计算函数

  **Acceptance Criteria**:
  
  **创建 `backend/ai/batch_test.go`**:
  ```go
  package ai

  import "testing"

  func TestCalculateBatches(t *testing.T) {
      tests := []struct {
          name           string
          comments       []CommentInput
          config         *BatchConfig
          expectedBatches int
      }{
          {
              name: "空评论列表",
              comments: []CommentInput{},
              config: nil,
              expectedBatches: 0,
          },
          {
              name: "单条评论",
              comments: []CommentInput{{ID: "1", Content: "测试"}},
              config: nil,
              expectedBatches: 1,
          },
          {
              name: "多条短评论合并",
              comments: []CommentInput{
                  {ID: "1", Content: "好用"},
                  {ID: "2", Content: "不错"},
                  {ID: "3", Content: "推荐"},
              },
              config: &BatchConfig{MaxCharsPerBatch: 100, MaxItemsPerBatch: 10, MinItemsPerBatch: 1},
              expectedBatches: 1,
          },
          {
              name: "超长评论单独一批",
              comments: []CommentInput{
                  {ID: "1", Content: string(make([]rune, 3000))}, // 3000 字符
                  {ID: "2", Content: "短评论"},
              },
              config: &BatchConfig{MaxCharsPerBatch: 3000, MaxItemsPerBatch: 15, MinItemsPerBatch: 1},
              expectedBatches: 2,
          },
          {
              name: "按条数限制分批",
              comments: func() []CommentInput {
                  var cs []CommentInput
                  for i := 0; i < 20; i++ {
                      cs = append(cs, CommentInput{ID: fmt.Sprintf("%d", i), Content: "短"})
                  }
                  return cs
              }(),
              config: &BatchConfig{MaxCharsPerBatch: 10000, MaxItemsPerBatch: 5, MinItemsPerBatch: 1},
              expectedBatches: 4,
          },
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              batches := CalculateBatches(tt.comments, tt.config)
              if len(batches) != tt.expectedBatches {
                  t.Errorf("期望 %d 批，实际 %d 批", tt.expectedBatches, len(batches))
              }
          })
      }
  }
  ```

  **运行测试**:
  ```bash
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  go test ./backend/ai/... -v
  ```

  **端到端验证**:
  1. 重启后端服务
  2. 创建新的分析任务
  3. 观察日志，确认批次数量减少
  4. 对比分析耗时

  **Commit**: YES
  - Message: `test(ai): add batch calculation tests`
  - Files: `backend/ai/batch_test.go`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `feat(ai): add batch calculation module` | backend/ai/batch.go | go build |
| 2 | `feat(ai): implement batch merged analysis` | backend/ai/analysis.go | go build |
| 3 | `refactor(ai): use batch merged analysis` | backend/ai/analysis.go | go build |
| 4 | `test(ai): add batch calculation tests` | backend/ai/batch_test.go | go test |

---

## Success Criteria

### Verification Commands
```bash
# 编译检查
cd /Users/daozhang/Trae_AI/B站商品评论解析
go build ./backend/...

# 运行测试
go test ./backend/ai/... -v

# 启动后端
go run backend/main.go
```

### Performance Metrics
| 指标 | 优化前 | 优化后目标 |
|------|--------|-----------|
| 500 条评论 API 调用次数 | 500 次 | ~35 次 |
| 分析耗时 | 5-10 分钟 | 1-2 分钟 |
| Token 消耗 | 高 | 降低 80%+ |

### Final Checklist
- [x] `batch.go` 创建并通过编译
- [x] `AnalyzeCommentsBatchMerged` 函数实现
- [x] `AnalyzeCommentsWithRateLimit` 使用新的批量分析
- [x] 单元测试通过
- [x] 端到端验证：日志显示批次数量减少
- [x] 降级机制正常工作
