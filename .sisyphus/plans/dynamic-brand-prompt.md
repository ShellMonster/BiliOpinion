# 动态品牌识别 systemPrompt 优化

## TL;DR

> **Quick Summary**: 将品牌识别的 systemPrompt 从硬编码 switch-case 改为动态上下文注入，利用任务中已有的信息（商品类别、已知品牌、已发现品牌）生成针对性提示词。
> 
> **Deliverables**: 
> - 重构 `backend/ai/brand_identify.go` - 删除 switch-case，改用动态模板
> - 修改 `backend/task/executor.go` - 传递完整上下文
> 
> **Estimated Effort**: Small
> **Parallel Execution**: NO - 顺序执行
> **Critical Path**: Task 1 → Task 2 → Task 3

---

## Context

### 原始任务背景
B站商品评论分析工具，用于抓取评论、AI分析品牌/型号、生成评分报告。

### 当前问题
品牌识别使用硬编码 switch-case，每个行业一个 case：
- 不可扩展（每增加行业要加代码）
- 跨行业型号冲突（"V12" 在不同行业是不同品牌）
- 没有利用任务上下文

### 解决方案：方案C - 混合方案
- 删除所有 switch-case
- 通用模板 + 动态注入上下文
- 利用：商品类别、用户指定品牌、AI已发现品牌

---

## Work Objectives

### Core Objective
让品牌识别的 systemPrompt 根据任务上下文动态生成，而不是硬编码行业知识。

### Concrete Deliverables
- `backend/ai/brand_identify.go` - 新增 `BrandIdentifyContext` 结构体，重写 `IdentifyBrandsForModels` 函数
- `backend/task/executor.go` - 修改调用方式，传递完整上下文

### Definition of Done
- [x] 删除所有 switch-case 硬编码
- [x] 新增 `BrandIdentifyContext` 结构体
- [x] `IdentifyBrandsForModels` 接收上下文参数
- [x] `executor.go` 传递 knownBrands 和 discoveredBrands
- [x] `go build ./backend/...` 通过

---

## TODOs

### Task 1: 重构 brand_identify.go

- [x] 1. 重构品牌识别模块

  **What to do**:
  - 新增 `BrandIdentifyContext` 结构体
  - 修改 `IdentifyBrandsForModels` 函数签名，接收 `BrandIdentifyContext`
  - 新增 `buildDynamicBrandPrompt` 函数，动态生成 systemPrompt
  - 删除 `generateBrandIdentifyPrompt` 函数及所有 switch-case

  **References**:
  - `backend/ai/brand_identify.go` - 当前实现

  **New Code**:
  ```go
  // BrandIdentifyContext 品牌识别上下文
  type BrandIdentifyContext struct {
      Category         string   // 商品类别："自动猫砂盆"
      KnownBrands      []string // 用户指定的品牌：["小佩", "CATLINK"]
      DiscoveredBrands []string // AI已识别的品牌：["霍曼", "美的"]
  }

  // IdentifyBrandsForModels 批量识别型号对应的品牌
  func (c *Client) IdentifyBrandsForModels(ctx context.Context, models []string, identifyCtx BrandIdentifyContext) (map[string]string, error)

  // buildDynamicBrandPrompt 根据上下文动态构建提示词
  func buildDynamicBrandPrompt(ctx BrandIdentifyContext) string
  ```

  **Acceptance Criteria**:
  - [x] `BrandIdentifyContext` 结构体已定义
  - [x] `IdentifyBrandsForModels` 接收 `BrandIdentifyContext` 参数
  - [x] `buildDynamicBrandPrompt` 函数实现动态模板
  - [x] 删除所有 switch-case 代码

---

### Task 2: 修改 executor.go 调用方式

- [x] 2. 修改 executor.go 传递完整上下文

  **What to do**:
  - 在 `analyzeComments` 函数中，收集已发现的品牌
  - 构建 `BrandIdentifyContext`，传递给 `IdentifyBrandsForModels`

  **References**:
  - `backend/task/executor.go:430-447` - 当前调用位置

  **Current Code**:
  ```go
  modelToBrand, err = aiClient.IdentifyBrandsForModels(ctx, models, category)
  ```

  **New Code**:
  ```go
  // 收集已发现的品牌（从analysisResults中提取）
  discoveredBrands := collectDiscoveredBrands(analysisResults)
  
  identifyCtx := ai.BrandIdentifyContext{
      Category:         category,
      KnownBrands:      brands,
      DiscoveredBrands: discoveredBrands,
  }
  modelToBrand, err = aiClient.IdentifyBrandsForModels(ctx, models, identifyCtx)
  ```

  **Acceptance Criteria**:
  - [x] 新增 `collectDiscoveredBrands` 辅助函数
  - [x] 构建 `BrandIdentifyContext` 并传递
  - [x] 编译通过

---

### Task 3: 编译验证

- [x] 3. 编译验证

  **What to do**:
  - 运行 `go build ./backend/...`
  - 确保无编译错误

  **Acceptance Criteria**:
  - [x] `go build ./backend/...` 通过

---

## Dynamic Prompt Template

最终生成的 systemPrompt 示例：

```
你是一个专业的【自动猫砂盆】产品型号识别专家。

## 任务背景
- 商品类别：自动猫砂盆
- 用户关注的品牌：小佩、CATLINK
- 已识别到的同类品牌：霍曼、美的

## 识别规则
1. **优先匹配**：如果型号明显属于已知品牌或已识别品牌，直接返回该品牌
2. **行业推断**：根据商品类别和已知品牌，推断该行业的其他常见品牌
3. **命名规律**：分析型号的命名规律（如前缀、系列名）来判断品牌
4. **品牌格式**：
   - 纯字母品牌用全大写（如 OPPO、CATLINK、JBL）
   - 中文品牌保持原样（如 小米、华为、小佩）
5. **无法确定**：如果确实无法判断，返回"未知"

## 重要提示
- 这是【自动猫砂盆】行业的型号，请在该行业范围内识别
- 同一型号在不同行业可能属于不同品牌，请根据上下文判断
- 必须严格返回JSON格式
```

---

## Success Criteria

### Verification Commands
```bash
go build ./backend/...
```

### Final Checklist
- [x] 删除所有 switch-case 硬编码
- [x] systemPrompt 根据上下文动态生成
- [x] 包含商品类别、已知品牌、已发现品牌
- [x] 编译通过
