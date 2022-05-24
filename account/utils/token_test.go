package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestValidateIDToken(t *testing.T) {
	var idExp int64 = 15 * 60
	uid, _ := uuid.NewRandom()
	user := &model.User{
		UID:   uid,
		Email: "bob@bob.com",
	}

	privateKeyFromPem, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFromPem)
	if err != nil {
		privateKey, _ = GeneratePrivateKey(2048)
	}

	publicKeyFromPem, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFromPem)
	if err != nil {
		publicKey = &privateKey.PublicKey
	}

	inValidPrivateKeyFromPEM, _ := GeneratePrivateKey(2048)

	t.Run("Valid token", func(t *testing.T) {
		ss, _ := GenerateIDToken(user, privateKey, idExp)

		uFromToken, err := ValidateIDToken(ss, publicKey)
		assert.NoError(t, err)
		assert.ElementsMatch(
			t,
			[]interface{}{user.Email, user.Name, user.UID, user.Website, user.ImageURL},
			[]interface{}{uFromToken.User.Email, uFromToken.User.Name, uFromToken.User.UID, uFromToken.User.Website, uFromToken.User.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		ss, _ := GenerateIDToken(user, privateKey, -1)

		expectedErr := apperrors.NewAuthorization("token is expired by 1s")

		_, err := ValidateIDToken(ss, publicKey)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		ss, _ := GenerateIDToken(user, inValidPrivateKeyFromPEM, 999999999) // expires one second ago

		expectedErr := apperrors.NewAuthorization("crypto/rsa: verification error")

		_, err := ValidateIDToken(ss, publicKey)
		assert.EqualError(t, err, expectedErr.Message)
	})
}

func TestValidateRefreshToken(t *testing.T) {
	var refreshExp int64 = 3 * 24 * 2600
	secret := "anotsorandomtestsecret"

	uid, _ := uuid.NewRandom()
	user := &model.User{
		UID:   uid,
		Email: "bob@bob.com",
	}

	//inValidPrivateKeyFromPEM, _ := GeneratePrivateKey(2048)

	t.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := GenerateRefreshToken(user.UID, secret, refreshExp)

		validatedRefreshToken, err := ValidateRefreshToken(testRefreshToken.SignedStringToken, secret)
		assert.NoError(t, err)
		assert.Equal(t, user.UID, validatedRefreshToken.UID)
	})

	t.Run("invalid signed token", func(t *testing.T) {
		testRefreshToken, _ := GenerateRefreshToken(user.UID, "secret", refreshExp)

		expectedErr := apperrors.NewAuthorization("signature is invalid")

		_, err := ValidateRefreshToken(testRefreshToken.SignedStringToken, secret)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := GenerateRefreshToken(user.UID, secret, -1)

		expectedErr := apperrors.NewAuthorization("token is expired by 1s")

		_, err := ValidateRefreshToken(testRefreshToken.SignedStringToken, secret)
		assert.EqualError(t, err, expectedErr.Message)
	})
}
