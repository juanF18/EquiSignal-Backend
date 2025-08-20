package repository

import (
	"github.com/juanF18/EquiSignal-Backend/internal/domain/models"
	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) SaveStocks(stocks []models.Stock) error {
	return r.db.Create(&stocks).Error
}
