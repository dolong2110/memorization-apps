package model

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// Token used for returning pairs of id and refresh tokens
type Token struct {
	AccessToken
	RefreshToken
}

// AccessToken stores token properties that
// are accessed in multiple application layers
type AccessToken struct {
	SignedStringToken string `json:"id_token"`
}

// RefreshToken stores token properties that
// are accessed in multiple application layers
type RefreshToken struct {
	ID                uuid.UUID `json:"-"`
	UID               uuid.UUID `json:"-"`
	SignedStringToken string    `json:"refresh_token"`
}

// AccessTokenInfo stores access token's initialize information
type AccessTokenInfo struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	Expires    int64
}

// RefreshTokenInfo stores refresh token's initialize information
type RefreshTokenInfo struct {
	Secret  string
	Expires int64
}

// AccessTokenCustomClaims holds structure of jwt claims of idToken
type AccessTokenCustomClaims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}

// RefreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// RefreshTokenData holds the actual signed jwt string along with the ID
// We return the id, so it can be used without re-parsing the JWT from signed string
type RefreshTokenData struct {
	SignedStringToken string
	ID                uuid.UUID
	ExpiresIn         time.Duration
}
