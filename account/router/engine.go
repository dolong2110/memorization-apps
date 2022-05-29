package router

import (
	"github.com/dolong2110/memorization-apps/account/handler"
	"github.com/dolong2110/memorization-apps/account/repository"
	"github.com/dolong2110/memorization-apps/account/service"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Router is the struct used to init gin-gonic
type Router struct {
	config     *Config
	dataSource *DataSources
}

// NewRouters is method to init router struct
func NewRouters(config *Config, dataSources *DataSources) *Router {
	return &Router{
		config:     config,
		dataSource: dataSources,
	}
}

// InitGin is a method which receive attributes of router and return *gin.Engine and error
func (r *Router) InitGin() (*gin.Engine, error) {
	// initialize data sources

	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(r.dataSource.PostgreSQLDB)
	tokenRepository := repository.NewTokenRepository(r.dataSource.RedisClient)
	imageRepository := repository.NewImageRepository(r.dataSource.CloudStorageClient, r.config.DataSource.GCP.GCPImageBucket)

	/*
	 * service layer
	 */
	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		ImageRepository: imageRepository,
	})

	tokenConfig := r.config.Token
	accessTokenInfo, err := initAccessToken(tokenConfig.AccessToken)
	if err != nil {
		log.Fatalf("could not get access token information: %v\n", err)
	}
	refreshTokenInfo := initRefreshToken(tokenConfig.RefreshToken)

	tokenService := service.NewTokenService(&service.TokenServiceConfig{
		AccessTokenInfo:  *accessTokenInfo,
		RefreshTokenInfo: *refreshTokenInfo,
		TokenRepository:  tokenRepository,
	})

	// initialize gin.Engine
	router := gin.Default()

	handler.NewHandler(&handler.Config{
		Engine:          router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         r.config.ApiUrl,
		TimeoutDuration: time.Duration(r.config.HandlerTimeout) * time.Second,
		MaxBodyBytes:    r.config.MaxBodyBytes,
	})

	return router, nil
}
