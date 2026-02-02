# 报告可视化升级 (Report Visualization Upgrade)

## TL;DR

> **Quick Summary**: 升级报告页面的可视化能力，添加15个新功能，包括热力图、词云、交互筛选、导出等
> 
> **Deliverables**:
> - 重构后的 Report.tsx (拆分为20+独立组件)
> - 5个新图表组件 (热力图、词云、直方图、饼图、网络图)
> - 3个交互组件 (品牌卡片、筛选器、弹窗)
> - 2个导出功能 (图片、Excel)
> - 后端新增关键词和情感分布数据
> 
> **Estimated Effort**: XL (26个任务)
> **Parallel Execution**: YES - 4 waves
> **Critical Path**: 重构 → 依赖安装 → 后端扩展 → 前端实现 → 测试

---

## Context

### Original Request
用户要求对报告页面进行全面可视化升级，包含15个功能点。

### Interview Summary
**Key Discussions**:
- ECharts版本: 降级到v5.6.0 (echarts-wordcloud兼容性)
- 关键词数据: 后端报告生成时提取
- 情感分布: 基于评分阈值计算 (>=8好评, 5-8中性, <5差评)
- 品牌详情弹窗: 包含得分+典型评论
- 测试策略: Playwright自动化测试

**Research Findings**:
- Report.tsx 当前616行，是God Component，必须先重构
- 后端ReportData结构完善，但缺少keyword_frequency和sentiment_distribution
- ECharts v6.0.0 已安装，需降级到v5.6.0

### Metis Review
**Identified Gaps** (addressed):
- echarts-wordcloud不兼容v6: 降级到v5.6.0
- 后端缺少关键词数据: 新增KeywordFrequency字段
- 后端缺少情感分布: 新增SentimentDistribution字段

---

## Work Objectives

### Core Objective
升级报告页面，添加15个可视化和交互功能，提升用户体验和数据洞察能力。

### Concrete Deliverables
- 重构后的组件结构 (20+文件)
- 5个新图表: 热力图、词云、直方图、情感饼图、网络图
- 3个交互组件: 品牌卡片、维度筛选器、品牌详情弹窗
- 5个深度分析: 型号分析、竞品对比、决策树、AI总结增强
- 2个导出功能: 图片导出、Excel导出
- 后端数据扩展: keyword_frequency, sentiment_distribution

### Definition of Done
- [ ] 所有15个功能可用
- [ ] 前端构建成功: `cd frontend && npm run build` 退出码0
- [ ] 后端编译成功: `go build ./backend/...` 退出码0
- [ ] Playwright测试通过

### Must Have
- 热力图显示品牌×维度得分矩阵
- 词云显示前50个高频关键词
- 情感饼图显示好评/中性/差评占比
- 品牌卡片可点击查看详情
- 图片和Excel导出功能

### Must NOT Have (Guardrails)
- ❌ 价格区间参考
- ❌ 时间趋势图
- ❌ 用户画像分析
- ❌ 评论展开/折叠
- ❌ 移动端适配
- ❌ 超过200行的单个组件文件
- ❌ 调用额外AI API进行关键词提取

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: NO (需要设置Playwright)
- **User wants tests**: YES (自动化测试)
- **Framework**: Playwright

### Automated Verification

每个任务使用以下验证方式:

**前端组件**: Playwright浏览器自动化
**后端API**: curl命令验证JSON响应
**构建验证**: npm run build / go build

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: 创建类型定义文件
├── Task 2: 创建数据获取Hook
└── Task 3: 安装新依赖 (降级ECharts, 安装wordcloud/xlsx/html2canvas)

Wave 2 (After Wave 1):
├── Task 4: 拆分Report.tsx - 提取Header组件
├── Task 5: 拆分Report.tsx - 提取统计卡片
├── Task 6: 拆分Report.tsx - 提取图表组件
├── Task 7: 后端新增KeywordFrequency字段
└── Task 8: 后端新增SentimentDistribution字段

