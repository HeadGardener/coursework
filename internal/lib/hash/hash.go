package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 14
)

func GetStringHash(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(passwordHash)
}

func CompareHashAndString(passHash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(passHash, []byte(password))

	return err == nil
}
