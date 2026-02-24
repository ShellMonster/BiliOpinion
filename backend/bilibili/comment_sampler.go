package bilibili

// SampleComments 采样视频评论，只获取主评论内容
// 用于快速获取评论样本进行AI分析
//
// 参数：
//   - bvid: 视频BV号
//   - count: 需要采样的评论数量
//
// 返回：
//   - []string: 评论内容字符串切片
//   - error: 错误信息
//
// 示例：
//
//	comments, err := client.SampleComments("BV1mH4y1u7UA", 50)
//	// 返回 ["评论1内容", "评论2内容", ...]
func (c *Client) SampleComments(bvid string, count int) ([]string, error) {
	// 调用GetAllComments获取评论
	// fetchReplies=false 表示不获取楼中楼，提高速度
	comments, err := c.GetAllComments(bvid, count, false)
	if err != nil {
		return nil, err
	}

	// 提取评论内容（Message字段）
	result := make([]string, 0, len(comments))
	for _, comment := range comments {
		if comment.Content.Message != "" {
			result = append(result, comment.Content.Message)
		}
	}

	return result, nil
}
