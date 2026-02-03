# 移动导出按钮到标签栏

## TL;DR

> **Quick Summary**: 将导出图片、导出Excel、导出PDF三个按钮从页面顶部移到标签栏（深度总结右侧），并删除查看历史按钮
> 
> **Deliverables**: 
> - 修改后的 ReportHeader.tsx（简化，只保留标题）
> - 修改后的 Report.tsx（标签栏包含导出按钮，右对齐）
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: NO - sequential
> **Critical Path**: Task 1 → Task 2

---

## Context

### Original Request
用户要求：
1. 把"导出图片"、"导出Excel"、"导出PDF"三个按钮移到"深度总结"标签页的右边，同一行，向右对齐
2. 删除"查看历史"按钮

### 当前代码结构
- `ReportHeader.tsx`: 包含标题和4个按钮（导出图片、导出Excel、导出PDF、查看历史）
- `Report.tsx`: 包含标签栏（总览、图表、深度总结）

---

## Work Objectives

### Core Objective
调整报告页面的按钮布局，将导出按钮移至标签栏右侧

### Concrete Deliverables
- `frontend/src/components/Report/ReportHeader.tsx` - 简化为只显示标题
- `frontend/src/pages/Report.tsx` - 标签栏增加导出按钮

### Definition of Done
- [x] 导出图片、导出Excel、导出PDF按钮在"深度总结"标签右侧
- [x] 三个导出按钮向右对齐
- [x] 查看历史按钮已删除
- [x] 所有导出功能正常工作

### Must Have
- 三个导出按钮在标签栏同一行
- 按钮向右对齐（使用 flex + ml-auto 或 justify-between）
- 保留原有的导出功能代码

### Must NOT Have (Guardrails)
- 不要删除导出功能的实现代码
- 不要改变按钮的样式和颜色
- 不要修改其他页面或组件

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: NO
- **User wants tests**: Manual-only
- **Framework**: none

### Manual Verification
启动前端开发服务器后，在浏览器中验证：
1. 访问报告页面
2. 确认标签栏显示：总览 | 图表 | 深度总结 | [导出图片] [导出Excel] [导出PDF]
3. 确认导出按钮在右侧对齐
4. 确认"查看历史"按钮已不存在
5. 点击各导出按钮确认功能正常

---

## TODOs

- [x] 1. 修改 Report.tsx - 将导出按钮移至标签栏

  **What to do**:
  1. 将导出功能代码（exportImage, exportExcel 函数）从 ReportHeader.tsx 移到 Report.tsx
  2. 添加必要的 state（imageExporting, excelExporting）
  3. 添加必要的 import（html2canvas, xlsx）
  4. 修改标签栏的 JSX 结构：
     - 使用 `flex justify-between items-center` 包裹整行
     - 左侧放标签按钮
     - 右侧放三个导出按钮（使用 `flex gap-2`）
  5. 删除 ReportHeader 组件中传递的 onExport, onBack, exporting, reportData props

  **Must NOT do**:
  - 不要改变导出按钮的样式
  - 不要改变标签按钮的样式

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **References**:
  - `frontend/src/pages/Report.tsx:66-68` - 当前标签栏代码
  - `frontend/src/components/Report/ReportHeader.tsx:29-122` - 导出功能实现代码
  - `frontend/src/components/Report/ReportHeader.tsx:131-203` - 导出按钮 JSX

  **Acceptance Criteria**:
  - [ ] Report.tsx 包含 exportImage 和 exportExcel 函数
  - [ ] 标签栏使用 flex justify-between 布局
  - [ ] 三个导出按钮在标签栏右侧
  - [ ] 运行 `cd frontend && npm run build` 无错误

  **Commit**: YES
  - Message: `feat(report): move export buttons to tab bar`
  - Files: `frontend/src/pages/Report.tsx`

- [x] 2. 简化 ReportHeader.tsx

  **What to do**:
  1. 删除导出相关的 state 和函数
  2. 删除不再需要的 props（exporting, onExport, onBack, reportData）
  3. 删除不再需要的 import（useState, html2canvas, xlsx, useToast）
  4. 简化组件为只显示标题和副标题

  **Must NOT do**:
  - 不要删除整个组件文件
  - 不要改变标题的样式

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **References**:
  - `frontend/src/components/Report/ReportHeader.tsx` - 当前完整代码
  - `frontend/src/pages/Report.tsx:64` - ReportHeader 的调用位置

  **Acceptance Criteria**:
  - [ ] ReportHeader 只接收 category 和 reportId props
  - [ ] ReportHeader 只渲染标题和副标题
  - [ ] 运行 `cd frontend && npm run build` 无错误

  **Commit**: YES
  - Message: `refactor(report): simplify ReportHeader component`
  - Files: `frontend/src/components/Report/ReportHeader.tsx`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `feat(report): move export buttons to tab bar` | Report.tsx | npm run build |
| 2 | `refactor(report): simplify ReportHeader component` | ReportHeader.tsx | npm run build |

---

## Success Criteria

### Verification Commands
```bash
cd frontend && npm run build  # Expected: Build successful
```

### Final Checklist
- [ ] 导出按钮在"深度总结"右侧
- [ ] 按钮向右对齐
- [ ] "查看历史"按钮已删除
- [ ] 构建无错误
