package service

import (
	"context"
	"github.com/dolong2110/memorization-apps/account/model"
	"github.com/dolong2110/memorization-apps/account/model/apperrors"
	"github.com/dolong2110/memorization-apps/account/utils"

	"github.com/google/uuid"
	"log"
)

// tokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	AccessToken     model.AccessTokenInfo
	RefreshToken    model.RefreshTokenInfo
	TokenRepository model.TokenRepository
}

// TokenServiceConfig will hold repositories that will eventually be injected into
// this service layer
type TokenServiceConfig struct {
	AccessTokenInfo  model.AccessTokenInfo
	RefreshTokenInfo model.RefreshTokenInfo
	TokenRepository  model.TokenRepository
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &tokenService{
		AccessToken:     c.AccessTokenInfo,
		RefreshToken:    c.RefreshTokenInfo,
		TokenRepository: c.TokenRepository,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenService) NewPairFromUser(ctx context.Context, user *model.User, prevRefreshTokenID string) (*model.Token, error) {
	// delete user's current refresh token (used when refreshing idToken)
	if prevRefreshTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, user.UID.String(), prevRefreshTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", user.UID.String(), prevRefreshTokenID)
			return nil, err
		}
	}

	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := utils.GenerateIDToken(user, s.AccessToken.PrivateKey, s.AccessToken.Expires)
	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UID, s.RefreshToken.Secret, s.RefreshToken.Expires)
	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, user.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", user.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	return &model.Token{
		AccessToken:  model.AccessToken{SignedStringToken: idToken},
		RefreshToken: model.RefreshToken{SignedStringToken: refreshToken.SignedStringToken, ID: refreshToken.ID, UID: user.UID},
	}, nil
}

// Signout reaches out to the repository layer to delete all valid tokens for a user
func (s *tokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	return s.TokenRepository.DeleteUserRefreshToken(ctx, uid.String())
}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the AccessTokenCustomClaims
func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := utils.ValidateIDToken(tokenString, s.AccessToken.PublicKey) // uses public RSA key
	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := utils.ValidateRefreshToken(tokenString, s.RefreshToken.Secret)
	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.Id)
	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SignedStringToken: tokenString,
		ID:                tokenUUID,
		UID:               claims.UID,
	}, nil
}
