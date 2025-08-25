package application

import (
	"log"
	"strings"

	"github.com/juanF18/EquiSignal-Backend/internal/algorithms/stock"
	"github.com/juanF18/EquiSignal-Backend/internal/domain/models"
	"github.com/juanF18/EquiSignal-Backend/internal/infrastructure/db"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/external"
)

type StockService struct {
	api *external.ExternalAPI
}

func NewStockService(api *external.ExternalAPI) *StockService {
	return &StockService{api: api}
}

// Trae todas las páginas y guarda en Cockroach
func (s *StockService) UpdateStocks() error {
	nextPage := ""
	for {
		resp, err := s.api.FetchStocks(nextPage)
		if err != nil {
			return err
		}

		// Guardar items en DB
		for _, item := range resp.Items {
			stock := models.Stock{
				Ticker:     item.Ticker,
				Company:    item.Company,
				Brokerage:  item.Brokerage,
				Action:     item.Action,
				RatingFrom: item.RatingFrom,
				RatingTo:   item.RatingTo,
				TargetFrom: item.TargetFrom,
				TargetTo:   item.TargetTo,
				Time:       item.Time,
			}
			if err := db.DB.Create(&stock).Error; err != nil {
				log.Printf("⚠️ Error guardando %s: %v", stock.Ticker, err)
			}
		}

		if resp.NextPage == "" {
			break // ya no hay más páginas
		}
		nextPage = resp.NextPage
	}

	return nil
}

// GetStocks devuelve una lista de stocks con paginación
func (s *StockService) GetStocks(page, pageSize int, search string) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	query := db.DB.Model(&models.Stock{})

	// Si el usuario pasó un search, filtramos
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			db.DB.Where("LOWER(Ticker) LIKE ?", like).
				Or("LOWER(Company) LIKE ?", like).
				Or("LOWER(Brokerage) LIKE ?", like),
		)
	}

	// contar el total de registros filtrados
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// aplicar paginación + ordenamiento
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (s *StockService) GetRecommend(limit int) ([]stock.StockRecommendation, error) {
	var stocks []models.Stock
	if err := db.DB.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stock.RecommendStocks(stocks, limit), nil
}
