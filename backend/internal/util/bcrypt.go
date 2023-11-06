package util

import (
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 16

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
}

func CheckPasswordHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