Wave 3 (After Wave 2):
├── Task 9: 实现热力图 BrandHeatmap.tsx
├── Task 10: 实现词云 KeywordCloud.tsx
├── Task 11: 实现情感饼图 SentimentPie.tsx
├── Task 12: 实现品牌卡片 BrandCard.tsx
├── Task 13: 实现维度筛选器 DimensionFilter.tsx
└── Task 14: 实现品牌详情弹窗 BrandDetailModal.tsx

Wave 4 (After Wave 3):
├── Task 15: 实现雷达图选择器 RadarBrandSelector.tsx
├── Task 16: 实现评分直方图 ScoreHistogram.tsx
├── Task 17: 实现型号分析 ModelAnalysis.tsx
├── Task 18: 实现竞品对比 CompetitorCompare.tsx
├── Task 19: 实现品牌网络图 BrandNetwork.tsx
├── Task 20: 实现决策树 DecisionTree.tsx
├── Task 21: 实现AI总结增强 EnhancedSummary.tsx
├── Task 22: 实现图片导出
└── Task 23: 实现Excel导出

Wave 5 (Final):
├── Task 24: 集成所有组件到Report.tsx
├── Task 25: 设置Playwright测试
└── Task 26: 运行完整测试验证
```

### Critical Path
Task 1 → Task 4-6 → Task 9-14 → Task 24 → Task 26

---

## TODOs

### Phase 0: 基础准备

- [x] 1. 创建类型定义文件 ✓

  **What to do**:
  - 创建 `frontend/src/types/report.ts`
  - 从 Report.tsx 提取所有 interface 定义
  - 新增 KeywordItem, SentimentStats 接口

  **Must NOT do**:
  - 修改现有接口的必填字段为可选

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3)
  - **Blocks**: Tasks 4-6
  - **Blocked By**: None

  **References**:
  - `frontend/src/pages/Report.tsx:1-80` - 现有接口定义
  - `backend/report/generator.go:20-100` - Go结构体定义

  **Acceptance Criteria**:
  ```bash
  # 验证文件存在且包含关键接口
  grep -q "interface ReportData" frontend/src/types/report.ts && echo "PASS" || echo "FAIL"
  grep -q "interface KeywordItem" frontend/src/types/report.ts && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `refactor(frontend): extract report types to dedicated file`
  - Files: `frontend/src/types/report.ts`

---

- [x] 2. 创建数据获取Hook ✓

  **What to do**:
  - 创建 `frontend/src/hooks/useReportData.ts`
  - 从 Report.tsx 提取数据获取逻辑
  - 使用 useState, useEffect 管理状态

  **Must NOT do**:
  - 改变API调用逻辑

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3)
  - **Blocks**: Tasks 4-6
  - **Blocked By**: None

  **References**:
  - `frontend/src/pages/Report.tsx:100-150` - 现有数据获取逻辑
  - `frontend/src/api/client.ts` - API客户端

  **Acceptance Criteria**:
  ```bash
  grep -q "export function useReportData" frontend/src/hooks/useReportData.ts && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `refactor(frontend): extract report data fetching to custom hook`
  - Files: `frontend/src/hooks/useReportData.ts`

---

- [x] 3. 安装新依赖并降级ECharts ✓

  **What to do**:
  - 降级 echarts 到 5.6.0
  - 安装 echarts-wordcloud@2.1.0
  - 安装 xlsx
  - 安装 html2canvas
  - 验证现有图表正常工作

  **Must NOT do**:
  - 保持 ECharts v6 (不兼容wordcloud)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2)
  - **Blocks**: Tasks 9-14
  - **Blocked By**: None

  **References**:
  - `frontend/package.json` - 当前依赖

  **Acceptance Criteria**:
  ```bash
  cd frontend && npm install echarts@5.6.0 echarts-wordcloud@2.1.0 xlsx html2canvas
  # 验证安装
  grep '"echarts": "5.6.0"' frontend/package.json && echo "PASS" || echo "FAIL"
  grep '"echarts-wordcloud"' frontend/package.json && echo "PASS" || echo "FAIL"
  # 验证构建
  cd frontend && npm run build
  ```

  **Commit**: YES
  - Message: `chore(frontend): downgrade echarts to v5.6.0 and add visualization deps`
  - Files: `frontend/package.json`, `frontend/package-lock.json`

---

### Phase 1: 重构Report.tsx

- [ ] 4. 拆分Report.tsx - 提取Header组件

  **What to do**:
  - 创建 `frontend/src/components/Report/ReportHeader.tsx`
  - 提取标题、导出按钮、返回按钮
  - 在Report.tsx中引用新组件

  **Must NOT do**:
  - 改变现有功能

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 5, 6, 7, 8)
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **References**:
  - `frontend/src/pages/Report.tsx:200-250` - Header部分

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/ReportHeader.tsx && echo "PASS" || echo "FAIL"
  cd frontend && npm run build
  ```

  **Commit**: YES
  - Message: `refactor(frontend): extract ReportHeader component`
  - Files: `frontend/src/components/Report/ReportHeader.tsx`, `frontend/src/pages/Report.tsx`

