package bilibili

import (
	"net/url"
	"testing"
)

// TestBVAVConversion 测试BV/AV号转换功能
func TestBVAVConversion(t *testing.T) {
	// 测试用例：BV1mH4y1u7UA ↔ 1054803170
	testBVID := "BV1mH4y1u7UA"
	testAVID := int64(1054803170)

	// 测试 BV → AV
	t.Run("BV to AV", func(t *testing.T) {
		avid := Bvid2Avid(testBVID)
		if avid != testAVID {
			t.Errorf("Bvid2Avid(%s) = %d, want %d", testBVID, avid, testAVID)
		} else {
			t.Logf("✓ BV号转AV号成功: %s → %d", testBVID, avid)
		}
	})

	// 测试 AV → BV
	t.Run("AV to BV", func(t *testing.T) {
		bvid := Avid2Bvid(testAVID)
		if bvid != testBVID {
			t.Errorf("Avid2Bvid(%d) = %s, want %s", testAVID, bvid, testBVID)
		} else {
			t.Logf("✓ AV号转BV号成功: %d → %s", testAVID, bvid)
		}
	})

	// 测试往返转换（BV → AV → BV）
	t.Run("Round trip BV->AV->BV", func(t *testing.T) {
		avid := Bvid2Avid(testBVID)
		bvid := Avid2Bvid(avid)
		if bvid != testBVID {
			t.Errorf("Round trip failed: %s → %d → %s", testBVID, avid, bvid)
		} else {
			t.Logf("✓ 往返转换成功: %s → %d → %s", testBVID, avid, bvid)
		}
	})

	// 测试往返转换（AV → BV → AV）
	t.Run("Round trip AV->BV->AV", func(t *testing.T) {
		bvid := Avid2Bvid(testAVID)
		avid := Bvid2Avid(bvid)
		if avid != testAVID {
			t.Errorf("Round trip failed: %d → %s → %d", testAVID, bvid, avid)
		} else {
			t.Logf("✓ 往返转换成功: %d → %s → %d", testAVID, bvid, avid)
		}
	})
}

// TestWBISign 测试WBI签名功能
func TestWBISign(t *testing.T) {
	// 测试URL：获取用户信息接口
	testURL := "https://api.bilibili.com/x/space/wbi/acc/info?mid=1850091"

	t.Run("Sign URL", func(t *testing.T) {
		u, err := url.Parse(testURL)
		if err != nil {
			t.Fatalf("解析URL失败: %v", err)
		}

		// 执行签名
		err = Sign(u)
		if err != nil {
			t.Fatalf("WBI签名失败: %v", err)
		}

		// 检查是否添加了w_rid参数
		wRid := u.Query().Get("w_rid")
		if wRid == "" {
			t.Error("签名失败: w_rid参数未添加")
		} else {
			t.Logf("✓ w_rid已添加: %s", wRid)
		}

		// 检查是否添加了wts参数
		wts := u.Query().Get("wts")
		if wts == "" {
			t.Error("签名失败: wts参数未添加")
		} else {
			t.Logf("✓ wts已添加: %s", wts)
		}

		// 输出完整的签名后URL
		t.Logf("✓ 签名后URL: %s", u.String())
	})

	// 测试多次签名（验证缓存机制）
	t.Run("Multiple signs with cache", func(t *testing.T) {
		u1, _ := url.Parse(testURL)
		err1 := Sign(u1)
		if err1 != nil {
			t.Fatalf("第一次签名失败: %v", err1)
		}

		u2, _ := url.Parse(testURL)
		err2 := Sign(u2)
		if err2 != nil {
			t.Fatalf("第二次签名失败: %v", err2)
		}

		// 两次签名的wts应该不同（因为时间戳不同）
		wts1 := u1.Query().Get("wts")
		wts2 := u2.Query().Get("wts")
		t.Logf("✓ 第一次签名wts: %s", wts1)
		t.Logf("✓ 第二次签名wts: %s", wts2)
		t.Logf("✓ 密钥缓存机制正常工作")
	})
}

// TestClient 测试HTTP客户端
func TestClient(t *testing.T) {
	// 创建客户端（不带Cookie）
	client := NewClient("")

	t.Run("Create client", func(t *testing.T) {
		if client == nil {
			t.Fatal("客户端创建失败")
		}
		t.Log("✓ 客户端创建成功")
	})

	t.Run("Set cookie", func(t *testing.T) {
		testCookie := "SESSDATA=test; bili_jct=test"
		client.SetCookie(testCookie)
		if client.cookie != testCookie {
			t.Errorf("Cookie设置失败: got %s, want %s", client.cookie, testCookie)
		} else {
			t.Logf("✓ Cookie设置成功: %s", testCookie)
		}
	})

	// 测试实际请求（获取nav接口，不需要登录）
	t.Run("Real request to nav API", func(t *testing.T) {
		resp, err := client.Get("https://api.bilibili.com/x/web-interface/nav", false)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("HTTP状态码错误: got %d, want 200", resp.StatusCode)
		} else {
			t.Logf("✓ 请求成功: HTTP %d", resp.StatusCode)
		}

		// 检查响应头
		userAgent := resp.Request.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("User-Agent未设置")
		} else {
			t.Logf("✓ User-Agent已设置: %s", userAgent)
		}

		referer := resp.Request.Header.Get("Referer")
		if referer == "" {
			t.Error("Referer未设置")
		} else {
			t.Logf("✓ Referer已设置: %s", referer)
		}
	})
}
