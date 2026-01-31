# B站商品评论分析工具

通过抓取B站视频评论，使用AI分析商品的真实用户评价，生成可视化报告帮助用户做购买决策。

## 功能特性

- 🤖 **AI智能解析**：自动提取品牌、评价维度和搜索关键词
- 📊 **多维度分析**：6个评价维度，全面了解商品表现
- 📈 **可视化报告**：雷达图、柱状图、对比表格
- 🔄 **实时进度**：SSE推送任务状态，实时查看进度
- 📝 **历史记录**：保存分析历史，随时查看报告
- ⚙️ **灵活配置**：支持自定义AI服务和B站Cookie

## 技术栈

### 后端
- **Go 1.24** - 高性能后端服务
- **Gin** - Web框架
- **GORM** - ORM框架
- **SQLite** - 轻量级数据库

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Tailwind CSS v4** - 样式框架
- **Recharts** - 图表库

## 快速开始

### 前置要求

- Go 1.24+
- Node.js 18+
- npm 或 pnpm

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd B站商品评论解析
```

2. **安装后端依赖**
```bash
go mod download
```

3. **安装前端依赖**
```bash
cd frontend
npm install
```

4. **启动后端服务**
```bash
go run backend/main.go
# 服务运行在 http://localhost:8080
```

5. **启动前端开发服务器**
```bash
cd frontend
npm run dev
# 访问 http://localhost:5173
```

## 使用说明

### 1. 配置设置

首次使用需要配置AI服务和B站Cookie：

1. 访问 http://localhost:5173/settings
2. 填写AI配置：
   - API Base URL（默认：https://api.openai.com/v1）
   - API Key
   - Model（如：gpt-3.5-turbo）
3. 填写B站Cookie（从浏览器开发者工具复制完整Cookie字符串）
4. 点击"保存设置"

### 2. 开始分析

1. 在首页输入商品类目（如："无线吸尘器"）
2. AI解析后展示品牌、维度和关键词
3. 确认后开始抓取和分析
4. 实时查看进度
5. 查看分析报告

### 3. 查看历史

访问 http://localhost:5173/history 查看所有历史分析记录。

## API文档

### 解析API
```http
POST /api/parse
Content-Type: application/json

{
  "category": "无线吸尘器"
}
```

**响应示例：**
```json
{
  "brands": ["戴森", "小米", "石头"],
  "dimensions": [
    {"name": "吸力", "description": "吸尘效果和清洁能力"},
    {"name": "续航", "description": "电池续航时间"}
  ],
  "keywords": [
    {"brand": "戴森", "keyword": "戴森吸尘器"},
    {"brand": "小米", "keyword": "小米吸尘器"}
  ]
}
```

### 确认API
```http
POST /api/confirm
Content-Type: application/json

{
  "category": "无线吸尘器",
  "brands": ["戴森", "小米"],
  "dimensions": [
    {"name": "吸力", "description": "吸尘效果和清洁能力"}
  ],
  "keywords": [
    {"brand": "戴森", "keyword": "戴森吸尘器"}
  ]
}
```

**响应示例：**
```json
{
  "task_id": "task_1738425600_abc123"
}
```

### SSE进度推送
```http
GET /api/sse?task_id=<task_id>
```

**事件流示例：**
```
data: {"status":"fetching","message":"正在抓取戴森评论...","progress":10}

data: {"status":"analyzing","message":"正在分析评论...","progress":50}

data: {"status":"completed","message":"分析完成","progress":100,"report_id":1}
```

### 历史记录API

**获取历史列表**
```http
GET /api/history
```

**获取历史详情**
```http
GET /api/history/:id
```

**删除历史记录**
```http
DELETE /api/history/:id
```

### 配置API

**获取配置**
```http
GET /api/config
```

**保存配置**
```http
POST /api/config
Content-Type: application/json

