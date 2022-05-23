package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, user *User) error
	Signin(ctx context.Context, user *User) error
	UpdateDetails(ctx context.Context, user *User) error
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWTs as string
type TokenService interface {
	NewPairFromUser(ctx context.Context, user *User, prevRefreshTokenID string) (*Token, error)
	Signout(ctx context.Context, uid uuid.UUID) error
	ValidateIDToken(idTokenString string) (*User, error)                   // jwt not require context, and we not do anything in repository or db that cancel or modify context
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error) // not need context because not reach DB or other layer.
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

// TokenRepository defines methods it expects a repository
// it interacts with to implement
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshToken(ctx context.Context, userID string) error
}

// ImageRepository defines methods it expects a repository
// it interacts with to implement
type ImageRepository interface {
	//DeleteProfile(ctx context.Context, objName string) error
	//UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}
