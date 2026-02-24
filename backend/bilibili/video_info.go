package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// VideoDetail 视频详细信息
// 包含视频的完整元数据，用于展示视频来源和分析
// 注意：与 search.go 中的 VideoInfo 不同，这个是详细信息版本
type VideoDetail struct {
	BVID         string `json:"bvid"`          // 视频BV号，如 BV1mH4y1u7UA
	Title        string `json:"title"`         // 视频标题
	Author       string `json:"author"`        // UP主名称
	PlayCount    int    `json:"play_count"`    // 播放量
	CommentCount int    `json:"comment_count"` // 评论数
	PubDate      string `json:"pub_date"`      // 发布时间（格式化后的字符串）
	Cover        string `json:"cover"`         // 视频封面图URL
	Description  string `json:"description"`   // 视频简介
}

// videoDetailAPIResponse B站视频信息API的响应结构
type videoDetailAPIResponse struct {
	Code    int    `json:"code"`    // 响应码，0 表示成功
	Message string `json:"message"` // 错误信息
	Data    struct {
		BVID  string `json:"bvid"`  // 视频BV号
		AID   int64  `json:"aid"`   // 视频AV号
		Title string `json:"title"` // 视频标题
		Desc  string `json:"desc"`  // 视频简介
		Owner struct {
			MID  int64  `json:"mid"`  // UP主ID
			Name string `json:"name"` // UP主名称
		} `json:"owner"`
		Stat struct {
			View    int `json:"view"`    // 播放量
			Danmaku int `json:"danmaku"` // 弹幕数
			Reply   int `json:"reply"`   // 评论数
		} `json:"stat"`
		PubDate int64  `json:"pubdate"` // 发布时间戳（秒）
		Pic     string `json:"pic"`     // 封面图URL
	} `json:"data"`
}

// GetVideoInfo 获取B站视频详细信息
// 参数：
//   - bvid: 视频的BV号，如 "BV1mH4y1u7UA"
//
// 返回：
//   - *VideoInfo: 视频信息结构体
//   - error: 视频不存在或请求失败时返回错误
//
// 示例：
//
//	info, err := client.GetVideoInfo("BV1mH4y1u7UA")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("标题: %s, UP主: %s, 播放量: %d\n", info.Title, info.Author, info.PlayCount)
func (c *Client) GetVideoInfo(bvid string) (*VideoDetail, error) {
	// 构建视频信息API URL（该API不需要WBI签名）
	u := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)

	// 发送GET请求（不需要WBI签名）
	resp, err := c.Get(u, false)
	if err != nil {
		return nil, fmt.Errorf("获取视频信息失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var apiResp videoDetailAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API响应码
	// -400: 请求错误（如BV号格式错误）
	// -404: 视频不存在
	if apiResp.Code != 0 {
		if apiResp.Code == -404 {
			return nil, fmt.Errorf("视频不存在: %s", bvid)
		}
		return nil, fmt.Errorf("API错误: %s (code: %d)", apiResp.Message, apiResp.Code)
	}

	// 将时间戳转换为可读格式
	pubDate := time.Unix(apiResp.Data.PubDate, 0).Format("2006-01-02 15:04:05")

	// 构建并返回 VideoDetail 结构体
	return &VideoDetail{
		BVID:         apiResp.Data.BVID,
		Title:        apiResp.Data.Title,
		Author:       apiResp.Data.Owner.Name,
		PlayCount:    apiResp.Data.Stat.View,
		CommentCount: apiResp.Data.Stat.Reply,
		PubDate:      pubDate,
		Cover:        apiResp.Data.Pic,
		Description:  apiResp.Data.Desc,
	}, nil
}
