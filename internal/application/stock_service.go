package application

import (
	"log"

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