{
  "ai_base_url": "https://api.openai.com/v1",
  "ai_api_key": "sk-xxx",
  "ai_model": "gpt-3.5-turbo",
  "bilibili_cookie": "SESSDATA=xxx; buvid3=xxx"
}
```

## 项目结构

```
B站商品评论解析/
├── backend/              # Go后端
│   ├── ai/              # AI服务集成
│   │   └── service.go   # OpenAI API调用
│   ├── api/             # API路由
│   │   ├── parse.go     # 解析接口
│   │   ├── confirm.go   # 确认接口
│   │   ├── history.go   # 历史记录接口
│   │   ├── config.go    # 配置接口
│   │   └── sse.go       # SSE接口
│   ├── bilibili/        # B站API集成
│   │   └── client.go    # B站评论抓取
│   ├── database/        # 数据库初始化
│   │   └── db.go        # SQLite初始化
│   ├── models/          # 数据模型
│   │   ├── config.go    # 配置模型
│   │   ├── history.go   # 历史记录模型
│   │   └── report.go    # 报告模型
│   ├── report/          # 报告生成
│   │   └── generator.go # 报告生成逻辑
│   ├── sse/             # SSE服务端
│   │   └── manager.go   # SSE连接管理
│   └── main.go          # 入口文件
├── frontend/            # React前端
│   ├── src/
│   │   ├── api/        # API客户端
│   │   │   └── client.ts
│   │   ├── components/ # 组件
│   │   │   ├── AnalysisForm.tsx      # 分析表单
│   │   │   ├── ConfirmDialog.tsx     # 确认对话框
│   │   │   ├── ProgressMonitor.tsx   # 进度监控
│   │   │   ├── ReportView.tsx        # 报告展示
│   │   │   └── HistoryList.tsx       # 历史列表
│   │   ├── hooks/      # 自定义hooks
│   │   │   └── useSSE.ts             # SSE连接hook
│   │   ├── pages/      # 页面
│   │   │   ├── Home.tsx              # 首页
│   │   │   ├── History.tsx           # 历史页
│   │   │   └── Settings.tsx          # 设置页
│   │   ├── App.tsx     # 应用入口
│   │   └── main.tsx    # 主入口
│   ├── package.json
│   └── vite.config.ts
├── data/                # SQLite数据库
│   └── bilibili.db
├── go.mod
├── go.sum
└── README.md
```

## 数据清理

- **原始评论数据**：自动清理3天前的数据（节省存储空间）
- **分析历史和报告**：永久保存（可手动删除）

## 开发

### 运行测试
```bash
# 后端测试
go test ./backend/...

# 前端测试
cd frontend
npm test
```

### 构建生产版本

**后端构建：**
```bash
go build -o bin/bilibili-analyzer ./backend
./bin/bilibili-analyzer
```

**前端构建：**
```bash
cd frontend
npm run build
# 构建产物在 frontend/dist 目录
```

### 代码规范

**后端（Go）：**
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加中文注释说明关键逻辑

**前端（TypeScript）：**
- 遵循 ESLint 规则
- 使用 Prettier 格式化代码
- 组件使用函数式组件 + Hooks

## 常见问题

### 1. B站Cookie如何获取？

1. 打开浏览器，访问 https://www.bilibili.com
2. 登录你的B站账号
3. 按 F12 打开开发者工具
4. 切换到 "Network" 标签
5. 刷新页面，找到任意请求
6. 在请求头中找到 "Cookie" 字段
7. 复制完整的 Cookie 字符串（包含 SESSDATA 和 buvid3）

### 2. AI API调用失败怎么办？

- 检查 API Key 是否正确
- 检查 API Base URL 是否可访问
- 检查网络连接是否正常
- 查看后端日志获取详细错误信息

### 3. 评论抓取失败怎么办？

- 检查 B站 Cookie 是否过期（需要重新获取）
- 检查搜索关键词是否有相关视频
- 检查网络连接是否正常
- B站可能有反爬限制，建议降低并发数

### 4. 如何自定义评价维度？

在解析结果确认页面，可以：
- 添加新维度
- 删除不需要的维度
- 修改维度名称和描述

### 5. 报告数据不准确怎么办？

可能原因：
- 评论数量太少（建议至少50条）
- 关键词不够精准（调整搜索关键词）
- AI模型理解偏差（尝试更换模型）

## 注意事项

1. **B站Cookie**：需要从浏览器复制完整Cookie字符串（包含SESSDATA和buvid3），Cookie会过期需要定期更新
2. **API限流**：建议设置合理的并发数，避免触发B站API限流（默认并发数为3）
3. **AI成本**：每次分析会调用多次AI API，注意成本控制（建议使用 gpt-3.5-turbo）
4. **数据隐私**：所有数据存储在本地SQLite数据库，不会上传到任何服务器
5. **评论质量**：分析结果依赖评论质量，建议选择热门视频的评论

## 技术亮点

1. **SSE实时推送**：使用Server-Sent Events实现实时进度推送，用户体验流畅
2. **并发控制**：使用Goroutine池控制并发数，避免API限流
3. **错误重试**：网络请求失败自动重试，提高稳定性
4. **数据清理**：自动清理过期数据，节省存储空间
5. **类型安全**：前端使用TypeScript，后端使用Go，类型安全有保障
6. **响应式设计**：前端使用Tailwind CSS，支持移动端和桌面端

## 未来规划

- [ ] 支持更多视频平台（抖音、小红书等）
- [ ] 支持导出PDF报告
- [ ] 支持多用户和权限管理
- [ ] 支持自定义AI提示词
- [ ] 支持评论情感分析
- [ ] 支持品牌对比分析
- [ ] 支持数据可视化大屏

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交 Issue
- 发送邮件

---

**感谢使用 B站商品评论分析工具！**
