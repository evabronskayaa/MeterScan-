package util

import (
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 16

func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	return bytes, err
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), hash)
	return err == nil
}
