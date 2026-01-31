# B站商品评论分析工具 - 工作计划

## TL;DR

> **Quick Summary**: 构建一个B站商品评论分析工具，通过抓取B站视频评论，使用AI分析商品评价，生成可视化报告帮助用户做购买决策。
> 
> **Deliverables**:
> - Go后端API服务 (B站抓取 + AI分析 + 数据存储)
> - React前端Web应用 (输入 + 报告展示 + 历史记录)
> - SQLite数据库 (配置 + 历史 + 报告数据)
> 
> **Estimated Effort**: Large (2-3周)
> **Parallel Execution**: YES - 4 waves
> **Critical Path**: 项目初始化 → B站API集成 → AI分析 → 前端报告

---

## Context

### Original Request
用户希望构建一个工具，用于抓取B站视频评论区，分析商品的真实用户评价，生成可视化报告。

### Interview Summary
**Key Discussions**:
- 技术栈: React + Go + SQLite，先做Web版，后续Tauri打包
- AI服务: OpenAI兼容接口，用户配置API Key/Base/Model
- 评价维度: 6个动态维度，AI根据品类生成
- B站API: 用户提供完整文档，包含WBI签名算法
- 数据策略: 历史记录永久保存，临时数据3天清理
- MVP范围: 核心流程 + 历史记录，导出功能第二期

**Research Findings**:
- B站API文档完整，WBI签名Go代码可直接使用
- 评论API需要AV号，需实现BV→AV转换
- WBI密钥每日更新，建议缓存1小时

### Metis Review
**Identified Gaps** (addressed):
- Cookie获取方式: 用户手动粘贴完整Cookie (包含SESSDATA + buvid3)
- 抓取数量限制: 50视频/500评论每视频
- 评分计算方式: AI打分(1-10)取平均

### Momus Review (高精度审查通过)
**非阻塞提醒** (已纳入计划):
- Cookie需求: 搜索接口需要buvid3，用户需粘贴完整Cookie字符串
- 数据清理: 明确raw_comments为临时数据，reports为永久数据
- Git初始化: Task 1需包含git init

---

## Work Objectives

### Core Objective
构建完整的B站商品评论分析流程：用户输入 → AI解析并澄清 → 用户确认 → B站搜索抓取 → AI评论分析 → 可视化报告

### Concrete Deliverables
- `/backend/` - Go后端服务 (含SSE任务状态推送)
- `/frontend/` - React前端应用 (含确认页面和进度展示)
- `/data/` - SQLite数据库文件
- 完整的API接口文档

### Definition of Done
- [ ] 用户可输入商品类目，获取AI解析的关键词
- [ ] AI返回澄清描述、品牌列表、对比维度供用户确认
- [ ] 用户确认后系统开始抓取和分析
- [ ] 后端通过SSE推送任务进度状态
- [ ] 系统可搜索B站视频并抓取评论
- [ ] AI可分析评论并生成各维度得分
- [ ] 前端可展示雷达图、柱状图、对比表格
- [ ] 历史记录可保存和查看

### Must Have
- OpenAI兼容API配置界面
- B站Cookie配置 (完整Cookie字符串，包含SESSDATA + buvid3)
- WBI签名实现
- **AI澄清确认页面** (展示澄清描述、品牌列表、对比维度，用户确认后继续)
- **SSE任务状态推送** (基于HTTP，简单稳定，浏览器自动重连)
- 6维度评价分析
- 雷达图 + 柱状图 + 对比表格
- 历史记录存储和查看

### Must NOT Have (Guardrails)
- 不实现B站登录流程 (用户手动粘贴Cookie)
- 不实现PDF/图片导出 (第二期)
- 不实现Tauri桌面打包 (第二期)
- 不过度抽象 (保持代码直接可读)
- 不添加非必要的中间件

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: NO (新项目)
- **User wants tests**: 先实现后测试
- **Framework**: Go: go test / React: vitest (后补)

### Automated Verification

每个TODO包含可执行的验证步骤：

**For Backend API** (using curl):
- 发送请求，验证响应JSON结构
- 检查HTTP状态码

