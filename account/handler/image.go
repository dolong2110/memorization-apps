package handler

import (
	"fmt"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/dolong2110/Memoirization-Apps/account/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Image handler
func (h *Handler) Image(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	// limit overly large request bodies
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	imageFileHeader, err := c.FormFile("image_file")

	// check for error before checking for non-nil header
	if err != nil {
		// should be a validation error
		log.Printf("Unable parse mage from multipart/form-data: %+v", err)

		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
			})
			return
		}
		e := apperrors.NewBadRequest("Unable to parse image from multipart/form-data")
		c.JSON(e.Status(), apperrors.Response{Error: e})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")

	// Validate image mime-type is allowable
	if valid := utils.IsAllowedImageType(mimeType); !valid {
		log.Println("Image is not an allowable mime-type")
		e := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx := c.Request.Context()

	updatedUser, err := h.UserService.SetProfileImage(ctx, authUser.UID, imageFileHeader)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image_url": updatedUser.ImageURL,
		"message":   "success",
	})
}
