package comment

import "testing"

func TestCleanBrandName_Normal(t *testing.T) {
	out := CleanBrandName("小佩", nil)
	if out != "小佩" {
		t.Fatalf("expected %q, got %q", "小佩", out)
	}
}

func TestCleanBrandName_WithSlash(t *testing.T) {
	out := CleanBrandName("喵洁易/Catlink", nil)
	if out != "喵洁易" {
		t.Fatalf("expected %q, got %q", "喵洁易", out)
	}
}

func TestCleanBrandName_MatchKnown(t *testing.T) {
	out := CleanBrandName("喵洁易/Catlink", []string{"Catlink", "小佩"})
	if out != "Catlink" {
		t.Fatalf("expected %q, got %q", "Catlink", out)
	}
}

func TestCleanBrandName_MultipleSlash(t *testing.T) {
	out := CleanBrandName("有陪/小佩/小米/糯雪", []string{"小米"})
	if out != "小米" {
		t.Fatalf("expected %q, got %q", "小米", out)
	}
}

func TestCleanBrandName_Empty(t *testing.T) {
	out := CleanBrandName("", nil)
	if out != "" {
		t.Fatalf("expected %q, got %q", "", out)
	}
}

func TestCleanBrandName_Unknown(t *testing.T) {
	out := CleanBrandName("未知", nil)
	if out != "未知" {
		t.Fatalf("expected %q, got %q", "未知", out)
	}
}

func TestCleanModelName_Normal(t *testing.T) {
	out := CleanModelName("V12")
	if out != "V12" {
		t.Fatalf("expected %q, got %q", "V12", out)
	}
}

func TestCleanModelName_WithSlash(t *testing.T) {
	out := CleanModelName("二代/通用")
	if out != "二代" {
		t.Fatalf("expected %q, got %q", "二代", out)
	}
}

func TestCleanModelName_Descriptive(t *testing.T) {
	out := CleanModelName("新款")
	if out != "通用" {
		t.Fatalf("expected %q, got %q", "通用", out)
	}
}
