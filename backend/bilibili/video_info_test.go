package bilibili

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockTransport struct {
	responseBody string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(m.responseBody)),
	}, nil
}

func TestGetVideoInfo_Success(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: &mockTransport{
				responseBody: `{
					"code": 0,
					"message": "0",
					"data": {
						"bvid": "BV1mH4y1u7UA",
						"aid": 1054803170,
						"title": "测试视频标题",
						"owner": {
							"mid": 123456,
							"name": "测试UP主"
						},
						"stat": {
							"view": 100000,
							"danmaku": 1000,
							"reply": 500
						},
						"pubdate": 1700000000,
						"pic": "https://example.com/cover.jpg"
					}
				}`,
			},
		},
	}

	info, err := client.GetVideoInfo("BV1mH4y1u7UA")
	if err != nil {
		t.Fatalf("GetVideoInfo() error = %v", err)
	}

	if info.BVID != "BV1mH4y1u7UA" {
		t.Errorf("BVID = %q, want %q", info.BVID, "BV1mH4y1u7UA")
	}
	if info.Title != "测试视频标题" {
		t.Errorf("Title = %q, want %q", info.Title, "测试视频标题")
	}
	if info.Author != "测试UP主" {
		t.Errorf("Author = %q, want %q", info.Author, "测试UP主")
	}
	if info.PlayCount != 100000 {
		t.Errorf("PlayCount = %d, want %d", info.PlayCount, 100000)
	}
	if info.CommentCount != 500 {
		t.Errorf("CommentCount = %d, want %d", info.CommentCount, 500)
	}
}

func TestGetVideoInfo_NotFound(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: &mockTransport{
				responseBody: `{"code": -404, "message": "啥都木有", "data": null}`,
			},
		},
	}

	_, err := client.GetVideoInfo("BV1notexist00")
	if err == nil {
		t.Fatal("GetVideoInfo() expected error for non-existent video")
	}
}

func TestGetVideoInfo_APIError(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: &mockTransport{
				responseBody: `{"code": -400, "message": "请求错误", "data": null}`,
			},
		},
	}

	_, err := client.GetVideoInfo("invalid")
	if err == nil {
		t.Fatal("GetVideoInfo() expected error for API error")
	}
}
