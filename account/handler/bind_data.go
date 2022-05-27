package handler

import (
	"fmt"
	"github.com/dolong2110/memorization-apps/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// bindData is helper function, returns false if data is not bound
func bindData(ctx *gin.Context, req interface{}) bool {
	// send error if Content-Type != application/json
	if ctx.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", ctx.FullPath())

		err := apperrors.NewUnsupportedMediaType(msg)

		ctx.JSON(err.Status(), apperrors.Response{Error: err})
		return false
	}

	// Bind incoming json to struct and check for validation errors
	if err := ctx.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

			ctx.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		// if we aren't able to properly extract validation errors,
		// we'll fall back and return an internal server error
		fallBack := apperrors.NewInternal()

		ctx.JSON(fallBack.Status(), apperrors.Response{Error: fallBack})
		return false
	}

	return true
}
