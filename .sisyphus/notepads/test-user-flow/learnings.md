# Playwright 测试学习记录

## 测试时间
2026-02-01 15:30

## 测试目标
验证 B站商品评论分析工具的完整用户流程（自由需求输入功能）

## 成功要点

### 1. 按钮定位策略
- **问题**：首页提交按钮是纯图标按钮（箭头图标），没有文本
- **解决方案**：使用 `page.locator('button:has(svg)')` 定位包含 SVG 图标的按钮
- **经验**：不要假设所有按钮都有文本，优先使用结构化选择器

### 2. 异步加载等待
- **问题**：确认页需要等待 AI API 返回数据，加载时间较长（30-60秒）
- **解决方案**：
  - 等待加载动画消失：`page.waitForSelector('.animate-spin', { state: 'detached', timeout: 60000 })`
  - 等待关键内容出现：`page.waitForSelector('text=/商品类型|品牌|评价维度/', { timeout: 10000 })`
  - 额外等待2秒确保内容完全渲染
- **经验**：对于 AI API 调用，需要设置足够长的超时时间（60秒）

### 3. 浏览器配置
- **问题**：Playwright 内置浏览器未安装
- **解决方案**：使用系统 Chrome 浏览器 `{ channel: 'chrome' }`
- **经验**：生产环境建议使用 headless 模式，开发调试时使用 headless: false

## 测试结果

### 全部通过 ✓
1. 首页访问正常
2. 输入框为自由需求输入（提示文字："描述你的需求，比如：想买个吸尘器，预算2000，家里有宠物..."）
3. 快捷示例为需求描述（"想买个机械键盘，主要打游戏用" 等）
4. 需求提交成功
5. 页面跳转到确认页
6. AI 理解描述显示正确（"我理解您想购买一款适合通勤使用的蓝牙耳机..."）
7. 商品类型显示正确（蓝牙耳机）
8. 品牌标签显示正确（索尼、BOSE、Apple、小米、华为）
9. 评价维度卡片显示正确（降噪效果、音质、佩戴舒适度、电池续航、连接稳定性、便携性）
10. 搜索关键词显示正确（通勤蓝牙耳机 | 降噪耳机推荐 | 蓝牙耳机评测）
11. 截图保存成功

## 关键代码片段

```javascript
// 等待加载动画消失（最多等待60秒）
await page.waitForSelector('.animate-spin', { state: 'detached', timeout: 60000 });

// 定位图标按钮
const analyzeButton = page.locator('button:has(svg)').first();
await analyzeButton.click();

// 等待 URL 跳转
await page.waitForURL(/\/confirm\?requirement=/, { timeout: 10000 });
```

## 项目特点
- 前端：React + TypeScript + Vite + Tailwind CSS v4
- 后端：Go + Gin + GORM + SQLite
- AI 集成：OpenAI API
- 实时推送：SSE (Server-Sent Events)

