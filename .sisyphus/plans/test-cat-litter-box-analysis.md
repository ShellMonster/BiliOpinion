# 测试自动猫砂盆分析任务

## TL;DR

> **Quick Summary**: 创建并执行一个完整的"自动猫砂盆"分析任务，验证当前系统的端到端功能
> 
> **Deliverables**: 
> - 完整的分析报告（包含品牌排名、型号排名、典型评论等）
> - 验证所有新功能正常工作
> 
> **Estimated Effort**: Quick (10-15分钟)
> **Parallel Execution**: NO - sequential
> **Critical Path**: 启动服务 → 创建任务 → 等待完成 → 验证结果

---

## Context

### 用户需求
用户已经在前端设置中配置好了：
- AI API Base URL
- AI API Key
- AI Model
- B站 Cookie

现在需要运行一个完整的分析任务来验证系统功能。

### 测试用例
**商品类目**: 自动猫砂盆

**预期结果**:
- AI 解析出品牌（如：小佩、CATLINK、拓牛等）
- AI 解析出评价维度（如：清洁效果、噪音、气味控制等）
- 抓取 B站 视频评论
- AI 分析评论并提取品牌和型号
- 生成完整报告，包含：
  - 品牌排名
  - 型号排名（新功能）
  - 数据统计
  - 典型评论
  - 优劣势分析
  - AI 购买建议

---

## Work Objectives

### Core Objective
验证系统的完整分析流程，确保所有功能正常工作。

### Concrete Deliverables
- 成功创建的分析任务
- 完整的分析报告
- 验证报告截图或数据

### Definition of Done
- [ ] 后端服务正常运行
- [ ] 前端服务正常运行
- [ ] 成功创建分析任务
- [ ] 任务完成并生成报告
- [ ] 报告包含所有预期内容

### Must Have
- 后端和前端服务都在运行
- 配置已保存（用户已完成）
- 完整的分析流程

### Must NOT Have (Guardrails)
- 不修改任何代码
- 不修改配置（用户已设置）
- 不中断正在运行的任务

---

## TODOs

- [ ] 1. 检查并启动后端服务

  **What to do**:
  - 检查后端服务是否在运行（端口 8080）
  - 如果未运行，启动后端服务
  
  **Commands**:
  ```bash
  # 检查后端是否运行
  curl -s http://localhost:8080/api/config > /dev/null && echo "后端运行中" || echo "后端未运行"
  
  # 如果未运行，启动后端（在后台）
  cd /Users/daozhang/Trae_AI/B站商品评论解析 && go run backend/main.go &
  ```
  
  **Acceptance Criteria**:
  - [ ] `curl http://localhost:8080/api/config` 返回 200
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 2
  - **Blocked By**: None

---

- [ ] 2. 检查并启动前端服务

  **What to do**:
  - 检查前端服务是否在运行（端口 5173 或 5177）
  - 如果未运行，启动前端服务
  
  **Commands**:
  ```bash
  # 检查前端是否运行
  curl -s http://localhost:5173 > /dev/null && echo "前端运行中(5173)" || \
  curl -s http://localhost:5177 > /dev/null && echo "前端运行中(5177)" || \
  echo "前端未运行"
  
  # 如果未运行，启动前端（在后台）
  cd /Users/daozhang/Trae_AI/B站商品评论解析/frontend && npm run dev &
  ```
  
  **Acceptance Criteria**:
  - [ ] 前端页面可访问
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

---

- [ ] 3. 创建分析任务

  **What to do**:
  - 使用 API 创建一个"自动猫砂盆"分析任务
  - 记录返回的 task_id
  
  **Commands**:
  ```bash
  # 步骤1: 调用解析 API
  curl -X POST http://localhost:8080/api/parse \
    -H "Content-Type: application/json" \
    -d '{"category":"自动猫砂盆"}' \
    -o /tmp/parse_result.json
  
  # 查看解析结果
  cat /tmp/parse_result.json | jq .
  
  # 步骤2: 使用解析结果调用确认 API
  # 注意：需要从解析结果中提取 brands, dimensions, keywords
  curl -X POST http://localhost:8080/api/confirm \
    -H "Content-Type: application/json" \
    -d @/tmp/parse_result.json \
    -o /tmp/confirm_result.json
  
  # 查看任务ID
  cat /tmp/confirm_result.json | jq .
  ```
  
  **Acceptance Criteria**:
  - [ ] 解析 API 返回品牌、维度、关键词
  - [ ] 确认 API 返回 task_id
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 4
  - **Blocked By**: Task 2

---

