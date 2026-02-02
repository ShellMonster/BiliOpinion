# 修复型号提取和小数点精度问题

## TL;DR

> **Quick Summary**: 修复两个问题：1) 型号排名数据为空 2) 分数显示小数位过多
> 
> **Deliverables**: 
> - 修复型号提取逻辑
> - 统一小数点精度为1位
> 
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 2 waves

---

## 问题分析

### 问题1: 型号排名数据为空

**现象**: `model_rankings` 字段不存在或为空数组

**可能原因**:
1. AI 没有正确提取型号（返回空字符串或"未知"）
2. 后端运行的是旧版本代码（已通过重新编译解决）
3. 正则后备提取没有匹配到型号

**验证方法**:
```bash
# 查看后端日志，检查 AI 返回的 model 字段
tail -100 /tmp/backend.log | grep -i model
```

### 问题2: 小数点精度不一致

**现象**: 分数显示为 `5.0446539684344565` 而不是 `5.0`

**原因**:
1. 后端 `generator.go` 计算平均分时没有做精度控制
2. 前端部分地方用 `toFixed(2)` 而不是 `toFixed(1)`

**需要修复的位置**:
- 后端: `backend/report/generator.go` - 计算平均分时四舍五入
- 前端: `frontend/src/pages/Report.tsx:446` - 改为 `toFixed(1)`

---

## TODOs

### Wave 1: 后端修复

- [ ] 1. 后端添加小数点精度控制

  **What to do**:
  修改 `backend/report/generator.go`，在计算平均分时使用 `math.Round` 保留1位小数
  
  **修改位置1** (约第 115-120 行，品牌得分计算):
  ```go
  // 计算平均分并四舍五入到1位小数
  avg := sum / float64(len(scores))
  avg = math.Round(avg*10) / 10  // 保留1位小数
  ```
  
  **修改位置2** (约第 236-237 行，型号得分计算):
  ```go
  avg := sum / float64(len(scores))
  avg = math.Round(avg*10) / 10  // 保留1位小数
  avgScores[dimName] = avg
  ```
  
  **修改位置3** (约第 243-244 行，综合得分计算):
  ```go
  overallScore := 0.0
  if dimCount > 0 {
      overallScore = math.Round((total/float64(dimCount))*10) / 10
  }
  ```
  
  **修改位置4** (约第 305-310 行，品牌综合得分):
  ```go
  overallScore := 0.0
  if count > 0 {
      overallScore = math.Round((total/float64(count))*10) / 10
  }
  ```
  
  **需要添加 import**:
  ```go
  import "math"
  ```
  
  **Acceptance Criteria**:
  - [ ] `go build ./backend/...` 通过
  - [ ] API 返回的分数都是1位小数
  
  **Commit**: YES
  - Message: `fix(report): round scores to 1 decimal place`
  - Files: `backend/report/generator.go`

---

- [ ] 2. 增强型号提取的 AI 提示词

  **What to do**:
  修改 `backend/ai/analysis.go` 的系统提示词，更强调型号提取
  
  **修改位置** (约第 67-90 行):
  在系统提示词中添加更多型号示例和强调：
  
  ```go
  systemPrompt := fmt.Sprintf(`你是一个专业的商品评论分析助手。你的任务是：

  1. 从视频标题和评论内容中识别：
     - 品牌名称（如"戴森"、"小米"、"苹果"、"Sony"、"小佩"、"CATLINK"）
     - 具体型号（如"V12"、"iPhone 15 Pro"、"G10"、"Pura X"、"SCOOPER"）

  重要：型号提取规则（必须严格遵守）：
  - 型号是产品的具体版本标识，通常包含字母和数字
  - 常见型号格式：
    * 字母+数字：V12, G10, T3, T4, X5
    * 品牌+型号：iPhone 15, Galaxy S24, Pura X
    * 系列名+版本：Pro Max, Ultra, Plus, Lite
    * 完整型号名：WH-1000XM5, AirPods Pro 2
  - 如果评论明确提到型号，必须提取出来
  - 如果评论只提到品牌没有型号，model 字段返回空字符串 ""
  - 不要猜测型号，只提取明确提到的
  
  猫砂盆常见型号示例：
  - 小佩：Pura X, Pura Max, T3, T4
  - CATLINK：SCOOPER, SCOOPER SE, Young
  - 霍曼：M1, M1 Pro
  - 拓牛：T1, T1S, T Air
  ...`)
  ```
  
  **Acceptance Criteria**:
  - [ ] `go build ./backend/...` 通过
  - [ ] AI 能正确提取评论中的型号
  
  **Commit**: YES
  - Message: `feat(ai): enhance model extraction prompt with more examples`
  - Files: `backend/ai/analysis.go`

