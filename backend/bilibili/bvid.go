package bilibili

import (
	"strings"
)

// BV/AV号转换相关常量
var (
	XOR_CODE = int64(23442827791579)                                        // XOR混淆码
	MAX_CODE = int64(2251799813685247)                                      // 最大编码值
	CHARTS   = "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf" // Base58字符表
	PAUL_NUM = int64(58)                                                    // 进制数（Base58）
)

// swapString 交换字符串中两个位置的字符
// 参数：
//   - s: 原始字符串
//   - x, y: 需要交换的两个位置索引
//
// 返回：
//   - string: 交换后的字符串
func swapString(s string, x, y int) string {
	chars := []rune(s)
	chars[x], chars[y] = chars[y], chars[x]
	return string(chars)
}

// Bvid2Avid 将BV号转换为AV号
// 算法流程：
//  1. 交换字符位置（3↔9, 4↔7）- 反混淆
//  2. 去除"BV1"前缀
//  3. Base58解码
//  4. XOR解密
//
// 参数：
//   - bvid: BV号字符串（如 "BV1mH4y1u7UA"）
//
// 返回：
//   - avid: AV号（如 1054803170），如果bvid无效则返回0
//
// 示例：
//
//	Bvid2Avid("BV1mH4y1u7UA") // 返回 1054803170
func Bvid2Avid(bvid string) (avid int64) {
	// 验证BVID格式：必须以"BV"开头且长度至少为12
	if len(bvid) < 12 || !strings.HasPrefix(bvid, "BV") {
		return 0
	}

	// 反向交换字符位置
	s := swapString(swapString(bvid, 3, 9), 4, 7)
	// 去除"BV1"前缀，保留后面的Base58编码部分
	bv1 := string([]rune(s)[3:])

	// Base58解码
	temp := int64(0)
	for _, c := range bv1 {
		idx := strings.IndexRune(CHARTS, c)
		if idx < 0 {
			return 0 // 无效字符
		}
		temp = temp*PAUL_NUM + int64(idx)
	}

	// XOR解密
	avid = (temp & MAX_CODE) ^ XOR_CODE
	return
}

// Avid2Bvid 将AV号转换为BV号
// 算法流程：
//  1. XOR加密
//  2. 添加高位标志位
//  3. Base58编码
//  4. 添加"BV1"前缀
//  5. 交换字符位置（3↔9, 4↔7）- 混淆
//
// 参数：
//   - avid: AV号（如 1054803170）
//
// 返回：
//   - bvid: BV号字符串（如 "BV1mH4y1u7UA"）
//
// 示例：
//
//	Avid2Bvid(1054803170) // 返回 "BV1mH4y1u7UA"
func Avid2Bvid(avid int64) (bvid string) {
	// 初始化数组，前3位固定为"B", "V", "1"
	arr := [12]string{"B", "V", "1"}
	bvIdx := len(arr) - 1

	// XOR加密并添加高位标志
	temp := (avid | (MAX_CODE + 1)) ^ XOR_CODE

	// Base58编码（从后往前填充）
	for temp > 0 {
		idx := temp % PAUL_NUM
		arr[bvIdx] = string(CHARTS[idx])
		temp /= PAUL_NUM
		bvIdx--
	}

	// 拼接字符串
	raw := strings.Join(arr[:], "")

	// 交换字符位置进行混淆
	bvid = swapString(swapString(raw, 3, 9), 4, 7)
	return
}
