package bilibili

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// WbiKeys 存储WBI密钥和缓存时间
type WbiKeys struct {
	Img            string    // img_key 图片密钥
	Sub            string    // sub_key 副密钥
	Mixin          string    // 混合后的最终密钥
	lastUpdateTime time.Time // 上次更新时间（用于1小时缓存）
}

// 全局WBI密钥缓存
var wbiKeys WbiKeys

// Sign 对URL进行WBI签名（公开函数）
// 参数：
//   - u: 需要签名的URL对象
//
// 返回：
//   - error: 签名失败时返回错误
func Sign(u *url.URL) error {
	return wbiKeys.Sign(u)
}

// Sign WBI签名实现
// 核心流程：
//  1. 更新密钥（如果缓存过期）
//  2. 移除特殊字符
//  3. 添加时间戳wts
//  4. 参数排序 + mixin key → MD5 → w_rid
func (wk *WbiKeys) Sign(u *url.URL) (err error) {
	// 更新密钥（如果需要）- 缓存1小时
	if err = wk.update(false); err != nil {
		return err
	}

	values := u.Query()
	// 移除可能干扰签名的特殊字符
	values = removeUnwantedChars(values, '!', '\'', '(', ')', '*')

	// 添加当前时间戳
	values.Set("wts", strconv.FormatInt(time.Now().Unix(), 10))

	// 编码参数（自动排序键）并添加盐值
	hash := md5.Sum([]byte(values.Encode() + wk.Mixin))

	// 设置签名
	values.Set("w_rid", hex.EncodeToString(hash[:]))
	u.RawQuery = values.Encode()
	return nil
}

// update 如果缓存过期则获取新密钥
// 参数：
//   - purge: 是否强制更新（忽略缓存）
//
// 返回：
//   - error: 更新失败时返回错误
func (wk *WbiKeys) update(purge bool) error {
	// 如果未强制更新且缓存未过期（1小时内），直接返回
	if !purge && time.Since(wk.lastUpdateTime) < time.Hour {
		return nil
	}

	// 从B站nav接口获取密钥
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应JSON
	nav := Nav{}
	err = json.Unmarshal(body, &nav)
	if err != nil {
		return err
	}

	// 检查响应码（0表示成功，-101表示未登录但仍可获取密钥）
	if nav.Code != 0 && nav.Code != -101 {
		return fmt.Errorf("unexpected code: %d", nav.Code)
	}

	img := nav.Data.WbiImg.ImgUrl
	sub := nav.Data.WbiImg.SubUrl

	// 提取密钥文件名（去除扩展名）
	imgParts := strings.Split(img, "/")
	subParts := strings.Split(sub, "/")
	wk.Img = strings.TrimSuffix(imgParts[len(imgParts)-1], ".png")
	wk.Sub = strings.TrimSuffix(subParts[len(subParts)-1], ".png")

	// 生成混合密钥
	wk.mixin()
	wk.lastUpdateTime = time.Now()
	return nil
}

// mixin 使用查找表生成最终盐值
// 算法：将img_key和sub_key拼接后，按照mixinKeyEncTab的顺序重新排列
func (wk *WbiKeys) mixin() {
	var mixin [32]byte
	wbi := wk.Img + wk.Sub
	for i := range mixin {
		mixin[i] = wbi[mixinKeyEncTab[i]]
	}
	wk.Mixin = string(mixin[:])
}

// mixinKeyEncTab 混合密钥查找表
// 用于打乱img_key和sub_key的顺序，生成最终的mixin key
var mixinKeyEncTab = [...]int{
	46, 47, 18, 2, 53, 8, 23, 32,
	15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19,
	29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61,
	26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63,
	57, 62, 11, 36, 20, 34, 44, 52,
}

// removeUnwantedChars 清理参数中的特殊字符
// 参数：
//   - v: URL参数
//   - chars: 需要移除的字符列表
//
// 返回：
//   - url.Values: 清理后的参数
func removeUnwantedChars(v url.Values, chars ...byte) url.Values {
	b := []byte(v.Encode())
	for _, c := range chars {
		b = bytes.ReplaceAll(b, []byte{c}, nil)
	}
	s, err := url.ParseQuery(string(b))
	if err != nil {
		panic(err)
	}
	return s
}

// Nav B站nav接口响应结构
type Nav struct {
	Code int `json:"code"`
	Data struct {
		WbiImg struct {
			ImgUrl string `json:"img_url"` // img_key URL
			SubUrl string `json:"sub_url"` // sub_key URL
		} `json:"wbi_img"`
	} `json:"data"`
}
