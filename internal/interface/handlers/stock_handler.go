package handlers

import (
	"net/http"
	"strconv"

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

func (h *StockHandler) GetStocks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	stocks, total, err := h.service.GetStocks(page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":     "Error fetching stocks",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        stocks,
		"total":       total,
		"page":        page,
		"pageSize":    pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})

}
