package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// HashPassword receive a string password and hash it with salt
func HashPassword(password string) (string, error) {
	// example for making salt - https://play.golang.org/p/_Aw6WeWC42I
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	sHash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(sHash), hex.EncodeToString(salt))

	return hashedPW, nil
}

// ComparePasswords get the string password and hashed it with salt
// get from stored password and compare both hashed passwords
func ComparePasswords(storedPassword string, suppliedPassword string) (bool, error) {
	pwSalt := strings.Split(storedPassword, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwSalt[1])
	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	sHash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, fmt.Errorf("unable to hash user password")
	}

	return hex.EncodeToString(sHash) == pwSalt[0], nil
}