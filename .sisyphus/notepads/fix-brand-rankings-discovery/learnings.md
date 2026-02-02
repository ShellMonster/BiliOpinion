## 修复品牌排名未包含AI发现品牌的问题

### 完成时间
2026-02-02

### 修改内容
1. **Task 1**: `generateRankings` 函数 - 遍历 `scores` map 而非 `brands` 参数
   - 文件: `backend/report/generator.go:318-320`
   - 添加注释说明新逻辑
   - 添加空品牌过滤

2. **Task 2**: `brandAnalysis` 调用 - 使用所有发现的品牌
   - 文件: `backend/report/generator.go:186-193`
   - 新增 `allBrands` 变量收集所有品牌

3. **Task 3**: `ReportData.Brands` - 返回所有品牌
   - 文件: `backend/report/generator.go:199-207`
   - 新增 `allBrandNames` 变量从排名中提取品牌

### 核心问题
`generateRankings` 只遍历 `input.Brands`（用户指定的品牌），忽略了 AI 在评论分析中发现的新品牌。

### 解决方案
改为遍历 `scores` map，其中包含所有有得分数据的品牌（包括AI发现的）。

### 验证
- `go build ./backend/...` 通过
- 代码审查确认修改正确

### 影响
- 新创建的报告将包含所有发现的品牌排名
- 现有报告（已保存到数据库的）不受影响
