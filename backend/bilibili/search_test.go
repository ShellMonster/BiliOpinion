package bilibili

import "testing"

func TestStripHTMLTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "B站标准高亮标签",
			input:    "普通人用了一年半的<em class=\"keyword\">吸尘器</em>测评",
			expected: "普通人用了一年半的吸尘器测评",
		},
		{
			name:     "多个标签",
			input:    "<em class=\"keyword\">戴森</em>V12 vs <em class=\"keyword\">小米</em>G10",
			expected: "戴森V12 vs 小米G10",
		},
		{
			name:     "无标签",
			input:    "没有标签的普通标题",
			expected: "没有标签的普通标题",
		},
		{
			name:     "其他HTML标签",
			input:    "<b>加粗</b>和<i>斜体</i>",
			expected: "加粗和斜体",
		},
		{
			name:     "自闭合标签",
			input:    "图片<img src=\"test.jpg\"/>在这里",
			expected: "图片在这里",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有标签",
			input:    "<em></em>",
			expected: "",
		},
		{
			name:     "嵌套标签",
			input:    "<div><span>内容</span></div>",
			expected: "内容",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripHTMLTags(tt.input)
			if result != tt.expected {
				t.Errorf("stripHTMLTags(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
