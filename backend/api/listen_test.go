package api

import (
	"net"
	"testing"
)

func TestListenAddr_IsLoopbackOnly(t *testing.T) {
	host, port, err := net.SplitHostPort(ListenAddr)
	if err != nil {
		t.Fatalf("ListenAddr must be host:port, got %q: %v", ListenAddr, err)
	}
	if port != "8080" {
		t.Fatalf("port must be 8080, got %q", port)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		t.Fatalf("host must be a loopback IP, got %q", host)
	}
}
