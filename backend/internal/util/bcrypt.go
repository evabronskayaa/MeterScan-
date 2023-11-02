package util

import (
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = 16

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
}

func CheckPasswordHash(password string, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}
