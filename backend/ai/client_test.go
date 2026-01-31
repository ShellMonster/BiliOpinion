package ai

import (
	"context"
	"testing"
	"time"
)

// TestNewClient 测试创建AI客户端
func TestNewClient(t *testing.T) {
	client := NewClient(Config{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})

	// 验证默认API Base URL
	if client.apiBase != "https://api.openai.com/v1" {
		t.Errorf("Expected default API base, got %s", client.apiBase)
	}

	// 验证API Key
	if client.apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", client.apiKey)
	}

	// 验证Model
	if client.model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got %s", client.model)
	}
}

// TestCustomAPIBase 测试自定义API Base URL
func TestCustomAPIBase(t *testing.T) {
	client := NewClient(Config{
		APIBase: "https://custom.api.com/v1",
		APIKey:  "test-key",
		Model:   "gpt-4",
	})

	// 验证自定义API Base URL
	if client.apiBase != "https://custom.api.com/v1" {
		t.Errorf("Expected custom API base, got %s", client.apiBase)
	}
}

// TestConcurrencyControl 测试并发控制
func TestConcurrencyControl(t *testing.T) {
	client := NewClient(Config{
		APIKey:        "test-key",
		Model:         "gpt-3.5-turbo",
		MaxConcurrent: 2, // 设置最大并发数为2
	})

	// 测试并发控制（模拟）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试获取2个信号量（应该成功）
	err1 := client.sem.Acquire(ctx, 1)
	err2 := client.sem.Acquire(ctx, 1)

	if err1 != nil || err2 != nil {
		t.Error("Should be able to acquire 2 semaphores")
	}

	// 第3个应该阻塞（在测试中我们设置短超时来验证）
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()
	err3 := client.sem.Acquire(ctx2, 1)

	if err3 == nil {
		t.Error("Should not be able to acquire 3rd semaphore")
		client.sem.Release(1) // 如果意外获取成功，释放它
	}

	// 释放前2个信号量
	client.sem.Release(2)
}

// TestDefaultMaxConcurrent 测试默认最大并发数
func TestDefaultMaxConcurrent(t *testing.T) {
	client := NewClient(Config{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
		// 不设置MaxConcurrent，应该使用默认值5
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试获取5个信号量（应该成功）
	for i := 0; i < 5; i++ {
		if err := client.sem.Acquire(ctx, 1); err != nil {
			t.Errorf("Should be able to acquire semaphore %d", i+1)
		}
	}

	// 第6个应该阻塞
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()
	err := client.sem.Acquire(ctx2, 1)

	if err == nil {
		t.Error("Should not be able to acquire 6th semaphore")
		client.sem.Release(1)
	}

	// 释放所有信号量
	client.sem.Release(5)
}

// TestHTTPClientTimeout 测试HTTP客户端超时设置
func TestHTTPClientTimeout(t *testing.T) {
	client := NewClient(Config{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})

	// 验证HTTP客户端超时为60秒
	expectedTimeout := 60 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}