- [ ] 4. 监控任务进度

  **What to do**:
  - 使用 SSE 接口监控任务进度
  - 等待任务完成
  - 记录 report_id
  
  **Commands**:
  ```bash
  # 从确认结果中提取 task_id
  TASK_ID=$(cat /tmp/confirm_result.json | jq -r .task_id)
  echo "任务ID: $TASK_ID"
  
  # 监控 SSE 进度（使用 curl）
  curl -N http://localhost:8080/api/sse?task_id=$TASK_ID
  
  # 或者使用 tail 查看后端日志
  # 后端日志会显示任务进度
  ```
  
  **Acceptance Criteria**:
  - [ ] SSE 推送进度消息
  - [ ] 最终收到 "completed" 状态和 report_id
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: Task 5
  - **Blocked By**: Task 3
  
  **Note**: 这个任务可能需要 3-10 分钟，取决于视频和评论数量

---

- [ ] 5. 验证报告内容

  **What to do**:
  - 获取生成的报告
  - 验证报告包含所有预期内容
  
  **Commands**:
  ```bash
  # 从 SSE 结果中提取 report_id（假设为 5）
  REPORT_ID=5  # 替换为实际的 report_id
  
  # 获取报告
  curl http://localhost:8080/api/report/$REPORT_ID -o /tmp/report.json
  
  # 验证报告内容
  cat /tmp/report.json | jq '{
    category: .category,
    brand_count: (.rankings | length),
    model_count: (.model_rankings | length),
    dimensions: (.dimensions | map(.name)),
    has_stats: (.stats != null),
    has_top_comments: (.top_comments != null),
    has_bad_comments: (.bad_comments != null),
    has_brand_analysis: (.brand_analysis != null),
    has_recommendation: (.recommendation != null)
  }'
  ```
  
  **Acceptance Criteria**:
  - [ ] 报告包含品牌排名（rankings）
  - [ ] 报告包含型号排名（model_rankings）
  - [ ] 报告包含数据统计（stats）
  - [ ] 报告包含典型评论（top_comments, bad_comments）
  - [ ] 报告包含优劣势分析（brand_analysis）
  - [ ] 报告包含 AI 购买建议（recommendation）
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: None
  - **Blocked By**: Task 4

---

- [ ] 6. 在浏览器中查看报告

  **What to do**:
  - 打开浏览器访问报告页面
  - 验证前端展示正常
  
  **URL**:
  ```
  http://localhost:5173/report/{report_id}
  或
  http://localhost:5177/report/{report_id}
  ```
  
  **验证内容**:
  - [ ] 品牌排名雷达图显示正常
  - [ ] 品牌综合得分柱状图显示正常
  - [ ] 型号排名表格显示正常（新功能）
  - [ ] 数据统计卡片显示正常
  - [ ] 典型评论展示正常
  - [ ] 优劣势分析展示正常
  - [ ] AI 购买建议展示正常
  - [ ] 详细得分表格显示正常（无空值错误）
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
  
  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Blocks**: None
  - **Blocked By**: Task 5

---

## Execution Strategy

### 依赖关系
```
Task 1 (启动后端) → Task 2 (启动前端) → Task 3 (创建任务) → Task 4 (监控进度) → Task 5 (验证报告) → Task 6 (浏览器查看)
```

### 预计时间
- Task 1-2: 1-2 分钟（启动服务）
- Task 3: 30 秒（创建任务）
- Task 4: 3-10 分钟（等待分析完成）
- Task 5: 30 秒（验证报告）
- Task 6: 1 分钟（浏览器查看）

**总计**: 约 5-15 分钟

---

## Success Criteria

### 验证命令
```bash
# 1. 检查服务状态
curl http://localhost:8080/api/config
curl http://localhost:5173

# 2. 检查报告
curl http://localhost:8080/api/report/{report_id} | jq .

# 3. 检查历史记录
curl http://localhost:8080/api/history | jq .
```

### Final Checklist
- [ ] 后端服务运行正常
- [ ] 前端服务运行正常
- [ ] 成功创建分析任务
- [ ] 任务完成并生成报告
- [ ] 报告包含品牌排名
- [ ] 报告包含型号排名（新功能）
- [ ] 报告包含数据统计
- [ ] 报告包含典型评论
- [ ] 报告包含优劣势分析
- [ ] 报告包含 AI 购买建议
- [ ] 前端展示无错误（包括空值检查修复）

---

## Notes

### 测试数据预期
根据之前的经验，"自动猫砂盆"分析可能会发现：
- **品牌**: 小佩、CATLINK、拓牛、Petree 等
- **型号**: Pura X、SCOOPER、T3、T4 等
- **维度**: 清洁效果、噪音、气味控制、易用性、耐用性等

### 故障排查
如果任务失败，检查：
1. B站 Cookie 是否过期（重新获取）
2. AI API Key 是否有效
3. 网络连接是否正常
4. 后端日志中的错误信息

### 数据库位置
- SQLite 数据库: `/Users/daozhang/Trae_AI/B站商品评论解析/data/bilibili.db`
- 可以使用 SQLite 客户端查看数据
