package api

// MaskSecret hides all but the last 4 runes of a secret.
// Empty stays empty. Values of 4 runes or fewer become "****".
func MaskSecret(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= 4 {
		return "****"
	}
	return "****" + string(runes[len(runes)-4:])
}

func shouldPreserveSecret(incoming, existing string) bool {
	if incoming == "" {
		return true
	}
	if existing != "" && incoming == MaskSecret(existing) {
		return true
	}
	return false
}
