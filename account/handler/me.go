package handler

import (
	"github.com/dolong2110/memorization-apps/account/model"
	"github.com/dolong2110/memorization-apps/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Me handler calls services for getting
// current login user's details
func (h *Handler) Me(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	user, exists := c.Get("user")

	// This shouldn't happen, as our middleware ought to throw an error.
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*model.User).UID

	// use the Request Context
	ctx := c.Request.Context()
	user, err := h.UserService.Get(ctx, uid)
	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
