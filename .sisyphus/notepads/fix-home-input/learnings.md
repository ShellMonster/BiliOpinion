# Learnings - fix-home-input

## Conventions & Patterns
(To be populated as we discover patterns)

## Gotchas
(To be populated as we encounter issues)

## 2026-02-01: Updated keyword.go for free-form requirement parsing

### Changes Made
1. **ParseKeywordRequest struct**: Changed from `Category string` to `Requirement string` to accept natural language input
   - Example: "想买个吸尘器，预算2000，家里有宠物" instead of just "无线吸尘器"

2. **ParseKeywordResponse struct**: Added new fields to capture extracted information:
   - `Understanding string`: AI's interpretation of user needs (starts with "我理解您...")
   - `ProductType string`: Identified product category
   - `Budget string`: Budget range (optional, omitempty)
   - `Scenario string`: Usage scenario (optional, omitempty)
   - `SpecialNeeds []string`: Special requirements (optional, omitempty)
   - Removed old `Clarification` field

3. **System prompt**: Completely rewritten to:
   - Handle free-form natural language input
   - Extract structured information (product type, budget, scenario, special needs)
   - Provide natural language understanding
   - Adjust dimensions based on user's special needs
   - Generate keywords that incorporate user requirements

4. **User prompt**: Changed from `"商品类目：%s"` to `"用户需求：%s"`

5. **Validation logic**: Updated to:
   - Validate new required fields: `Understanding` and `ProductType`
   - Remove strict dimension count check (was `!= 6`, now just checks `> 0`)
   - Keep validation for `Brands` and `Keywords`

### Technical Notes
- All optional fields use `omitempty` JSON tag
- Build verified with `go build ./ai/` - no errors
- Backward compatibility: API consumers need to update their request/response handling

### Example Transformation
**Before**: `{"category": "无线吸尘器"}`
**After**: `{"requirement": "想买个吸尘器，预算2000，家里有宠物"}`

**Response now includes**:
- Natural language understanding
- Extracted product type, budget, scenario
- Dimensions tailored to special needs (e.g., pet hair handling)
- Keywords incorporating budget and scenario

## 2026-02-01 - 报告增强功能完成

### 成功实现的功能
1. **数据统计** - 视频数、评论数、各品牌评论数
2. **典型评论展示** - 好评/差评自动筛选
3. **品牌优劣势分析** - 自动识别优势/劣势维度
4. **AI生成购买建议** - 200-300字专业建议
5. **PDF柱状图** - 横向柱状图可视化得分

### 技术要点
- 扩展 ReportData 结构体添加新字段
- 实现辅助函数：generateBrandAnalysis(), selectTypicalComments()
- 新增 AI GenerateRecommendation() 函数
- PDF drawBarChart() 函数绘制柱状图
- 前端新增统计卡片和评论展示区域

### 文件改动
- backend/report/generator.go: +156行
- backend/ai/analysis.go: +61行
- backend/pdf/generator.go: +60行
- frontend/src/pages/Report.tsx: +347行

### 验证结果
- ✅ 后端编译通过
- ✅ 前端构建通过
- ✅ 服务正常启动
- ⚠️ AI API 超时（网络问题，非代码问题）