**For Frontend** (using playwright):
- 导航到页面，验证元素存在
- 填写表单，验证交互

---

## Execution Strategy

### Parallel Execution Waves

Wave 1 (Start Immediately):
├── Task 1: 项目初始化 (Go + React)
└── Task 2: 数据库设计与实现

Wave 2 (After Wave 1):
├── Task 3: B站API集成 (WBI签名)
├── Task 4: AI服务集成
└── Task 5: 前端基础框架

Wave 3 (After Wave 2):
├── Task 6: 搜索与评论抓取
├── Task 7: AI关键词解析
└── Task 8: AI评论分析

Wave 4 (After Wave 3):
├── Task 9: 报告生成逻辑
├── Task 10: 前端报告页面
└── Task 11: 历史记录功能

Wave 5 (Final):
└── Task 12: 集成测试与优化

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 3,4,5 | 2 |
| 2 | None | 6,7,8 | 1 |
| 3 | 1 | 6 | 4,5 |
| 4 | 1 | 7,8 | 3,5 |
| 5 | 1 | 10,11 | 3,4 |
| 6 | 2,3 | 9 | 7,8 |
| 7 | 2,4 | 9 | 6,8 |
| 8 | 2,4 | 9 | 6,7 |
| 9 | 6,7,8 | 10 | None |
| 10 | 5,9 | 12 | 11 |
| 11 | 5,2 | 12 | 10 |
| 12 | 10,11 | None | None |

---

## TODOs

- [x] 1. 项目初始化

  **What to do**:
  - 初始化Git仓库: `git init`
  - 创建Go模块: `go mod init bilibili-analyzer`
  - 创建React项目: `npm create vite@latest frontend -- --template react-ts`
  - 设置项目目录结构
  - 配置Go依赖: gin, gorm, sqlite
  - 配置React依赖: axios, recharts, tailwindcss

  **Must NOT do**:
  - 不添加不必要的依赖
  - 不过度设计目录结构

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - 简单的项目初始化任务

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 2)
  - **Blocks**: Tasks 3, 4, 5
  - **Blocked By**: None

  **References**:
  - Go官方文档: https://go.dev/doc/
  - Vite文档: https://vitejs.dev/guide/

  **Acceptance Criteria**:
  - [ ] `git status` 显示已初始化
  - [ ] `go build ./...` 成功
  - [ ] `cd frontend && npm run dev` 启动成功
  - [ ] 目录结构符合规划

  **Commit**: YES
  - Message: `feat: initialize project structure`
  - Files: `go.mod, go.sum, frontend/package.json`

---

- [x] 2. 数据库设计与实现

  **What to do**:
  - 设计SQLite表结构:
    - `settings`: AI配置、B站Cookie (完整Cookie字符串)
    - `analysis_history`: 分析历史记录 (永久保存)
    - `reports`: 报告数据JSON (永久保存)
    - `raw_comments`: 原始评论数据 (临时，3天清理)
  - 实现GORM模型
  - 实现数据库初始化和迁移
  - 实现启动时3天临时数据清理逻辑 (仅清理raw_comments表)

  **Must NOT do**:
  - 不使用复杂的ORM关系
  - 不过度规范化表结构
  - 不清理reports和analysis_history表

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocks**: Tasks 6, 7, 8
  - **Blocked By**: None

  **References**:
  - GORM文档: https://gorm.io/docs/

  **Acceptance Criteria**:
  - [ ] 数据库文件创建成功
  - [ ] 4个表结构正确创建 (settings, analysis_history, reports, raw_comments)
  - [ ] CRUD操作测试通过
  - [ ] 清理逻辑仅删除raw_comments中超过3天的数据

  **Commit**: YES
  - Message: `feat: implement database models and migrations`
  - Files: `backend/models/*.go, backend/database/*.go`

---

