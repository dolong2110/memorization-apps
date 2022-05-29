package router

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolong2110/memorization-apps/account/model"
	"io/ioutil"
)

func initAccessToken(accessTokenConfig AccessToken) (*model.AccessTokenInfo, error) {
	// load rsa keys
	publicKeyByte, err := ioutil.ReadFile(accessTokenConfig.PublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	privateKeyByte, err := ioutil.ReadFile(accessTokenConfig.PrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	return &model.AccessTokenInfo{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Expires:    accessTokenConfig.AccessTokenExpire,
	}, nil
}

func initRefreshToken(refreshTokenConfig RefreshToken) *model.RefreshTokenInfo {
	return &model.RefreshTokenInfo{
		Secret:  refreshTokenConfig.RefreshTokenSecret,
		Expires: refreshTokenConfig.RefreshTokenExpire,
	}
}
