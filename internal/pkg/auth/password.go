package auth

import (
	"crypto/md5"
	"io"
)

func HashPassword(password string) (string, error) {
	hash := md5.New()
	if _, err := io.WriteString(hash, password); err != nil {
		return "", err
	}
	hashed := hash.Sum(nil)
	return string(hashed), nil
}

func EqualPassword(password, hashed string) (bool, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return false, err
	}

	return hashedPassword == hashed, nil
}
