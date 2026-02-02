# 品牌/型号发现优化计划

## TL;DR

> **Quick Summary**: 将分析系统从"确认模式"改为"发现模式"，收集所有AI识别的品牌和型号，而不仅仅是用户指定的。同时优化搜索和评论抓取策略，提高数据覆盖率。
> 
> **Deliverables**: 
> - 品牌发现模式（不再丢弃未指定品牌）
> - 通用搜索关键词生成
> - 增加视频和评论抓取数量
> - 增强型号提取能力
> 
> **Estimated Effort**: Large
> **Parallel Execution**: YES - 3 waves
> **Critical Path**: Task 1.1 → Task 1.2 → Task 1.3 → 验证

---

## Context

### Original Request
用户反馈分析时会漏掉很多品牌和型号，需要从整个流程上进行改进。

### Interview Summary
**Key Discussions**:
- 当前系统是"确认模式"：只保留用户指定的品牌，其他品牌被丢弃
- 搜索关键词太具体，限制了发现新品牌的机会
- 视频和评论数量限制过低
- 型号提取率低

**Research Findings**:
- `executor.go:360-401` 存在品牌白名单过滤，是数据丢失的主要原因
- `DefaultTaskConfig` 限制每关键词10个视频、每视频100条评论
- AI提示词对型号提取的要求不够明确

---

## Work Objectives

### Core Objective
将分析系统从"确认模式"改为"发现模式"，最大化品牌和型号的发现能力。

### Concrete Deliverables
- 修改后的 `backend/task/executor.go`（品牌发现模式）
- 修改后的 `backend/ai/keyword.go`（通用关键词生成）
- 修改后的 `backend/ai/analysis.go`（增强型号提取）
- 更新的配置参数（视频/评论数量）
- 前端显示发现的品牌

### Definition of Done
- [ ] 运行分析任务时，日志显示"发现新品牌: [...]"
- [ ] 报告中包含用户未指定但AI发现的品牌
- [ ] 型号排名表格有更多数据
- [ ] `go build ./backend/...` 通过
- [ ] `cd frontend && npm run build` 通过

### Must Have
- 品牌发现模式（不丢弃数据）
- 增加数据抓取量
- 向后兼容（不破坏现有功能）

### Must NOT Have (Guardrails)
- 不修改数据库结构
- 不改变API接口格式
- 不删除现有功能

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (go build, npm run build)
- **User wants tests**: Manual verification
- **Framework**: Go build + npm build + 运行分析任务

### Automated Verification
```bash
# 后端编译
go build ./backend/...

# 前端编译
cd frontend && npm run build

# 运行分析任务测试
curl -X POST http://localhost:8080/api/confirm \
  -H "Content-Type: application/json" \
  -d '{"requirement":"自动猫砂盆","brands":["小佩"],"dimensions":[{"name":"清洁效果","description":"清洁能力"}],"keywords":["自动猫砂盆评测"]}'
```

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1.1: 移除品牌白名单过滤
├── Task 2.1: AI解析生成通用关键词
└── Task 2.2: 增加视频数量限制

Wave 2 (After Wave 1):
├── Task 1.2: 添加品牌归一化逻辑
├── Task 3.1: 增加评论数量限制
└── Task 4.1: 增强AI提示词提取型号

Wave 3 (After Wave 2):
├── Task 1.3: 前端显示发现的品牌
├── Task 2.3: 添加视频过滤
└── Task 3.2: 评论抓取包含子评论

