package ai

import (
	"strings"
	"testing"
)

// TestBuildDynamicBrandPrompt 测试动态品牌提示词生成
func TestBuildDynamicBrandPrompt(t *testing.T) {
	tests := []struct {
		name         string
		ctx          BrandIdentifyContext
		wantContains []string // 期望提示词中包含的内容
	}{
		{
			name: "猫砂盆场景 - 完整上下文",
			ctx: BrandIdentifyContext{
				Category:         "自动猫砂盆",
				KnownBrands:      []string{"小佩", "CATLINK"},
				DiscoveredBrands: []string{"霍曼", "美的"},
			},
			wantContains: []string{
				"自动猫砂盆",
				"小佩、CATLINK",
				"霍曼、美的",
				"优先匹配",
				"行业推断",
				"命名规律",
			},
		},
		{
			name: "耳机场景 - 只有已知品牌",
			ctx: BrandIdentifyContext{
				Category:         "无线耳机",
				KnownBrands:      []string{"OPPO", "小米", "华为"},
				DiscoveredBrands: []string{},
			},
			wantContains: []string{
				"无线耳机",
				"OPPO、小米、华为",
				"无", // 已发现品牌为"无"
			},
		},
		{
			name: "吸尘器场景 - 只有已发现品牌",
			ctx: BrandIdentifyContext{
				Category:         "无线吸尘器",
				KnownBrands:      []string{},
				DiscoveredBrands: []string{"戴森", "石头", "追觅"},
			},
			wantContains: []string{
				"无线吸尘器",
				"无", // 已知品牌为"无"
				"戴森、石头、追觅",
			},
		},
		{
			name: "通用场景 - 无品牌信息",
			ctx: BrandIdentifyContext{
				Category:         "智能手表",
				KnownBrands:      []string{},
				DiscoveredBrands: []string{},
			},
			wantContains: []string{
				"智能手表",
				"无", // 两个都是"无"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildDynamicBrandPrompt(tt.ctx)

			// 验证提示词包含期望的内容
			for _, want := range tt.wantContains {
				if !strings.Contains(prompt, want) {
					t.Errorf("buildDynamicBrandPrompt() 生成的提示词不包含 %q\n生成的提示词:\n%s", want, prompt)
				}
			}

			// 验证提示词结构完整性
			if !strings.Contains(prompt, "## 任务背景") {
				t.Error("提示词缺少 '## 任务背景' 部分")
			}
			if !strings.Contains(prompt, "## 识别规则") {
				t.Error("提示词缺少 '## 识别规则' 部分")
			}
			if !strings.Contains(prompt, "## 重要提示") {
				t.Error("提示词缺少 '## 重要提示' 部分")
			}
		})
	}
}

// TestBuildDynamicBrandPromptFormat 测试品牌格式规则
func TestBuildDynamicBrandPromptFormat(t *testing.T) {
	ctx := BrandIdentifyContext{
		Category:         "测试类别",
		KnownBrands:      []string{"品牌A", "BRAND_B"},
		DiscoveredBrands: []string{"品牌C"},
	}

	prompt := buildDynamicBrandPrompt(ctx)

	// 验证品牌格式规则说明
	if !strings.Contains(prompt, "纯字母品牌用全大写") {
		t.Error("提示词缺少纯字母品牌格式说明")
	}
	if !strings.Contains(prompt, "中文品牌保持原样") {
		t.Error("提示词缺少中文品牌格式说明")
	}
}

// TestBrandIdentifyContextEmpty 测试空上下文处理
func TestBrandIdentifyContextEmpty(t *testing.T) {
	ctx := BrandIdentifyContext{
		Category:         "",
		KnownBrands:      nil,
		DiscoveredBrands: nil,
	}

	prompt := buildDynamicBrandPrompt(ctx)

	// 即使类别为空，也应该生成有效的提示词
	if !strings.Contains(prompt, "【】") && !strings.Contains(prompt, "【") {
		t.Log("警告: 类别为空时提示词可能不完整")
	}

	// 验证空品牌列表显示为"无"
	if !strings.Contains(prompt, "无") {
		t.Error("空品牌列表应该显示为'无'")
	}
}

// TestBrandIdentifyContextSpecialChars 测试特殊字符处理
func TestBrandIdentifyContextSpecialChars(t *testing.T) {
	ctx := BrandIdentifyContext{
		Category:         "智能家电/空调",
		KnownBrands:      []string{"美的", "格力", "Haier"},
		DiscoveredBrands: []string{"小米", "华凌"},
	}

	prompt := buildDynamicBrandPrompt(ctx)

	// 验证特殊字符正确处理
	if !strings.Contains(prompt, "智能家电/空调") {
		t.Error("提示词应该包含带斜杠的类别名称")
	}

	// 验证品牌正确连接
	if !strings.Contains(prompt, "美的、格力、Haier") {
		t.Error("品牌应该用顿号连接")
	}
}
