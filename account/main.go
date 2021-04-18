package main

import (
	"context"
	"log"

)

func main() {
	// you could insert your favorite logger here for structured or leveled logging
	log.Println("Starting server...")
	
	router := gin.Default()

	router.GET("/api/account", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	srv := &http.Server{}
}