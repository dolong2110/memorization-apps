package utils

import (
	"crypto/rsa"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/google/uuid"
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
	log.Printf("KEY STR: " + ss)
	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

// GenerateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func GenerateRefreshToken(uid uuid.UUID, key string, exp int64) (*model.RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom()         // v4 uuid in the google uuid lib

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

	return &model.RefreshToken{
		SignedTokenString: signedToken,
		ID:                tokenID.String(),
		ExpiresIn:         tokenExp.Sub(currentTime),
	}, nil
}