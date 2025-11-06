package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8081",
			"http://localhost:8082",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:8082",
			"http://192.168.1.82:8081",
			"http://192.168.1.82:8082",
			// tus puertos de Expo dev si aplica (19000–19006)
			"http://localhost:19000", "http://localhost:19001",
			"http://localhost:19002", "http://localhost:19003",
			"http://localhost:19004", "http://localhost:19005", "http://localhost:19006",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // ← si no usás cookies, mejor false
		MaxAge:           12 * time.Hour,
	})
}
