package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

var ErrMismatchedHashAndPassword = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")

func GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
