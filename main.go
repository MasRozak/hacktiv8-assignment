package main

import (
	"log"

	"social-media-api/config"
	"social-media-api/database"
	"social-media-api/middleware"
	"social-media-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg.GetDatabaseURL())
	defer database.CloseDB()

	database.CreateTables()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Timestamping())

	routes.SetupRoutes(r)

	log.Printf("Server running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
