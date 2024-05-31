package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