---

- [ ] 5. 拆分Report.tsx - 提取统计卡片

  **What to do**:
  - 创建 `frontend/src/components/Report/Overview/KeyStatsCards.tsx`
  - 提取视频数、评论数、品牌数统计卡片

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **References**:
  - `frontend/src/pages/Report.tsx:300-350` - 统计卡片部分

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Overview/KeyStatsCards.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `refactor(frontend): extract KeyStatsCards component`

---

- [ ] 6. 拆分Report.tsx - 提取图表组件

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/BrandRadarChart.tsx`
  - 创建 `frontend/src/components/Report/Charts/BrandScoreChart.tsx`
  - 提取现有雷达图和柱状图

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: Tasks 9-14
  - **Blocked By**: Tasks 1, 2, 3

  **References**:
  - `frontend/src/pages/Report.tsx:400-500` - 图表配置

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/BrandRadarChart.tsx && echo "PASS" || echo "FAIL"
  test -f frontend/src/components/Report/Charts/BrandScoreChart.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `refactor(frontend): extract chart components`

---

### Phase 2: 后端扩展

- [ ] 7. 后端新增KeywordFrequency字段

  **What to do**:
  - 在 `backend/report/generator.go` 中新增 KeywordItem 结构体
  - 在 ReportData 中新增 KeywordFrequency 字段
  - 实现从评论中提取关键词的逻辑 (简单词频统计)
  - 过滤停用词和单字词

  **Must NOT do**:
  - 调用额外的AI API
  - 破坏现有API响应格式

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4-6, 8)
  - **Blocks**: Task 10 (词云)
  - **Blocked By**: None

  **References**:
  - `backend/report/generator.go:20-50` - ReportData结构体
  - `backend/report/generator.go:200-300` - 报告生成逻辑

  **Acceptance Criteria**:
  ```bash
  # 验证编译
  cd /Users/daozhang/Trae_AI/B站商品评论解析 && go build ./backend/...
  # 验证结构体存在
  grep -q "KeywordFrequency" backend/report/generator.go && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(backend): add keyword frequency extraction to report`
  - Files: `backend/report/generator.go`

---

- [ ] 8. 后端新增SentimentDistribution字段

  **What to do**:
  - 在 ReportData 中新增 SentimentDistribution 字段
  - 实现情感分布计算: >=8好评, 5-8中性, <5差评
  - 返回各类别的数量和百分比

  **Must NOT do**:
  - 使用AI进行情感分析

  **Recommended Agent Profile**:
  - **Category**: `business-logic`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 11 (情感饼图)
  - **Blocked By**: None

  **References**:
  - `backend/report/generator.go` - 报告生成逻辑

  **Acceptance Criteria**:
  ```bash
  go build ./backend/...
  grep -q "SentimentDistribution" backend/report/generator.go && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(backend): add sentiment distribution to report`
  - Files: `backend/report/generator.go`

---

### Phase 3: 新图表组件

- [ ] 9. 实现热力图 BrandHeatmap.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/BrandHeatmap.tsx`
  - 使用ECharts heatmap显示品牌×维度得分矩阵
  - 颜色编码: 绿色(>=8), 橙色(6-8), 红色(<6)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 3, 6

  **References**:
  - `frontend/src/components/Report/Charts/BrandRadarChart.tsx` - ECharts使用模式
  - ECharts heatmap文档

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/BrandHeatmap.tsx && echo "PASS" || echo "FAIL"
  cd frontend && npm run build
  ```

  **Commit**: YES
  - Message: `feat(frontend): add brand heatmap visualization`

---

- [ ] 10. 实现词云 KeywordCloud.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/KeywordCloud.tsx`
  - 使用echarts-wordcloud显示前50个高频关键词
  - 词大小与频率成正比

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 3, 7

  **References**:
  - echarts-wordcloud文档

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/KeywordCloud.tsx && echo "PASS" || echo "FAIL"
  grep -q "wordCloud" frontend/src/components/Report/Charts/KeywordCloud.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add keyword word cloud visualization`

---

- [ ] 11. 实现情感饼图 SentimentPie.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/SentimentPie.tsx`
  - 使用ECharts pie显示好评/中性/差评占比
  - 颜色: 绿色(好评), 灰色(中性), 红色(差评)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 3, 8

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/SentimentPie.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add sentiment distribution pie chart`

---

- [ ] 12. 实现品牌卡片 BrandCard.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/BrandCard.tsx`
  - 卡片显示: 品牌名、综合得分、排名、优劣势标签
  - 点击卡片打开详情弹窗

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/BrandCard.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add brand card component`

---

- [ ] 13. 实现维度筛选器 DimensionFilter.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/DimensionFilter.tsx`
  - 多选checkbox允许用户选择关注的维度
  - 选择后重新计算排名

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/DimensionFilter.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add dimension filter component`

---

- [ ] 14. 实现品牌详情弹窗 BrandDetailModal.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/BrandDetailModal.tsx`
  - 显示: 品牌详细得分、各维度分数、典型好评/差评
  - 使用现有Modal组件模式

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **References**:
  - `frontend/src/components/common/Modal.tsx` - 弹窗模式 (如果存在)

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/BrandDetailModal.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add brand detail modal`

---

### Phase 4: 高级功能

- [ ] 15. 实现雷达图选择器 RadarBrandSelector.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/RadarBrandSelector.tsx`
  - 允许用户选择2-4个品牌进行雷达图对比
  - 使用checkbox或dropdown选择

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 6

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/RadarBrandSelector.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add radar chart brand selector`

---

- [ ] 16. 实现评分直方图 ScoreHistogram.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/ScoreHistogram.tsx`
  - 显示各品牌的评分分布 (好评/中评/差评数量)
  - 使用ECharts bar chart

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 8

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/ScoreHistogram.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add score distribution histogram`

---

- [ ] 17. 实现型号分析 ModelAnalysis.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/ModelAnalysis.tsx`
  - 显示型号级别的详细分析
  - 包含型号雷达图对比

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 6

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/ModelAnalysis.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add model-level analysis component`

---

- [ ] 18. 实现竞品对比 CompetitorCompare.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/CompetitorCompare.tsx`
  - 允许选择2个品牌进行1v1深度对比
  - 并排显示各维度得分

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/CompetitorCompare.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add competitor comparison component`

---

- [ ] 19. 实现品牌网络图 BrandNetwork.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/Charts/BrandNetwork.tsx`
  - 使用ECharts graph显示品牌-型号关系
  - 节点大小与评论数成正比

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 3

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/Charts/BrandNetwork.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add brand network graph`

---

- [ ] 20. 实现决策树 DecisionTree.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/DecisionTree.tsx`
  - 交互式问答推荐品牌
  - 基于用户选择的维度优先级

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/DecisionTree.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add purchase decision tree`

---

- [ ] 21. 实现AI总结增强 EnhancedSummary.tsx

  **What to do**:
  - 创建 `frontend/src/components/Report/EnhancedSummary.tsx`
  - 显示分场景推荐、避坑指南
  - 使用react-markdown渲染

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  ```bash
  test -f frontend/src/components/Report/EnhancedSummary.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add enhanced AI summary component`

---

### Phase 5: 导出功能

- [ ] 22. 实现图片导出

  **What to do**:
  - 在ReportHeader中添加"导出图片"按钮
  - 使用html2canvas截取报告区域
  - 下载为PNG文件

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 3

  **Acceptance Criteria**:
  ```bash
  grep -q "html2canvas" frontend/src/components/Report/ReportHeader.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add image export functionality`

---

- [ ] 23. 实现Excel导出

  **What to do**:
  - 在ReportHeader中添加"导出Excel"按钮
  - 使用xlsx库生成Excel文件
  - 包含品牌排名、型号排名、维度得分

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 24
  - **Blocked By**: Task 3

  **Acceptance Criteria**:
  ```bash
  grep -q "xlsx" frontend/src/components/Report/ReportHeader.tsx && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `feat(frontend): add Excel export functionality`

