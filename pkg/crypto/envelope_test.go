package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"reflect"
	"sync"
	"testing"
)

func resetKEKForTest() {
	defaultKEKOnce = sync.Once{}
	cachedKEK = nil
	kekError = nil
}

func TestEnvelopeEncryptionRoundTrip(t *testing.T) {
	resetKEKForTest()
	os.Setenv(EnvKeyName, "0123456789abcdef0123456789abcdef") // 32 chars raw
	defer os.Unsetenv(EnvKeyName)

	plaintext := map[string]interface{}{
		"username": "root",
		"password": "superSecretPassword123!",
		"ssh_key":  "-----BEGIN RSA PRIVATE KEY-----\n...",
		"port":     float64(22),
	}

	encrypted, err := EncryptAuthData(plaintext)
	if err != nil {
		t.Fatalf("EncryptAuthData failed: %v", err)
	}

	// 验证已加密字段标记
	if isEnc, ok := encrypted[EncryptedMarker].(bool); !ok || !isEnc {
		t.Fatalf("Expected $encrypted: true marker in ciphertext map")
	}
	if alg, ok := encrypted["alg"].(string); !ok || alg != AlgorithmMarker {
		t.Fatalf("Expected alg: %s", AlgorithmMarker)
	}
	if _, ok := encrypted["encrypted_dek"].(string); !ok {
		t.Fatalf("Missing encrypted_dek field")
	}
	if _, ok := encrypted["ciphertext"].(string); !ok {
		t.Fatalf("Missing ciphertext field")
	}

	// 幂等测试：再次加密应当直接返回原加密内容
	reEncrypted, err := EncryptAuthData(encrypted)
	if err != nil {
		t.Fatalf("Re-encrypting already encrypted map failed: %v", err)
	}
	if !reflect.DeepEqual(encrypted, reEncrypted) {
		t.Fatalf("EncryptAuthData is not idempotent on encrypted map")
	}

	// 解密测试
	decrypted, err := DecryptAuthData(encrypted)
	if err != nil {
		t.Fatalf("DecryptAuthData failed: %v", err)
	}

	if !reflect.DeepEqual(plaintext, decrypted) {
		t.Fatalf("Decrypted map does not match original. Got: %+v, want: %+v", decrypted, plaintext)
	}
}

func TestLegacyUnencryptedDataPassthrough(t *testing.T) {
	resetKEKForTest()

	legacyData := map[string]interface{}{
		"username": "admin",
		"password": "legacy_password",
	}

	decrypted, err := DecryptAuthData(legacyData)
	if err != nil {
		t.Fatalf("DecryptAuthData on legacy data failed: %v", err)
	}
	if !reflect.DeepEqual(legacyData, decrypted) {
		t.Fatalf("Expected legacy data to pass through unchanged")
	}
}

func TestGetKEKFormats(t *testing.T) {
	tests := []struct {
		name      string
		envVal    string
		wantBytes int
		wantErr   bool
	}{
		{
			name:      "Default dev fallback",
			envVal:    "",
			wantBytes: 32,
			wantErr:   false,
		},
		{
			name:      "Raw 32 bytes string",
			envVal:    "12345678901234567890123456789012",
			wantBytes: 32,
			wantErr:   false,
		},
		{
			name:      "Hex encoded 32 bytes (64 chars)",
			envVal:    hex.EncodeToString([]byte("12345678901234567890123456789012")),
			wantBytes: 32,
			wantErr:   false,
		},
		{
			name:      "Base64 encoded 32 bytes",
			envVal:    base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012")),
			wantBytes: 32,
			wantErr:   false,
		},
		{
			name:    "Invalid length string",
			envVal:  "tooshort",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetKEKForTest()
			if tt.envVal != "" {
				os.Setenv(EnvKeyName, tt.envVal)
			} else {
				os.Unsetenv(EnvKeyName)
			}
			defer os.Unsetenv(EnvKeyName)

			key, err := GetKEK()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetKEK() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(key) != tt.wantBytes {
				t.Fatalf("GetKEK() key length = %d, want %d", len(key), tt.wantBytes)
			}
		})
	}
}

