# 修复品牌/型号解析问题

## TL;DR

> **Quick Summary**: 修复AI解析评论时品牌包含"/"分隔符和型号为"通用"的问题，通过优化AI提示词+添加后处理逻辑双管齐下。
> 
> **Deliverables**:
> - 优化AI提示词，明确品牌必须单一
> - 添加品牌清洗函数，拆分"/"分隔的品牌
> - 添加单元测试验证清洗逻辑
> 
> **Estimated Effort**: Medium (2-3小时)
> **Parallel Execution**: YES - 2 waves
> **Critical Path**: Task 1 (提示词) + Task 2 (清洗函数) → Task 3 (集成)

---

## Context

### 问题背景
从报告 ID=6 的数据分析发现：
- 总型号数：76
- "通用"型号数：38 (50%)
- 包含"/"的异常品牌：8个

### 异常数据示例
```
品牌: "喵洁易/Catlink"     → 实际是2个品牌
品牌: "有陪/小佩/小米/糯雪" → 实际是4个品牌
型号: "通用"              → 未能提取具体型号
```

### 根本原因
1. AI提示词没有明确"品牌必须单一"
2. 缺少后处理逻辑清洗异常数据

---

## Work Objectives

### Core Objective
确保AI返回的品牌是单一品牌名称，型号尽可能具体。

### Concrete Deliverables
- `backend/ai/analysis.go` - 优化AI提示词
- `backend/comment/brand_cleaner.go` - 新建品牌清洗模块
- `backend/comment/brand_cleaner_test.go` - 单元测试
- `backend/task/executor.go` - 集成清洗逻辑

### Definition of Done
- [x] AI提示词明确品牌必须单一
- [x] 品牌清洗函数能正确拆分"/"分隔的品牌
- [x] 单元测试覆盖各种边界情况
- [x] 集成到评论分析流程中

### Must Have
- 品牌清洗：拆分"/"分隔的品牌，返回第一个或匹配已知品牌
- AI提示词优化：明确品牌单一、提供已知品牌参考
- 单元测试覆盖

### Must NOT Have (Guardrails)
- ❌ 不修复历史数据（只影响新任务）
- ❌ 不改变数据库结构
- ❌ 不修改前端代码
- ❌ 不实现复杂的型号提取逻辑（方案C）

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (Go testing)
- **User wants tests**: YES
- **Framework**: `go test`

### Test Coverage Requirements
品牌清洗模块必须包含以下测试用例：
1. 正常品牌名不变
2. "/"分隔的品牌返回第一个
3. 匹配已知品牌优先返回
4. 空字符串处理
5. 多个"/"的情况

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: 优化AI提示词 [no dependencies]
└── Task 2: 创建品牌清洗模块 [no dependencies]

Wave 2 (After Wave 1):
└── Task 3: 集成清洗逻辑到执行器 [depends: 1, 2]

Critical Path: Task 1 + Task 2 → Task 3
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 3 | 2 |
| 2 | None | 3 | 1 |
| 3 | 1, 2 | None | None (final) |

---

## TODOs

- [x] 1. 优化AI提示词，明确品牌必须单一

  **What to do**:
  - 修改 `backend/ai/analysis.go` 中的批量分析提示词（约第411-425行）
  - 添加以下规则到提示词：
    ```
    - 品牌必须是单一品牌名称，绝对不能包含"/"或其他分隔符
    - 如果评论对比多个品牌（如"A比B好"），只提取评论主要评价的那个品牌
    - 如果评论同时讨论多个品牌且无法确定主要品牌，brand填"未知"
    - 型号必须是具体型号名（如"V12"、"Max"、"Pro"），不能是描述性文字（如"新款"、"基础款"）
    ```
  - 同样修改单条分析提示词（约第93-103行）

  **Must NOT do**:
  - 不改变返回的JSON格式
  - 不修改评分逻辑

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 2)
  - **Blocks**: Task 3
  - **Blocked By**: None

  **References**:
  - `backend/ai/analysis.go:85-103` - 单条分析提示词
  - `backend/ai/analysis.go:411-425` - 批量分析提示词

  **Acceptance Criteria**:

  ```bash
  # AC1: 后端编译通过
  cd backend && go build ./...
  # Expected: 编译成功

  # AC2: 提示词包含新规则
  grep -n "不能包含" backend/ai/analysis.go
  # Expected: 找到相关规则
  ```

  **Commit**: YES
  - Message: `fix(ai): clarify brand must be single in analysis prompts`
  - Files: `backend/ai/analysis.go`

---

