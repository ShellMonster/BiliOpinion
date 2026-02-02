package ai

// BatchConfig 批次配置
type BatchConfig struct {
	MaxCharsPerBatch int // 每批最大字符数（默认 3000）
	MaxItemsPerBatch int // 每批最大条数（默认 15）
	MinItemsPerBatch int // 每批最小条数（默认 1）
}

// DefaultBatchConfig 默认批次配置
func DefaultBatchConfig() BatchConfig {
	return BatchConfig{
		MaxCharsPerBatch: 3000,
		MaxItemsPerBatch: 15,
		MinItemsPerBatch: 1,
	}
}

// CalculateBatches 按字符数动态计算批次
// 返回分好批的评论列表
func CalculateBatches(comments []CommentInput, config *BatchConfig) [][]CommentInput {
	if config == nil {
		cfg := DefaultBatchConfig()
		config = &cfg
	}

	var batches [][]CommentInput
	var currentBatch []CommentInput
	currentChars := 0

	for _, c := range comments {
		// 计算当前评论的字符数（内容 + 视频标题）
		commentLen := len([]rune(c.Content)) + len([]rune(c.VideoTitle))

		// 如果当前批次加上这条评论会超限，且当前批次不为空，则开始新批次
		shouldStartNewBatch := (currentChars+commentLen > config.MaxCharsPerBatch ||
			len(currentBatch) >= config.MaxItemsPerBatch) &&
			len(currentBatch) >= config.MinItemsPerBatch

		if shouldStartNewBatch {
			batches = append(batches, currentBatch)
			currentBatch = nil
			currentChars = 0
		}

		currentBatch = append(currentBatch, c)
		currentChars += commentLen
	}

	// 添加最后一批
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}
