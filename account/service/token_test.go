package service

import (
	"context"
	"fmt"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/dolong2110/Memoirization-Apps/account/model/mocks"
	"github.com/dolong2110/Memoirization-Apps/account/utils"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func TestNewPairFromUser(t *testing.T) {
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600
	privateKeyFromPem, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFromPem)
	if err != nil {
		privateKey, _ = utils.GeneratePrivateKey(2048)
	}
	publicKeyFromPem, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFromPem)
	if err != nil {
		publicKey = &privateKey.PublicKey
	}
	secret := "anotsorandomtestsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository:       mockTokenRepository,
		PrivateKey:            privateKey,
		PublicKey:             publicKey,
		RefreshSecret:         secret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	user := &model.User{
		UID:      uid,
		Email:    "long@do.com",
		Password: "blarghedymcblarghface",
	}

	// Setup mock call responses in setup before t.Run statements
	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &model.User{
		UID:      uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uErrorCase.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UID.String(),
		prevID,
	}

	// mock call argument/responses
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()
		tokenPair, err := tokenService.NewPairFromUser(ctx, user, prevID)
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SignedStringToken)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &model.IDTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken.SignedStringToken, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			user.UID,
			user.Email,
			user.Name,
			user.ImageURL,
			user.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &model.RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SignedStringToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SignedStringToken)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, user.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})
	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, uErrorCase, "")
		assert.Error(t, err) // should return an error

		// SetRefreshToken should be called with setErrorArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, user, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
	t.Run("Prev token not in repository", func(t *testing.T) {
		ctx := context.Background()
		uid, _ := uuid.NewRandom()
		user := &model.User{
			UID: uid,
		}

		tokenIDNotInRepo := "not_in_token_repo"

		deleteArgs := mock.Arguments{
			ctx,
			user.UID.String(),
			tokenIDNotInRepo,
		}

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		mockTokenRepository.
			On("DeleteRefreshToken", deleteArgs...).
			Return(mockError)

		_, err := tokenService.NewPairFromUser(ctx, user, tokenIDNotInRepo)
		assert.Error(t, err)

		appError, ok := err.(*apperrors.Error)

		assert.True(t, ok)
		assert.Equal(t, apperrors.Authorization, appError.Type)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteArgs...)
		mockTokenRepository.AssertNotCalled(t, "SetRefreshToken")
	})
}

func TestSignout(t *testing.T) {
	mockTokenRepository := new(mocks.MockTokenRepository)
	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository: mockTokenRepository,
	})

	t.Run("No error", func(t *testing.T) {
		uidSuccess, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshToken", mock.AnythingOfType("*context.emptyCtx"), uidSuccess.String()).
			Return(nil)

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidSuccess)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		uidError, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshToken", mock.AnythingOfType("*context.emptyCtx"), uidError.String()).
			Return(apperrors.NewInternal())

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidError)

		assert.Error(t, err)

		apperr, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperr.Type, apperrors.Internal)
	})
}

func TestValidateIDToken(t *testing.T) {
	var idExp int64 = 15 * 60

	privateKeyFromPem, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFromPem)
	if err != nil {
		privateKey, _ = utils.GeneratePrivateKey(2048)
	}
	publicKeyFromPem, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFromPem)
	if err != nil {
		publicKey = &privateKey.PublicKey
	}

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TokenServiceConfig{
		PrivateKey:       privateKey,
		PublicKey:        publicKey,
		IDExpirationSecs: idExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := utils.GenerateIDToken(u, privateKey, idExp)

		uFromToken, err := tokenService.ValidateIDToken(ss)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{u.Email, u.Name, u.UID, u.Website, u.ImageURL},
			[]interface{}{uFromToken.Email, uFromToken.Name, uFromToken.UID, uFromToken.Website, uFromToken.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := utils.GenerateIDToken(u, privateKey, -1) // expires one second ago

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := utils.GenerateIDToken(u, privateKey, -1) // expires one second ago

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	// TODO - Add other invalid token types
}

func TestValidateRefreshToken(t *testing.T) {
	var refreshExp int64 = 3 * 24 * 2600
	secret := "anotsorandomtestsecret"

	tokenService := NewTokenService(&TokenServiceConfig{
		RefreshSecret:         secret,
		RefreshExpirationSecs: refreshExp,
	})

	uid, _ := uuid.NewRandom()
	user := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := utils.GenerateRefreshToken(user.UID, secret, refreshExp)

		validatedRefreshToken, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedStringToken)
		assert.NoError(t, err)

		assert.Equal(t, user.UID, validatedRefreshToken.UID)
		assert.Equal(t, testRefreshToken.SignedStringToken, validatedRefreshToken.SignedStringToken)
		assert.Equal(t, user.UID, validatedRefreshToken.UID)
	})

	t.Run("invalid signed token", func(t *testing.T) {
		testRefreshToken, _ := utils.GenerateRefreshToken(user.UID, "secret", refreshExp)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedStringToken)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := utils.GenerateRefreshToken(user.UID, secret, -1)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedStringToken)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Error Token ID: Not uuid type in uid field of token", func(t *testing.T) {
		testRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxIiwiZXhwIjoxNjUzNTkyMjk2LCJqdGkiOiI1ZGQzNjg3Ny05MTZlLTQ2MTUtOThjNC0zYTllNzVjYjAwNTgiLCJpYXQiOjk5OTk5OTk5OTk5OX0.YXJLGwK8kYqGzSgfIZOog-kTuW4fLgRTlFN-lhEEX0g"

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Error Token ID: int type in uid field of token", func(t *testing.T) {
		testRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsImV4cCI6MTY1MzU5MjI5NiwianRpIjoiNWRkMzY4NzctOTE2ZS00NjE1LTk4YzQtM2E5ZTc1Y2IwMDU4IiwiaWF0Ijo5OTk5OTk5OTk5OTl9.zqaapBDRioxF5V5FzN8cXRWNxYRKllXQ91pjsRMGzA0"

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Error User ID: not uid type in jti field of token", func(t *testing.T) {
		testRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3ZDIwYTMzZi1mZTE1LTRmZmQtOTBlZS1kNDRkMDEzYzI2MGUiLCJleHAiOjE2NTM1OTIyOTYsImp0aSI6IjEiLCJpYXQiOjk5OTk5OTk5OTk5OX0.AUzh2t9RHaKrJl9cEqoSxhO5nFvAtVzISd7c-6AFowk"

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Error Token ID: int type in jti field of token", func(t *testing.T) {
		testRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3ZDIwYTMzZi1mZTE1LTRmZmQtOTBlZS1kNDRkMDEzYzI2MGUiLCJleHAiOjE2NTM1OTIyOTYsImp0aSI6MSwiaWF0Ijo5OTk5OTk5OTk5OTl9.ts6ZmNnTyaAKXsetR53-bV42q51Z3PoL0ozzshUT-Vw"

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken)
		assert.EqualError(t, err, expectedErr.Message)
	})
}
