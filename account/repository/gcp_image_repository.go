package repository

import (
	"cloud.google.com/go/storage"
	"github.com/dolong2110/Memoirization-Apps/account/model"
)

type gcpImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

// NewImageRepository is a factory for initializing User Repositories
func NewImageRepository(gcClient *storage.Client, bucketName string) model.ImageRepository {
	return &gcpImageRepository{
		Storage:    gcClient,
		BucketName: bucketName,
	}
}
