package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateIDToken generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateIDTokenString, but the signature makes this fairly clear
func GenerateIDToken(user *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := model.IDTokenCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

// GenerateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func GenerateRefreshToken(uid uuid.UUID, key string, exp int64) (*model.RefreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib
	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := model.RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &model.RefreshTokenData{
		SignedStringToken: signedToken,
		ID:                tokenID,
		ExpiresIn:         tokenExp.Sub(currentTime),
	}, nil
}

// ValidateIDToken returns the token's claims if the token is valid
func ValidateIDToken(tokenString string, key *rsa.PublicKey) (*model.IDTokenCustomClaims, error) {
	claims := &model.IDTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*model.IDTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

// ValidateRefreshToken uses the secret key to validate a refresh token
func ValidateRefreshToken(tokenString string, key string) (*model.RefreshTokenCustomClaims, error) {
	claims := &model.RefreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("refresh token is invalid")
	}

	claims, ok := token.Claims.(*model.RefreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("refresh token valid but couldn't parse claims")
	}

	return claims, nil
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(rsaPublicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(rsaPublicKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}
