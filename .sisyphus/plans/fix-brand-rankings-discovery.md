# 修复品牌排名未包含AI发现品牌的问题

## TL;DR

> **Quick Summary**: 修复报告生成逻辑，使品牌排名包含所有AI发现的品牌，而不仅仅是用户指定的品牌。
> 
> **Deliverables**: 
> - 修改 `backend/report/generator.go` 中的三处逻辑
> 
> **Estimated Effort**: Small
> **Parallel Execution**: NO - 顺序执行
> **Critical Path**: Task 1 → Task 2 → Task 3 → Task 4

---

## Context

### 问题现象
报告页面显示：
- **顶部品牌排名**: 只有 1 个品牌（用户指定的）
- **型号排名**: 有 170 个型号，涉及 50+ 个品牌

### 根本原因
`generateRankings` 函数只遍历 `input.Brands`（用户指定的品牌），而不是 `scores` map 中所有AI发现的品牌。

```
用户输入: brands = ["苹果"]  ← 用户只指定了1个品牌

AI分析后: scores = {
    "苹果": {...},
    "华为": {...},      ← AI发现的新品牌被忽略
    "索尼": {...},
    "三星": {...},
    ... 50+ 品牌
}
```

### 解决方案
修改 `generateRankings` 遍历 `scores` 而非 `brands`，同时更新相关函数。

---

## Work Objectives

### Core Objective
让品牌排名包含所有AI发现的品牌，与型号排名保持一致。

### Concrete Deliverables
- `backend/report/generator.go` - 修改三处逻辑

### Definition of Done
- [x] `generateRankings` 遍历 `scores` 中的所有品牌
- [x] `brandAnalysis` 包含所有发现的品牌
- [x] `ReportData.Brands` 返回所有品牌（按排名顺序）
- [x] `go build ./backend/...` 通过

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: NO (无单元测试)
- **User wants tests**: Manual-only
- **QA approach**: 编译验证 + 手动检查报告数据

---

## TODOs

### Task 1: 修改 generateRankings 函数

- [x] 1. 修改 generateRankings 遍历所有品牌

  **What to do**:
  - 修改函数遍历 `scores` map 而非 `brands` 参数
  - 添加空品牌过滤

  **References**:
  - `backend/report/generator.go:304-311` - 当前实现

  **Current Code** (第304-311行):
  ```go
  // generateRankings 生成品牌排名
  // 根据各维度得分计算综合得分，并按综合得分排序
  func generateRankings(brands []string, dimensions []ai.Dimension, scores map[string]map[string]float64) []BrandRanking {
  	rankings := make([]BrandRanking, 0, len(brands))

  	// 为每个品牌计算综合得分
  	for _, brand := range brands {
  		brandScores := scores[brand]
  ```

  **New Code**:
  ```go
  // generateRankings 生成品牌排名
  // 根据各维度得分计算综合得分，并按综合得分排序
  // 注意：遍历 scores 中的所有品牌（包括AI发现的新品牌），而不仅仅是用户指定的品牌
  func generateRankings(brands []string, dimensions []ai.Dimension, scores map[string]map[string]float64) []BrandRanking {
  	rankings := make([]BrandRanking, 0, len(scores))

  	// 为每个品牌计算综合得分（遍历所有发现的品牌）
  	for brand, brandScores := range scores {
  		if brand == "" {
  			continue
  		}
  ```

  **Acceptance Criteria**:
  - [x] 函数遍历 `scores` map 而非 `brands` 参数
  - [x] 添加了空品牌过滤 `if brand == "" { continue }`
  - [x] 注释已更新说明新逻辑

---

### Task 2: 修改 GenerateReportWithInput 中的 brandAnalysis 调用

- [x] 2. 修改 brandAnalysis 使用所有发现的品牌

  **What to do**:
  - 在调用 `generateBrandAnalysis` 前收集所有品牌
  - 传递完整品牌列表

  **References**:
  - `backend/report/generator.go:183-186` - 当前实现

  **Current Code** (第183-186行):
  ```go
  	rankings := generateRankings(input.Brands, input.Dimensions, scores)
  	recommendation := generateRecommendation(rankings, input.Dimensions)
  	brandAnalysis := generateBrandAnalysis(input.Brands, input.Dimensions, scores)
  ```

  **New Code**:
  ```go
  	rankings := generateRankings(input.Brands, input.Dimensions, scores)
  	recommendation := generateRecommendation(rankings, input.Dimensions)

  	// 收集所有发现的品牌（用于品牌分析）
  	allBrands := make([]string, 0, len(scores))
  	for brand := range scores {
  		if brand != "" {
  			allBrands = append(allBrands, brand)
  		}
  	}
  	brandAnalysis := generateBrandAnalysis(allBrands, input.Dimensions, scores)
  ```

  **Acceptance Criteria**:
  - [x] 新增 `allBrands` 变量收集所有品牌
  - [x] `generateBrandAnalysis` 使用 `allBrands` 而非 `input.Brands`

---

### Task 3: 修改 ReportData 返回值中的 Brands 字段

- [x] 3. 修改 Brands 字段返回所有品牌

  **What to do**:
  - 从排名结果中提取所有品牌名称（已按得分排序）
  - 更新 `ReportData.Brands` 字段

  **References**:
  - `backend/report/generator.go:191-195` - 当前实现

  **Current Code** (第191-195行):
  ```go
  	return &ReportData{
  		Category:       input.Category,
  		Brands:         input.Brands,
  		Dimensions:     input.Dimensions,
  ```

  **New Code**:
  ```go
  	// 收集所有品牌名称用于报告（按排名顺序）
  	allBrandNames := make([]string, 0, len(rankings))
  	for _, r := range rankings {
  		allBrandNames = append(allBrandNames, r.Brand)
  	}

  	return &ReportData{
  		Category:       input.Category,
  		Brands:         allBrandNames,
  		Dimensions:     input.Dimensions,
  ```

  **Acceptance Criteria**:
  - [x] 新增 `allBrandNames` 变量从排名中提取品牌
  - [x] `Brands` 字段使用 `allBrandNames` 而非 `input.Brands`

---

### Task 4: 编译验证

- [x] 4. 编译验证

  **What to do**:
  - 运行 `go build ./backend/...`
  - 确保无编译错误

  **Acceptance Criteria**:
  - [x] `go build ./backend/...` 通过

---

## Success Criteria

### Verification Commands
```bash
go build ./backend/...
```

### Final Checklist
- [x] `generateRankings` 遍历所有发现的品牌
- [x] `brandAnalysis` 包含所有发现的品牌
- [x] `ReportData.Brands` 返回所有品牌（按排名顺序）
- [x] 编译通过
