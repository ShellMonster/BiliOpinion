package comment

import (
	"bilibili-analyzer/backend/bilibili"
	"math"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Comment 复用 bilibili.Comment，避免在 comment 包里重复定义。
type Comment = bilibili.Comment

// FilterConfig 过滤与排序配置。
// 说明：MaxComments 与任务配置中的 MaxComments 含义保持一致。
type FilterConfig struct {
	MaxComments int
	Keywords    []string

	// MinLength 最小有效字符数（按 rune 计数），默认 10。
	MinLength int
	// FilterEmoji 是否启用“纯表情/符号”过滤，默认 true。
	FilterEmoji bool
}

// FilterAndRank 按规则过滤评论，并按质量分数降序排序，最后截断到 MaxComments。
// 注意：不会修改输入切片 comments。
func FilterAndRank(comments []Comment, config FilterConfig) []Comment {
	if len(comments) == 0 {
		return []Comment{}
	}

	minLen := config.MinLength
	if minLen <= 0 {
		minLen = 10
	}
	// 需求要求必须过滤纯表情/符号评论；因此这里默认开启。
	// 注意：即使 config.FilterEmoji 为零值 false，也会启用过滤。
	filterEmoji := true
	if config.FilterEmoji {
		filterEmoji = true
	}

	// 评分阶段：先过滤，再计算分数。
	type scored struct {
		c     Comment
		score float64
	}
	kept := make([]scored, 0, len(comments))
	for _, c := range comments {
		if !isValidComment(c, minLen, filterEmoji) {
			continue
		}
		kept = append(kept, scored{c: c, score: scoreComment(c, config.Keywords)})
	}

	if len(kept) == 0 {
		return []Comment{}
	}

	// 排序：分数降序；分数相同用时间、点赞做稳定的确定性排序。
	sort.SliceStable(kept, func(i, j int) bool {
		if kept[i].score != kept[j].score {
			return kept[i].score > kept[j].score
		}
		if kept[i].c.Ctime != kept[j].c.Ctime {
			return kept[i].c.Ctime > kept[j].c.Ctime
		}
		if kept[i].c.Like != kept[j].c.Like {
			return kept[i].c.Like > kept[j].c.Like
		}
		return kept[i].c.RPID > kept[j].c.RPID
	})

	max := config.MaxComments
	if max > 0 && len(kept) > max {
		kept = kept[:max]
	}

	out := make([]Comment, 0, len(kept))
	for _, s := range kept {
		out = append(out, s.c)
	}
	return out
}

// scoreComment 计算单条评论质量分（0-100）。
// 总分 = 热度(0-40) + 长度(0-30) + 关键词(0-30)
func scoreComment(c Comment, keywords []string) float64 {
	msg := strings.TrimSpace(c.Content.Message)
	charCount := utf8.RuneCountInString(msg)

	// 热度分（0-40）：Like 与 Count（回复数）
	likeScore := math.Min(float64(c.Like)/100.0, 20)
	replyScore := math.Min(float64(c.Count)/10.0, 20)
	popularity := likeScore + replyScore

	// 长度分（0-30）：按 rune 计数
	lengthScore := math.Min(float64(charCount)/10.0, 30)

	keywordScore := 0.0
	if len(keywords) > 0 {
		lowerMsg := strings.ToLower(msg)
		for _, kw := range keywords {
			kw = strings.TrimSpace(kw)
			if kw == "" {
				continue
			}
			if strings.Contains(lowerMsg, strings.ToLower(kw)) {
				keywordScore += 10
				if keywordScore >= 30 {
					keywordScore = 30
					break
				}
			}
		}
	}

	total := popularity + lengthScore + keywordScore
	if total < 0 {
		return 0
	}
	if total > 100 {
		return 100
	}
	return total
}

// isValidComment 判断评论是否满足最小长度与“纯表情/符号”过滤规则。
func isValidComment(c Comment, minLength int, filterEmoji bool) bool {
	msg := strings.TrimSpace(c.Content.Message)
	if utf8.RuneCountInString(msg) < minLength {
		return false
	}

	if !filterEmoji {
		return true
	}

	// 移除表情与符号后，如果可读内容不足 minLength，则视为“纯表情/符号”。
	clean := strings.TrimSpace(removeEmojiAndSymbols(msg))
	if utf8.RuneCountInString(clean) < minLength {
		return false
	}
	return true
}

// removeEmojiAndSymbols 移除 emoji 与各类符号/标点，保留字母、数字与空白。
// 目的：用于识别“纯表情/符号评论”，不是做通用的文本清洗。
func removeEmojiAndSymbols(text string) string {
	var b strings.Builder
	// 预估：ASCII/UTF-8 情况下 Grow 近似即可。
	b.Grow(len(text))
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
			continue
		}
		// 其余全部视为符号/表情/标点，直接丢弃。
	}
	return b.String()
}
