package main

import (
	"context"
	"fmt"
	

	"bilibili-analyzer/backend/ai"
)

func main() {
	client := ai.NewClient(ai.Config{
		APIBase:       "https://yunwu.ai/v1",
		APIKey:        "sk-VuMgwjl1N8Xuy490KON5JHRx46WayvLtmu3ELIYFH2dZz6bL",
		Model:         "gemini-3-flash-preview",
		MaxConcurrent: 5,
	})

	// 准备20条测试评论
	comments := []ai.CommentInput{}
	for i := 1; i <= 20; i++ {
		comments = append(comments, ai.CommentInput{
			ID:         fmt.Sprintf("comment_%d", i),
			Content:    fmt.Sprintf("这是第%d条评论，产品很好用，音质不错", i),
			VideoTitle: fmt.Sprintf("视频%d", i),
		})
	}

	dimensions := []ai.Dimension{
		{Name: "音质", Description: "音质表现"},
		{Name: "性价比", Description: "价格与性能比"},
	}

	// 测试批次计算
	config := ai.DefaultBatchConfig()
	batches := ai.CalculateBatches(comments, &config)

	fmt.Printf("📦 批次计算结果:\n")
	fmt.Printf("   总评论数: %d\n", len(comments))
	fmt.Printf("   批次数量: %d\n", len(batches))
	fmt.Printf("   平均每批: %.1f 条\n", float64(len(comments))/float64(len(batches)))
	fmt.Println()

	for i, batch := range batches {
		totalChars := 0
		for _, c := range batch {
			totalChars += len([]rune(c.Content)) + len([]rune(c.VideoTitle))
		}
		fmt.Printf("   批次 %d: %d 条评论, 约 %d 字符\n", i+1, len(batch), totalChars)
	}

	// 测试批量分析
	fmt.Println("\n⏳ 开始批量分析...")
	results, err := client.AnalyzeCommentsWithRateLimit(context.Background(), comments, dimensions, 0)
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		return
	}

	fmt.Printf("\n✅ 分析完成！成功分析 %d 条评论\n", len(results))

	// 统计成功率
	successCount := 0
	for _, r := range results {
		if r.Error == "" && r.Scores != nil {
			successCount++
		}
	}
	fmt.Printf("📈 成功率: %d/%d (%.1f%%)\n", successCount, len(results), float64(successCount)*100/float64(len(results)))

	// 显示前3条结果
	fmt.Println("\n📝 前3条分析结果:")
	for i := 0; i < 3 && i < len(results); i++ {
		r := results[i]
		fmt.Printf("\n评论 %d:\n", i+1)
		fmt.Printf("  内容: %s...\n", r.Content[:min(30, len(r.Content))])
		fmt.Printf("  品牌: %s, 型号: %s\n", r.Brand, r.Model)
		if r.Error != "" {
			fmt.Printf("  错误: %s\n", r.Error)
		} else {
			fmt.Printf("  得分: ")
			for dim, score := range r.Scores {
				if score != nil {
					fmt.Printf("%s=%.1f ", dim, *score)
				}
			}
			fmt.Println()
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
