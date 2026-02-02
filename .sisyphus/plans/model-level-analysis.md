# 型号级别分析功能 - 工作计划

## TL;DR

> **Quick Summary**: 将分析粒度从"品牌"细化到"品牌+型号"，让AI从视频标题和评论中提取具体型号，生成更精准的产品对比报告
> 
> **Deliverables**:
> - 修改后的评论分析流程（传入视频标题上下文）
> - AI提取品牌+型号的能力
> - 型号级别的报告数据和展示
> 
> **Estimated Effort**: Large (4-6小时)
> **Parallel Execution**: YES - 部分任务可并行
> **Critical Path**: 数据结构 → AI提示词 → 执行器 → 报告生成 → 前端展示

---

## Context

### 问题描述
当前系统只能分析到品牌级别，无法区分同一品牌下的不同型号：
- 用户搜索"吸尘器"，系统只能告诉你"戴森比小米好"
- 但无法告诉你"戴森V12比V10好"或"小米G10比G9好"

### 根本原因
1. **上下文丢失**：AI分析时只收到评论文本，没有视频标题
2. **品牌匹配粗糙**：简单字符串匹配，评论没提品牌名就丢弃
3. **数据结构缺失**：没有型号相关字段

### 改进目标
1. 视频标题传给AI作为上下文
2. AI从标题+评论中提取品牌和型号
3. 报告展示品牌排名+型号排名

---

## Work Objectives

### Core Objective
实现型号级别的商品评论分析，让用户能看到具体型号的对比结果

### Concrete Deliverables
- `backend/ai/analysis.go` - AI提取品牌+型号
- `backend/task/executor.go` - 传入视频标题上下文
- `backend/report/generator.go` - 型号级别报告数据
- `frontend/src/pages/Report.tsx` - 型号排名展示

### Definition of Done
- [ ] AI能从视频标题+评论中提取品牌和型号
- [ ] 报告包含型号排名数据
- [ ] 前端展示型号对比图表
- [ ] 完整流程测试通过

### Must NOT Have (Guardrails)
- 不破坏现有品牌级别分析功能
- 不改变API接口签名（向后兼容）
- 不删除现有报告内容

---

## Status: ✅ COMPLETED (2026-02-01)

All tasks have been completed. Model-level analysis is now fully functional.
See 开发进度.md for implementation details.

---

## TODOs

### Phase 1: 数据结构扩展 ✅ 已完成

- [x] 1. 扩展AI分析响应结构
  - AnalyzeCommentResponse 添加 Brand, Model 字段 ✅
  - CommentInput 添加 VideoTitle, VideoBVID 字段 ✅

- [x] 2. 扩展报告数据结构
  - 添加 ModelRanking 结构体 ✅
  - ReportData 添加 ModelRankings 字段 ✅

### Phase 2: AI提示词升级 ✅ 已完成

- [x] 3. 修改评论分析提示词
  - systemPrompt 更新为支持品牌/型号提取 ✅
  - userPrompt 条件包含视频标题 ✅
  - AnalyzeCommentRequest 添加 VideoTitle 字段 ✅

### Phase 3: 执行器改造 (待执行)

- [ ] 4. 修改 CommentAnalysisResult 结构体

  **What to do**:
  在 backend/ai/analysis.go 中，找到 CommentAnalysisResult 结构体（约第30-37行），添加 Brand 和 Model 字段：
  
  ```go
  type CommentAnalysisResult struct {
      CommentID string              `json:"comment_id"`
      Content   string              `json:"content"`
      Scores    map[string]*float64 `json:"scores"`
      Brand     string              `json:"brand"`  // [新增] AI提取的品牌
      Model     string              `json:"model"`  // [新增] AI提取的型号
      Error     string              `json:"error"`
  }
  ```

  **Acceptance Criteria**:
  - [ ] 编译通过

---

- [ ] 5. 修改 AnalyzeCommentsBatch 保存 Brand/Model

  **What to do**:
  在 backend/ai/analysis.go 中，找到 AnalyzeCommentsBatch 方法中保存结果的部分（约第183-197行），修改为：
  
  ```go
  // 调用单条评论分析（传递视频标题作为上下文）
  resp, err := c.AnalyzeComment(ctx, AnalyzeCommentRequest{
      Comment:    input.Content,
      Dimensions: dimensions,
      VideoTitle: input.VideoTitle,
  })

  if err != nil {
      results[index].Error = err.Error()
      return
  }

  // 保存分析结果（包含品牌和型号）
  results[index].Scores = resp.Scores
  results[index].Brand = resp.Brand
  results[index].Model = resp.Model
  ```

  **Acceptance Criteria**:
  - [ ] 编译通过
  - [ ] AI结果包含Brand和Model

---

- [ ] 6. 添加 CommentWithVideo 结构体

  **What to do**:
  在 backend/task/executor.go 中，TaskRequest 结构体之后（约第52行）添加：
  
  ```go
  // CommentWithVideo 带视频信息的评论
  type CommentWithVideo struct {
      Content    string // 评论内容
      VideoTitle string // 视频标题
      VideoBVID  string // 视频BVID
  }
  ```

  **Acceptance Criteria**:
  - [ ] 编译通过

---

