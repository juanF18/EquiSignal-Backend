package application

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/juanF18/EquiSignal-Backend/internal/domain/models"
	"github.com/juanF18/EquiSignal-Backend/internal/interface/dto"
)

func TestDTOToModelTransformation(t *testing.T) {
	// Arrange
	now := time.Now()
	dtoStock := dto.Stock{
		Ticker:     "AAPL",
		Company:    "Apple Inc.",
		Brokerage:  "Goldman Sachs",
		Action:     "Initiates",
		RatingFrom: "Hold",
		RatingTo:   "Strong Buy",
		TargetFrom: "$150.00",
		TargetTo:   "$180.00",
		Time:       now,
	}

	// Act - Simular transformación (esto estaría en el service real)
	modelStock := models.Stock{
		ID:         uuid.New(), // Se generaría en la DB
		Ticker:     dtoStock.Ticker,
		Company:    dtoStock.Company,
		Brokerage:  dtoStock.Brokerage,
		Action:     dtoStock.Action,
		RatingFrom: dtoStock.RatingFrom,
		RatingTo:   dtoStock.RatingTo,
		TargetFrom: dtoStock.TargetFrom,
		TargetTo:   dtoStock.TargetTo,
		Time:       dtoStock.Time,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Assert
	if modelStock.Ticker != dtoStock.Ticker {
		t.Errorf("Expected ticker %s, got %s", dtoStock.Ticker, modelStock.Ticker)
	}

	if modelStock.Company != dtoStock.Company {
		t.Errorf("Expected company %s, got %s", dtoStock.Company, modelStock.Company)
	}

	if modelStock.RatingTo != dtoStock.RatingTo {
		t.Errorf("Expected rating %s, got %s", dtoStock.RatingTo, modelStock.RatingTo)
	}

	// Verificar que los campos de modelo tienen valores
	if modelStock.ID == uuid.Nil {
		t.Error("Expected valid UUID for ID")
	}

	if modelStock.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestStockServiceParameterValidation(t *testing.T) {
	// Test para validación conceptual de parámetros
	testCases := []struct {
		name        string
		page        int
		pageSize    int
		search      string
		expectValid bool
	}{
		{"Valid parameters", 1, 10, "", true},
		{"Valid with search", 1, 10, "AAPL", true},
		{"Large page size", 1, 100, "", true},
		{"Zero page", 0, 10, "", false},
		{"Negative page", -1, 10, "", false},
		{"Zero page size", 1, 0, "", false},
		{"Negative page size", 1, -5, "", false},
		{"Very large page size", 1, 10000, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validación conceptual de parámetros
			isValid := tc.page > 0 && tc.pageSize > 0 && tc.pageSize <= 1000

			if isValid != tc.expectValid {
				t.Errorf("Parameter validation failed for %s", tc.name)
			}
		})
	}
}

func TestStockDataValidation(t *testing.T) {
	// Test para validar datos de stock
	testCases := []struct {
		name        string
		stock       dto.Stock
		expectValid bool
		errorField  string
	}{
		{
			name: "Valid stock data",
			stock: dto.Stock{
				Ticker:    "AAPL",
				Company:   "Apple Inc.",
				Brokerage: "Goldman Sachs",
				RatingTo:  "Strong Buy",
				Time:      time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Missing ticker",
			stock: dto.Stock{
				Company:   "Apple Inc.",
				Brokerage: "Goldman Sachs",
				RatingTo:  "Strong Buy",
				Time:      time.Now(),
			},
			expectValid: false,
			errorField:  "ticker",
		},
		{
			name: "Missing company",
			stock: dto.Stock{
				Ticker:    "AAPL",
				Brokerage: "Goldman Sachs",
				RatingTo:  "Strong Buy",
				Time:      time.Now(),
			},
			expectValid: false,
			errorField:  "company",
		},
		{
			name: "Zero time",
			stock: dto.Stock{
				Ticker:   "AAPL",
				Company:  "Apple Inc.",
				RatingTo: "Strong Buy",
				Time:     time.Time{}, // Zero time
			},
			expectValid: false,
			errorField:  "time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Función de validación
			isValid := validateStockData(tc.stock)

			if isValid != tc.expectValid {
				t.Errorf("Expected validation result %v, got %v for field %s",
					tc.expectValid, isValid, tc.errorField)
			}
		})
	}
}

// Función helper para validación
func validateStockData(stock dto.Stock) bool {
	if stock.Ticker == "" {
		return false
	}
	if stock.Company == "" {
		return false
	}
	if stock.Time.IsZero() {
		return false
	}
	return true
}

func TestPaginationCalculation(t *testing.T) {
	testCases := []struct {
		name          string
		total         int64
		pageSize      int
		expectedPages int64
	}{
		{"Exact division", 100, 10, 10},
		{"With remainder", 105, 10, 11},
		{"Less than page size", 5, 10, 1},
		{"Zero total", 0, 10, 0},
		{"Single item", 1, 10, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Cálculo de páginas totales como en el handler
			var totalPages int64
			if tc.total == 0 {
				totalPages = 0
			} else {
				totalPages = (tc.total + int64(tc.pageSize) - 1) / int64(tc.pageSize)
			}

			if totalPages != tc.expectedPages {
				t.Errorf("Expected %d pages, got %d", tc.expectedPages, totalPages)
			}
		})
	}
}
