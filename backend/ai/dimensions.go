package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GenerateDimensions 根据视频信息和评论样本生成评价维度
// 使用AI分析视频标题、简介和评论内容，生成最适合用于分析该视频评论的评价维度
//
// 参数：
//   - ctx: 上下文（用于取消和超时控制）
//   - videoTitle: 视频标题
//   - videoDesc: 视频简介/描述
//   - comments: 评论样本列表（用于了解用户讨论的话题）
//
// 返回：
//   - []Dimension: 生成的评价维度列表（6个维度）
//   - error: 生成失败时返回错误
func (c *Client) GenerateDimensions(ctx context.Context, videoTitle, videoDesc string, comments []string) ([]Dimension, error) {
	// 构建评论样本文本（限制数量，避免内容过长）
	maxComments := 20 // 最多使用20条评论作为样本
	if len(comments) > maxComments {
		comments = comments[:maxComments]
	}

	// 将评论列表转换为文本
	commentsText := strings.Join(comments, "\n")

	// 构建系统提示词
	systemPrompt := `你是一个专业的视频评论分析师。根据以下视频信息和用户评论样本，生成6个最适合用于分析该视频评论的评价维度。

【要求】
1. 维度要反映用户在评论中实际讨论的话题
2. 每个维度名称控制在4-6个字
3. 每个维度需要有简要描述（10-20字）
4. 必须返回JSON格式，不要添加任何其他文字

【返回格式】返回JSON数组：
[{"name": "维度名称", "description": "维度描述"}, ...]`

	// 构建用户提示词
	userPrompt := fmt.Sprintf(`【视频信息】
标题：%s
简介：%s

【评论样本】
%s`, videoTitle, videoDesc, commentsText)

	// 构建消息列表
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// 调用AI接口
	response, err := c.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI请求失败: %w", err)
	}

	// 清理AI返回的响应（移除Markdown代码块标记）
	cleanedResponse := cleanJSONResponse(response)

	// 尝试解析JSON数组
	var dimensions []Dimension
	if err := json.Unmarshal([]byte(cleanedResponse), &dimensions); err != nil {
		// 如果直接解析失败，尝试从响应中提取JSON数组
		extracted := extractJSONArray(cleanedResponse)
		if extracted == "" {
			return nil, fmt.Errorf("无法从响应中提取JSON数组: %s", response[:min(len(response), 200)])
		}

		if err := json.Unmarshal([]byte(extracted), &dimensions); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w, 原始响应: %s", err, response[:min(len(response), 200)])
		}
	}

	// 验证结果
	if len(dimensions) == 0 {
		return nil, fmt.Errorf("AI未返回任何评价维度")
	}

	// 确保维度名称和描述不为空
	for i := range dimensions {
		if dimensions[i].Name == "" {
			return nil, fmt.Errorf("第%d个维度名称为空", i+1)
		}
		if dimensions[i].Description == "" {
			return nil, fmt.Errorf("第%d个维度描述为空", i+1)
		}
	}

	return dimensions, nil
}

// extractJSONArray 从文本中提取JSON数组
// 用于处理AI返回的包含额外文本的情况
func extractJSONArray(text string) string {
	// 查找第一个 [ 和最后一个 ]
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")

	if start == -1 || end == -1 || start >= end {
		return ""
	}

	return text[start : end+1]
}