- [ ] 3. B站API集成 - WBI签名

  **What to do**:
  - 从用户文档复制WBI签名Go代码
  - 实现WBI密钥获取和缓存(1小时)
  - 实现BV号转AV号函数
  - 封装HTTP客户端(带Cookie和签名)

  **Must NOT do**:
  - 不修改已验证的签名算法
  - 不忽略错误处理

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4, 5)
  - **Blocks**: Task 6
  - **Blocked By**: Task 1

  **References**:
  - WBI签名代码: `/Users/daozhang/Downloads/bilibili-API-collect/docs/misc/sign/wbi.md:388-566`
  - BV转AV代码: `/Users/daozhang/Downloads/bilibili-API-collect/docs/misc/bvid_desc.md:324-382`

  **Acceptance Criteria**:
  - [ ] WBI签名生成正确
  - [ ] BV→AV转换正确
  - [ ] 密钥缓存生效

  **Commit**: YES
  - Message: `feat: implement bilibili WBI signing and BV/AV conversion`
  - Files: `backend/bilibili/wbi.go, backend/bilibili/bvid.go`

---

- [ ] 4. AI服务集成

  **What to do**:
  - 实现OpenAI兼容客户端
  - 支持自定义API Base URL
  - 实现请求重试逻辑(1次)
  - 实现并发请求控制

  **Must NOT do**:
  - 不硬编码API配置
  - 不忽略API错误

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 5)
  - **Blocks**: Tasks 7, 8
  - **Blocked By**: Task 1

  **References**:
  - OpenAI API文档: https://platform.openai.com/docs/api-reference

  **Acceptance Criteria**:
  - [ ] 可配置API Base/Key/Model
  - [ ] 请求成功返回响应
  - [ ] 重试逻辑生效

  **Commit**: YES
  - Message: `feat: implement OpenAI-compatible AI client`
  - Files: `backend/ai/client.go`

---