Final:
└── Task 5.1: 编译验证和测试
```

---

## TODOs

### Phase 1: 核心修复（P0 - 最重要）

- [ ] 1.1. 移除品牌白名单过滤，改为品牌发现模式

  **What to do**:
  - 修改 `backend/task/executor.go` 的 `analyzeComments` 函数
  - 将品牌过滤逻辑改为品牌分类逻辑
  - 收集所有AI识别的品牌，分为"指定品牌"和"发现品牌"两类
  - 添加日志记录发现的新品牌

  **Must NOT do**:
  - 不改变函数签名
  - 不修改返回类型

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 1.2, Task 1.3
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go:360-401` - 当前品牌过滤逻辑（需要替换）
  - `backend/report/generator.go:CommentWithScore` - 评论数据结构

  **Acceptance Criteria**:
  - [ ] 替换 lines 360-401 的代码为新的发现模式逻辑
  - [ ] 新代码包含 `specifiedResults` 和 `discoveredResults` 两个 map
  - [ ] 日志输出包含 "发现新品牌" 信息
  - [ ] `go build ./backend/...` 通过

  **New Code to Insert** (replace lines 360-401):
  ```go
  // === DISCOVERY MODE: 收集所有AI识别的品牌，不仅仅是用户指定的 ===

  // 创建用户指定品牌的映射（用于分类，不是过滤）
  specifiedBrands := make(map[string]string) // lowercase -> original
  for _, brand := range brands {
      specifiedBrands[strings.ToLower(brand)] = brand
  }

  // 分类收集结果：指定品牌 vs 发现的新品牌
  specifiedResults := make(map[string][]report.CommentWithScore)
  discoveredResults := make(map[string][]report.CommentWithScore)

  for _, r := range analysisResults {
      if r.Error != "" || r.Scores == nil {
          continue
      }

      // 从AI结果获取品牌
      brand := strings.TrimSpace(r.Brand)
      if brand == "" || brand == "未知" {
          // 尝试从评论内容中匹配用户指定的品牌
          contentLower := strings.ToLower(r.Content)
          for specBrandLower, origBrand := range specifiedBrands {
              if strings.Contains(contentLower, specBrandLower) {
                  brand = origBrand
                  break
              }
          }
      }

      if brand == "" {
          continue // 仍然没有品牌则跳过
      }

      comment := report.CommentWithScore{
          Content: r.Content,
          Scores:  r.Scores,
          Brand:   r.Brand,
          Model:   r.Model,
      }

      // 分类：指定品牌还是发现的新品牌
      brandLower := strings.ToLower(brand)
      isSpecified := false
      for specBrandLower, origBrand := range specifiedBrands {
          if strings.Contains(brandLower, specBrandLower) || strings.Contains(specBrandLower, brandLower) {
              specifiedResults[origBrand] = append(specifiedResults[origBrand], comment)
              isSpecified = true
              break
          }
      }

      if !isSpecified {
          // 这是新发现的品牌 - 保留它！
          discoveredResults[brand] = append(discoveredResults[brand], comment)
      }
  }

  // 记录发现的新品牌
  if len(discoveredResults) > 0 {
      var discoveredBrandNames []string
      for brand := range discoveredResults {
          discoveredBrandNames = append(discoveredBrandNames, brand)
      }
      log.Printf("[Task %s] 🔍 发现新品牌: %v", taskID, discoveredBrandNames)
  }

  // 合并结果：先指定品牌，再发现的品牌
  results := make(map[string][]report.CommentWithScore)
  for brand, comments := range specifiedResults {
      results[brand] = comments
      log.Printf("[Task %s] 指定品牌 %s: %d 条评论", taskID, brand, len(comments))
  }
  for brand, comments := range discoveredResults {
      results[brand] = comments
      log.Printf("[Task %s] 发现品牌 %s: %d 条评论", taskID, brand, len(comments))
  }
  ```

  **Commit**: YES
  - Message: `feat(executor): 改为品牌发现模式，收集所有AI识别的品牌`
  - Files: `backend/task/executor.go`

---

- [ ] 1.2. 添加品牌归一化逻辑（处理别名）

  **What to do**:
  - 在 `backend/task/executor.go` 添加品牌别名映射
  - 处理常见别名：Apple=苹果=iPhone, Dyson=戴森, Xiaomi=小米
  - 在分类时使用归一化后的品牌名

  **Must NOT do**:
  - 不硬编码所有品牌别名（只处理最常见的）

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 1.3
  - **Blocked By**: Task 1.1

  **References**:
  - `backend/task/executor.go` - 品牌分类逻辑

  **Acceptance Criteria**:
  - [ ] 添加 `normalizeBrand` 函数
  - [ ] 常见品牌别名被正确归一化
  - [ ] `go build ./backend/...` 通过

  **New Code to Add** (在 executor.go 文件末尾添加):
  ```go
  // normalizeBrand 品牌名称归一化
  // 处理常见的品牌别名，返回统一的品牌名称
  func normalizeBrand(brand string) string {
      brandLower := strings.ToLower(strings.TrimSpace(brand))
      
      // 品牌别名映射
      aliases := map[string][]string{
          "苹果":   {"apple", "iphone", "ipad", "mac", "airpods"},
          "戴森":   {"dyson"},
          "小米":   {"xiaomi", "mi", "redmi"},
          "华为":   {"huawei", "honor", "荣耀"},
          "三星":   {"samsung", "galaxy"},
          "索尼":   {"sony", "playstation", "ps5"},
          "小佩":   {"petkit"},
          "CATLINK": {"catlink", "猫猫狗狗"},
      }
      
      for normalized, aliasList := range aliases {
          for _, alias := range aliasList {
              if strings.Contains(brandLower, alias) || strings.Contains(alias, brandLower) {
                  return normalized
              }
          }
          // 也检查归一化名称本身
          if strings.Contains(brandLower, strings.ToLower(normalized)) {
              return normalized
          }
      }
      
      return brand // 返回原始名称
  }
  ```

  **Commit**: YES
  - Message: `feat(executor): 添加品牌名称归一化，处理常见别名`
  - Files: `backend/task/executor.go`

