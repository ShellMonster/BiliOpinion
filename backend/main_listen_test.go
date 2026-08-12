package main

import (
	"os"
	"strings"
	"testing"
)

func TestMainRunsOnAPIListenAddr(t *testing.T) {
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	text := string(src)
	if !strings.Contains(text, "r.Run(api.ListenAddr)") {
		t.Fatal("main must bind via r.Run(api.ListenAddr)")
	}
	if strings.Contains(text, `r.Run(":8080")`) || strings.Contains(text, `r.Run("0.0.0.0`) {
		t.Fatal("main must not bind all interfaces")
	}
}