---

### Phase 6: 集成与测试

- [ ] 24. 集成所有组件到Report.tsx

  **What to do**:
  - 更新 `frontend/src/pages/Report.tsx`
  - 引入所有新组件
  - 确保Report.tsx不超过150行
  - 添加Tab或Section切换不同视图

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (Sequential)
  - **Blocks**: Tasks 25, 26
  - **Blocked By**: Tasks 4-23

  **Acceptance Criteria**:
  ```bash
  # 验证Report.tsx行数
  wc -l < frontend/src/pages/Report.tsx | awk '{if($1<=200) print "PASS"; else print "FAIL"}'
  # 验证构建
  cd frontend && npm run build
  ```

  **Commit**: YES
  - Message: `feat(frontend): integrate all report components`

---

- [ ] 25. 设置Playwright测试

  **What to do**:
  - 安装 @playwright/test
  - 创建 `frontend/tests/report.spec.ts`
  - 编写基本测试: 页面加载、图表渲染、导出功能

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`playwright`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (Sequential)
  - **Blocks**: Task 26
  - **Blocked By**: Task 24

  **Acceptance Criteria**:
  ```bash
  test -f frontend/tests/report.spec.ts && echo "PASS" || echo "FAIL"
  ```

  **Commit**: YES
  - Message: `test(frontend): add Playwright tests for report page`

