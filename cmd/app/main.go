package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/juanF18/EquiSignal-Backend/internal/config"
	"github.com/juanF18/EquiSignal-Backend/internal/infrastructure/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se encontró .env, usando variables del sistema")
	}

	cfg := config.LoadConfig()

	db.ConnectCockroachDB(cfg)

	r := gin.Default()

	var now time.Time
	db.DB.Raw("SELECT NOW()").Scan(&now)
	fmt.Println("⏰ DB Time:", now)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	addr := fmt.Sprintf(":%s", cfg.HttpPort)
	log.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