- [x] 2. 创建品牌清洗模块

  **What to do**:
  - 创建 `backend/comment/brand_cleaner.go`
  - 实现 `CleanBrandName(brand string, knownBrands []string) string` 函数
  - 实现 `CleanModelName(model string) string` 函数
  - 创建 `backend/comment/brand_cleaner_test.go` 单元测试

  **清洗逻辑**:
  ```go
  func CleanBrandName(brand string, knownBrands []string) string {
      brand = strings.TrimSpace(brand)
      if brand == "" || brand == "未知" {
          return brand
      }
      
      // 如果不包含"/"，直接返回
      if !strings.Contains(brand, "/") {
          return brand
      }
      
      // 拆分品牌
      parts := strings.Split(brand, "/")
      
      // 优先返回匹配已知品牌的部分
      for _, part := range parts {
          part = strings.TrimSpace(part)
          for _, known := range knownBrands {
              if strings.EqualFold(part, known) {
                  return known // 返回已知品牌的标准名称
              }
          }
      }
      
      // 都不匹配，返回第一个
      return strings.TrimSpace(parts[0])
  }
  
  func CleanModelName(model string) string {
      model = strings.TrimSpace(model)
      
      // 如果包含"/"，取第一个
      if strings.Contains(model, "/") {
          parts := strings.Split(model, "/")
          model = strings.TrimSpace(parts[0])
      }
      
      // 过滤掉描述性文字，保留"通用"
      descriptive := []string{"新款", "旧款", "基础款", "升级款", "标准版"}
      for _, d := range descriptive {
          if model == d {
              return "通用"
          }
      }
      
      return model
  }
  ```

  **Must NOT do**:
  - 不访问数据库
  - 不修改其他文件

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocks**: Task 3
  - **Blocked By**: None

  **References**:
  - `backend/comment/filter.go` - 同包下的现有模块，参考代码风格

  **Acceptance Criteria**:

  ```bash
  # AC1: 单元测试全部通过
  cd backend && go test ./comment/... -v -run TestClean
  # Expected: 所有测试通过

  # AC2: 测试覆盖率 >= 80%
  cd backend && go test ./comment/... -cover
  # Expected: coverage >= 80%
  ```

  **必须包含的测试用例**:
  ```go
  func TestCleanBrandName_Normal(t *testing.T)
  // Input: "小佩", [] → Output: "小佩"
  
  func TestCleanBrandName_WithSlash(t *testing.T)
  // Input: "喵洁易/Catlink", [] → Output: "喵洁易"
  
  func TestCleanBrandName_MatchKnown(t *testing.T)
  // Input: "喵洁易/Catlink", ["Catlink", "小佩"] → Output: "Catlink"
  
  func TestCleanBrandName_MultipleSlash(t *testing.T)
  // Input: "有陪/小佩/小米/糯雪", ["小米"] → Output: "小米"
  
  func TestCleanBrandName_Empty(t *testing.T)
  // Input: "", [] → Output: ""
  
  func TestCleanBrandName_Unknown(t *testing.T)
  // Input: "未知", [] → Output: "未知"
  
  func TestCleanModelName_Normal(t *testing.T)
  // Input: "V12" → Output: "V12"
  
  func TestCleanModelName_WithSlash(t *testing.T)
  // Input: "二代/通用" → Output: "二代"
  
  func TestCleanModelName_Descriptive(t *testing.T)
  // Input: "新款" → Output: "通用"
  ```

  **Commit**: YES
  - Message: `feat(comment): add brand and model name cleaner`
  - Files: `backend/comment/brand_cleaner.go`, `backend/comment/brand_cleaner_test.go`

---

- [x] 3. 集成清洗逻辑到执行器

  **What to do**:
  - 修改 `backend/task/executor.go`
  - 在 `analyzeComments` 函数中，AI分析返回后调用清洗函数
  - 传递已知品牌列表（来自任务配置的 `req.Brands`）

  **集成位置**（约在第469-470行之后）:
  ```go
  // AI分析返回后，清洗品牌和型号
  results[i].Brand = comment.CleanBrandName(r.Brand, req.Brands)
  results[i].Model = comment.CleanModelName(r.Model)
  ```

  **Must NOT do**:
  - 不修改AI调用逻辑
  - 不修改报告生成逻辑

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (final)
  - **Blocks**: None
  - **Blocked By**: Task 1, Task 2

  **References**:
  - `backend/task/executor.go:460-475` - AI分析结果处理位置
  - `backend/comment/brand_cleaner.go` - Task 2 创建的清洗模块

  **Acceptance Criteria**:

  ```bash
  # AC1: 后端编译通过
  cd backend && go build ./...
  # Expected: 编译成功

  # AC2: 验证集成
  grep -n "CleanBrandName\|CleanModelName" backend/task/executor.go
  # Expected: 找到调用
  ```

  **Commit**: YES
  - Message: `fix(executor): integrate brand/model cleaning after AI analysis`
  - Files: `backend/task/executor.go`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `fix(ai): clarify brand must be single in prompts` | analysis.go | `go build ./...` |
| 2 | `feat(comment): add brand/model cleaner` | brand_cleaner.go, brand_cleaner_test.go | `go test ./comment/...` |
| 3 | `fix(executor): integrate brand/model cleaning` | executor.go | `go build ./...` |

---

## Success Criteria

### Verification Commands
```bash
# 后端编译
cd backend && go build ./...
# Expected: 成功

# 单元测试
cd backend && go test ./comment/... -v
# Expected: 全部通过

# 验证提示词更新
grep -A5 "品牌必须" backend/ai/analysis.go
# Expected: 找到新规则
```

### Final Checklist
- [x] AI提示词明确品牌必须单一
- [x] 品牌清洗函数正确处理"/"分隔
- [x] 型号清洗函数过滤描述性文字
- [x] 单元测试覆盖率 >= 80%
- [x] 集成到执行器流程中
