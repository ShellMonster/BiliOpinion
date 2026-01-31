
## B站API Go语言SDK调研结果

### 发现的主要SDK库

#### 1. CuteReimu/bilibili (★184, 最推荐)
- **仓库**: https://github.com/CuteReimu/bilibili
- **最新提交**: f9f2a81d9a777f54f66134518d69486c3a9ffd82
- **特点**:
  - 功能最完整,接口最全面
  - 支持Go 1.23+
  - 基于resty/v2 HTTP客户端
  - 完整的WBI签名实现
  - 完善的评论API支持
  - 良好的Cookie管理
  - 详细的中文文档

- **核心依赖**:
  ```
  github.com/go-resty/resty/v2 v2.16.5
  github.com/pkg/errors v0.9.1
  golang.org/x/sync v0.16.0
  ```

- **评论API示例**:
  ```go
  // 获取评论区明细
  func (c *Client) GetCommentsDetail(param GetCommentsDetailParam) (*CommentsDetail, error)
  
  // 参数结构
  type GetCommentsDetailParam struct {
      Type  int // 评论区类型代码
      Oid   int // 目标评论区id
      Sort  int // 排序方式: 0按时间, 1按点赞, 2按回复
      Ps    int // 每页项数,默认20,范围1-20
      Pn    int // 页码,默认1
  }
  ```

- **搜索API示例**:
  ```go
  // 综合搜索(web端)
  func (c *Client) IntergratedSearch(param SearchParam) (*SearchRespData, error)
  
  // 使用WBI签名
  url: "https://api.bilibili.com/x/web-interface/wbi/search/all/v2"
  ```

#### 2. iyear/biligo (★89)
- **仓库**: https://github.com/iyear/biligo
- **最新提交**: cc9f58c336def0a7f3889cc6501934a0d20ca3c7
- **特点**:
  - 轻量级设计
  - 支持Go 1.25
  - 使用标准库http.Client
  - 简洁的WBI签名实现
  - 支持直播消息流监听

- **核心依赖**:
  ```
  github.com/tidwall/gjson v1.18.0
  github.com/coder/websocket v1.8.14
  ```

#### 3. WhiteBlue/bilibili-sdk-go (★322, 较旧)
- **仓库**: https://github.com/WhiteBlue/bilibili-sdk-go
- **特点**: 较早期的SDK,可能不支持最新API

#### 4. FKU-bilimall/bilibili-ticket-go
- **仓库**: https://github.com/FKU-bilimall/bilibili-ticket-go
- **最新提交**: d0d7c0d610920ff435e3046c21c344e82224a25c
- **特点**: 专注于WBI签名实现

### WBI签名实现对比

#### CuteReimu/bilibili实现
**文件**: [wbi.go](https://github.com/CuteReimu/bilibili/blob/f9f2a81d9a777f54f66134518d69486c3a9ffd82/wbi.go)

```go
// WBI签名核心逻辑
type WBI struct {
    cookies        []*http.Cookie
    mixinKeyEncTab []int
    updateCheckerInterval time.Duration
    lastInitTime   time.Time
    storage        Storage
    sfg            singleflight.Group
}

// 签名步骤:
// 1. 获取imgKey和subKey (从 https://api.bilibili.com/x/web-interface/nav)
// 2. 生成mixinKey (通过mixinKeyEncTab打乱顺序)
// 3. 添加wts时间戳
// 4. 排序参数并移除特殊字符
// 5. 计算MD5: md5(params + mixinKey)
// 6. 添加w_rid参数

// mixinKeyEncTab固定值
var _defaultMixinKeyEncTab = []int{
    46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
    33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
    61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
    36, 20, 34, 44, 52,
}
```

#### bilibili-ticket-go实现
**文件**: [bili/sign.go](https://github.com/FKU-bilimall/bilibili-ticket-go/blob/d0d7c0d610920ff435e3046c21c344e82224a25c/bili/sign.go)

```go
// 包含App签名和WBI签名两种方式
const appKey = "1d8b6e7d45233436"
const appSec = "560c52ccd288fed045859ed18bffd973"

// WBI签名
func (c *Client) getSignedParameterWithAbi(forceUpdate bool, u *url.URL) error {
    if c.wbi == nil || c.wbi.isExpired() || forceUpdate {
        err := c.refreshWbiToken()
    }
    values := u.Query()
    values.Set("wts", fmt.Sprintf("%d", time.Now().Unix()))
    wbi := md5.Sum([]byte(values.Encode() + c.wbi.mixin))
    values.Set("w_rid", hex.EncodeToString(wbi[:]))
    u.RawQuery = values.Encode()
}

// App签名
func (c *Client) getSignedParameterWithApp(params map[string]any) url.Values {
    values.Set("appkey", appKey)
    values.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
    sign := md5.Sum([]byte(values.Encode() + appSec))
    values.Set("sign", hex.EncodeToString(sign[:]))
}
```

### HTTP客户端封装模式

#### CuteReimu/bilibili模式
```go
// 使用resty封装
type Client struct {
    wbi   *WBI
    resty *resty.Client
}

// 默认配置
restyClient := resty.New().
    SetRedirectPolicy(resty.NoRedirectPolicy()).
    SetTimeout(20*time.Second).
    SetHeader("User-Agent", "Mozilla/5.0...")
```

#### iyear/biligo模式
```go
// 使用标准库http.Client
type client struct {
    *http.Client
    Headers http.Header
}

// 自定义Request包装
type Request struct {
    *http.Request
    querys *netUrl.Values
    wbi    bool // wbi签名标记
}
```

### Cookie处理方式

#### CuteReimu/bilibili
```go
// Cookie字符串存储
cookiesString := client.GetCookiesString()
client.SetCookiesString(cookiesString)

// 原始Cookie设置
client.SetRawCookies("cookie1=xxx; cookie2=xxx")

// 匿名访客Cookie
client := bilibili.NewAnonymousClient()
```

#### iyear/biligo
```go
// 自动获取访客Cookie
type cookieManager struct {
    cookie        *http.Cookie
    guestFetching sync.Mutex
}

func (c *cookieManager) fetchGuestCookie() error {
    resp, err := httpClient.Do(NewGet(URL_MAIN_PAGE))
    // 从Set-Cookie header提取
}
```

### 推荐架构选择

**推荐使用 CuteReimu/bilibili 作为基础SDK**

理由:
1. ✅ 功能最完整,接口覆盖全面
2. ✅ 活跃维护,支持最新Go版本
3. ✅ 完善的WBI签名实现
4. ✅ 良好的错误处理和类型定义
5. ✅ 详细的中文文档和示例
6. ✅ 基于resty,易于扩展
7. ✅ 完整的评论API支持

### 关键API端点

```
# 搜索商品
https://api.bilibili.com/x/web-interface/wbi/search/all/v2
参数: keyword, search_type
需要: WBI签名

# 获取评论
https://api.bilibili.com/x/v2/reply
参数: type, oid, pn, ps, sort
不需要: WBI签名(但需要Cookie)

# 获取评论回复
https://api.bilibili.com/x/v2/reply/reply
参数: type, oid, root, pn, ps

# 获取WBI密钥
https://api.bilibili.com/x/web-interface/nav
返回: wbi_img.img_url, wbi_img.sub_url
```

### 评论区类型代码(type参数)

```
1: 视频
11: 图文
12: 专栏
17: 动态
```

