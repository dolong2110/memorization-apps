package service

import (
	"context"
	"crypto/rsa"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/dolong2110/Memoirization-Apps/account/utils"
	"log"
)

// tokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TokenServiceConfig will hold repositories that will eventually be injected into
// this service layer
type TokenServiceConfig struct {
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivateKey:            c.PrivateKey,
		PublicKey:             c.PublicKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, user *model.User, prevTokenID string) (*model.Token, error) {
	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := utils.GenerateIDToken(user, s.PrivateKey, s.IDExpirationSecs)
	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UID, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, user.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, user.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", user.UID.String(), prevTokenID)
		}
	}

	return &model.Token{
		IDToken:      idToken,
		RefreshToken: refreshToken.SignedTokenString,
	}, nil
}
