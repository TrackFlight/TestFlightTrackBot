package utils

import "strings"

const base36Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase36(n int64) string {
	if n == 0 {
		return "0"
	}
	var encoded strings.Builder
	for n > 0 {
		rem := n % 36
		encoded.WriteByte(base36Chars[rem])
		n /= 36
	}
	runes := []rune(encoded.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func padLeft(s string, minLen int) string {
	if len(s) >= minLen {
		return s
	}
	return strings.Repeat("0", minLen-len(s)) + s
}

func EncodeTag(id int64) string {
	obfuscated := id ^ 0x64
	encoded := encodeBase36(obfuscated)
	return padLeft(encoded, 3)
}
