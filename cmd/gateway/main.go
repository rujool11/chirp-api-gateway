package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rujool11/chirp-api-gateway/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	AUTH_URL := os.Getenv("AUTH_URL")
	CORE_URL := os.Getenv("CORE_URL")
	if AUTH_URL == "" {
		log.Println("Defaulting AUTH_URL to http://localhost:8001")
		AUTH_URL = "http://localhost:8001"
	}
	if CORE_URL == "" {
		log.Println("Defaulting CORE_URL to http://localhost:8002")
		CORE_URL = "http://localhost:8002"
	}

	r := gin.Default()

	// redirect /auth to auth service
	r.Any("/auth/*path", utils.ReverseProxy(AUTH_URL, "/auth"))

	// redirect /core to core service
	r.Any("/core/*path", utils.ReverseProxy(CORE_URL, "/core"))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from chirp-api-gateway",
		})
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
		log.Println("Defaulting PORT to 8000")
	}

	r.Run(":" + PORT)
}
