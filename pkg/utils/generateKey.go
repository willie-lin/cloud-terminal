package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomKey 生成随机密钥的函数
func GenerateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
