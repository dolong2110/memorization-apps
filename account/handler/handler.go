package handler

import (
	"github.com/dolong2110/Memoirization-Apps/account/handler/middleware"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"time"
)

// Handler struct holds required services for handler to function
type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
	MaxBodyBytes int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Engine          *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	// Create an account group
	// Create a handler (which will later have injected services)
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	// Create a group, or base url for all routes
	g := c.Engine.Group(c.BaseURL) // Create a handler (which will later have injected services)
	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.Details)
		g.POST("/image", middleware.AuthUser(h.TokenService), h.Image)
		g.DELETE("/image", middleware.AuthUser(h.TokenService), h.DeleteImage)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
		g.PUT("/details", h.Details)
		g.POST("/image", h.Image)
		g.DELETE("/image", h.DeleteImage)
	}

	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/tokens", h.Tokens)
}