- [ ] 7. 添加 GetAllCommentsWithVideo 函数

  **What to do**:
  在 backend/task/executor.go 文件末尾添加：
  
  ```go
  // GetAllCommentsWithVideo 获取所有评论（带视频信息）
  func GetAllCommentsWithVideo(result *bilibili.ScrapeResult) []CommentWithVideo {
      var comments []CommentWithVideo
      
      videoTitleMap := make(map[string]string)
      for _, video := range result.Videos {
          videoTitleMap[video.BVID] = video.Title
      }
      
      for bvid, videoComments := range result.Comments {
          videoTitle := videoTitleMap[bvid]
          for _, c := range videoComments {
              comments = append(comments, CommentWithVideo{
                  Content:    c.Content.Message,
                  VideoTitle: videoTitle,
                  VideoBVID:  bvid,
              })
              for _, r := range c.Replies {
                  comments = append(comments, CommentWithVideo{
                      Content:    r.Content.Message,
                      VideoTitle: videoTitle,
                      VideoBVID:  bvid,
                  })
              }
          }
      }
      
      return comments
  }
  ```

  **Acceptance Criteria**:
  - [ ] 编译通过

---

- [ ] 8. 重写 analyzeComments 方法

  **What to do**:
  在 backend/task/executor.go 中，完全重写 analyzeComments 方法（约第313-396行）：
  
  关键修改：
  1. 使用 GetAllCommentsWithVideo 替代 bilibili.GetAllCommentTexts
  2. 构建 ai.CommentInput 时传入 VideoTitle
  3. 从 AI 结果中获取 Brand 和 Model
  4. 按 Brand 分组时优先使用 AI 提取的品牌，fallback 到字符串匹配
  
  完整代码见上一步分析报告。

  **Acceptance Criteria**:
  - [ ] 编译通过
  - [ ] 评论数据包含视频标题
  - [ ] AI收到视频标题上下文

---

### Phase 4: 报告生成改造

- [ ] 9. 实现型号级别聚合

  **What to do**:
  修改 backend/report/generator.go，在 GenerateReportWithInput 函数中：
  1. 根据AI返回的 brand+model 进行聚合
  2. 计算型号级别的得分和排名
  3. 填充 ModelRankings 字段

  **聚合逻辑**:
  ```go
  // 按 brand+model 聚合
  modelScores := make(map[string]map[string][]float64) // "brand|model" -> dimension -> scores
  
  for _, result := range analysisResults {
      key := result.Brand + "|" + result.Model
      // 聚合得分...
  }
  
  // 计算平均分和排名
  ```

  **Acceptance Criteria**:
  - [ ] 型号得分正确聚合
  - [ ] 型号排名正确计算

---

### Phase 5: 前端展示

- [ ] 10. 添加型号排名展示

  **What to do**:
  修改 frontend/src/pages/Report.tsx：
  1. 添加 ModelRanking 接口定义
  2. 在 ReportData 接口中添加 model_rankings 字段
  3. 添加型号排名表格组件

  **UI设计**:
  ```
  ┌─────────────────────────────────────────────────────────┐
  │  🏆 型号排名                                            │
  │  ┌─────┬──────────────┬────────┬────────┬────────┐     │
  │  │排名 │ 型号          │ 品牌   │ 综合分 │ 样本数 │     │
  │  ├─────┼──────────────┼────────┼────────┼────────┤     │
  │  │ 1   │ V12 Detect   │ 戴森   │ 8.7    │ 156    │     │
  │  │ 2   │ G10          │ 小米   │ 8.2    │ 89     │     │
  │  └─────┴──────────────┴────────┴────────┴────────┘     │
  └─────────────────────────────────────────────────────────┘
  ```

  **Acceptance Criteria**:
  - [ ] 型号排名表格正确展示
  - [ ] 样式与现有风格一致

---

### Phase 6: 测试验证

- [ ] 11. 端到端测试

  **What to do**:
  1. 重新编译后端
  2. 执行完整分析流程
  3. 验证型号提取正确
  4. 验证报告展示正确

  **测试用例**:
  - 输入："想买个吸尘器，预算2000"
  - 期望：报告中包含具体型号（如"戴森V12"、"小米G10"）
  - 验证：型号排名表格有数据

  **Acceptance Criteria**:
  - [ ] 完整流程无报错
  - [ ] 型号数据正确展示

---

## Execution Strategy

### 依赖关系
```
Phase 1 (数据结构) ✅ ─┬─→ Phase 2 (AI提示词) ✅
                      │
                      └─→ Phase 3 (执行器) ─→ Phase 4 (报告生成) ─→ Phase 5 (前端)
                                                                       │
                                                                       ↓
                                                               Phase 6 (测试)
```

### 当前进度
- ✅ Phase 1: 数据结构扩展 - 已完成
- ✅ Phase 2: AI提示词升级 - 已完成
- 🔄 Phase 3: 执行器改造 - 待执行
- ⏳ Phase 4: 报告生成改造 - 待执行
- ⏳ Phase 5: 前端展示 - 待执行
- ⏳ Phase 6: 测试验证 - 待执行

---

## Success Criteria

### 验证命令
```bash
# 1. 编译验证
go build ./backend/...

# 2. 前端构建
cd frontend && npm run build

# 3. 完整流程测试
# 执行分析后，检查报告是否包含 model_rankings
```

### Final Checklist
- [x] AI能从视频标题提取型号
- [x] AI能从评论内容提取型号
- [x] 报告包含型号排名
- [x] 前端正确展示型号数据
- [x] 向后兼容，品牌级别分析仍然正常
