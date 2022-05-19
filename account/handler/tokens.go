package handler

import (
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type tokensReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Tokens handler
func (h *Handler) Tokens(c *gin.Context) {
	// bind JSON to req of type tokensRew
	var req tokensReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	// verify refresh JWT
	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// get up-to-date user
	user, err := h.UserService.Get(ctx, refreshToken.UID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create fresh pair of tokens
	tokens, err := h.TokenService.NewPairFromUser(ctx, user, refreshToken.ID.String())
	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", user, err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
