# 型号归一化与未知品牌识别

## TL;DR

> **Quick Summary**: 对品牌为"未知"的型号批量发送AI识别品牌，同时对型号进行归一化处理（去空格、统一大小写），合并同一型号的数据。
> 
> **Deliverables**: 
> - AI批量识别未知品牌功能
> - 型号归一化函数（TWS5、TWS 5、Tws5 → TWS5）
> - 品牌格式化函数（oppo → OPPO）
> 
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 2 waves
> **Critical Path**: Task 1 → Task 2 → Task 3 → Task 4

---

## Context

### Original Request
用户希望：
1. 对品牌为"未知"的型号，批量发送给AI识别品牌（一次性提交，JSON返回）
2. 型号归一化：处理大小写、空格差异，合并同一型号（如 TWS5、TWS 5、Tws5）
3. 品牌格式化：纯字母品牌全大写（oppo → OPPO, vivo → VIVO）

### Interview Summary
**Key Discussions**:
- 归一化时必须带上品牌，避免不同品牌的同名型号被合并
- 品牌格式：纯字母 → 全大写，中文 → 保持原样

**Research Findings**:
- `executor.go:434-450` 处理品牌识别，品牌为"未知"时尝试从评论匹配
- `report/generator.go:182-254` 处理型号聚合，使用 "brand|model" 作为key
- `ai/analysis.go` 包含AI分析接口

---

## Work Objectives

### Core Objective
增强型号和品牌的识别与归一化能力，提高数据聚合的准确性。

### Concrete Deliverables
- `backend/ai/brand_identify.go` - AI批量识别未知品牌
- `backend/task/executor.go` - 调用品牌识别，应用品牌格式化
- `backend/report/generator.go` - 型号归一化函数

### Definition of Done
- [x] 品牌为"未知"的型号被AI批量识别
- [x] TWS5、TWS 5、Tws5 被合并为同一型号
- [x] oppo、OPPO、Oppo 被统一为 OPPO
- [x] `go build ./backend/...` 通过
- [x] 运行分析任务验证功能

### Must Have
- AI批量识别未知品牌
- 型号归一化（去空格、统一格式）
- 品牌格式化（字母全大写）
- 归一化时保留品牌区分

### Must NOT Have (Guardrails)
- 不修改数据库结构
- 不改变API接口格式
- 不合并不同品牌的同名型号
- 不过度调用AI（只对未知品牌调用一次）

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (go build)
- **User wants tests**: Manual verification
- **Framework**: Go build + 运行分析任务

### Automated Verification
```bash
# 后端编译
go build ./backend/...

# 运行分析任务，检查日志
# 应该看到：
# - "🔍 AI识别未知品牌: [型号列表]"
# - "✅ 品牌识别完成: TWS5 -> OPPO"
# - 型号排名中相似型号被合并
```

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: 品牌格式化函数 (formatBrandName)
└── Task 2: 型号归一化函数 (normalizeModelKey, getDisplayModel)

Wave 2 (After Wave 1):
├── Task 3: AI批量识别未知品牌 (IdentifyBrandsForModels)
└── Task 4: 集成到executor和generator

