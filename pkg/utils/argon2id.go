package utils

// Argon2id 参数配置
//const (
//	time      = 1
//	memory    = 64 * 1024
//	threads   = 4
//	keyLength = 32
//)
//
//// 生成随机salt函数
//func generateSalt(length int) []byte {
//	salt := make([]byte, length)
//	_, err := rand.Read(salt)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return salt
//}

// 生成密码哈希函数
//func hashPassword(password string) (string, error) {
//	salt := generateSalt(16)
//	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLength)
//	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
//	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
//	return fmt.Sprintf("%s:%s", b64Salt, b64Hash), nil
//}
//
//// 验证密码函数
//func verifyPassword(storedPassword, providedPassword string) bool {
//	parts := strings.Split(storedPassword, ":")
//	if len(parts) != 2 {
//		return false
//	}
//
//	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	storedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	newHash := argon2.IDKey([]byte(providedPassword), salt, time, memory, threads, keyLength)
//	return subtle.ConstantTimeCompare(storedHash, newHash) == 1
//}
