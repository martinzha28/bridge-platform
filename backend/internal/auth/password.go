package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ValidatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	var hasUpper, hasSpecial bool
	for _, c := range p {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			hasSpecial = true
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
