package service

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"strings"
)

func hashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func comparePasswords(storedPassword, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	if len(pwsalt) < 2 {
		return false, fmt.Errorf("did not provide a valid hash")
	}

	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
