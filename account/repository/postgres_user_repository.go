package repository

import (
	"context"
	"github.com/dolong2110/Memoirization-Apps/account/model"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type PGUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &PGUserRepository{
		DB: db,
	}
}

// Create reaches out to database SQLX api
func (r *PGUserRepository) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.Get(user, query, user.Email, user.Password); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, err.Code.Name())
			return apperrors.NewConflict("email", user.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByID fetches user by id
func (r *PGUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Get(user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}