- [ ] 5. 前端基础框架

  **What to do**:
  - 配置Tailwind CSS
  - 创建基础布局组件
  - 创建路由结构(首页/确认页/进度页/报告页/历史/设置)
  - 实现设置页面(AI配置/Cookie配置)
  - 实现API请求封装
  - **实现SSE连接封装** (使用EventSource接收任务状态推送)

  **Must NOT do**:
  - 不过度组件化
  - 不添加复杂状态管理

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 4)
  - **Blocks**: Tasks 10, 11
  - **Blocked By**: Task 1

  **References**:
  - Tailwind CSS: https://tailwindcss.com/docs
  - React Router: https://reactrouter.com/
  - MDN EventSource: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
  - **UI风格参考项目**: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/`

  **UI风格规范 (MUST FOLLOW)**:
  
  参考项目风格特点：现代毛玻璃效果、大圆角、渐变按钮、柔和阴影
  
  **背景色**:
  ```css
  /* 页面背景 */
  bg-[#f8fafc]  /* 浅灰蓝色 */
  ```
  
  **卡片容器**:
  ```css
  /* 毛玻璃卡片 */
  bg-white/70 backdrop-blur-xl rounded-3xl shadow-sm border border-white/40
  
  /* 内容区域 */
  bg-white/50 backdrop-blur-sm rounded-2xl
  ```
  
  **按钮样式**:
  ```css
  /* 主按钮 - 渐变蓝紫色 */
  bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700
  text-white font-bold rounded-2xl px-6 py-3
  transition-all duration-200 shadow-lg shadow-blue-500/25
  
  /* 次要按钮 */
  bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-2xl
  
  /* 危险按钮 */
  bg-gradient-to-r from-red-500 to-pink-500 text-white rounded-2xl
  ```
  
  **输入框样式**:
  ```css
  /* 输入框 */
  bg-slate-100 rounded-2xl px-4 py-3 
  border-2 border-transparent focus:border-blue-500 focus:bg-white
  transition-all duration-200 outline-none
  
  /* 文本域 */
  bg-slate-100 rounded-2xl px-4 py-3 resize-none
  ```
  
  **文字样式**:
  ```css
  /* 标题 */
  text-2xl font-black text-slate-800
  
  /* 副标题 */
  text-lg font-bold text-slate-700
  
  /* 正文 */
  text-sm text-slate-600
  
  /* 标签 */
  text-xs font-medium text-slate-500
  ```
  
  **布局组件参考**:
  - 主布局: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/Layout/MainLayout.tsx`
  - 设置弹窗: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/Settings/SettingsModal.tsx`
  - 按钮组件: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/common/Button.tsx`
  - 输入框组件: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/common/Input.tsx`
  - 进度条组件: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/GenerateArea/ProgressBar.tsx`

  **Acceptance Criteria**:
  - [ ] 页面背景使用 `bg-[#f8fafc]`
  - [ ] 卡片使用毛玻璃效果 `bg-white/70 backdrop-blur-xl rounded-3xl`
  - [ ] 按钮使用渐变色 `bg-gradient-to-r from-blue-600 to-indigo-600`
  - [ ] 输入框使用 `bg-slate-100 rounded-2xl`
  - [ ] 路由切换正常
  - [ ] 设置可保存
  - [ ] EventSource可连接SSE端点

  **Commit**: YES
  - Message: `feat: implement frontend base layout and settings`
  - Files: `frontend/src/components/*.tsx, frontend/src/pages/*.tsx`

---

- [ ] 6. 搜索与评论抓取 + SSE状态推送

  **What to do**:
  - 实现B站视频搜索API调用
  - 实现评论列表获取(分页)
  - 实现楼中楼评论获取
  - 实现并发抓取(可配置并发数)
  - 实现抓取数量限制(50视频/500评论)
  - **实现SSE服务端** (Go标准库，无需第三方依赖)
  - **实现任务状态推送**，状态包括:
    - `parsing`: AI正在解析关键词
    - `waiting_confirm`: 等待用户确认
    - `searching`: 正在搜索视频 (进度: x/50)
    - `scraping`: 正在抓取评论 (进度: x/500)
    - `analyzing`: AI正在分析评论 (进度: x/total)
    - `generating`: 正在生成报告
    - `completed`: 完成
    - `error`: 出错 (附带错误信息)

  **SSE实现示例** (Go标准库):
  ```go
  func handleSSE(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "text/event-stream")
      w.Header().Set("Cache-Control", "no-cache")
      w.Header().Set("Connection", "keep-alive")
      
      flusher, _ := w.(http.Flusher)
      for status := range taskStatusChan {
          fmt.Fprintf(w, "data: %s\n\n", status)
          flusher.Flush()
      }
  }
  ```

  **SSE消息格式**:
  ```json
  {
    "task_id": "uuid",
    "status": "scraping",
    "progress": {"current": 150, "total": 500},
    "message": "正在抓取评论..."
  }
  ```

  **Must NOT do**:
  - 不超过配置的抓取限制
  - 不忽略API错误
  - 不在没有进度更新时频繁推送

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 7, 8)
  - **Blocks**: Task 9
  - **Blocked By**: Tasks 2, 3

  **References**:
  - 搜索接口: `/Users/daozhang/Downloads/bilibili-API-collect/搜索与评论接口汇总.md:61-104`
  - 评论接口: `/Users/daozhang/Downloads/bilibili-API-collect/搜索与评论接口汇总.md:108-188`
  - MDN SSE: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events

  **Acceptance Criteria**:
  - [ ] 搜索返回视频列表
  - [ ] 评论获取包含楼中楼
  - [ ] 并发控制生效
  - [ ] SSE端点可连接并接收状态推送
  - [ ] 进度信息实时更新

  **Commit**: YES
  - Message: `feat: implement bilibili scraping with SSE progress`
  - Files: `backend/bilibili/search.go, backend/bilibili/comment.go, backend/sse/handler.go`

---

