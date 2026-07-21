package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	// EncryptedMarker 用于标识 auth_data 是否经过信封加密
	EncryptedMarker = "$encrypted"
	// AlgorithmMarker 加密算法标识
	AlgorithmMarker = "AES-256-GCM-ENVELOPE"
	// EnvKeyName 存放 KEK (Key Encryption Key) 的环境变量名
	EnvKeyName = "RESOURCE_ENCRYPTION_KEY"
)

var (
	defaultKEKOnce sync.Once
	cachedKEK      []byte
	kekError       error
)

// GetKEK 从环境变量获取或初始化主密钥（KEK）
// 支持 32 字节十六进制（64字符）、Base64 编码或原生 32 字节字符串
func GetKEK() ([]byte, error) {
	defaultKEKOnce.Do(func() {
		val := os.Getenv(EnvKeyName)
		if val == "" {
			log.Printf("WARNING: %s is not set. Using temporary dev key for auth_data envelope encryption.", EnvKeyName)
			// 开发默认 KEK (仅供开发测试，生产环境强烈建议配置环境变量)
			cachedKEK = []byte("cloud-terminal-dev-kek-32bytes-k")
			return
		}

		// 尝试十六进制解码 (64字符对应32字节)
		if len(val) == 64 {
			if decoded, err := hex.DecodeString(val); err == nil && len(decoded) == 32 {
				cachedKEK = decoded
				return
			}
		}

		// 尝试 Base64 解码
		if decoded, err := base64.StdEncoding.DecodeString(val); err == nil && len(decoded) == 32 {
			cachedKEK = decoded
			return
		}

		// 原生字符串长度校验
		if len(val) == 32 {
			cachedKEK = []byte(val)
			return
		}

		kekError = fmt.Errorf("invalid RESOURCE_ENCRYPTION_KEY format: must be 32 bytes (raw), 64 chars (hex), or base64 encoded 32 bytes")
	})

	if kekError != nil {
		return nil, kekError
	}
	return cachedKEK, nil
}

// encryptGCM 使用 AES-256-GCM 加密数据
func encryptGCM(key []byte, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptGCM 使用 AES-256-GCM 解密数据
func decryptGCM(key []byte, encodedCiphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptAuthData 使用信封加密算法（DEK + KEK）加密 auth_data 字典
// 最终生成带 "$encrypted": true 标识的加密字典
func EncryptAuthData(plaintext map[string]interface{}) (map[string]interface{}, error) {
	if plaintext == nil {
		return nil, nil
	}

	// 如果已经被加密，直接返回
	if isEnc, ok := plaintext[EncryptedMarker].(bool); ok && isEnc {
		return plaintext, nil
	}

	kek, err := GetKEK()
	if err != nil {
		return nil, fmt.Errorf("get KEK failed: %w", err)
	}

	// 1. 序列化明文
	jsonBytes, err := json.Marshal(plaintext)
	if err != nil {
		return nil, fmt.Errorf("marshal auth_data failed: %w", err)
	}

	// 2. 生成 32 字节随机 DEK (Data Encryption Key)
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, fmt.Errorf("generate DEK failed: %w", err)
	}

	// 3. 使用 DEK 加密明文 payload
	ciphertext, err := encryptGCM(dek, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("encrypt payload with DEK failed: %w", err)
	}

	// 4. 使用 KEK 加密 DEK
	encryptedDEK, err := encryptGCM(kek, dek)
	if err != nil {
		return nil, fmt.Errorf("encrypt DEK with KEK failed: %w", err)
	}

	return map[string]interface{}{
		EncryptedMarker: true,
		"alg":           AlgorithmMarker,
		"encrypted_dek": encryptedDEK,
		"ciphertext":    ciphertext,
	}, nil
}

// DecryptAuthData 解密信封加密的 auth_data 字典
// 如果数据未经过加密（兼容历史存量明文数据），则直接原样返回
func DecryptAuthData(authData map[string]interface{}) (map[string]interface{}, error) {
	if authData == nil {
		return nil, nil
	}

	// 判断是否为加密数据
	isEnc, ok := authData[EncryptedMarker].(bool)
	if !ok || !isEnc {
		// 存量明文直接兼容返回
		return authData, nil
	}

	encDEKStr, ok1 := authData["encrypted_dek"].(string)
	ciphertextStr, ok2 := authData["ciphertext"].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("invalid encrypted auth_data structure: missing encrypted_dek or ciphertext")
	}

	kek, err := GetKEK()
	if err != nil {
		return nil, fmt.Errorf("get KEK failed: %w", err)
	}

	// 1. 使用 KEK 解密出 DEK
	dek, err := decryptGCM(kek, encDEKStr)
	if err != nil {
		return nil, fmt.Errorf("decrypt DEK failed: %w", err)
	}

	// 2. 使用 DEK 解密 payload
	jsonBytes, err := decryptGCM(dek, ciphertextStr)
	if err != nil {
		return nil, fmt.Errorf("decrypt payload failed: %w", err)
	}

	// 3. 反序列化 JSON 到字典
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal decrypted payload failed: %w", err)
	}

	return result, nil
}

