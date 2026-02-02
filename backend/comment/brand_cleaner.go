package comment

import (
	"strings"
)

// CleanBrandName 清洗品牌字段。
//
// 规则：
// - 去除首尾空白
// - 空字符串或“未知”原样返回
// - 不包含“/”直接返回
// - 包含“/”时，优先返回匹配 knownBrands 的部分（忽略大小写）；否则返回第一个分段
func CleanBrandName(brand string, knownBrands []string) string {
	brand = strings.TrimSpace(brand)
	if brand == "" || brand == "未知" {
		return brand
	}

	// 如果不包含"/"，直接返回
	if !strings.Contains(brand, "/") {
		return brand
	}

	// 拆分品牌
	parts := strings.Split(brand, "/")

	// 优先返回匹配已知品牌的部分
	for _, part := range parts {
		part = strings.TrimSpace(part)
		for _, known := range knownBrands {
			if strings.EqualFold(part, known) {
				return known // 返回已知品牌的标准名称
			}
		}
	}

	// 都不匹配，返回第一个
	return strings.TrimSpace(parts[0])
}

// CleanModelName 清洗型号字段。
//
// 规则：
// - 去除首尾空白
// - 如果包含“/”，取第一个分段
// - 描述性文案（如“新款/旧款”等）统一映射为“通用”
func CleanModelName(model string) string {
	model = strings.TrimSpace(model)

	// 如果包含"/"，取第一个
	if strings.Contains(model, "/") {
		parts := strings.Split(model, "/")
		model = strings.TrimSpace(parts[0])
	}

	// 过滤掉描述性文字，保留"通用"
	descriptive := []string{"新款", "旧款", "基础款", "升级款", "标准版"}
	for _, d := range descriptive {
		if model == d {
			return "通用"
		}
	}

	return model
}
