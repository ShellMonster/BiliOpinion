package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// 注意：Dimension 类型已在 keyword.go 中定义，此处复用

// AnalyzeCommentRequest 分析评论请求
// 包含待分析的评论内容和评价维度列表
type AnalyzeCommentRequest struct {
	Comment    string      // 评论内容
	Dimensions []Dimension // 评价维度列表
}

// AnalyzeCommentResponse 分析评论响应
// 包含各维度的得分结果
type AnalyzeCommentResponse struct {
	Scores map[string]*float64 `json:"scores"` // 维度名 -> 得分(1-10)，nil表示未提及
}

// CommentAnalysisResult 评论分析结果（包含原始评论信息）
// 用于批量分析时返回完整的分析结果
type CommentAnalysisResult struct {
	CommentID string              `json:"comment_id"` // 评论ID
	Content   string              `json:"content"`    // 评论内容
	Scores    map[string]*float64 `json:"scores"`     // 各维度得分
	Error     string              `json:"error"`      // 分析错误信息（如有）
}

// AnalyzeComment 分析单条评论
// 使用AI对评论内容进行多维度打分
// 参数：
//   - ctx: 上下文（用于取消和超时控制）
//   - req: 分析请求，包含评论内容和评价维度
//
// 返回：
//   - *AnalyzeCommentResponse: 各维度得分结果
//   - error: 分析失败时返回错误
func (c *Client) AnalyzeComment(ctx context.Context, req AnalyzeCommentRequest) (*AnalyzeCommentResponse, error) {
	// 检查评论内容是否为空
	if strings.TrimSpace(req.Comment) == "" {
		return nil, fmt.Errorf("评论内容不能为空")
	}

	// 检查维度列表是否为空
	if len(req.Dimensions) == 0 {
		return nil, fmt.Errorf("评价维度不能为空")
	}

	// 构建维度列表字符串
	var dimList []string
	for _, dim := range req.Dimensions {
		dimList = append(dimList, fmt.Sprintf("- %s：%s", dim.Name, dim.Description))
	}

	// 构建系统提示词（System Prompt）
	// 指导AI如何分析评论并返回结构化的JSON结果
	systemPrompt := fmt.Sprintf(`你是一个专业的商品评论分析助手。你的任务是根据用户评论内容，对以下维度进行打分（1-10分）：

%s

评分标准：
- 1-3分：差评/负面评价
- 4-5分：一般/中性评价
- 6-7分：较好/正面评价
- 8-10分：优秀/强烈好评

重要规则：
1. 只根据评论中明确提及的内容打分
2. 如果评论完全未提及某个维度，该维度返回null
3. 必须严格返回JSON格式，不要添加任何其他文字

返回格式示例：
{"scores":{"维度1":8.5,"维度2":null,"维度3":7.0}}`, strings.Join(dimList, "\n"))

	// 构建用户提示词
	userPrompt := fmt.Sprintf("请分析以下评论：\n\n%s", req.Comment)

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

	// 解析AI返回的JSON响应
	result, err := parseAnalysisResponse(response)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	return result, nil
}

// parseAnalysisResponse 解析AI返回的分析响应
// 从AI返回的文本中提取JSON并解析为结构化数据
func parseAnalysisResponse(response string) (*AnalyzeCommentResponse, error) {
	// 尝试直接解析
	var result AnalyzeCommentResponse
	if err := json.Unmarshal([]byte(response), &result); err == nil {
		return &result, nil
	}

	// 如果直接解析失败，尝试从响应中提取JSON
	// AI有时会在JSON前后添加额外的文字说明
	jsonPattern := regexp.MustCompile(`\{[\s\S]*"scores"[\s\S]*\}`)
	match := jsonPattern.FindString(response)
	if match == "" {
		return nil, fmt.Errorf("无法从响应中提取JSON: %s", response)
	}

	if err := json.Unmarshal([]byte(match), &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w, 原始响应: %s", err, response)
	}

	return &result, nil
}

