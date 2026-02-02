# 测试HTML标签过滤功能

## TL;DR

> **Quick Summary**: 验证 `stripHTMLTags` 函数能正确过滤B站搜索API返回的HTML标签
> 
> **Deliverables**:
> - 单元测试文件 `backend/bilibili/search_test.go`
> - 验证各种HTML标签场景
> 
> **Estimated Effort**: Quick (15分钟)
> **Parallel Execution**: NO

---

## Context

### 问题背景
B站搜索API返回的视频标题包含HTML高亮标签：
```
普通人用了一年半的<em class="keyword">吸尘器</em>测评
```

### 当前实现
使用正则表达式 `<[^>]*>` 匹配并移除所有HTML标签：
```go
func stripHTMLTags(s string) string {
    re := regexp.MustCompile(`<[^>]*>`)
    return re.ReplaceAllString(s, "")
}
```

---

## Work Objectives

### Core Objective
验证 `stripHTMLTags` 函数能正确处理各种HTML标签场景。

### 测试用例

| 输入 | 期望输出 | 说明 |
|------|----------|------|
| `普通人用了一年半的<em class="keyword">吸尘器</em>测评` | `普通人用了一年半的吸尘器测评` | 标准B站高亮标签 |
| `<em class="keyword">戴森</em>V12 vs <em class="keyword">小米</em>G10` | `戴森V12 vs 小米G10` | 多个标签 |
| `没有标签的普通标题` | `没有标签的普通标题` | 无标签 |
| `<b>加粗</b>和<i>斜体</i>` | `加粗和斜体` | 其他HTML标签 |
| `<script>alert('xss')</script>正常内容` | `alert('xss')正常内容` | 安全测试 |
| `空字符串` | `空字符串` | 边界情况 |

---

## TODOs

- [ ] 1. 创建单元测试文件验证stripHTMLTags函数

  **What to do**:
  - 创建 `backend/bilibili/search_test.go`
  - 编写 `TestStripHTMLTags` 测试函数
  - 覆盖上述所有测试用例
  - 运行测试验证结果

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **References**:
  - `backend/bilibili/search.go:234-240` - stripHTMLTags 函数定义

  **Acceptance Criteria**:

  ```bash
  # AC1: 单元测试全部通过
  cd backend && go test ./bilibili/... -v -run TestStripHTMLTags
  # Expected: 所有测试用例通过

  # AC2: 验证B站标准格式
  # 输入: 普通人用了一年半的<em class="keyword">吸尘器</em>测评
  # 输出: 普通人用了一年半的吸尘器测评
  ```

  **Commit**: YES
  - Message: `test(bilibili): add unit tests for stripHTMLTags function`
  - Files: `backend/bilibili/search_test.go`

---

## Success Criteria

### Verification Commands
```bash
# 运行测试
cd backend && go test ./bilibili/... -v -run TestStripHTMLTags

# 检查覆盖率
cd backend && go test ./bilibili/... -cover -run TestStripHTMLTags
```

### Final Checklist
- [ ] 所有测试用例通过
- [ ] B站标准高亮标签正确过滤
- [ ] 多标签场景正确处理
- [ ] 无标签输入不受影响
- [ ] 边界情况正确处理
