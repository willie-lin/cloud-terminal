package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	cotp "github.com/pquerna/otp"
	ctotp "github.com/pquerna/otp/totp"
	"time"
)

type GenerateOpts ctotp.GenerateOpts

// NewTOPO 生成新的totp
func NewTOPO(opt GenerateOpts) (*cotp.Key, error) {

	return ctotp.Generate(ctotp.GenerateOpts(opt))

}

// Validate 验证totp
func Validate(code string, secret string) bool {
	if secret == "" {
		return true
	}
	return ctotp.Validate(code, secret)
}

func Hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

func Totp(key []byte, t time.Time, digits int) int {
	return Hotp(key, uint64(t.Unix())/30, digits)
}
