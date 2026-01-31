package bilibili

import (
	"net/http"
	"net/url"
	"time"
)

// Client B站API客户端
// 封装HTTP请求，自动处理Cookie和WBI签名
type Client struct {
	httpClient *http.Client // HTTP客户端
	cookie     string       // 用户Cookie（用于需要登录的接口）
}

// NewClient 创建新的B站API客户端
// 参数：
//   - cookie: 用户的完整Cookie字符串（可选，未登录接口可传空字符串）
//
// 返回：
//   - *Client: 客户端实例
//
// 示例：
//
//	client := NewClient("SESSDATA=xxx; bili_jct=xxx")
func NewClient(cookie string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second, // 设置20秒超时
		},
		cookie: cookie,
	}
}

// Get 发送GET请求（自动添加WBI签名）
// 参数：
//   - urlStr: 请求URL字符串
//   - needSign: 是否需要WBI签名（大部分B站API都需要）
//
// 返回：
//   - *http.Response: HTTP响应对象
//   - error: 请求失败时返回错误
//
// 示例：
//
//	// 需要签名的请求
//	resp, err := client.Get("https://api.bilibili.com/x/space/wbi/acc/info?mid=1850091", true)
//
//	// 不需要签名的请求
//	resp, err := client.Get("https://api.bilibili.com/x/web-interface/nav", false)
func (c *Client) Get(urlStr string, needSign bool) (*http.Response, error) {
	// 解析URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// 如果需要签名，添加WBI签名
	if needSign {
		if err := Sign(u); err != nil {
			return nil, err
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// 设置必需的请求头
	// User-Agent: 模拟Chrome浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	// Referer: B站主站（防止反爬）
	req.Header.Set("Referer", "https://www.bilibili.com/")

	// 设置Cookie（如果提供）
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	// 发送请求
	return c.httpClient.Do(req)
}

// SetCookie 设置Cookie
// 参数：
//   - cookie: 新的Cookie字符串
//
// 用途：
//   - 动态更新Cookie（如用户重新登录）
//
// 示例：
//
//	client.SetCookie("SESSDATA=new_value; bili_jct=new_value")
func (c *Client) SetCookie(cookie string) {
	c.cookie = cookie
}
