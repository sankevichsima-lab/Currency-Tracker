package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash with salt from a plain text password.
// HashPassword генерирует хэш с солью на основе открытого пароля через bcrypt.
func HashPassword(password string) (string, error) {
	pas_byte := []byte(password)
	passBytes, err := bcrypt.GenerateFromPassword(pas_byte, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("problem with generate password: %w", err)
	}
	return string(passBytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// CheckPassword сравнивает захэшированный пароль из БД с открытым паролем от пользователя.
func CheckPassword(password, hash string) error {
	passByte, hashByte := []byte(password), []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashByte, passByte)
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err)
	}
	return nil
}
