# 设置页面改为弹窗形式

## TL;DR

> **Quick Summary**: 将设置从独立页面改为弹窗形式，点击导航栏"设置"按钮弹出设置弹窗
> 
> **Deliverables**:
> - 新建 Modal 通用组件
> - 新建 SettingsModal 组件
> - 修改 Layout 添加弹窗状态和触发按钮
> - 移除 /settings 路由
> 
> **Estimated Effort**: Quick (30-45分钟)
> **Parallel Execution**: YES - 2 waves
> **Critical Path**: Modal组件 → SettingsModal → Layout集成

---

## Context

### 用户需求
将设置页面从独立路由页面改为弹窗形式，点击导航栏"设置"时弹出弹窗而非跳转页面。

### 当前实现
- `frontend/src/pages/Settings.tsx` - 独立页面，路由 `/settings`
- 使用 localStorage 存储配置
- 项目缺少 Modal 组件

### 参考实现
- `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/common/Modal.tsx` - 通用弹窗组件
- 使用 createPortal 渲染到 body
- 支持毛玻璃效果和动画

---

## Work Objectives

### Core Objective
将设置从页面改为弹窗，保持现有功能不变

### Concrete Deliverables
- `frontend/src/components/common/Modal.tsx` - 通用弹窗组件
- `frontend/src/components/Settings/SettingsModal.tsx` - 设置弹窗
- 修改后的 `frontend/src/components/Layout/Layout.tsx` - 集成弹窗
- 修改后的 `frontend/src/App.tsx` - 移除 /settings 路由

### Definition of Done
- [ ] Modal 组件可正常打开/关闭
- [ ] 设置弹窗显示所有配置项
- [ ] 点击导航栏"设置"弹出弹窗
- [ ] 保存功能正常工作
- [ ] /settings 路由已移除

### Must NOT Have (Guardrails)
- 不改变设置的存储方式（保持 localStorage）
- 不添加复杂的状态管理（如 Zustand）
- 不添加 i18n 国际化
- 不添加 Tabs 布局（保持简单）

---

## Status: ✅ COMPLETED (2026-02-01)

All tasks have been completed. This plan is now archived.

---

## TODOs

