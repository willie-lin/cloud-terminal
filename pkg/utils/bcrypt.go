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

//
//func GenerateFromPassword(password []byte, cost int) ([]byte, error) {
//	p, err := newFromPassword(password, cost)
//	if err != nil {
//		return nil, err
//	}
//	return p.Hash(), nil
//}
//
//// CompareHashAndPassword compares a bcrypt hashed password with its possible
//// plaintext equivalent. Returns nil on success, or an error on failure.
//func CompareHashAndPassword(hashedPassword, password []byte) error {
//	p, err := newFromHash(hashedPassword)
//	if err != nil {
//		return err
//	}
//
//	otherHash, err := bcrypt(password, p.cost, p.salt)
//	if err != nil {
//		return err
//	}
//
//	otherP := &hashed{otherHash, p.salt, p.cost, p.major, p.minor}
//	if subtle.ConstantTimeCompare(p.Hash(), otherP.Hash()) == 1 {
//		return nil
//	}
//
//	return ErrMismatchedHashAndPassword
//}

// Cost returns the hashing cost used to create the given hashed
// password. When, in the future, the hashing cost of a password system needs
// to be increased in order to adjust for greater computational power, this
// function allows one to establish which passwords need to be updated.
//func Cost(hashedPassword []byte) (int, error) {
//	p, err := newFromHash(hashedPassword)
//	if err != nil {
//		return 0, err
//	}
//	return p.cost, nil
//}
