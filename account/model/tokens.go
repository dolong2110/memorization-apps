package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// Token used for returning pairs of id and refresh tokens
type Token struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// IDTokenCustomClaims holds structure of jwt claims of idToken
type IDTokenCustomClaims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}

// RefreshToken holds the actual signed jwt string along with the ID
// We return the id, so it can be used without re-parsing the JWT from signed string
type RefreshToken struct {
	SignedTokenString string
	ID                string
	ExpiresIn   	  time.Duration
}

// RefreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}