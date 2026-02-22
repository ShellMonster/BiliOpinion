package bilibili

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ParseVideoURL 解析B站视频URL，提取BV号
// 支持的URL格式：
//   - bilibili.com/video/BVxxx
//   - bilibili.com/video/avxxx (自动转换为BV号)
//   - m.bilibili.com/video/BVxxx
//   - 带参数链接（忽略 ? 后的参数）
//
// 参数：
//   - url: B站视频链接
//
// 返回：
//   - bvid: BV号（如 "BV1mH4y1u7UA"）
//   - err: 错误信息，如果URL无效则返回错误
//
// 示例：
//
//	ParseVideoURL("https://www.bilibili.com/video/BV1mH4y1u7UA")
//	  → "BV1mH4y1u7UA", nil
//	ParseVideoURL("https://www.bilibili.com/video/av1054803170")
//	  → "BV1mH4y1u7UA", nil (转换后)
func ParseVideoURL(url string) (bvid string, err error) {
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	url = strings.TrimSuffix(url, "/")

	bvPattern := regexp.MustCompile(`bilibili\.com/video/(BV[a-zA-Z0-9]{10})`)
	if matches := bvPattern.FindStringSubmatch(url); len(matches) > 1 {
		bvid = matches[1]
		if !strings.HasPrefix(bvid, "BV") || len(bvid) != 12 {
			return "", errors.New("无效的BV号格式")
		}
		return bvid, nil
	}

	avPattern := regexp.MustCompile(`bilibili\.com/video/av(\d+)`)
	if matches := avPattern.FindStringSubmatch(url); len(matches) > 1 {
		avid, parseErr := strconv.ParseInt(matches[1], 10, 64)
		if parseErr != nil {
			return "", errors.New("无效的AV号格式")
		}
		bvid = Avid2Bvid(avid)
		if bvid == "" {
			return "", errors.New("AV号转换失败")
		}
		return bvid, nil
	}

	if strings.Contains(url, "b23.tv") {
		return "", errors.New("不支持b23.tv短链接，请使用完整的bilibili.com链接")
	}

	return "", errors.New("无效的B站视频链接")
}