// AnalyzeCommentsBatch 批量分析评论（并发）
// 使用goroutine并发分析多条评论，提高处理效率
// 参数：
//   - ctx: 上下文（用于取消和超时控制）
//   - comments: 评论列表，每个元素包含评论ID和内容
//   - dimensions: 评价维度列表
//
// 返回：
//   - []CommentAnalysisResult: 所有评论的分析结果
//   - error: 如果所有评论都分析失败则返回错误
func (c *Client) AnalyzeCommentsBatch(ctx context.Context, comments []CommentInput, dimensions []Dimension) ([]CommentAnalysisResult, error) {
	// 检查输入参数
	if len(comments) == 0 {
		return nil, fmt.Errorf("评论列表不能为空")
	}
	if len(dimensions) == 0 {
		return nil, fmt.Errorf("评价维度不能为空")
	}

	// 创建结果切片
	results := make([]CommentAnalysisResult, len(comments))

	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup

	// 并发分析每条评论
	for i, comment := range comments {
		wg.Add(1)

		// 启动goroutine处理单条评论
		go func(index int, input CommentInput) {
			defer wg.Done()

			// 初始化结果
			results[index] = CommentAnalysisResult{
				CommentID: input.ID,
				Content:   input.Content,
			}

			// 调用单条评论分析
			resp, err := c.AnalyzeComment(ctx, AnalyzeCommentRequest{
				Comment:    input.Content,
				Dimensions: dimensions,
			})

			if err != nil {
				// 记录错误但不中断其他评论的分析
				results[index].Error = err.Error()
				return
			}

			// 保存分析结果
			results[index].Scores = resp.Scores
		}(i, comment)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 统计成功和失败的数量
	successCount := 0
	for _, r := range results {
		if r.Error == "" {
			successCount++
		}
	}

	// 如果所有评论都分析失败，返回错误
	if successCount == 0 {
		return results, fmt.Errorf("所有评论分析均失败")
	}

	return results, nil
}

// CommentInput 评论输入
// 用于批量分析时传入评论信息
type CommentInput struct {
	ID      string // 评论ID
	Content string // 评论内容
}

// AnalyzeCommentsWithRateLimit 带速率限制的批量分析
// 控制并发数量，避免API限流
// 参数：
//   - ctx: 上下文
//   - comments: 评论列表
//   - dimensions: 评价维度
//   - batchSize: 每批处理的评论数量
//
// 返回：
//   - []CommentAnalysisResult: 分析结果
//   - error: 错误信息
func (c *Client) AnalyzeCommentsWithRateLimit(ctx context.Context, comments []CommentInput, dimensions []Dimension, batchSize int) ([]CommentAnalysisResult, error) {
	// 设置默认批次大小
	if batchSize <= 0 {
		batchSize = 5
	}

	var allResults []CommentAnalysisResult

	// 分批处理评论
	for i := 0; i < len(comments); i += batchSize {
		// 计算当前批次的结束索引
		end := i + batchSize
		if end > len(comments) {
			end = len(comments)
		}

		// 获取当前批次的评论
		batch := comments[i:end]

		// 分析当前批次
		results, err := c.AnalyzeCommentsBatch(ctx, batch, dimensions)
		if err != nil {
			// 记录错误但继续处理下一批
			// 部分失败不影响整体流程
			for _, comment := range batch {
				allResults = append(allResults, CommentAnalysisResult{
					CommentID: comment.ID,
					Content:   comment.Content,
					Error:     err.Error(),
				})
			}
			continue
		}

		// 添加到总结果
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// RecommendationInput AI生成购买建议的输入数据
type RecommendationInput struct {
	Category      string
	Rankings      []BrandRankingInfo
	BrandAnalysis map[string]BrandStrengthWeakness
}

// BrandRankingInfo 品牌排名信息（用于AI生成建议）
type BrandRankingInfo struct {
	Brand        string
	OverallScore float64
	Rank         int
}

// BrandStrengthWeakness 品牌优劣势（用于AI生成建议）
type BrandStrengthWeakness struct {
	Strengths  []string
	Weaknesses []string
}

// GenerateRecommendation 使用AI生成专业的购买建议
func (c *Client) GenerateRecommendation(ctx context.Context, input RecommendationInput) (string, error) {
	if len(input.Rankings) == 0 {
		return "暂无足够数据生成购买建议", nil
	}

	var rankingText string
	for _, r := range input.Rankings {
		analysis := input.BrandAnalysis[r.Brand]
		rankingText += fmt.Sprintf("第%d名：%s（%.1f分）", r.Rank, r.Brand, r.OverallScore)
		if len(analysis.Strengths) > 0 {
			rankingText += fmt.Sprintf("，优势：%v", analysis.Strengths)
		}
		if len(analysis.Weaknesses) > 0 {
			rankingText += fmt.Sprintf("，劣势：%v", analysis.Weaknesses)
		}
		rankingText += "\n"
	}

	systemPrompt := `你是一位专业的商品评测专家。请根据以下品牌评分和优劣势分析，生成一段200-300字的专业购买建议。
要求：
1. 客观分析各品牌的优缺点
2. 针对不同用户需求给出具体建议
3. 语言专业但易懂
4. 不要使用markdown格式，直接输出纯文本`

	userPrompt := fmt.Sprintf("商品类别：%s\n\n品牌排名及分析：\n%s\n请生成购买建议：", input.Category, rankingText)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := c.ChatCompletion(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("AI生成建议失败: %w", err)
	}

	return strings.TrimSpace(response), nil
}
