package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ParseKeywordRequest 解析关键词请求
// 用户用自然语言描述购买需求，AI会帮助解析出品牌、维度和搜索关键词
type ParseKeywordRequest struct {
	Requirement string // 用户的需求描述（自由格式），例如"想买个吸尘器，预算2000，家里有宠物"
}

// ParseKeywordResponse 解析关键词响应
// AI返回的完整解析结果，包含需求理解、商品类型、预算、场景、特殊需求、品牌列表、评价维度和搜索关键词
type ParseKeywordResponse struct {
	Understanding string      `json:"understanding"`           // AI对用户需求的理解描述（用"我理解您..."开头）
	ProductType   string      `json:"product_type"`            // 识别出的商品类型
	Budget        string      `json:"budget,omitempty"`        // 预算范围（可选）
	Scenario      string      `json:"scenario,omitempty"`      // 使用场景（可选）
	SpecialNeeds  []string    `json:"special_needs,omitempty"` // 特殊需求（可选）
	Brands        []string    `json:"brands"`                  // 品牌列表，该类目的主流品牌
	Dimensions    []Dimension `json:"dimensions"`              // 评价维度，用于后续分析的维度
	Keywords      []string    `json:"keywords"`                // 搜索关键词，用于在B站搜索相关视频
}

// Dimension 评价维度
// 每个维度包含名称和描述，用于指导AI分析评论
type Dimension struct {
	Name        string `json:"name"`        // 维度名称，例如"吸力性能"
	Description string `json:"description"` // 维度描述，例如"评估吸尘器的吸力大小和清洁效果"
}

// cleanJSONResponse 清理AI返回的JSON响应
// 移除Markdown代码块标记（```json 和 ```）
func cleanJSONResponse(response string) string {
	// 去除首尾空白
	response = strings.TrimSpace(response)

	// 移除 ```json 或 ``` 开头
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
	}

	// 移除 ``` 结尾
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	// 再次去除首尾空白
	response = strings.TrimSpace(response)

	// 如果还有问题，尝试用正则提取JSON对象
	if !strings.HasPrefix(response, "{") {
		// 尝试提取 {...} 部分
		re := regexp.MustCompile(`(?s)\{.*\}`)
		if match := re.FindString(response); match != "" {
			response = match
		}
	}

	return response
}

// ParseKeyword 解析关键词
// 这是核心方法，调用AI来解析用户输入的商品类目
func (c *Client) ParseKeyword(ctx context.Context, req ParseKeywordRequest) (*ParseKeywordResponse, error) {
	// 构建系统提示词，告诉AI它的角色和任务
	systemPrompt := `你是一个商品分析助手。用户会用自然语言描述他们的购买需求，你需要：

1. 理解用户的真实意图，提取关键信息：
   - 商品类型（必须）
   - 预算范围（如果提到）
   - 使用场景（如果提到）
   - 特殊需求（如果提到）

2. 用通俗易懂的语言描述你的理解（用"我理解您..."开头）

3. 根据用户需求推荐5个左右的主流品牌（按市场份额和用户需求匹配度排序）

4. 提出6个针对性的评价维度（根据商品特点和用户特殊需求调整）

5. 生成B站搜索关键词，包含两类：
   a) 品牌特定关键词（3-5个）：每个主流品牌的"品牌名+商品类型"组合
      例如："戴森吸尘器"、"小米吸尘器"
   
   b) 通用发现关键词（4个）：不包含品牌名，用于发现市场上所有品牌
      必须包含以下4种类型：
      - "商品类型+评测"（如"自动猫砂盆评测"）
      - "商品类型+推荐"（如"自动猫砂盆推荐"）
      - "商品类型+横评"（如"自动猫砂盆横评"）
      - "商品类型+对比"（如"自动猫砂盆对比"）

直接返回JSON格式，不要使用Markdown代码块：
{
  "understanding": "我理解您想购买...",
  "product_type": "商品类型",
  "budget": "预算范围（如果用户提到）",
  "scenario": "使用场景（如果用户提到）",
  "special_needs": ["特殊需求1", "特殊需求2"],
  "brands": ["品牌1", "品牌2", "品牌3", "品牌4", "品牌5"],
  "dimensions": [
    {"name": "维度名", "description": "维度说明（结合用户需求）"}
  ],
  "keywords": [
    "品牌1+商品类型",
    "品牌2+商品类型",
    "品牌3+商品类型",
    "商品类型+评测",
    "商品类型+推荐",
    "商品类型+横评",
    "商品类型+对比"
  ]
}

注意：
- 如果用户没有提到预算/场景/特殊需求，对应字段可以为空或省略
- 品牌名称要准确，使用官方中文名
- 维度要针对用户的特殊需求调整（如用户提到宠物，维度描述要体现对宠物毛发的处理能力）
- 品牌特定关键词：结合用户需求生成（如预算、场景等）
- 通用发现关键词：必须包含"评测"、"推荐"、"横评"、"对比"这4种类型，用于发现所有品牌
- 重要：直接返回JSON，不要用代码块包裹`

	// 构建用户提示词，传入用户输入的商品类目
	userPrompt := fmt.Sprintf("用户需求：%s", req.Requirement)

	// 调用AI接口
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// 发送请求到AI服务
	response, err := c.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI请求失败: %w", err)
	}

	// 清理AI返回的响应（移除Markdown代码块标记）
	cleanedResponse := cleanJSONResponse(response)

	// 解析AI返回的JSON响应
	var result ParseKeywordResponse
	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		// 如果解析失败，记录原始响应以便调试
		return nil, fmt.Errorf("解析AI响应失败: %w (原始响应: %s)", err, response[:min(len(response), 200)])
	}

	// 验证结果的完整性
	if result.Understanding == "" {
		return nil, fmt.Errorf("AI未返回需求理解描述")
	}
	if result.ProductType == "" {
		return nil, fmt.Errorf("AI未返回商品类型")
	}
	if len(result.Brands) == 0 {
		return nil, fmt.Errorf("AI未返回品牌列表")
	}
	if len(result.Dimensions) == 0 {
		return nil, fmt.Errorf("AI未返回评价维度")
	}
	if len(result.Keywords) == 0 {
		return nil, fmt.Errorf("AI未返回搜索关键词")
	}

	return &result, nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
