package main

import (
	"fmt"
	"notes-api/models"
	"notes-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"

)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	models.ConnectDatabase()
	r := gin.Default()
	// Set up CORS before running the server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	routes.RegisterNoteRoutes(r)
	r.Run(":3000")

}