- [ ] 7. AI关键词解析与澄清

  **What to do**:
  - 设计关键词解析Prompt，要求AI返回:
    - `clarification`: 对用户需求的澄清描述 (例如"您想了解家用无线吸尘器的评价，我将为您分析以下品牌...")
    - `brands`: 品牌列表 (例如["戴森", "小米", "追觅", "石头"])
    - `dimensions`: 6个评价维度 (例如[{name: "吸力", desc: "..."}, ...])
    - `keywords`: 搜索关键词列表
  - 实现用户输入→AI解析流程
  - 返回JSON供前端展示确认页面
  - **实现确认API**: 用户确认后才开始后续抓取

  **AI返回JSON示例**:
  ```json
  {
    "clarification": "您想了解家用无线吸尘器的真实用户评价。我将为您分析市面上主流品牌的用户反馈，重点关注以下6个维度...",
    "brands": ["戴森", "小米追觅", "石头", "添可", "美的"],
    "dimensions": [
      {"name": "吸力", "description": "吸尘效果强弱"},
      {"name": "续航", "description": "单次充电使用时长"},
      {"name": "静音", "description": "工作噪音大小"},
      {"name": "便携", "description": "重量和操作便利性"},
      {"name": "清洁效果", "description": "地面/缝隙清洁能力"},
      {"name": "性价比", "description": "价格与功能的匹配度"}
    ],
    "keywords": ["无线吸尘器测评", "吸尘器推荐", "戴森vs小米"]
  }
  ```

  **Must NOT do**:
  - 不硬编码评价维度
  - 不忽略AI解析失败
  - 不跳过用户确认直接开始抓取

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 6, 8)
  - **Blocks**: Task 9
  - **Blocked By**: Tasks 2, 4

  **References**:
  - OpenAI Chat API: https://platform.openai.com/docs/guides/chat

  **Acceptance Criteria**:
  - [ ] 输入"吸尘器"返回包含clarification的JSON
  - [ ] 返回品牌列表和6个评价维度
  - [ ] JSON格式正确，前端可直接渲染
  - [ ] 确认API可触发后续流程

  **Commit**: YES
  - Message: `feat: implement AI keyword parsing with clarification`
  - Files: `backend/ai/keyword.go, backend/api/confirm.go`

---

- [ ] 8. AI评论分析

  **What to do**:
  - 设计评论分析Prompt
  - 实现批量评论分析(并发)
  - 对每条评论的每个维度打分(1-10)
  - 实现分析结果存储

  **Must NOT do**:
  - 不一次性发送过多评论
  - 不忽略分析失败的评论

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 6, 7)
  - **Blocks**: Task 9
  - **Blocked By**: Tasks 2, 4

  **References**:
  - OpenAI Chat API: https://platform.openai.com/docs/guides/chat

  **Acceptance Criteria**:
  - [ ] 评论分析返回各维度得分
  - [ ] 并发分析生效
  - [ ] 结果正确存储

  **Commit**: YES
  - Message: `feat: implement AI comment analysis`
  - Files: `backend/ai/analysis.go`

---

- [ ] 9. 报告生成逻辑

  **What to do**:
  - 实现各维度得分计算(平均值)
  - 实现品牌/产品排名
  - 实现购买建议生成
  - 生成结构化报告JSON

  **Must NOT do**:
  - 不使用复杂的统计算法
  - 不生成过于冗长的报告

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 10
  - **Blocked By**: Tasks 6, 7, 8

  **References**:
  - 无外部参考

  **Acceptance Criteria**:
  - [ ] 各维度得分正确计算
  - [ ] 排名逻辑正确
  - [ ] JSON结构完整

  **Commit**: YES
  - Message: `feat: implement report generation logic`
  - Files: `backend/report/generator.go`

---

