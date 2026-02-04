package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

// 注意：Dimension 类型已在 keyword.go 中定义，此处复用

// ProgressCallback 进度回调函数类型
// 参数：
//   - stage: 当前阶段
//   - current: 当前进度
//   - total: 总数
//   - message: 状态消息
type ProgressCallback func(stage string, current, total int, message string)

// AnalyzeCommentRequest 分析评论请求
// 包含待分析的评论内容和评价维度列表
type AnalyzeCommentRequest struct {
	Comment    string      // 评论内容
	Dimensions []Dimension // 评价维度列表
	VideoTitle string      // 视频标题，作为上下文
}

// AnalyzeCommentResponse 分析评论响应
// 包含各维度的得分结果
type AnalyzeCommentResponse struct {
	Scores map[string]*float64 `json:"scores"` // 维度名 -> 得分(1-10)，nil表示未提及
	Brand  string              `json:"brand"`  // 提取的品牌名称
	Model  string              `json:"model"`  // 提取的具体型号
}

// CommentAnalysisResult 评论分析结果（包含原始评论信息）
// 用于批量分析时返回完整的分析结果
type CommentAnalysisResult struct {
	CommentID string              `json:"comment_id"` // 评论ID
	Content   string              `json:"content"`    // 评论内容
	Scores    map[string]*float64 `json:"scores"`     // 各维度得分
	Brand     string              `json:"brand"`      // AI提取的品牌
	Model     string              `json:"model"`      // AI提取的型号
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

	systemPrompt := fmt.Sprintf(`你是一个专业的商品评论分析助手。你的任务是：

1. 从视频标题和评论内容中识别：
   - 品牌名称（如"戴森"、"小米"、"苹果"、"Sony"）
   - 具体型号（如"V12"、"iPhone 15 Pro"、"G10"、"WH-1000XM5"）

重要：型号提取规则：
- 优先从评论内容中提取具体型号（评论比标题更准确）
- 常见型号格式：
  * 字母+数字组合（V12, G10, S23, XM5）
  * 品牌+型号（iPhone 15, Galaxy S23, Watch GT3）
  * 系列+后缀（Pro, Max, Plus, Ultra, Lite）
- 如果评论提到多个型号，选择评论主要讨论的那个
- 如果无法确定具体型号但能确定系列，填写系列名（如"V系列"、"Pro系列"）
- 注意区分型号和代数（"第二代"不是型号，"V2"才是）

2. 对以下维度进行打分（1-10分）：
%s

评分标准：
- 1-3分：差评/负面评价
- 4-5分：一般/中性评价
- 6-7分：较好/正面评价
- 8-10分：优秀/强烈好评

重要规则：
1. **必须从评论内容中提取品牌**（视频标题仅供参考上下文）
2. 如果评论中没有明确提及任何品牌，brand字段必须填"未知"
3. 品牌必须是单一品牌名称，绝对不能包含"/"或其他分隔符
4. 如果评论对比多个品牌（如"A比B好"），只提取评论主要评价的那个品牌
5. 如果评论同时讨论多个品牌且无法确定主要品牌，brand填"未知"
6. 型号必须是具体型号名（如"V12"、"Max"、"Pro"），不能是描述性文字（如"新款"、"基础款"）
7. 如果无法确定型号，model字段填"通用"
8. 如果无法确定品牌，brand字段填"未知"
9. 只根据评论中明确提及的内容打分
10. 如果评论完全未提及某个维度，该维度返回null
11. 必须严格返回JSON格式，不要添加任何其他文字

返回JSON格式：
{"brand":"品牌名","model":"型号名","scores":{"维度1":8.5,"维度2":null}}`, strings.Join(dimList, "\n"))

	var userPrompt string
	if req.VideoTitle != "" {
		userPrompt = fmt.Sprintf("视频标题：%s\n\n评论内容：%s", req.VideoTitle, req.Comment)
	} else {
		userPrompt = fmt.Sprintf("评论内容：%s", req.Comment)
	}

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

		go func(index int, input CommentInput) {
			defer wg.Done()

			results[index] = CommentAnalysisResult{
				CommentID: input.ID,
				Content:   input.Content,
			}

			resp, err := c.AnalyzeComment(ctx, AnalyzeCommentRequest{
				Comment:    input.Content,
				Dimensions: dimensions,
				VideoTitle: input.VideoTitle,
			})

			if err != nil {
				results[index].Error = err.Error()
				return
			}

			results[index].Scores = resp.Scores
			results[index].Brand = resp.Brand
			results[index].Model = resp.Model
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
	ID         string // 评论ID
	Content    string // 评论内容
	VideoTitle string // 视频标题，作为上下文
	VideoBVID  string // 视频BVID
}

// AnalyzeCommentsWithRateLimit 带速率限制的批量分析（优化版）
// 使用批量合并策略和并发控制，大幅减少 API 调用次数并提高处理速度
// 参数：
//   - ctx: 上下文
//   - comments: 评论列表
//   - dimensions: 评价维度
//   - concurrency: 并发数（控制同时发送的AI请求数）
//
// 返回：
//   - []CommentAnalysisResult: 分析结果
//   - error: 错误信息
func (c *Client) AnalyzeCommentsWithRateLimit(ctx context.Context, comments []CommentInput, dimensions []Dimension, concurrency int) ([]CommentAnalysisResult, error) {
	// 使用动态批次计算
	config := DefaultBatchConfig()
	batches := CalculateBatches(comments, &config)
	log.Printf("[AI] 评论分析：共 %d 条评论，分为 %d 批处理，并发数: %d", len(comments), len(batches), concurrency)

	// 并发控制
	sem := semaphore.NewWeighted(int64(concurrency))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 预分配结果存储（按批次索引）
	batchResults := make([][]CommentAnalysisResult, len(batches))
	var completedCount int

	for i, batch := range batches {
		// 检查上下文取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 获取信号量
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(idx int, b []CommentInput) {
			defer wg.Done()
			defer sem.Release(1)

			log.Printf("[AI] 正在分析第 %d/%d 批（%d 条评论）...", idx+1, len(batches), len(b))

			// 尝试批量合并分析
			results, err := c.AnalyzeCommentsBatchMerged(ctx, b, dimensions)
			if err != nil {
				// 降级：使用原有的并发单条分析
				log.Printf("[AI] 批量分析失败，降级到单条分析: %v", err)
				results, err = c.AnalyzeCommentsBatch(ctx, b, dimensions)
				if err != nil {
					// 如果单条分析也失败，记录错误但继续
					log.Printf("[AI] 单条分析也失败: %v", err)
					results = make([]CommentAnalysisResult, len(b))
					for j, comment := range b {
						results[j] = CommentAnalysisResult{
							CommentID: comment.ID,
							Content:   comment.Content,
							Error:     err.Error(),
						}
					}
				}
			}

			// 线程安全地存储结果和更新进度
			mu.Lock()
			batchResults[idx] = results
			completedCount++
			c.reportProgress("analyzing", completedCount, len(batches),
				fmt.Sprintf("正在分析第 %d/%d 批（%d 条评论）", completedCount, len(batches), len(b)))
			mu.Unlock()

			log.Printf("[AI] 完成第 %d/%d 批分析", completedCount, len(batches))
		}(i, batch)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 按顺序合并结果
	var allResults []CommentAnalysisResult
	for _, results := range batchResults {
		allResults = append(allResults, results...)
	}

	log.Printf("[AI] 所有批次分析完成，共 %d 条结果", len(allResults))
	return allResults, nil
}

// RecommendationInput AI生成购买建议的输入数据
type RecommendationInput struct {
	Category      string
	Rankings      []BrandRankingInfo
	BrandAnalysis map[string]BrandStrengthWeakness
	ModelRankings []ModelRankingInfo
}

// BrandRankingInfo 品牌排名信息（用于AI生成建议）
type BrandRankingInfo struct {
	Brand        string
	OverallScore float64
	Rank         int
}

// ModelRankingInfo 型号排名信息（用于AI生成建议）
type ModelRankingInfo struct {
	Model        string
	Brand        string
	OverallScore float64
	Rank         int
	CommentCount int
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
4. 使用Markdown格式输出，包括：
   - 使用 ## 作为小标题
   - 使用 **加粗** 强调重点
   - 使用 - 列表展示要点
   - 使用 > 引用块突出关键建议`

	var modelText string
	if len(input.ModelRankings) > 0 {
		modelText = "\n\n型号排名：\n"
		for _, m := range input.ModelRankings {
			modelText += fmt.Sprintf("第%d名：%s %s（%.1f分，%d条评论）\n",
				m.Rank, m.Brand, m.Model, m.OverallScore, m.CommentCount)
		}
	}

	userPrompt := fmt.Sprintf("商品类别：%s\n\n品牌排名及分析：\n%s%s\n请生成购买建议：",
		input.Category, rankingText, modelText)

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

// BatchAnalysisResult 批量分析结果（用于 JSON 解析）
type BatchAnalysisResult struct {
	Results []struct {
		ID     string              `json:"id"`
		Brand  string              `json:"brand"`
		Model  string              `json:"model"`
		Scores map[string]*float64 `json:"scores"`
	} `json:"results"`
}

// AnalyzeCommentsBatchMerged 真正的批量分析（多条评论合并到一个请求）
// 将多条评论合并到一个 API 请求中，大幅减少请求次数
func (c *Client) AnalyzeCommentsBatchMerged(ctx context.Context, comments []CommentInput, dimensions []Dimension) ([]CommentAnalysisResult, error) {
	if len(comments) == 0 {
		return nil, fmt.Errorf("评论列表不能为空")
	}
	if len(dimensions) == 0 {
		return nil, fmt.Errorf("评价维度不能为空")
	}

	// 构建维度列表
	var dimList []string
	for _, dim := range dimensions {
		dimList = append(dimList, fmt.Sprintf("- %s：%s", dim.Name, dim.Description))
	}

	// 构建评论列表文本
	var commentList []string
	for i, c := range comments {
		if c.VideoTitle != "" {
			commentList = append(commentList, fmt.Sprintf("[%d] 视频：%s | 内容：%s", i+1, c.VideoTitle, c.Content))
		} else {
			commentList = append(commentList, fmt.Sprintf("[%d] 内容：%s", i+1, c.Content))
		}
	}

	systemPrompt := fmt.Sprintf(`你是商品评论分析助手。分析以下多条评论，为每条评论：
1. 提取品牌名称和具体型号
2. 对以下维度打分（1-10分，未提及则为null）：
%s

评分标准：1-3差评，4-5一般，6-7较好，8-10优秀

重要规则：
- 每条评论独立分析，用评论编号[1][2]等标识
- 品牌必须是单一品牌名称，绝对不能包含"/"或其他分隔符
- 如果评论对比多个品牌（如"A比B好"），只提取评论主要评价的那个品牌
- 如果评论同时讨论多个品牌且无法确定主要品牌，brand填"未知"
- 型号必须是具体型号名（如"V12"、"Max"、"Pro"），不能是描述性文字（如"新款"、"基础款"）
- 无法确定品牌填"未知"，无法确定型号填"通用"
- 必须返回JSON格式，不要添加任何其他文字
- results数组的顺序必须与输入评论顺序一致

返回格式：
{"results":[{"id":"1","brand":"品牌","model":"型号","scores":{"维度1":8.5,"维度2":null}},{"id":"2",...}]}`, strings.Join(dimList, "\n"))

	userPrompt := fmt.Sprintf("评论列表（共%d条）：\n%s", len(comments), strings.Join(commentList, "\n"))

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := c.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("AI请求失败: %w", err)
	}

	// 清理响应（移除可能的 markdown 代码块）
	cleanedResponse := cleanJSONResponse(response)

	// 解析批量结果
	var batchResult BatchAnalysisResult
	if err := json.Unmarshal([]byte(cleanedResponse), &batchResult); err != nil {
		// 尝试提取 JSON
		jsonPattern := regexp.MustCompile(`\{[\s\S]*"results"[\s\S]*\}`)
		match := jsonPattern.FindString(response)
		if match != "" {
			if err := json.Unmarshal([]byte(match), &batchResult); err != nil {
				return nil, fmt.Errorf("解析批量响应失败: %w, 原始响应: %s", err, response[:min(len(response), 500)])
			}
		} else {
			return nil, fmt.Errorf("无法从响应中提取JSON: %s", response[:min(len(response), 500)])
		}
	}

	// 转换为 CommentAnalysisResult 格式
	results := make([]CommentAnalysisResult, len(comments))
	for i, c := range comments {
		results[i] = CommentAnalysisResult{
			CommentID: c.ID,
			Content:   c.Content,
		}

		// 查找对应的分析结果
		for _, r := range batchResult.Results {
			// 支持 "1" 或 "comment_0" 格式的 ID
			if r.ID == fmt.Sprintf("%d", i+1) || r.ID == c.ID {
				results[i].Brand = r.Brand
				results[i].Model = r.Model
				results[i].Scores = r.Scores
				break
			}
		}

		// 如果没找到对应结果，标记错误
		if results[i].Scores == nil {
			results[i].Error = "未在批量响应中找到对应结果"
		}
	}

	return results, nil
}
