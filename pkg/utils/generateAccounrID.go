package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
)

func generateID() string {
	id := uuid.New().String()
	hash := sha256.Sum256([]byte(id))
	num := binary.BigEndian.Uint64(hash[:8])       // 取前 8 字节转换为 uint64
	return fmt.Sprintf("%012d", num%1000000000000) //取模并格式化
}