- [x] 1. 创建 Modal 通用组件

  **What to do**:
  - 创建 `frontend/src/components/common/Modal.tsx`
  - 使用 createPortal 渲染到 body
  - 支持 isOpen, onClose, title, children props
  - 实现毛玻璃背景和动画效果
  - 点击背景关闭弹窗
  - 打开时禁止背景滚动

  **Must NOT do**:
  - 不添加复杂的 density/scrollable 选项（简化版）

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 2
  - **Blocked By**: None

  **References**:
  - 参考实现: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/common/Modal.tsx`

  **Acceptance Criteria**:
  ```bash
  cd frontend && npm run build
  # 期望: 编译成功，无类型错误
  ```

  **Commit**: YES
  - Message: `feat(frontend): add Modal component`
  - Files: `frontend/src/components/common/Modal.tsx`

---

- [x] 2. 创建 SettingsModal 组件

  **What to do**:
  - 创建 `frontend/src/components/Settings/SettingsModal.tsx`
  - 复用现有 Settings.tsx 的逻辑和 UI
  - 接收 isOpen 和 onClose props
  - 使用 Modal 组件包裹内容
  - 保存成功后自动关闭弹窗

  **UI结构**:
  ```tsx
  <Modal isOpen={isOpen} onClose={onClose} title="系统设置">
    {/* AI配置区域 */}
    <div className="space-y-4 mb-8">
      <h3>AI配置</h3>
      <Input label="API Base URL" ... />
      <Input label="API Key" ... />
      <Input label="Model" ... />
    </div>
    
    {/* B站Cookie区域 */}
    <div className="space-y-4 mb-8">
      <h3>B站Cookie</h3>
      <textarea ... />
    </div>
    
    {/* 保存按钮 */}
    <Button onClick={handleSave}>保存设置</Button>
  </Modal>
  ```

  **Must NOT do**:
  - 不添加 Tabs 布局
  - 不改变存储逻辑

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

  **References**:
  - 现有设置页面: `frontend/src/pages/Settings.tsx`
  - Modal组件: `frontend/src/components/common/Modal.tsx`

  **Acceptance Criteria**:
  ```bash
  cd frontend && npm run build
  # 期望: 编译成功，无类型错误
  ```

  **Commit**: YES
  - Message: `feat(frontend): add SettingsModal component`
  - Files: `frontend/src/components/Settings/SettingsModal.tsx`

---

- [x] 3. 修改 Layout 集成弹窗

  **What to do**:
  - 修改 `frontend/src/components/Layout/Layout.tsx`
  - 添加 isSettingsOpen 状态
  - 将"设置"链接改为按钮，点击打开弹窗
  - 在 Layout 中渲染 SettingsModal

  **代码结构**:
  ```tsx
  export default function Layout({ children }: LayoutProps) {
    const [isSettingsOpen, setIsSettingsOpen] = useState(false)
    
    return (
      <div className="min-h-screen bg-[#f8fafc]">
        <header>
          <nav>
            <Link to="/">首页</Link>
            <Link to="/history">历史记录</Link>
            <button onClick={() => setIsSettingsOpen(true)}>设置</button>
          </nav>
        </header>
        
        <main>{children}</main>
        
        <SettingsModal 
          isOpen={isSettingsOpen} 
          onClose={() => setIsSettingsOpen(false)} 
        />
      </div>
    )
  }
  ```

  **Must NOT do**:
  - 不改变其他导航链接
  - 不改变整体布局结构

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 4
  - **Blocked By**: Task 2

  **References**:
  - 当前 Layout: `frontend/src/components/Layout/Layout.tsx`

  **Acceptance Criteria**:
  ```bash
  cd frontend && npm run build
  # 期望: 编译成功
  ```

  **Commit**: YES
  - Message: `feat(frontend): integrate SettingsModal into Layout`
  - Files: `frontend/src/components/Layout/Layout.tsx`

---

- [x] 4. 移除 /settings 路由和页面

  **What to do**:
  - 修改 `frontend/src/App.tsx` 移除 /settings 路由
  - 删除 `frontend/src/pages/Settings.tsx` 文件
  - 移除 Settings 的 import

  **Must NOT do**:
  - 不改变其他路由

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: None
  - **Blocked By**: Task 3

  **References**:
  - 当前 App.tsx: `frontend/src/App.tsx`

  **Acceptance Criteria**:
  ```bash
  cd frontend && npm run build
  # 期望: 编译成功，无未使用的 import 警告
  ```

  **Commit**: YES
  - Message: `refactor(frontend): remove Settings page route`
  - Files: `frontend/src/App.tsx`

---

## Commit Strategy

| After Task | Message | Files |
|------------|---------|-------|
| 1 | `feat(frontend): add Modal component` | frontend/src/components/common/Modal.tsx |
| 2 | `feat(frontend): add SettingsModal component` | frontend/src/components/Settings/SettingsModal.tsx |
| 3 | `feat(frontend): integrate SettingsModal into Layout` | frontend/src/components/Layout/Layout.tsx |
| 4 | `refactor(frontend): remove Settings page route` | frontend/src/App.tsx |

---

## Success Criteria

### Verification Commands
```bash
# 前端构建
cd frontend && npm run build
# Expected: 编译成功

# 启动前端
cd frontend && npm run dev
# Expected: Vite dev server on :5173
```

### Final Checklist
- [x] 点击导航栏"设置"弹出弹窗
- [x] 弹窗显示 AI 配置和 B站 Cookie 配置
- [x] 保存功能正常工作
- [x] 点击背景或关闭按钮可关闭弹窗
- [x] /settings 路由已移除
- [x] 前端编译无错误
