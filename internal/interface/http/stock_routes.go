package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/handlers"
)

func RegisterStockRoutes(r *gin.RouterGroup, h *handlers.StockHandler) {
	{
		r.GET("/stocks", h.GetStocks)
		r.GET("/stocks/recommend", h.GetRecommend)
	}
}
