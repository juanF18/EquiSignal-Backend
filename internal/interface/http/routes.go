package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/handlers"
)

func SetupRoutes(r *gin.Engine, stockHandler *handlers.StockHandler) {
	api := r.Group("/api")
	{
		RegisterExternalAPIRoutes(api, stockHandler)
	}
}