Final:
└── Task 5: 编译验证
```

---

## TODOs

### Task 1: 品牌格式化函数

- [x] 1. 添加品牌格式化函数 `formatBrandName`

  **What to do**:
  - 在 `backend/task/executor.go` 添加 `formatBrandName` 函数
  - 纯字母品牌 → 全大写（oppo → OPPO）
  - 中文品牌 → 保持原样
  - 混合品牌 → 保持原样（避免破坏格式）

  **Must NOT do**:
  - 不改变已有的 `normalizeBrand` 函数逻辑

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 4
  - **Blocked By**: None

  **References**:
  - `backend/task/executor.go:559-600` - 现有的 `normalizeBrand` 函数

  **Acceptance Criteria**:
  - [ ] `formatBrandName("oppo")` → `"OPPO"`
  - [ ] `formatBrandName("vivo")` → `"VIVO"`
  - [ ] `formatBrandName("sony")` → `"SONY"`
  - [ ] `formatBrandName("小米")` → `"小米"`
  - [ ] `formatBrandName("华为")` → `"华为"`
  - [ ] `go build ./backend/...` 通过

  **New Code to Add** (在 executor.go 文件中添加):
  ```go
  // formatBrandName 格式化品牌名称
  // 纯字母品牌转全大写，中文品牌保持原样
  func formatBrandName(brand string) string {
      brand = strings.TrimSpace(brand)
      if brand == "" {
          return brand
      }
      
      // 检查是否为纯字母（ASCII字母）
      isPureAlpha := true
      for _, r := range brand {
          if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
              isPureAlpha = false
              break
          }
      }
      
      if isPureAlpha {
          return strings.ToUpper(brand)
      }
      
      return brand
  }
  ```

  **Commit**: YES
  - Message: `feat(executor): 添加品牌格式化函数，纯字母品牌转全大写`
  - Files: `backend/task/executor.go`

---

### Task 2: 型号归一化函数

- [x] 2. 添加型号归一化函数

  **What to do**:
  - 在 `backend/report/generator.go` 添加型号归一化函数
  - `normalizeModelKey(brand, model)` - 生成归一化的key用于比对
  - `getDisplayModel(models)` - 从多个变体中选择最佳显示名称
  - 归一化规则：去除空格 + 转小写 + 带品牌前缀

  **Must NOT do**:
  - 不合并不同品牌的同名型号

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 4
  - **Blocked By**: None

  **References**:
  - `backend/report/generator.go:182-254` - 现有型号聚合逻辑

  **Acceptance Criteria**:
  - [ ] `normalizeModelKey("OPPO", "TWS5")` → `"oppo|tws5"`
  - [ ] `normalizeModelKey("OPPO", "TWS 5")` → `"oppo|tws5"`
  - [ ] `normalizeModelKey("OPPO", "Tws5")` → `"oppo|tws5"`
  - [ ] `normalizeModelKey("小米", "Air 2")` → `"小米|air2"`
  - [ ] 不同品牌的同名型号不会被合并
  - [ ] `go build ./backend/...` 通过

  **New Code to Add** (在 generator.go 文件中添加):
  ```go
  // normalizeModelKey 生成归一化的型号key用于比对
  // 规则：品牌小写 + "|" + 型号去空格小写
  // 例如：("OPPO", "TWS 5") -> "oppo|tws5"
  func normalizeModelKey(brand, model string) string {
      // 品牌转小写
      brandKey := strings.ToLower(strings.TrimSpace(brand))
      
      // 型号：去除所有空格，转小写
      modelKey := strings.ToLower(strings.TrimSpace(model))
      modelKey = strings.ReplaceAll(modelKey, " ", "")
      modelKey = strings.ReplaceAll(modelKey, "-", "")
      modelKey = strings.ReplaceAll(modelKey, "_", "")
      
      return brandKey + "|" + modelKey
  }

  // getDisplayModel 从多个型号变体中选择最佳显示名称
  // 优先选择：有空格分隔的 > 首字母大写的 > 第一个出现的
  func getDisplayModel(variants []string) string {
      if len(variants) == 0 {
          return ""
      }
      if len(variants) == 1 {
          return variants[0]
      }
      
      // 优先选择有空格的版本（如 "TWS 5" 优于 "TWS5"）
      for _, v := range variants {
          if strings.Contains(v, " ") {
              return v
          }
      }
      
      // 其次选择首字母大写的版本
      for _, v := range variants {
          if len(v) > 0 && v[0] >= 'A' && v[0] <= 'Z' {
              return v
          }
      }
      
      // 默认返回第一个
      return variants[0]
  }
  ```

  **Commit**: YES
  - Message: `feat(report): 添加型号归一化函数，合并相似型号`
  - Files: `backend/report/generator.go`

---

### Task 3: AI批量识别未知品牌

- [x] 3. 添加AI批量识别未知品牌功能

  **What to do**:
  - 在 `backend/ai/brand_identify.go` 新建文件
  - 添加 `IdentifyBrandsForModels` 函数
  - 一次性提交所有未知品牌的型号，JSON返回型号→品牌映射
  - 使用简洁的提示词，减少token消耗

  **Must NOT do**:
  - 不为每个型号单独调用AI（太慢太贵）
  - 不修改现有的 `AnalyzeComment` 函数

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 4
  - **Blocked By**: None

  **References**:
  - `backend/ai/analysis.go` - AI调用模式
  - `backend/ai/keyword.go` - JSON解析模式

  **Acceptance Criteria**:
  - [ ] 新建 `backend/ai/brand_identify.go` 文件
  - [ ] `IdentifyBrandsForModels(ctx, []string{"TWS5", "Air 2", "V12"})` 返回 `map[string]string{"TWS5": "OPPO", "Air 2": "小米", "V12": "戴森"}`
  - [ ] 一次API调用处理所有型号
  - [ ] `go build ./backend/...` 通过

  **New File Content** (`backend/ai/brand_identify.go`):
  ```go
  package ai

  import (
      "context"
      "encoding/json"
      "fmt"
      "log"
      "strings"
  )

  // BrandIdentifyRequest 品牌识别请求
  type BrandIdentifyRequest struct {
      Models []string // 需要识别品牌的型号列表
  }

  // BrandIdentifyResponse 品牌识别响应
  type BrandIdentifyResponse struct {
      Results map[string]string `json:"results"` // 型号 -> 品牌
  }

  // IdentifyBrandsForModels 批量识别型号对应的品牌
  // 一次性提交所有未知品牌的型号，返回型号→品牌映射
  func (c *Client) IdentifyBrandsForModels(ctx context.Context, models []string) (map[string]string, error) {
      if len(models) == 0 {
          return make(map[string]string), nil
      }

      // 去重
      uniqueModels := make([]string, 0, len(models))
      seen := make(map[string]bool)
      for _, m := range models {
          m = strings.TrimSpace(m)
          if m != "" && !seen[strings.ToLower(m)] {
              uniqueModels = append(uniqueModels, m)
              seen[strings.ToLower(m)] = true
          }
      }

      if len(uniqueModels) == 0 {
          return make(map[string]string), nil
      }

      log.Printf("[AI] 🔍 批量识别 %d 个型号的品牌...", len(uniqueModels))

      systemPrompt := `你是一个专业的电子产品型号识别专家。
  根据型号名称识别对应的品牌。

  常见品牌和型号特征：
  - OPPO: TWS系列、Enco系列、Find系列
  - VIVO: TWS系列、X系列
  - 小米/Xiaomi: Air系列、Buds系列、Pro系列
  - 华为/Huawei: FreeBuds系列、Mate系列
  - 苹果/Apple: AirPods系列、iPhone系列、MacBook系列
  - 戴森/Dyson: V系列(V10/V12/V15)、Airwrap系列
  - 索尼/Sony: WH系列、WF系列、XM系列
  - 三星/Samsung: Galaxy系列、Buds系列
  - JBL: Tour系列、Live系列、Tune系列
  - Bose: QuietComfort系列、SoundSport系列
  - 漫步者/Edifier: TWS系列、LolliPods系列
  - 1MORE/万魔: ComfoBuds系列、EVO系列

  规则：
  1. 如果能确定品牌，返回品牌名称（中文优先，如"小米"而非"Xiaomi"）
  2. 纯字母品牌用全大写（如OPPO、VIVO、JBL、BOSE）
  3. 如果无法确定，返回"未知"
  4. 必须严格返回JSON格式`

      userPrompt := fmt.Sprintf(`请识别以下型号对应的品牌，返回JSON格式：

  型号列表：
  %s

  返回格式示例：
  {"results": {"TWS5": "OPPO", "Air 2": "小米", "V12": "戴森"}}`, strings.Join(uniqueModels, "\n"))

      resp, err := c.chat(ctx, systemPrompt, userPrompt)
      if err != nil {
          return nil, fmt.Errorf("AI调用失败: %w", err)
      }

      // 解析JSON响应
      var result BrandIdentifyResponse
      
      // 尝试提取JSON
      jsonStr := resp
      if idx := strings.Index(resp, "{"); idx != -1 {
          if endIdx := strings.LastIndex(resp, "}"); endIdx != -1 {
              jsonStr = resp[idx : endIdx+1]
          }
      }

      if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
          log.Printf("[AI] 品牌识别JSON解析失败: %v, 原始响应: %s", err, resp)
          return make(map[string]string), nil
      }

      log.Printf("[AI] ✅ 品牌识别完成: %v", result.Results)
      return result.Results, nil
  }
  ```

  **Commit**: YES
  - Message: `feat(ai): 添加AI批量识别未知品牌功能`
  - Files: `backend/ai/brand_identify.go`

---

### Task 4: 集成到executor和generator

- [x] 4. 集成品牌识别和型号归一化

  **What to do**:
  - 在 `executor.go` 的 `analyzeComments` 函数中：
    1. 收集所有品牌为"未知"但有型号的评论
    2. 调用 `IdentifyBrandsForModels` 批量识别
    3. 更新评论的品牌字段
    4. 应用 `formatBrandName` 格式化品牌
  - 在 `generator.go` 的型号聚合逻辑中：
    1. 使用 `normalizeModelKey` 生成归一化key
    2. 使用 `getDisplayModel` 选择显示名称
    3. 合并相似型号的数据

  **Must NOT do**:
  - 不改变函数签名
  - 不破坏现有的品牌发现模式

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (after Task 1, 2, 3)
  - **Blocks**: Task 5
  - **Blocked By**: Task 1, Task 2, Task 3

  **References**:
  - `backend/task/executor.go:429-500` - 品牌分类逻辑
  - `backend/report/generator.go:182-254` - 型号聚合逻辑

  **Acceptance Criteria**:
  - [ ] 品牌为"未知"的型号被AI识别后更新
  - [ ] 日志显示 "🔍 AI识别未知品牌: [型号列表]"
  - [ ] 型号排名中 TWS5、TWS 5、Tws5 被合并
  - [ ] 品牌显示为全大写（OPPO、VIVO）
  - [ ] `go build ./backend/...` 通过

  **Code Changes for executor.go** (在 analyzeComments 函数中，约 line 415 后添加):
  ```go
  // === 批量识别未知品牌 ===
  // 收集品牌为"未知"但有型号的评论
  unknownBrandModels := make(map[string]bool) // 型号 -> 是否需要识别
  for _, r := range analysisResults {
      if r.Error != "" || r.Scores == nil {
          continue
      }
      brand := strings.TrimSpace(r.Brand)
      model := strings.TrimSpace(r.Model)
      if (brand == "" || brand == "未知") && model != "" && model != "未知" && model != "通用" {
          unknownBrandModels[model] = true
      }
  }

  // 批量调用AI识别品牌
  var modelToBrand map[string]string
  if len(unknownBrandModels) > 0 {
      models := make([]string, 0, len(unknownBrandModels))
      for m := range unknownBrandModels {
          models = append(models, m)
      }
      log.Printf("[Task %s] 🔍 AI识别未知品牌: %v", taskID, models)
      
      var err error
      modelToBrand, err = aiClient.IdentifyBrandsForModels(ctx, models)
      if err != nil {
          log.Printf("[Task %s] ⚠️ 品牌识别失败: %v", taskID, err)
          modelToBrand = make(map[string]string)
      }
  }

  // 更新分析结果中的品牌
  for i := range analysisResults {
      r := &analysisResults[i]
      brand := strings.TrimSpace(r.Brand)
      model := strings.TrimSpace(r.Model)
      
      // 如果品牌未知，尝试从AI识别结果获取
      if (brand == "" || brand == "未知") && model != "" {
          if identifiedBrand, ok := modelToBrand[model]; ok && identifiedBrand != "" && identifiedBrand != "未知" {
              r.Brand = identifiedBrand
              brand = identifiedBrand
          }
      }
      
      // 格式化品牌名称（纯字母转大写）
      if brand != "" {
          r.Brand = formatBrandName(brand)
      }
  }
  ```

  **Code Changes for generator.go** (修改型号聚合逻辑，约 line 182-212):
  ```go
  // 按"品牌+型号"聚合，生成型号排名（使用归一化key）
  modelScores := make(map[string]map[string][]float64)    // normalizedKey -> 维度 -> 分数列表
  modelCommentCounts := make(map[string]int)              // normalizedKey -> 评论数
  modelVariants := make(map[string][]string)              // normalizedKey -> 原始型号变体列表
  modelBrands := make(map[string]string)                  // normalizedKey -> 品牌

  for brandKey, results := range input.AnalysisResults {
      for _, result := range results {
          model := strings.TrimSpace(result.Model)
          if model == "" || model == "未知" || model == "通用" {
              continue
          }
          brand := strings.TrimSpace(result.Brand)
          if brand == "" {
              brand = strings.TrimSpace(brandKey)
          }
          if brand == "" || brand == "未知" {
              continue
          }

          // 使用归一化key进行聚合
          normalizedKey := normalizeModelKey(brand, model)
          
          if modelScores[normalizedKey] == nil {
              modelScores[normalizedKey] = make(map[string][]float64)
              modelBrands[normalizedKey] = brand
          }
          
          // 记录型号变体
          modelVariants[normalizedKey] = append(modelVariants[normalizedKey], model)

          for dimName, score := range result.Scores {
              if score != nil {
                  modelScores[normalizedKey][dimName] = append(modelScores[normalizedKey][dimName], *score)
              }
          }
          modelCommentCounts[normalizedKey]++
      }
  }

  modelRankings := make([]ModelRanking, 0, len(modelScores))
  for normalizedKey, dimScores := range modelScores {
      commentCount := modelCommentCounts[normalizedKey]
      if commentCount < 1 {
          continue
      }
      
      brand := modelBrands[normalizedKey]
      // 选择最佳显示名称
      displayModel := getDisplayModel(modelVariants[normalizedKey])

      // ... 后续计算平均分的逻辑保持不变 ...
  ```

  **Commit**: YES
  - Message: `feat: 集成品牌识别和型号归一化，合并相似型号`
  - Files: `backend/task/executor.go`, `backend/report/generator.go`

---

### Task 5: 编译验证

- [x] 5. 编译验证和测试

  **What to do**:
  - 运行后端编译测试
  - 检查代码逻辑正确性

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
  - [ ] 无编译错误或警告

  **Commit**: NO (verification only)

---

## Commit Strategy

| After Task | Message | Files |
|------------|---------|-------|
| 1 | `feat(executor): 添加品牌格式化函数` | executor.go |
| 2 | `feat(report): 添加型号归一化函数` | generator.go |
| 3 | `feat(ai): 添加AI批量识别未知品牌` | brand_identify.go |
| 4 | `feat: 集成品牌识别和型号归一化` | executor.go, generator.go |

---

## Success Criteria

### Verification Commands
```bash
# 后端编译
go build ./backend/...
```

### Final Checklist
- [x] 品牌为"未知"的型号被AI批量识别
- [x] TWS5、TWS 5、Tws5 被合并为同一型号
- [x] oppo、OPPO、Oppo 被统一为 OPPO
- [x] 不同品牌的同名型号不会被合并
- [x] 所有编译测试通过
