package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/sync/semaphore"
)

// Client AI客户端
// 用于与OpenAI兼容的API进行交互
type Client struct {
	apiBase          string              // API基础URL（如：https://api.openai.com/v1）
	apiKey           string              // API密钥
	model            string              // 使用的模型名称（如：gemini-3-flash-preview）
	httpClient       *http.Client        // HTTP客户端
	sem              *semaphore.Weighted // 并发控制信号量
	progressCallback ProgressCallback    // 进度回调函数
}

// Config AI客户端配置
type Config struct {
	APIBase       string // API Base URL（默认：https://api.openai.com/v1）
	APIKey        string // API Key
	Model         string // 模型名称
	MaxConcurrent int64  // 最大并发数（默认：5）
}

// NewClient 创建新的AI客户端
// 参数：
//   - cfg: 客户端配置
//
// 返回：
//   - *Client: AI客户端实例
func NewClient(cfg Config) *Client {
	// 设置默认API Base URL
	if cfg.APIBase == "" {
		cfg.APIBase = "https://api.openai.com/v1"
	}
	// 设置默认最大并发数
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 10
	}

	return &Client{
		apiBase: cfg.APIBase,
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 60秒超时（AI请求可能较慢）
		},
		sem: semaphore.NewWeighted(cfg.MaxConcurrent), // 创建并发控制信号量
	}
}

// ChatCompletionRequest Chat Completion请求结构
type ChatCompletionRequest struct {
	Model    string    `json:"model"`    // 模型名称
	Messages []Message `json:"messages"` // 消息列表
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // 角色：system（系统）/user（用户）/assistant（助手）
	Content string `json:"content"` // 消息内容
}

// ChatCompletionResponse Chat Completion响应结构
type ChatCompletionResponse struct {
	ID      string   `json:"id"`      // 响应ID
	Object  string   `json:"object"`  // 对象类型
	Created int64    `json:"created"` // 创建时间戳
	Model   string   `json:"model"`   // 使用的模型
	Choices []Choice `json:"choices"` // 选择项列表
}

// Choice 选择项结构
type Choice struct {
	Index        int     `json:"index"`         // 索引
	Message      Message `json:"message"`       // 消息内容
	FinishReason string  `json:"finish_reason"` // 完成原因
}

// ChatCompletion 发送Chat Completion请求
// 参数：
//   - ctx: 上下文（用于取消和超时控制）
//   - messages: 消息列表
//
// 返回：
//   - string: AI返回的文本内容
//   - error: 请求失败时返回错误
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, error) {
	// 并发控制：获取信号量（限制同时进行的请求数）
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("acquire semaphore failed: %w", err)
	}
	defer c.sem.Release(1) // 请求完成后释放信号量

	// 构建请求
	req := ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	// 重试逻辑：最多重试1次（总共尝试2次）
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		resp, err := c.doRequest(ctx, req)
		if err == nil {
			return resp, nil // 请求成功，返回结果
		}
		lastErr = err

		// 第一次失败后等待1秒再重试
		if attempt == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return "", fmt.Errorf("request failed after 2 attempts: %w", lastErr)
}

// doRequest 执行HTTP请求
// 参数：
//   - ctx: 上下文
//   - req: Chat Completion请求
//
// 返回：
//   - string: AI返回的文本内容
//   - error: 请求失败时返回错误
func (c *Client) doRequest(ctx context.Context, req ChatCompletionRequest) (string, error) {
	// 序列化请求体为JSON
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.apiBase+"/chat/completions", // OpenAI Chat Completion端点
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey) // Bearer Token认证

	// 发送HTTP请求
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	// 检查HTTP状态码
	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", httpResp.StatusCode, string(respBody))
	}

	// 解析JSON响应
	var resp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %w", err)
	}

	// 提取AI返回的文本内容
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}

// SetProgressCallback 设置进度回调
//
// 参数：
//   - callback: 进度回调函数
//
// 示例：
//
//	client.SetProgressCallback(func(stage string, current, total int, message string) {
//	    fmt.Printf("[%s] %d/%d: %s\n", stage, current, total, message)
//	})
func (c *Client) SetProgressCallback(callback ProgressCallback) {
	c.progressCallback = callback
}

// reportProgress 报告进度
func (c *Client) reportProgress(stage string, current, total int, message string) {
	if c.progressCallback != nil {
		c.progressCallback(stage, current, total, message)
	}
}
