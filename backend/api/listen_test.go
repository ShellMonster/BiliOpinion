package api

import (
	"strings"
	"testing"
)

func TestListenAddr_IsLoopbackOnly(t *testing.T) {
	if ListenAddr != "127.0.0.1:8080" {
		t.Fatalf("server must listen on 127.0.0.1:8080, got %q", ListenAddr)
	}
	if strings.HasPrefix(ListenAddr, ":") || strings.HasPrefix(ListenAddr, "0.0.0.0") {
		t.Fatalf("must not bind all interfaces, got %q", ListenAddr)
	}
}
