package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// GetCommentsRequest 获取评论请求参数
type GetCommentsRequest struct {
	BVID     string // BV号（视频唯一标识）
	Page     int    // 页码（从1开始）
	PageSize int    // 每页数量（最大20）
	Sort     int    // 排序方式：0=时间，1=点赞，2=回复数
}

// CommentsResponse B站评论API响应结构
type CommentsResponse struct {
	Code    int    `json:"code"`    // 响应码（0表示成功）
	Message string `json:"message"` // 响应消息
	Data    struct {
		Page struct {
			Num   int `json:"num"`   // 当前页码
			Size  int `json:"size"`  // 每页数量
			Count int `json:"count"` // 评论总数
		} `json:"page"`
		Replies []Comment `json:"replies"` // 评论列表
	} `json:"data"`
}

// Comment 评论结构
// 包含评论的完整信息，包括作者、内容、互动数据等
type Comment struct {
	RPID       int64     `json:"rpid"`    // 评论ID（唯一标识）
	OID        int64     `json:"oid"`     // 目标ID（视频AV号）
	Type       int       `json:"type"`    // 评论类型（1=视频）
	Mid        int64     `json:"mid"`     // 评论者UID
	Root       int64     `json:"root"`    // 根评论ID（0表示是根评论）
	Parent     int64     `json:"parent"`  // 父评论ID（0表示是根评论）
	Dialog     int64     `json:"dialog"`  // 对话ID
	Count      int       `json:"count"`   // 回复数
	RCount     int       `json:"rcount"`  // 回复数（另一个字段）
	Like       int       `json:"like"`    // 点赞数
	Ctime      int64     `json:"ctime"`   // 创建时间戳
	Content    Content   `json:"content"` // 评论内容
	Member     Member    `json:"member"`  // 评论者信息
	Replies    []Comment `json:"replies"` // 楼中楼评论（预加载的前3条）
	ReplyCount int       `json:"-"`       // 实际回复数（用于判断是否需要获取更多）
}

// Content 评论内容结构
type Content struct {
	Message string `json:"message"` // 评论文本内容
	Emote   any    `json:"emote"`   // 表情信息
}

// Member 评论者信息结构
type Member struct {
	Mid    string `json:"mid"`                      // 用户UID（字符串格式）
	Uname  string `json:"uname"`                    // 用户昵称
	Sex    string `json:"sex"`                      // 性别
	Sign   string `json:"sign"`                     // 个性签名
	Avatar string `json:"avatar"`                   // 头像URL
	Level  int    `json:"level_info.current_level"` // 用户等级
}