// EncryptString 使用信封加密（KEK）加密任意字符串（如 totp_secret）
// 返回格式：base64(nonce || GCM_ciphertext_of_DEK_encrypted_plaintext)
// 实际上是两层加密：plaintext → DEK(AES-256-GCM) → encrypted_dek(KEK) → 存 "$enc:<encDEK>:<ciphertext>"
func EncryptString(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	kek, err := GetKEK()
	if err != nil {
		return "", fmt.Errorf("get KEK failed: %w", err)
	}

	// 生成 32 字节随机 DEK
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return "", fmt.Errorf("generate DEK failed: %w", err)
	}

	// 用 DEK 加密明文
	ciphertext, err := encryptGCM(dek, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("encrypt with DEK failed: %w", err)
	}

	// 用 KEK 加密 DEK
	encDEK, err := encryptGCM(kek, dek)
	if err != nil {
		return "", fmt.Errorf("encrypt DEK with KEK failed: %w", err)
	}

	// 拼装成可存储的格式：$enc:<base64_encDEK>:<base64_ciphertext>
	return fmt.Sprintf("$enc:%s:%s", encDEK, ciphertext), nil
}

// DecryptString 解密由 EncryptString 生成的加密字符串
// 兼容未加密的明文（不以 "$enc:" 开头时直接返回，用于迁移旧数据）
func DecryptString(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}
	// 兼容历史明文存量数据
	if len(encrypted) < 5 || encrypted[:5] != "$enc:" {
		return encrypted, nil
	}

	// 解析 $enc:<encDEK>:<ciphertext>
	rest := encrypted[5:] // 去掉 "$enc:"
	// 找到第一个 ":" 分隔 encDEK 和 ciphertext
	// encDEK 是 base64，不包含 ":"，但 ciphertext 也是 base64 不含 ":"，直接 SplitN 2
	var encDEK, ciphertext string
	for i := 0; i < len(rest); i++ {
		if rest[i] == ':' {
			encDEK = rest[:i]
			ciphertext = rest[i+1:]
			break
		}
	}
	if encDEK == "" || ciphertext == "" {
		return "", errors.New("invalid encrypted string format")
	}

	kek, err := GetKEK()
	if err != nil {
		return "", fmt.Errorf("get KEK failed: %w", err)
	}

	// 用 KEK 解密 DEK
	dek, err := decryptGCM(kek, encDEK)
	if err != nil {
		return "", fmt.Errorf("decrypt DEK failed: %w", err)
	}

	// 用 DEK 解密明文
	plainBytes, err := decryptGCM(dek, ciphertext)
	if err != nil {
		return "", fmt.Errorf("decrypt plaintext failed: %w", err)
	}

	return string(plainBytes), nil
}
