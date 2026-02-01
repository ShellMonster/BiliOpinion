package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// SearchVideosRequest 搜索视频请求参数
// 用于指定搜索关键词和分页信息
type SearchVideosRequest struct {
	Keyword  string // 搜索关键词（必填）
	Page     int    // 页码（从1开始，默认1）
	PageSize int    // 每页数量（默认20，最大50）
}

// SearchVideosResponse B站搜索API响应结构
type SearchVideosResponse struct {
	Code    int    `json:"code"`    // 响应码（0表示成功）
	Message string `json:"message"` // 响应消息
	Data    struct {
		NumResults int         `json:"numResults"` // 搜索结果总数
		NumPages   int         `json:"numPages"`   // 总页数
		Result     []VideoInfo `json:"result"`     // 视频列表
	} `json:"data"`
}

// VideoInfo 视频信息结构
// 包含视频的基本信息，用于展示搜索结果
type VideoInfo struct {
	BVID        string `json:"bvid"`         // BV号（视频唯一标识）
	AID         int64  `json:"aid"`          // AV号（旧版视频ID）
	Title       string `json:"title"`        // 视频标题（可能包含HTML高亮标签）
	Author      string `json:"author"`       // UP主昵称
	Mid         int64  `json:"mid"`          // UP主UID
	Play        int    `json:"play"`         // 播放量
	VideoReview int    `json:"video_review"` // 评论数
	Favorites   int    `json:"favorites"`    // 收藏数
	Duration    string `json:"duration"`     // 视频时长（格式：mm:ss）
	Pic         string `json:"pic"`          // 封面图URL
	Description string `json:"description"`  // 视频简介
	Pubdate     int64  `json:"pubdate"`      // 发布时间戳
}

// SearchVideos 搜索B站视频
// 使用B站搜索API，支持关键词搜索和分页
//
// 参数：
//   - req: 搜索请求参数
//
// 返回：
//   - []VideoInfo: 视频列表
//   - int: 搜索结果总数
//   - error: 错误信息
//
// 示例：
//
//	videos, total, err := client.SearchVideos(SearchVideosRequest{
//	    Keyword: "iPhone 15 评测",
//	    Page: 1,
//	    PageSize: 20,
//	})
func (c *Client) SearchVideos(req SearchVideosRequest) ([]VideoInfo, int, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	// 限制每页最大数量为50
	if req.PageSize > 50 {
		req.PageSize = 50
	}

	// 构建搜索URL
	// search_type=video 表示搜索视频
	// keyword 需要URL编码
	u := fmt.Sprintf(
		"https://api.bilibili.com/x/web-interface/wbi/search/type?search_type=video&keyword=%s&page=%d&page_size=%d",
		url.QueryEscape(req.Keyword),
		req.Page,
		req.PageSize,
	)

	// 发送请求（需要WBI签名）
	resp, err := c.Get(u, true)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析JSON响应
	var searchResp SearchVideosResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, 0, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API响应码
	if searchResp.Code != 0 {
		return nil, 0, fmt.Errorf("API错误: %s (code: %d)", searchResp.Message, searchResp.Code)
	}

	// 过滤掉BVID为空的视频
	var validVideos []VideoInfo
	for _, v := range searchResp.Data.Result {
		if v.BVID != "" {
			validVideos = append(validVideos, v)
		}
	}

	return validVideos, searchResp.Data.NumResults, nil
}

// SearchVideosWithLimit 搜索视频（带数量限制）
// 自动分页获取视频，直到达到指定数量或没有更多结果
//
// 参数：
//   - keyword: 搜索关键词
//   - maxVideos: 最大视频数量（默认50）
//
// 返回：
//   - []VideoInfo: 视频列表
//   - error: 错误信息
//
// 示例：
//
//	videos, err := client.SearchVideosWithLimit("iPhone 15 评测", 50)
func (c *Client) SearchVideosWithLimit(keyword string, maxVideos int) ([]VideoInfo, error) {
	// 默认限制50个视频
	if maxVideos <= 0 {
		maxVideos = 50
	}

	var allVideos []VideoInfo
	page := 1
	pageSize := 20 // 每页20个

	for len(allVideos) < maxVideos {
		// 搜索当前页
		videos, _, err := c.SearchVideos(SearchVideosRequest{
			Keyword:  keyword,
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("搜索第%d页失败: %w", page, err)
		}

		// 没有更多结果
		if len(videos) == 0 {
			break
		}

		// 添加到结果列表
		allVideos = append(allVideos, videos...)
		page++

		// 防止无限循环（最多搜索10页）
		if page > 10 {
			break
		}
	}

	// 截取到指定数量
	if len(allVideos) > maxVideos {
		allVideos = allVideos[:maxVideos]
	}

	return allVideos, nil
}