// GetComments 获取视频评论列表
// 获取指定视频的评论，支持分页和排序
//
// 参数：
//   - req: 获取评论请求参数
//
// 返回：
//   - []Comment: 评论列表
//   - int: 评论总数
//   - error: 错误信息
//
// 示例：
//
//	comments, total, err := client.GetComments(GetCommentsRequest{
//	    BVID: "BV1mH4y1u7UA",
//	    Page: 1,
//	    PageSize: 20,
//	    Sort: 1, // 按点赞排序
//	})
func (c *Client) GetComments(req GetCommentsRequest) ([]Comment, int, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 20 {
		req.PageSize = 20 // B站API限制每页最多20条
	}

	// 将BV号转换为AV号（评论API使用AV号）
	avid := Bvid2Avid(req.BVID)

	// 构建评论API URL
	// type=1 表示视频评论
	// oid 是视频的AV号
	// pn 是页码，ps 是每页数量
	// sort 是排序方式
	u := fmt.Sprintf(
		"https://api.bilibili.com/x/v2/reply?type=1&oid=%d&pn=%d&ps=%d&sort=%d",
		avid,
		req.Page,
		req.PageSize,
		req.Sort,
	)

	// 发送请求（评论API不需要WBI签名）
	resp, err := c.Get(u, false)
	if err != nil {
		return nil, 0, fmt.Errorf("获取评论失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var commentResp CommentsResponse
	if err := json.Unmarshal(body, &commentResp); err != nil {
		return nil, 0, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API响应码
	if commentResp.Code != 0 {
		return nil, 0, fmt.Errorf("API错误: %s (code: %d)", commentResp.Message, commentResp.Code)
	}

	// 设置回复数（用于后续判断是否需要获取楼中楼）
	comments := commentResp.Data.Replies
	for i := range comments {
		comments[i].ReplyCount = comments[i].RCount
	}

	return comments, commentResp.Data.Page.Count, nil
}

// GetReplies 获取楼中楼评论（评论的回复）
// 获取指定评论下的所有回复
//
// 参数：
//   - bvid: 视频BV号
//   - rootRPID: 根评论ID
//   - page: 页码（从1开始）
//   - pageSize: 每页数量（最大20）
//
// 返回：
//   - []Comment: 回复列表
//   - error: 错误信息
//
// 示例：
//
//	replies, err := client.GetReplies("BV1mH4y1u7UA", 123456789, 1, 20)
func (c *Client) GetReplies(bvid string, rootRPID int64, page, pageSize int) ([]Comment, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 20 {
		pageSize = 20
	}

	// 将BV号转换为AV号
	avid := Bvid2Avid(bvid)

	// 构建楼中楼API URL
	// root 是根评论ID
	u := fmt.Sprintf(
		"https://api.bilibili.com/x/v2/reply/reply?type=1&oid=%d&root=%d&pn=%d&ps=%d",
		avid,
		rootRPID,
		page,
		pageSize,
	)

	// 发送请求
	resp, err := c.Get(u, false)
	if err != nil {
		return nil, fmt.Errorf("获取楼中楼评论失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var replyResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Page struct {
				Count int `json:"count"` // 回复总数
			} `json:"page"`
			Replies []Comment `json:"replies"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &replyResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if replyResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s (code: %d)", replyResp.Message, replyResp.Code)
	}

	return replyResp.Data.Replies, nil
}

// GetAllReplies 获取评论的所有楼中楼回复
// 自动分页获取所有回复
//
// 参数：
//   - bvid: 视频BV号
//   - rootRPID: 根评论ID
//   - maxReplies: 最大回复数量（0表示不限制）
//
// 返回：
//   - []Comment: 所有回复列表
//   - error: 错误信息
func (c *Client) GetAllReplies(bvid string, rootRPID int64, maxReplies int) ([]Comment, error) {
	var allReplies []Comment
	page := 1
	pageSize := 20

	for {
		replies, err := c.GetReplies(bvid, rootRPID, page, pageSize)
		if err != nil {
			return nil, err
		}

		// 没有更多回复
		if len(replies) == 0 {
			break
		}

		allReplies = append(allReplies, replies...)
		page++

		// 达到最大数量限制
		if maxReplies > 0 && len(allReplies) >= maxReplies {
			allReplies = allReplies[:maxReplies]
			break
		}

		// 防止无限循环（最多获取10页）
		if page > 10 {
			break
		}

		// 添加延迟，避免请求过快
		time.Sleep(100 * time.Millisecond)
	}

	return allReplies, nil
}

// GetCommentsWithReplies 获取评论及其楼中楼
// 获取评论列表，并自动获取每条评论的楼中楼回复
//
// 参数：
//   - req: 获取评论请求参数
//   - fetchReplies: 是否获取楼中楼
//   - maxRepliesPerComment: 每条评论最大回复数（0表示不限制）
//
// 返回：
//   - []Comment: 评论列表（包含楼中楼）
//   - int: 评论总数
//   - error: 错误信息
func (c *Client) GetCommentsWithReplies(req GetCommentsRequest, fetchReplies bool, maxRepliesPerComment int) ([]Comment, int, error) {
	// 获取评论列表
	comments, total, err := c.GetComments(req)
	if err != nil {
		return nil, 0, err
	}

	// 如果不需要获取楼中楼，直接返回
	if !fetchReplies {
		return comments, total, nil
	}

	// 获取每条评论的楼中楼
	for i := range comments {
		// 只有当评论有回复时才获取
		if comments[i].ReplyCount > 0 {
			replies, err := c.GetAllReplies(req.BVID, comments[i].RPID, maxRepliesPerComment)
			if err != nil {
				// 获取楼中楼失败不影响主流程，记录错误继续
				continue
			}
			comments[i].Replies = replies
		}

		// 添加延迟，避免请求过快
		time.Sleep(50 * time.Millisecond)
	}

	return comments, total, nil
}

// GetAllComments 获取视频的所有评论（带数量限制）
// 自动分页获取评论，直到达到指定数量或没有更多结果
//
// 参数：
//   - bvid: 视频BV号
//   - maxComments: 最大评论数量（默认500）
//   - fetchReplies: 是否获取楼中楼
//
// 返回：
//   - []Comment: 评论列表
//   - error: 错误信息
func (c *Client) GetAllComments(bvid string, maxComments int, fetchReplies bool) ([]Comment, error) {
	// 默认限制500条评论
	if maxComments <= 0 {
		maxComments = 500
	}

	var allComments []Comment
	page := 1
	pageSize := 20

	for len(allComments) < maxComments {
		// 获取当前页评论
		comments, _, err := c.GetCommentsWithReplies(GetCommentsRequest{
			BVID:     bvid,
			Page:     page,
			PageSize: pageSize,
			Sort:     1, // 按点赞排序，获取高质量评论
		}, fetchReplies, 10) // 每条评论最多获取10条楼中楼

		if err != nil {
			return nil, fmt.Errorf("获取第%d页评论失败: %w", page, err)
		}

		// 没有更多评论
		if len(comments) == 0 {
			break
		}

		allComments = append(allComments, comments...)
		page++

		// 防止无限循环（最多获取50页）
		if page > 50 {
			break
		}

		// 添加延迟，避免请求过快
		time.Sleep(200 * time.Millisecond)
	}

	// 截取到指定数量
	if len(allComments) > maxComments {
		allComments = allComments[:maxComments]
	}

	return allComments, nil
}
