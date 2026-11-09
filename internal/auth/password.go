package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pas_byte := []byte(password)
	passBytes, err := bcrypt.GenerateFromPassword(pas_byte, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("problem with generate password: %w", err)
	}
	return string(passBytes), nil
}

func CheckPassword(password, hash string) error {
	passByte, hashByte := []byte(password), []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashByte, passByte)
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err)
	}
	return nil
}
