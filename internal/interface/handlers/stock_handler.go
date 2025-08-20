package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/application"
)

type StockHandler struct {
	service *application.StockService
}

func NewStockHandler(service *application.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) UpdateStocks(c *gin.Context) {
	err := h.service.UpdateStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stocks updated successfully"})
}