---

- [ ] 1.3. 前端显示发现的其他品牌

  **What to do**:
  - 修改 `frontend/src/pages/Report.tsx`
  - 在品牌排名表格中，用不同样式标记"发现的品牌"
  - 添加提示说明这些品牌是AI自动发现的

  **Must NOT do**:
  - 不改变现有的排名逻辑
  - 不修改API响应格式

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: None
  - **Blocked By**: Task 1.1, Task 1.2

  **References**:
  - `frontend/src/pages/Report.tsx:397-426` - 品牌排名表格

  **Acceptance Criteria**:
  - [ ] 发现的品牌在表格中有特殊标记（如 🔍 图标）
  - [ ] 添加 tooltip 说明"此品牌由AI自动发现"
  - [ ] `npm run build` 通过

  **Commit**: YES
  - Message: `feat(report): 前端标记AI发现的品牌`
  - Files: `frontend/src/pages/Report.tsx`

---

### Phase 2: 搜索优化（P1）

- [ ] 2.1. AI解析时生成通用搜索关键词

  **What to do**:
  - 修改 `backend/ai/keyword.go` 的 AI 提示词
  - 除了品牌+品类关键词，还要生成通用关键词
  - 例如："自动猫砂盆评测"、"自动猫砂盆推荐"、"猫砂盆横评"

  **Must NOT do**:
  - 不改变返回的数据结构

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/ai/keyword.go` - 关键词生成逻辑

  **Acceptance Criteria**:
  - [ ] AI 提示词要求生成通用关键词
  - [ ] 返回的关键词包含"XX评测"、"XX推荐"、"XX横评"等
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES
  - Message: `feat(ai): 生成通用搜索关键词，增加品牌发现机会`
  - Files: `backend/ai/keyword.go`

---

- [ ] 2.2. 增加视频数量限制到20

  **What to do**:
  - 修改 `backend/task/executor.go` 的 `DefaultTaskConfig`
  - 将 `MaxVideosPerKeyword` 从 10 改为 20

  **Must NOT do**:
  - 不改变其他配置

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go:36-43` - DefaultTaskConfig

  **Acceptance Criteria**:
  - [ ] `MaxVideosPerKeyword: 20`
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES (groups with 3.1)
  - Message: `perf(executor): 增加视频和评论抓取数量`
  - Files: `backend/task/executor.go`

---

- [ ] 2.3. 添加视频过滤（时长>1分钟）

  **What to do**:
  - 修改 `backend/bilibili/search.go` 的 `SearchVideosWithLimit` 函数
  - 过滤掉时长小于60秒的视频（短视频通常没有深度评测）

  **Must NOT do**:
  - 不改变API调用逻辑

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/bilibili/search.go:136-178` - SearchVideosWithLimit
  - `backend/bilibili/search.go:40` - Duration 字段

  **Acceptance Criteria**:
  - [ ] 添加 `parseDuration` 辅助函数
  - [ ] 过滤时长 < 60秒的视频
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES
  - Message: `feat(search): 过滤短视频，只保留时长>1分钟的评测视频`
  - Files: `backend/bilibili/search.go`

---

### Phase 3: 评论优化（P2）

- [ ] 3.1. 增加评论数量限制到200

  **What to do**:
  - 修改 `backend/task/executor.go` 的 `DefaultTaskConfig`
  - 将 `MaxCommentsPerVideo` 从 100 改为 200

  **Must NOT do**:
  - 不改变其他配置

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go:36-43` - DefaultTaskConfig

  **Acceptance Criteria**:
  - [ ] `MaxCommentsPerVideo: 200`
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES (groups with 2.2)
  - Message: `perf(executor): 增加视频和评论抓取数量`
  - Files: `backend/task/executor.go`

---

