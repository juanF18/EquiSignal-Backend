package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanF18/EquiSignal-Backend/internal/infrastructure/external"
)

type StockHandler struct {
	api *external.ExternalAPI
}

func NewStockHandler(api *external.ExternalAPI) *StockHandler {
	return &StockHandler{api: api}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	nextPage := c.Query("next_page")

	data, err := h.api.FetchStocks(nextPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
