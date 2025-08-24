package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/juanF18/EquiSignal-Backend/internal/application"
	"github.com/juanF18/EquiSignal-Backend/internal/config"
	"github.com/juanF18/EquiSignal-Backend/internal/infrastructure/db"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/external"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/handlers"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se encontró .env, usando variables del sistema")
	}

	cfg := config.LoadConfig()

	db.ConnectCockroachDB(cfg)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontEndURL}, // cambia según tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	var now time.Time
	db.DB.Raw("SELECT NOW()").Scan(&now)
	fmt.Println("⏰ DB Time:", now)

	// API externa
	externalAPI := external.NewExternalAPI(cfg)
	stockService := application.NewStockService(externalAPI)
	stockHandler := handlers.NewStockHandler(stockService)

	http.SetupRoutes(r, stockHandler)

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