- [ ] 3.2. 评论抓取包含子评论（回复）

  **What to do**:
  - 检查 `backend/task/executor.go` 的 `GetAllCommentsWithVideo` 函数
  - 确保子评论（Replies）被正确收集

  **Must NOT do**:
  - 不改变评论数据结构

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go:484-510` - GetAllCommentsWithVideo
  - `backend/bilibili/comment.go` - 评论数据结构

  **Acceptance Criteria**:
  - [ ] 子评论被正确收集到结果中
  - [ ] 日志显示包含子评论的总数
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES
  - Message: `feat(executor): 确保子评论被正确收集`
  - Files: `backend/task/executor.go`

---

### Phase 4: 型号提取优化（P2-P3）

- [ ] 4.1. 增强AI提示词提取型号

  **What to do**:
  - 修改 `backend/ai/analysis.go` 的 `AnalyzeComment` 函数
  - 增强系统提示词，强调从评论内容中提取具体型号
  - 添加型号示例帮助AI理解

  **Must NOT do**:
  - 不改变返回的JSON格式

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/ai/analysis.go:67-92` - 系统提示词

  **Acceptance Criteria**:
  - [ ] 提示词强调从评论内容提取型号
  - [ ] 添加型号示例（如 "V12", "Pro Max", "G10"）
  - [ ] `go build ./backend/...` 通过

  **New Prompt Enhancement**:
  ```
  1. 从视频标题和评论内容中识别：
     - 品牌名称（如"戴森"、"小米"、"苹果"、"Sony"）
     - 具体型号（如"V12"、"iPhone 15 Pro"、"G10"、"WH-1000XM5"）
     
  重要：型号提取规则：
  - 优先从评论内容中提取具体型号
  - 常见型号格式：字母+数字（V12, G10）、品牌+型号（iPhone 15）、系列名（Pro, Max, Plus）
  - 如果评论提到多个型号，选择评论主要讨论的那个
  - 如果无法确定具体型号，尝试提取系列名（如"V系列"、"Pro系列"）
  ```

  **Commit**: YES
  - Message: `feat(ai): 增强型号提取提示词`
  - Files: `backend/ai/analysis.go`

---

- [ ] 4.2. 添加型号正则匹配（P3 - 可选）

  **What to do**:
  - 在 `backend/task/executor.go` 添加型号正则匹配作为后备
  - 如果AI没有提取到型号，尝试用正则从评论中提取

  **Must NOT do**:
  - 不覆盖AI提取的结果

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go` - 评论处理逻辑

  **Acceptance Criteria**:
  - [ ] 添加 `extractModelFromContent` 函数
  - [ ] 正则匹配常见型号格式
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES
  - Message: `feat(executor): 添加型号正则匹配作为后备`
  - Files: `backend/task/executor.go`

---

- [ ] 4.3. 型号归一化（P3 - 可选）

  **What to do**:
  - 在 `backend/report/generator.go` 添加型号归一化逻辑
  - 合并相似型号名称（如 "V12 Detect" = "V12"）

  **Must NOT do**:
  - 不过度合并（保留重要区分）

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: None
  - **Blocked By**: None

  **References**:
  - `backend/report/generator.go:152-238` - 型号聚合逻辑

  **Acceptance Criteria**:
  - [ ] 添加 `normalizeModel` 函数
  - [ ] 相似型号被合并
  - [ ] `go build ./backend/...` 通过

  **Commit**: YES
  - Message: `feat(report): 添加型号归一化，合并相似型号`
  - Files: `backend/report/generator.go`

---

### Phase 5: 验证

- [ ] 5.1. 编译验证和测试

  **What to do**:
  - 运行后端编译测试
  - 运行前端编译测试
  - 运行一个完整的分析任务验证功能

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Final
  - **Blocks**: None
  - **Blocked By**: All previous tasks

  **Acceptance Criteria**:
  - [ ] `go build ./backend/...` 通过
  - [ ] `cd frontend && npm run build` 通过
  - [ ] 运行分析任务，日志显示"发现新品牌"
  - [ ] 报告中包含发现的品牌

  **Commit**: NO (verification only)

---

## Commit Strategy

| After Task | Message | Files |
|------------|---------|-------|
| 1.1 | `feat(executor): 改为品牌发现模式` | executor.go |
| 1.2 | `feat(executor): 添加品牌归一化` | executor.go |
| 1.3 | `feat(report): 前端标记发现的品牌` | Report.tsx |
| 2.1 | `feat(ai): 生成通用搜索关键词` | keyword.go |
| 2.2+3.1 | `perf(executor): 增加抓取数量` | executor.go |
| 2.3 | `feat(search): 过滤短视频` | search.go |
| 3.2 | `feat(executor): 收集子评论` | executor.go |
| 4.1 | `feat(ai): 增强型号提取` | analysis.go |
| 4.2 | `feat(executor): 型号正则匹配` | executor.go |
| 4.3 | `feat(report): 型号归一化` | generator.go |

---

## Success Criteria

### Verification Commands
```bash
# 后端编译
go build ./backend/...

# 前端编译
cd frontend && npm run build

# 运行分析任务
curl -X POST http://localhost:8080/api/confirm \
  -H "Content-Type: application/json" \
  -d '{"requirement":"自动猫砂盆","brands":["小佩"],"dimensions":[{"name":"清洁效果","description":"清洁能力"}],"keywords":["自动猫砂盆评测"]}'

# 检查日志是否有"发现新品牌"
```

### Final Checklist
- [ ] 品牌发现模式生效（日志显示发现的品牌）
- [ ] 报告包含用户未指定的品牌
- [ ] 视频抓取数量增加到20
- [ ] 评论抓取数量增加到200
- [ ] 型号提取率提高
- [ ] 所有编译测试通过
