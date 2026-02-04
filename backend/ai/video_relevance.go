package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// VideoRelevanceChecker 视频相关性检查器
// 用于判断视频标题是否与用户需求相关，避免无关视频污染分析结果
type VideoRelevanceChecker struct {
	client *Client // AI客户端，用于调用大模型进行相关性判断
}

// NewVideoRelevanceChecker 创建视频相关性检查器
//
// 参数:
//   - client: AI客户端实例
//
// 返回:
//   - *VideoRelevanceChecker: 检查器实例
func NewVideoRelevanceChecker(client *Client) *VideoRelevanceChecker {
	return &VideoRelevanceChecker{client: client}
}

// CheckRelevanceRequest 相关性检查请求
// 用于封装检查所需的输入参数
type CheckRelevanceRequest struct {
	VideoTitle      string `json:"video_title"`      // 视频标题
	UserRequirement string `json:"user_requirement"` // 用户需求描述
}

// CheckRelevanceResponse 相关性检查响应
// AI返回的判断结果
type CheckRelevanceResponse struct {
	IsRelevant bool   `json:"is_relevant"` // 是否相关
	Reason     string `json:"reason"`      // 判断理由
}

// CheckRelevance 检查视频标题是否与用户需求相关
//
// 工作原理：
//  1. 构造系统提示词，定义判断标准和输出格式
//  2. 将视频标题和用户需求作为用户输入
//  3. 调用AI模型进行判断
//  4. 解析JSON响应，提取相关性和理由
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - videoTitle: 视频标题
//   - userRequirement: 用户需求/任务描述
//
// 返回:
//   - isRelevant: 是否相关（true=相关，false=不相关）
//   - reason: 判断理由（用于日志和调试）
//   - err: 错误信息（AI调用失败时返回）
func (c *VideoRelevanceChecker) CheckRelevance(
	ctx context.Context,
	videoTitle string,
	userRequirement string,
) (isRelevant bool, reason string, err error) {
	// 参数校验：标题或需求为空时直接返回不相关
	if videoTitle == "" || userRequirement == "" {
		return false, "视频标题或用户需求为空", nil
	}

	// 系统提示词：定义AI的角色、任务和判断标准
	systemPrompt := `你是一个视频内容相关性判断助手。

你的任务是判断视频标题是否与用户需求相关。

判断标准：
1. 视频标题应该直接涉及用户需求的主题
2. 如果视频标题涉及的是相关但不同的产品/主题，视为不相关
3. 考虑同义词和相关概念

输出格式（JSON）：
{"is_relevant": true/false, "reason": "判断理由"}

示例：
- 需求"猫砂盆"，标题"小佩猫砂盆测评" → {"is_relevant": true, "reason": "直接涉及猫砂盆产品"}
- 需求"猫砂盆"，标题"猫咪喂食器推荐" → {"is_relevant": false, "reason": "涉及的是喂食器而非猫砂盆"}
- 需求"吸尘器"，标题"戴森V12吸尘器测评" → {"is_relevant": true, "reason": "直接涉及吸尘器产品"}
- 需求"吸尘器"，标题"扫地机器人推荐" → {"is_relevant": false, "reason": "涉及的是扫地机器人而非吸尘器"}`

	// 用户输入：包含需求和视频标题
	userPrompt := fmt.Sprintf("用户需求：%s\n视频标题：%s", userRequirement, videoTitle)

	// 构造消息列表
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// 调用AI模型
	response, err := c.client.ChatCompletion(ctx, messages)
	if err != nil {
		return false, "", fmt.Errorf("AI请求失败: %w", err)
	}

	// 解析JSON响应
	var result CheckRelevanceResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// 如果解析失败，尝试从文本中提取（容错处理）
		// 检查响应中是否包含"true"或"相关"关键词
		isRelevant = strings.Contains(response, "true") || strings.Contains(strings.ToLower(response), "相关")
		reason = response
		return isRelevant, reason, nil
	}

	return result.IsRelevant, result.Reason, nil
}

// BatchCheckRelevance 批量检查视频相关性
//
// 工作原理：
//  1. 遍历所有视频标题
//  2. 对每个标题调用CheckRelevance进行判断
//  3. 收集相关视频的索引和不相关视频的信息
//  4. 返回过滤结果
//
// 参数:
//   - ctx: 上下文
//   - videoTitles: 视频标题列表
//   - userRequirement: 用户需求
//   - concurrency: 并发数（当前版本未使用，预留用于未来优化）
//
// 返回:
//   - relevantIndices: 相关视频的索引列表（用于过滤原始视频列表）
//   - irrelevantVideos: 不相关视频信息列表（包含标题和理由，用于日志）
//   - err: 错误信息
func (c *VideoRelevanceChecker) BatchCheckRelevance(
	ctx context.Context,
	videoTitles []string,
	userRequirement string,
	concurrency int,
) (relevantIndices []int, irrelevantVideos []map[string]string, err error) {
	// 并发数校验（预留参数，当前版本串行处理）
	if concurrency <= 0 {
		concurrency = 5
	}

	// 遍历所有视频标题
	for i, title := range videoTitles {
		// 检查单个视频的相关性
		isRelevant, reason, checkErr := c.CheckRelevance(ctx, title, userRequirement)
		if checkErr != nil {
			// 出错时默认保留该视频（避免误删）
			relevantIndices = append(relevantIndices, i)
			continue
		}

		if isRelevant {
			// 相关视频：记录索引
			relevantIndices = append(relevantIndices, i)
		} else {
			// 不相关视频：记录详细信息用于日志
			irrelevantVideos = append(irrelevantVideos, map[string]string{
				"title":  title,
				"reason": reason,
			})
		}
	}

	return relevantIndices, irrelevantVideos, nil
}