- [ ] 10. 前端报告页面 (含确认页和进度页)

  **What to do**:
  - **实现输入页面**: 商品类目输入框
  - **实现确认页面** (AI澄清后展示):
    - 显示AI的澄清描述 (clarification)
    - 显示品牌列表 (卡片/标签形式)
    - 显示6个评价维度 (名称+描述)
    - "确认开始分析"按钮
    - (第一期只读展示，不可修改)
  - **实现进度页面** (SSE驱动):
    - 显示当前状态文字
    - 显示进度条 (current/total)
    - 显示各阶段完成状态 (搜索✓ → 抓取✓ → 分析... → 生成)
  - **实现报告展示页面**:
    - 雷达图(6维度)
    - 柱状图(品牌对比)
    - 对比表格
    - 购买建议

  **页面流程**:
  ```
  输入页 → [AI解析] → 确认页 → [用户确认] → 进度页 → [完成] → 报告页
  ```

  **Must NOT do**:
  - 不过度动画效果
  - 不使用过于复杂的图表
  - 不在确认页实现编辑功能 (第一期)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Task 11)
  - **Blocks**: Task 12
  - **Blocked By**: Tasks 5, 9

  **References**:
  - Recharts: https://recharts.org/
  - Tailwind CSS: https://tailwindcss.com/
  - **UI风格参考项目**: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/`
  - 卡片组件参考: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/HistoryPanel/ImageCard.tsx`
  - 进度条参考: `/Users/daozhang/Trae_AI/文生图前后端/frontend/src/components/GenerateArea/ProgressBar.tsx`

  **UI风格规范 (MUST FOLLOW)**:
  
  继承Task 5的基础风格，以下是各页面特定样式：
  
  **输入页面** (OpenAI/ChatGPT 风格 - 居中简洁大输入框):
  ```tsx
  {/* 整体布局 - 垂直居中 */}
  <div className="min-h-screen bg-[#f8fafc] flex flex-col items-center justify-center px-4">
    
    {/* 标题区域 */}
    <div className="text-center mb-8">
      <h1 className="text-4xl md:text-5xl font-black text-slate-800 mb-3">
        B站商品评论分析
      </h1>
      <p className="text-lg text-slate-500">
        输入商品类目，AI 为您分析真实用户评价
      </p>
    </div>
    
    {/* 核心输入区域 - 类似 ChatGPT 的大输入框 */}
    <div className="w-full max-w-2xl">
      <div className="relative bg-white rounded-2xl shadow-lg shadow-slate-200/50 
                      border border-slate-200 overflow-hidden
                      focus-within:shadow-xl focus-within:shadow-blue-500/10 
                      focus-within:border-blue-300 transition-all duration-300">
        
        {/* 输入框 */}
        <input 
          className="w-full px-6 py-5 text-lg text-slate-800 
                     placeholder:text-slate-400 outline-none bg-transparent"
          placeholder="无线吸尘器、机械键盘、空气炸锅..."
        />
        
        {/* 发送按钮 - 右侧内嵌 */}
        <button className="absolute right-3 top-1/2 -translate-y-1/2
                           bg-gradient-to-r from-blue-600 to-indigo-600 
                           hover:from-blue-700 hover:to-indigo-700
                           text-white font-semibold rounded-xl px-5 py-2.5
                           transition-all duration-200 
                           shadow-md shadow-blue-500/25 hover:shadow-lg
                           disabled:opacity-50 disabled:cursor-not-allowed">
          开始分析
        </button>
      </div>
      
      {/* 快捷示例标签 */}
      <div className="flex flex-wrap justify-center gap-2 mt-4">
        {['无线吸尘器', '机械键盘', '空气炸锅', '投影仪'].map(example => (
          <button 
            key={example}
            className="px-4 py-2 text-sm text-slate-600 
                       bg-white/80 hover:bg-white rounded-full
                       border border-slate-200 hover:border-slate-300
                       transition-all duration-200 hover:shadow-sm"
          >
            {example}
          </button>
        ))}
      </div>
    </div>
    
    {/* 底部说明 */}
    <p className="mt-12 text-sm text-slate-400">
      基于 B站视频评论区的真实用户反馈
    </p>
  </div>
  ```
  
  **输入框交互细节**:
  - 默认状态: 白底 + 浅灰边框 + 柔和阴影
  - 聚焦状态: 边框变蓝 + 阴影加深带蓝色
  - 按钮内嵌在输入框右侧，保持整体简洁
  - 快捷示例标签可点击快速填入
  - 整体垂直居中，留白充足
  
  **确认页面**:
  ```tsx
  {/* 澄清描述卡片 */}
  <div className="bg-blue-50/80 backdrop-blur-sm rounded-2xl p-6 border border-blue-100">
    <p className="text-slate-700 leading-relaxed">{clarification}</p>
  </div>
  
  {/* 品牌标签 */}
  <div className="flex flex-wrap gap-2">
    {brands.map(brand => (
      <span className="bg-white/80 backdrop-blur-sm rounded-xl px-4 py-2 
                       text-sm font-medium text-slate-700 
                       border border-slate-200 shadow-sm">
        {brand}
      </span>
    ))}
  </div>
  
  {/* 维度卡片 */}
  <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
    {dimensions.map(dim => (
      <div className="bg-white/60 backdrop-blur-sm rounded-xl p-4 border border-white/40">
        <h4 className="font-bold text-slate-800">{dim.name}</h4>
        <p className="text-xs text-slate-500 mt-1">{dim.description}</p>
      </div>
    ))}
  </div>
  ```
  
  **进度页面**:
  ```tsx
  {/* 进度容器 */}
  <div className="bg-white/70 backdrop-blur-xl rounded-3xl shadow-sm border border-white/40 p-8">
    {/* 状态文字 */}
    <h2 className="text-xl font-bold text-slate-800 mb-2">{statusText}</h2>
    <p className="text-sm text-slate-500 mb-6">{message}</p>
    
    {/* 进度条 */}
    <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
      <div 
        className="h-full bg-gradient-to-r from-blue-500 to-indigo-500 rounded-full
                   transition-all duration-500 ease-out"
        style={{ width: `${progress}%` }}
      />
    </div>
    
    {/* 阶段指示器 */}
    <div className="flex justify-between mt-6">
      {stages.map((stage, i) => (
        <div className={`flex flex-col items-center ${stage.done ? 'text-green-600' : 'text-slate-400'}`}>
          <div className={`w-8 h-8 rounded-full flex items-center justify-center
                          ${stage.done ? 'bg-green-100' : stage.current ? 'bg-blue-100 animate-pulse' : 'bg-slate-100'}`}>
            {stage.done ? '✓' : i + 1}
          </div>
          <span className="text-xs mt-1">{stage.name}</span>
        </div>
      ))}
    </div>
  </div>
  ```
  
  **报告页面**:
  ```tsx
  {/* 报告标题区 */}
  <div className="bg-gradient-to-r from-blue-600 to-indigo-600 rounded-3xl p-6 text-white mb-6">
    <h1 className="text-2xl font-black">分析报告</h1>
    <p className="text-blue-100 mt-1">{category} · {brandCount}个品牌 · {commentCount}条评论</p>
  </div>
  
  {/* 图表卡片 */}
  <div className="bg-white/70 backdrop-blur-xl rounded-3xl shadow-sm border border-white/40 p-6">
    <h3 className="text-lg font-bold text-slate-800 mb-4">维度对比</h3>
    {/* Recharts 雷达图 */}
  </div>
  
  {/* 购买建议卡片 */}
  <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-6 border border-green-100">
    <h3 className="text-lg font-bold text-green-800 mb-2">💡 购买建议</h3>
    <p className="text-green-700 leading-relaxed">{recommendation}</p>
  </div>
  
  {/* 对比表格 */}
  <div className="bg-white/70 backdrop-blur-xl rounded-3xl shadow-sm border border-white/40 overflow-hidden">
    <table className="w-full">
      <thead className="bg-slate-50/80">
        <tr>
          <th className="px-4 py-3 text-left text-sm font-bold text-slate-700">品牌</th>
          {dimensions.map(d => (
            <th className="px-4 py-3 text-center text-sm font-bold text-slate-700">{d.name}</th>
          ))}
          <th className="px-4 py-3 text-center text-sm font-bold text-slate-700">综合</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-slate-100">
        {/* 表格行 */}
      </tbody>
    </table>
  </div>
  ```
  
  **Recharts 图表配色**:
  ```tsx
  // 雷达图配色
  const RADAR_COLORS = {
    fill: 'rgba(99, 102, 241, 0.2)',  // indigo-500 with opacity
    stroke: '#6366f1',                 // indigo-500
  }
  
  // 柱状图配色 (品牌区分)
  const BAR_COLORS = [
    '#3b82f6', // blue-500
    '#6366f1', // indigo-500
    '#8b5cf6', // violet-500
    '#a855f7', // purple-500
    '#d946ef', // fuchsia-500
  ]
  ```

  **Acceptance Criteria**:
  - [ ] 输入页使用毛玻璃卡片和渐变按钮
  - [ ] 确认页正确显示澄清描述、品牌标签、维度卡片
  - [ ] 进度页显示渐变进度条和阶段指示器
  - [ ] 进度页实时显示SSE推送的状态
  - [ ] 报告页使用渐变标题区
  - [ ] 雷达图正确显示6维度
  - [ ] 柱状图显示品牌对比
  - [ ] 表格使用圆角和分割线样式
  - [ ] 购买建议使用绿色渐变背景

  **Commit**: YES
  - Message: `feat: implement confirmation, progress and report pages`
  - Files: `frontend/src/pages/Confirm.tsx, frontend/src/pages/Progress.tsx, frontend/src/pages/Report.tsx`

