package bilibili

import (
	"testing"
)

func TestParseVideoURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantBvid  string
		wantError bool
	}{
		{
			name:      "标准BV号链接",
			url:       "https://www.bilibili.com/video/BV1mH4y1u7UA",
			wantBvid:  "BV1mH4y1u7UA",
			wantError: false,
		},
		{
			name:      "AV号链接(转换后)",
			url:       "https://www.bilibili.com/video/av1054803170",
			wantBvid:  "BV1mH4y1u7UA",
			wantError: false,
		},
		{
			name:      "移动端BV号链接",
			url:       "https://m.bilibili.com/video/BV1mH4y1u7UA",
			wantBvid:  "BV1mH4y1u7UA",
			wantError: false,
		},
		{
			name:      "带参数的BV号链接",
			url:       "https://www.bilibili.com/video/BV1mH4y1u7UA?p=1",
			wantBvid:  "BV1mH4y1u7UA",
			wantError: false,
		},
		{
			name:      "非B站链接",
			url:       "https://youtube.com/xxx",
			wantBvid:  "",
			wantError: true,
		},
		{
			name:      "b23短链接(不支持)",
			url:       "https://b23.tv/BV1mH4y1u7UA",
			wantBvid:  "",
			wantError: true,
		},
		{
			name:      "带多个参数的链接",
			url:       "https://www.bilibili.com/video/BV1mH4y1u7UA?p=1&spm_id_from=333.788",
			wantBvid:  "BV1mH4y1u7UA",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBvid, err := ParseVideoURL(tt.url)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseVideoURL() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if gotBvid != tt.wantBvid {
				t.Errorf("ParseVideoURL() = %v, want %v", gotBvid, tt.wantBvid)
			}
		})
	}
}
