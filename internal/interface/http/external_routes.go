package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/handlers"
)

func RegisterExternalAPIRoutes(r *gin.RouterGroup, h *handlers.StockHandler) {
	externalGroup := r.Group("/external")
	{
		externalGroup.GET("/update-stocks", h.UpdateStocks)
	}
}