---

- [ ] 26. 运行完整测试验证

  **What to do**:
  - 启动后端服务
  - 启动前端服务
  - 运行Playwright测试
  - 验证所有功能正常

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`playwright`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (Final)
  - **Blocks**: None
  - **Blocked By**: Task 25

  **Acceptance Criteria**:
  ```bash
  cd frontend && npx playwright test
  # Assert: All tests pass
  ```

  **Commit**: NO (测试任务)

---

## Commit Strategy

| After Task | Message | Files |
|------------|---------|-------|
| 1 | `refactor(frontend): extract report types` | types/report.ts |
| 2 | `refactor(frontend): extract data hook` | hooks/useReportData.ts |
| 3 | `chore(frontend): update dependencies` | package.json |
| 4-6 | `refactor(frontend): split Report.tsx` | components/Report/* |
| 7-8 | `feat(backend): add keyword and sentiment data` | generator.go |
| 9-14 | `feat(frontend): add visualization components` | components/Report/Charts/* |
| 15-23 | `feat(frontend): add advanced features` | components/Report/* |
| 24 | `feat(frontend): integrate all components` | pages/Report.tsx |
| 25 | `test(frontend): add Playwright tests` | tests/* |

---

## Success Criteria

### Verification Commands
```bash
# 后端编译
go build ./backend/...

# 前端构建
cd frontend && npm run build

# Playwright测试
cd frontend && npx playwright test
```

### Final Checklist
- [ ] 所有15个功能可用
- [ ] Report.tsx 不超过200行
- [ ] 后端API包含keyword_frequency和sentiment_distribution
- [ ] 图片和Excel导出功能正常
- [ ] Playwright测试通过
