package comment

import (
	"bilibili-analyzer/backend/bilibili"
	"testing"
)

func TestFilterAndRank_EmptyInput(t *testing.T) {
	out := FilterAndRank(nil, FilterConfig{MaxComments: 10})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}

func TestFilterAndRank_PureEmojiFiltered(t *testing.T) {
	comments := []Comment{
		{Content: bilibili.Content{Message: "😀😀😀😀😀😀😀😀😀😀"}},
	}

	out := FilterAndRank(comments, FilterConfig{MaxComments: 10, MinLength: 10, FilterEmoji: true})
	if len(out) != 0 {
		t.Fatalf("expected emoji-only comment to be filtered, got %d", len(out))
	}
}

func TestFilterAndRank_ShortCommentFiltered(t *testing.T) {
	comments := []Comment{
		{Content: bilibili.Content{Message: "太好"}},
	}

	out := FilterAndRank(comments, FilterConfig{MaxComments: 10, MinLength: 10, FilterEmoji: true})
	if len(out) != 0 {
		t.Fatalf("expected short comment to be filtered, got %d", len(out))
	}
}

func TestFilterAndRank_ValidCommentKept(t *testing.T) {
	comments := []Comment{
		{Like: 10, Count: 1, Content: bilibili.Content{Message: "这个吸尘器真的很好用，吸力很强，续航也不错"}},
	}

	out := FilterAndRank(comments, FilterConfig{MaxComments: 10, MinLength: 10, FilterEmoji: true})
	if len(out) != 1 {
		t.Fatalf("expected valid comment kept, got %d", len(out))
	}
}

func TestFilterAndRank_SortByScore(t *testing.T) {
	low := Comment{
		Like:    0,
		Count:   0,
		Content: bilibili.Content{Message: "这个东西一般般，没啥特别的地方"},
		Ctime:   1,
		RPID:    1,
	}
	high := Comment{
		Like:    5000,
		Count:   300,
		Content: bilibili.Content{Message: "戴森 Dyson 真的好用，吸力强，噪音小，清洁很彻底，推荐"},
		Ctime:   2,
		RPID:    2,
	}

	comments := []Comment{low, high}
	out := FilterAndRank(comments, FilterConfig{MaxComments: 10, MinLength: 10, FilterEmoji: true, Keywords: []string{"dyson"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 comments kept, got %d", len(out))
	}
	if out[0].RPID != high.RPID {
		t.Fatalf("expected highest score first, got rpid=%d", out[0].RPID)
	}
}

func TestFilterAndRank_LimitApplied(t *testing.T) {
	comments := []Comment{
		{RPID: 1, Like: 0, Count: 0, Content: bilibili.Content{Message: "这个产品还可以，符合预期，用起来挺顺手的"}},
		{RPID: 2, Like: 1000, Count: 20, Content: bilibili.Content{Message: "非常推荐，做工扎实，体验很好，性价比也高"}},
		{RPID: 3, Like: 2000, Count: 50, Content: bilibili.Content{Message: "用了一周感觉很棒，吸力强劲，清理很方便，续航也不错"}},
	}

	out := FilterAndRank(comments, FilterConfig{MaxComments: 2, MinLength: 10, FilterEmoji: true})
	if len(out) != 2 {
		t.Fatalf("expected limit applied to 2, got %d", len(out))
	}
}

func TestScoreComment_KeywordCaseInsensitive(t *testing.T) {
	c := Comment{Like: 0, Count: 0, Content: bilibili.Content{Message: "I like Dyson vacuum cleaners a lot"}}

	noKW := scoreComment(c, nil)
	withKW := scoreComment(c, []string{"dyson"})
	// 关键词匹配应至少增加 10 分（不依赖长度分/热度分）。
	if withKW-noKW < 9.999 {
		t.Fatalf("expected keyword score applied, got noKW=%v withKW=%v", noKW, withKW)
	}
}

func TestIsValidComment_WhitespaceOnly(t *testing.T) {
	c := Comment{Content: bilibili.Content{Message: "   \n\t  "}}

	if isValidComment(c, 10, true) {
		t.Fatalf("expected whitespace-only comment to be invalid")
	}
}
