package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// BrandIdentifyRequest 品牌识别请求
type BrandIdentifyRequest struct {
	Models []string // 需要识别品牌的型号列表
}

// BrandIdentifyResponse 品牌识别响应
type BrandIdentifyResponse struct {
	Results map[string]string `json:"results"` // 型号 -> 品牌
}

// BrandIdentifyContext 品牌识别上下文
type BrandIdentifyContext struct {
	Category         string   // 商品类别："自动猫砂盆"
	KnownBrands      []string // 用户指定的品牌：["小佩", "CATLINK"]
	DiscoveredBrands []string // AI已识别的品牌：["霍曼", "美的"]
}

// IdentifyBrandsForModels 批量识别型号对应的品牌
// 一次性提交所有未知品牌的型号，返回型号→品牌映射
// category: 商品类别（如"耳机"、"吸尘器"、"猫砂盆"），用于生成针对性的systemPrompt
func (c *Client) IdentifyBrandsForModels(ctx context.Context, models []string, identifyCtx BrandIdentifyContext) (map[string]string, error) {
	if len(models) == 0 {
		return make(map[string]string), nil
	}

	// 去重
	uniqueModels := make([]string, 0, len(models))
	seen := make(map[string]bool)
	for _, m := range models {
		m = strings.TrimSpace(m)
		if m != "" && !seen[strings.ToLower(m)] {
			uniqueModels = append(uniqueModels, m)
			seen[strings.ToLower(m)] = true
		}
	}

	if len(uniqueModels) == 0 {
		return make(map[string]string), nil
	}

	log.Printf("[AI] 🔍 批量识别 %d 个型号的品牌 (类别: %s)...", len(uniqueModels), identifyCtx.Category)

	// 根据商品类别生成针对性的systemPrompt
	systemPrompt := buildDynamicBrandPrompt(identifyCtx)

	userPrompt := fmt.Sprintf(`请识别以下型号对应的品牌，返回JSON格式：

型号列表：
%s

返回格式示例：
{"results": {"TWS5": "OPPO", "Air 2": "小米", "V12": "戴森"}}`, strings.Join(uniqueModels, "\n"))

	// 构建消息列表
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	resp, err := c.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI调用失败: %w", err)
	}

	// 解析JSON响应
	var result BrandIdentifyResponse

	// 尝试提取JSON
	jsonStr := resp
	if idx := strings.Index(resp, "{"); idx != -1 {
		if endIdx := strings.LastIndex(resp, "}"); endIdx != -1 {
			jsonStr = resp[idx : endIdx+1]
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("[AI] 品牌识别JSON解析失败: %v, 原始响应: %s", err, resp)
		return make(map[string]string), nil
	}

	log.Printf("[AI] ✅ 品牌识别完成: %v", result.Results)
	return result.Results, nil
}

// buildDynamicBrandPrompt 根据上下文动态构建提示词
func buildDynamicBrandPrompt(ctx BrandIdentifyContext) string {
	// 构建已知品牌列表
	knownBrandsStr := "无"
	if len(ctx.KnownBrands) > 0 {
		knownBrandsStr = strings.Join(ctx.KnownBrands, "、")
	}

	// 构建已发现品牌列表
	discoveredBrandsStr := "无"
	if len(ctx.DiscoveredBrands) > 0 {
		discoveredBrandsStr = strings.Join(ctx.DiscoveredBrands, "、")
	}

	return fmt.Sprintf(`你是一个专业的【%s】产品型号识别专家。

## 任务背景
- 商品类别：%s
- 用户关注的品牌：%s
- 已识别到的同类品牌：%s

## 识别规则
1. **优先匹配**：如果型号明显属于已知品牌或已识别品牌，直接返回该品牌
2. **行业推断**：根据商品类别和已知品牌，推断该行业的其他常见品牌
3. **命名规律**：分析型号的命名规律（如前缀、系列名）来判断品牌
4. **品牌格式**：
   - 纯字母品牌用全大写（如 OPPO、CATLINK、JBL）
   - 中文品牌保持原样（如 小米、华为、小佩）
5. **无法确定**：如果确实无法判断，返回"未知"

## 重要提示
- 这是【%s】行业的型号，请在该行业范围内识别
- 同一型号在不同行业可能属于不同品牌，请根据上下文判断
- 必须严格返回JSON格式`, ctx.Category, ctx.Category, knownBrandsStr, discoveredBrandsStr, ctx.Category)
}
