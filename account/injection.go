package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolong2110/Memoirization-Apps/account/handler"
	"github.com/dolong2110/Memoirization-Apps/account/repository"
	"github.com/dolong2110/Memoirization-Apps/account/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)

	bucketName := os.Getenv("GCP_IMAGE_BUCKET")
	imageRepository := repository.NewImageRepository(d.StorageClient, bucketName)

	/*
	 * service layer
	 */
	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		ImageRepository: imageRepository,
	})

	// load rsa keys
	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	privateKeyByte, err := ioutil.ReadFile(privateKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUBLIC_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenService := service.NewTokenService(&service.TokenServiceConfig{
		TokenRepository:       tokenRepository,
		PrivateKey:            privateKey,
		PublicKey:             pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// initialize gin.Engine
	router := gin.Default()

	// read in ACCOUNT_API_URL
	baseURL := os.Getenv("ACCOUNT_API_URL")
	// read in HANDLER_TIMEOUT
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	timeoutInt, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		Engine:          router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(timeoutInt) * time.Second,
	})

	return router, nil
}
