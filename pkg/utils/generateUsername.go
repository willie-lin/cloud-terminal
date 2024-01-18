package utils

import (
	"math/rand"
	"strings"
)

// GenerateUsername 生成唯一用户名
func GenerateUsername() string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var digits = []rune("0123456789")

	var b strings.Builder
	for i := 0; i < 6; i++ {
		b.WriteRune(letters[rand.Intn(len(letters))])
	}
	for i := 0; i < 10; i++ {
		b.WriteRune(digits[rand.Intn(len(digits))])
	}
	return b.String()
}
