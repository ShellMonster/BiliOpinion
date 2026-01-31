package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

// ParseKeywordRequest 解析关键词请求
// 用户输入商品类目，AI会帮助解析出品牌、维度和搜索关键词
type ParseKeywordRequest struct {
	Category string // 用户输入的商品类目，例如"无线吸尘器"
}

// ParseKeywordResponse 解析关键词响应
// AI返回的完整解析结果，包含澄清描述、品牌列表、评价维度和搜索关键词
type ParseKeywordResponse struct {
	Clarification string      `json:"clarification"` // 澄清描述，帮助用户确认理解是否正确
	Brands        []string    `json:"brands"`        // 品牌列表，该类目的主流品牌
	Dimensions    []Dimension `json:"dimensions"`    // 评价维度，用于后续分析的6个维度
	Keywords      []string    `json:"keywords"`      // 搜索关键词，用于在B站搜索相关视频
}

// Dimension 评价维度
// 每个维度包含名称和描述，用于指导AI分析评论
type Dimension struct {
	Name        string `json:"name"`        // 维度名称，例如"吸力性能"
	Description string `json:"description"` // 维度描述，例如"评估吸尘器的吸力大小和清洁效果"
}

// ParseKeyword 解析关键词
// 这是核心方法，调用AI来解析用户输入的商品类目
func (c *Client) ParseKeyword(ctx context.Context, req ParseKeywordRequest) (*ParseKeywordResponse, error) {
	// 构建系统提示词，告诉AI它的角色和任务
	systemPrompt := `你是一个商品分析助手。用户会输入一个商品类目，你需要：
1. 理解用户需求并给出澄清描述（用"您想了解..."开头）
2. 列出该类目的主流品牌（5个左右，按市场份额排序）
3. 提出6个评价维度（如性能、价格、外观、噪音、续航、易用性等，根据商品特点调整）
4. 生成3-5个B站搜索关键词（包含商品类目+评测/测评/开箱等词）

返回JSON格式：
{
  "clarification": "您想了解...",
  "brands": ["品牌1", "品牌2", "品牌3", "品牌4", "品牌5"],
  "dimensions": [
    {"name": "维度名", "description": "维度说明"}
  ],
  "keywords": ["关键词1", "关键词2", "关键词3"]
}

注意：
- 品牌名称要准确，使用官方中文名
- 维度要具体且可量化，避免模糊概念
- 关键词要符合B站用户搜索习惯`

	// 构建用户提示词，传入用户输入的商品类目
	userPrompt := fmt.Sprintf("商品类目：%s", req.Category)

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

	// 解析AI返回的JSON响应
	var result ParseKeywordResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 验证结果的完整性
	if len(result.Brands) == 0 {
		return nil, fmt.Errorf("AI未返回品牌列表")
	}
	if len(result.Dimensions) != 6 {
		return nil, fmt.Errorf("期望6个评价维度，实际得到%d个", len(result.Dimensions))
	}
	if len(result.Keywords) == 0 {
		return nil, fmt.Errorf("AI未返回搜索关键词")
	}

	return &result, nil
}