func TestTamperedCiphertext(t *testing.T) {
	resetKEKForTest()
	os.Setenv(EnvKeyName, "12345678901234567890123456789012")
	defer os.Unsetenv(EnvKeyName)

	plaintext := map[string]interface{}{"secret": "top-secret"}
	encrypted, _ := EncryptAuthData(plaintext)

	// 篡改 ciphertext
	encrypted["ciphertext"] = "aW52YWxpZGNpcGhlcnRleHQ=" // invalid base64/GCM

	_, err := DecryptAuthData(encrypted)
	if err == nil {
		t.Fatalf("Expected decryption error on tampered ciphertext, got nil")
	}
}

// ========== EncryptString / DecryptString Tests ==========

func TestEncryptDecryptStringRoundTrip(t *testing.T) {
	resetKEKForTest()
	os.Setenv(EnvKeyName, "12345678901234567890123456789012")
	defer os.Unsetenv(EnvKeyName)

	cases := []string{
		"JBSWY3DPEHPK3PXP",                             // 典型 TOTP base32 secret
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ234567",              // 长 secret
		"hello world",                                   // 普通字符串
		"中文内容",                                        // Unicode
		"特殊字符 !@#$%^&*()_+-=[]{}|;':\",./<>?",        // 特殊字符
	}

	for _, original := range cases {
		encrypted, err := EncryptString(original)
		if err != nil {
			t.Fatalf("EncryptString(%q) error: %v", original, err)
		}

		// 必须以 $enc: 开头
		if len(encrypted) < 5 || encrypted[:5] != "$enc:" {
			t.Fatalf("EncryptString(%q) result missing $enc: prefix: %q", original, encrypted)
		}

		decrypted, err := DecryptString(encrypted)
		if err != nil {
			t.Fatalf("DecryptString(%q) error: %v", encrypted, err)
		}
		if decrypted != original {
			t.Fatalf("Round-trip mismatch: got %q, want %q", decrypted, original)
		}
	}
}

func TestEncryptStringEmpty(t *testing.T) {
	resetKEKForTest()

	enc, err := EncryptString("")
	if err != nil {
		t.Fatalf("EncryptString('') should not error: %v", err)
	}
	if enc != "" {
		t.Fatalf("EncryptString('') should return empty string, got %q", enc)
	}

	dec, err := DecryptString("")
	if err != nil {
		t.Fatalf("DecryptString('') should not error: %v", err)
	}
	if dec != "" {
		t.Fatalf("DecryptString('') should return empty string, got %q", dec)
	}
}

func TestDecryptStringLegacyPlaintext(t *testing.T) {
	resetKEKForTest()

	// 历史明文（不以 $enc: 开头）应该原样返回
	legacy := "JBSWY3DPEHPK3PXP"
	dec, err := DecryptString(legacy)
	if err != nil {
		t.Fatalf("DecryptString(legacy) should not error: %v", err)
	}
	if dec != legacy {
		t.Fatalf("Legacy passthrough failed: got %q, want %q", dec, legacy)
	}
}

func TestDecryptStringTampered(t *testing.T) {
	resetKEKForTest()
	os.Setenv(EnvKeyName, "12345678901234567890123456789012")
	defer os.Unsetenv(EnvKeyName)

	original := "JBSWY3DPEHPK3PXP"
	encrypted, _ := EncryptString(original)

	// 篡改密文最后几位
	tampered := encrypted[:len(encrypted)-4] + "XXXX"
	_, err := DecryptString(tampered)
	if err == nil {
		t.Fatalf("Expected error when decrypting tampered ciphertext, got nil")
	}
}

func TestEncryptStringNonDeterministic(t *testing.T) {
	resetKEKForTest()
	os.Setenv(EnvKeyName, "12345678901234567890123456789012")
	defer os.Unsetenv(EnvKeyName)

	// 相同明文每次加密结果应不同（因为随机 nonce + DEK）
	plain := "JBSWY3DPEHPK3PXP"
	enc1, _ := EncryptString(plain)
	enc2, _ := EncryptString(plain)
	if enc1 == enc2 {
		t.Fatal("EncryptString should produce different ciphertext each time (non-deterministic)")
	}
}

