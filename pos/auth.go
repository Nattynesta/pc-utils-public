package main

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

func isSHA256Hash(hash string) bool {
	return len(hash) == 64 && !isBcryptHash(hash)
}

func isBcryptHash(hash string) bool {
	return len(hash) >= 60 && hash[:4] == "$2a$"
}

func migrateHash(plain, oldHash string) (string, error) {
	h := sha256.Sum256([]byte(plain))
	if fmt.Sprintf("%x", h) != oldHash {
		return "", fmt.Errorf("password does not match old hash")
	}
	return HashPassword(plain)
}