---

- [ ] 11. 历史记录功能

  **What to do**:
  - 实现历史记录列表页面
  - 实现历史报告查看
  - 实现历史记录删除
  - 实现启动时3天数据清理

  **Must NOT do**:
  - 不删除用户明确保存的报告
  - 不实现复杂的搜索过滤

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Task 10)
  - **Blocks**: Task 12
  - **Blocked By**: Tasks 5, 2

  **References**:
  - 无外部参考

  **Acceptance Criteria**:
  - [ ] 历史列表正确显示
  - [ ] 可查看历史报告
  - [ ] 3天清理逻辑生效

  **Commit**: YES
  - Message: `feat: implement history management`
  - Files: `frontend/src/pages/History.tsx, backend/api/history.go`

---

- [ ] 12. 集成测试与优化

  **What to do**:
  - 端到端流程测试
  - 错误处理优化
  - 性能优化(如需要)
  - 编写README文档

  **Must NOT do**:
  - 不添加新功能
  - 不过度优化

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Final
  - **Blocks**: None
  - **Blocked By**: Tasks 10, 11

  **References**:
  - 无外部参考

  **Acceptance Criteria**:
  - [ ] 完整流程可运行
  - [ ] 错误有友好提示
  - [ ] README完整

  **Commit**: YES
  - Message: `docs: add README and finalize project`
  - Files: `README.md`

