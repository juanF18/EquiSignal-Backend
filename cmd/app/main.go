package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	r := gin.Default()

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
