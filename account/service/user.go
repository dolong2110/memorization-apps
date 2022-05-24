package service

import (
	"context"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/dolong2110/Memoirization-Apps/account/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// USConfig will hold repositories that will eventually be injected into
// this service layer
type USConfig struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		ImageRepository: c.ImageRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user, err := s.UserRepository.FindByID(ctx, uid)

	return user, err
}

// Signup reaches out to a UserRepository to sign up the user.
// UserRepository Create should handle checking for user exists conflicts
func (s *userService) Signup(ctx context.Context, user *model.User) error {
	pwd, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return apperrors.NewInternal()
	}

	// now I realize why I originally used Signup(ctx, email, password)
	// then created a user. It's somewhat un-natural to mutate the user here
	user.Password = pwd
	if err := s.UserRepository.Create(ctx, user); err != nil {
		return err
	}

	// ...

	return nil
}

// Signin reaches our to a UserRepository check if the user exists
// and then compares the supplied password with the provided password.
// If a valid email/password combo is provided, u will hold all
// available user fields
func (s *userService) Signin(ctx context.Context, user *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := utils.ComparePasswords(uFetched.Password, user.Password)
	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*user = *uFetched
	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, user *model.User) error {
	// Update user in UserRepository
	err := s.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	// // Publish user updated
	// err = s.EventsBroker.PublishUserUpdated(user, false)
	// if err != nil {
	// 	return apperrors.NewInternal()
	// }

	return nil
}

func (s *userService) SetProfileImage(
	ctx context.Context,
	uid uuid.UUID,
	imageFileHeader *multipart.FileHeader,
) (*model.User, error) {
	user, err := s.UserRepository.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	objName, err := utils.ObjNameFromURL(user.ImageURL)
	if err != nil {
		return nil, err
	}

	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("Failed to open image file: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	// Upload user's image to ImageRepository
	// Possibly received updated imageURL
	imageURL, err := s.ImageRepository.UpdateProfile(ctx, objName, imageFile)
	if err != nil {
		log.Printf("Unable to upload image to cloud provider: %v\n", err)
		return nil, err
	}

	updatedUser, err := s.UserRepository.UpdateImage(ctx, user.UID, imageURL)
	if err != nil {
		log.Printf("Unable to update imageURL: %v\n", err)
		return nil, err
	}

	return updatedUser, nil
}
