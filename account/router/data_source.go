package router

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type DataSources struct {
	PostgreSQLDB       *sqlx.DB
	RedisClient        *redis.Client
	CloudStorageClient *storage.Client
}

// InitDS establishes connections to fields in DataSources
func InitDS(config *Config) (*DataSources, error) {
	log.Printf("Initializing data sources\n")
	pg := config.DataSource.PostGreSQL
	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pg.PostGresHost, pg.PostGresPort, pg.PostGresUser, pg.PostGresPassword, pg.PostGresDB, pg.PostGresSSL)

	log.Printf("Connecting to Postgresql\n")
	db, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	db.SetConnMaxLifetime(time.Duration(pg.PostGresConnectionTimeOut) * time.Minute)

	rd := config.DataSource.Redis
	log.Printf("Connecting to Redis\n")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rd.RedisHost, rd.RedisPort),
		Password: "",
		DB:       0,
	})

	// verify redis connection
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	// Initialize google storage client
	log.Printf("Connecting to Cloud Storage\n")
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(config.DataSource.GCP.CloudConnectionTimeout)*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	cloudStorage, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating cloud storage client: %w", err)
	}

	return &DataSources{
		PostgreSQLDB:       db,
		RedisClient:        rdb,
		CloudStorageClient: cloudStorage,
	}, nil
}

// Close to be used in graceful server shutdown
func (d *DataSources) Close() error {
	if err := d.PostgreSQLDB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}

	if err := d.CloudStorageClient.Close(); err != nil {
		return fmt.Errorf("error closing Cloud Storage client: %w", err)
	}

	return nil
}