---

## Commit Strategy

| After Task | Message | Files |
|------------|---------|-------|
| 1 | `feat: initialize project structure` | go.mod, frontend/package.json |
| 2 | `feat: implement database models` | backend/models/*.go |
| 3 | `feat: implement WBI signing` | backend/bilibili/wbi.go |
| 4 | `feat: implement AI client` | backend/ai/client.go |
| 5 | `feat: implement frontend base` | frontend/src/**/*.tsx |
| 6 | `feat: implement scraping` | backend/bilibili/*.go |
| 7 | `feat: implement keyword parsing` | backend/ai/keyword.go |
| 8 | `feat: implement comment analysis` | backend/ai/analysis.go |
| 9 | `feat: implement report generation` | backend/report/*.go |
| 10 | `feat: implement report pages` | frontend/src/pages/*.tsx |
| 11 | `feat: implement history` | frontend/src/pages/History.tsx |
| 12 | `docs: finalize project` | README.md |

---

## Success Criteria

### Verification Commands

```bash
# 后端启动
cd backend && go run main.go
# Expected: Server running on :8080

# 前端启动
cd frontend && npm run dev
# Expected: Vite dev server on :5173

# API测试
curl http://localhost:8080/api/health
# Expected: {"status": "ok"}
```

### Final Checklist
- [ ] 用户可配置AI和B站Cookie
- [ ] 输入商品类目可获取分析报告
- [ ] 报告包含雷达图、柱状图、表格
- [ ] 历史记录可保存和查看
- [ ] 3天数据自动清理
