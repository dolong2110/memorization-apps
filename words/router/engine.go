package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Router struct {
	config      *Config
	dataSources *dataSources
}

type dataSources struct {
	DB *sqlx.DB
}

func NewRouter(config *Config) *Router {
	return &Router{
		config: config,
	}
}

func (r *Router) InitGin() (*gin.Engine, error) {

	router := gin.Default()

	return router, nil
}

func InitDS(r *Router) error {
	log.Printf("Initializing data sources\n")

	pg := r.config.DataSource.PostGreSQL
	log.Printf("Connecting to Postgresql\n")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.PostGreSHost, pg.PostGreSPort, pg.PostGreSUser, pg.PostGreSPassword, pg.PostGreSDB, pg.PostGreSSSL)
	db, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}

	r.dataSources.DB = db

	return nil
}

func (r *Router) Close() error {
	if err := r.dataSources.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