---

- [ ] 3. 增强正则后备提取

  **What to do**:
  修改 `backend/task/executor.go` 的 `extractModelFromContent` 函数，添加更多型号模式
  
  **查找函数位置**:
  ```bash
  grep -n "extractModelFromContent" backend/task/executor.go
  ```
  
  **增强正则模式**:
  ```go
  func extractModelFromContent(content string) string {
      // 常见型号模式
      patterns := []string{
          // 字母+数字组合 (V12, G10, T3, M1)
          `\b([A-Z]\d{1,2})\b`,
          // 品牌+型号 (Pura X, Pura Max)
          `(?i)\b(Pura\s*(?:X|Max|Air))\b`,
          // CATLINK 型号
          `(?i)\b(SCOOPER(?:\s*SE)?|Young)\b`,
          // 霍曼型号
          `(?i)\b(M1(?:\s*Pro)?)\b`,
          // 拓牛型号
          `(?i)\b(T1(?:S)?|T\s*Air)\b`,
          // Pro/Max/Ultra 后缀
          `\b(\w+\s*(?:Pro|Max|Ultra|Plus|Lite))\b`,
          // 数字+字母组合 (15 Pro, S24)
          `\b(\d+\s*(?:Pro|Max|Ultra|Plus))\b`,
      }
      
      for _, pattern := range patterns {
          re := regexp.MustCompile(pattern)
          if match := re.FindString(content); match != "" {
              return strings.TrimSpace(match)
          }
      }
      return ""
  }
  ```
  
  **Acceptance Criteria**:
  - [ ] `go build ./backend/...` 通过
  - [ ] 正则能匹配常见型号格式
  
  **Commit**: YES
  - Message: `feat(executor): enhance model regex patterns for cat litter boxes`
  - Files: `backend/task/executor.go`

---

### Wave 2: 前端修复

- [ ] 4. 前端统一小数点精度

  **What to do**:
  修改 `frontend/src/pages/Report.tsx`，将所有 `toFixed(2)` 改为 `toFixed(1)`
  
  **修改位置** (第 446 行):
  ```tsx
  // 修改前
  {ranking.overall_score.toFixed(2)}
  
  // 修改后
  {ranking.overall_score.toFixed(1)}
  ```
  
  **Acceptance Criteria**:
  - [ ] `npm run build` 通过
  - [ ] 所有分数显示为1位小数
  
  **Commit**: YES
  - Message: `fix(report): use consistent 1 decimal place for scores`
  - Files: `frontend/src/pages/Report.tsx`

---

### Wave 3: 验证

- [ ] 5. 重新编译后端并测试

  **What to do**:
  1. 重新编译后端
  2. 重启后端服务
  3. 运行新的分析任务
  4. 验证型号排名和小数点精度
  
  **Commands**:
  ```bash
  # 编译
  cd /Users/daozhang/Trae_AI/B站商品评论解析
  go build -o /tmp/bilibili-server ./backend
  
  # 重启
  pkill -f bilibili-server
  nohup /tmp/bilibili-server > /tmp/backend.log 2>&1 &
  
  # 测试
  curl -X POST http://localhost:8080/api/parse \
    -H "Content-Type: application/json" \
    -d '{"category":"自动猫砂盆"}'
  ```
  
  **Acceptance Criteria**:
  - [ ] 后端启动成功
  - [ ] 新任务能正确提取型号
  - [ ] 分数显示为1位小数

---

## Success Criteria

### 验证命令
```bash
# 检查报告数据
curl -s http://localhost:8080/api/report/{id} | python3 -c "
import sys, json
data = json.load(sys.stdin)['data']
print('型号数量:', len(data.get('model_rankings', [])))
print('品牌分数示例:', data['rankings'][0]['overall_score'])
"
```

### Final Checklist
- [ ] 型号排名有数据
- [ ] 所有分数显示为1位小数
- [ ] 后端编译通过
- [ ] 前端编译通过
