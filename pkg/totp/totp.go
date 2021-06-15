package totp

import (
	cotp "github.com/pquerna/otp"
	ctotp "github.com/pquerna/otp/totp"
)

type GenerateOpts ctotp.GenerateOpts

// 生成新的to
func NewTOPO(opt GenerateOpts) (*cotp.Key, error) {

	return ctotp.Generate(ctotp.GenerateOpts(opt))

}

func Validate(code string, secret string) bool {
	if secret == "" {
		return true
	}
	return ctotp.Validate(code, secret)
}
