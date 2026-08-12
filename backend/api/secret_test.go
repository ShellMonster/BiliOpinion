package api

import "testing"

func TestMaskSecret_EmptyStaysEmpty(t *testing.T) {
	if got := MaskSecret(""); got != "" {
		t.Fatalf("empty secret should stay empty, got %q", got)
	}
}

func TestMaskSecret_ShortFullyMasked(t *testing.T) {
	if got := MaskSecret("ab"); got != "****" {
		t.Fatalf("short secret should be fully masked, got %q", got)
	}
}

func TestMaskSecret_KeepsLastFourNeverFullValue(t *testing.T) {
	secret := "sk-live-super-secret-key-1234"
	got := MaskSecret(secret)
	if got == secret {
		t.Fatal("masked secret must not equal the full value")
	}
	if got != "****1234" {
		t.Fatalf("expected ****1234, got %q", got)
	}
}
