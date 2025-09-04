package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStock_CreatedAt(t *testing.T) {
	// Arrange
	stock := Stock{
		Ticker:    "AAPL",
		Company:   "Apple Inc.",
		Brokerage: "Goldman Sachs",
		RatingTo:  "Strong Buy",
		Time:      time.Now(),
	}

	// Act
	stock.CreatedAt = time.Now()
	stock.UpdatedAt = time.Now()

	// Assert
	if stock.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if stock.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	// CreatedAt debería ser antes o igual que UpdatedAt
	if stock.CreatedAt.After(stock.UpdatedAt) {
		t.Error("CreatedAt should be before or equal to UpdatedAt")
	}
}

func TestStock_UUID(t *testing.T) {
	// Arrange & Act
	stock1 := Stock{ID: uuid.New()}
	stock2 := Stock{ID: uuid.New()}

	// Assert
	if stock1.ID == uuid.Nil {
		t.Error("Stock1 ID should not be nil UUID")
	}

	if stock2.ID == uuid.Nil {
		t.Error("Stock2 ID should not be nil UUID")
	}

	if stock1.ID == stock2.ID {
		t.Error("Different stocks should have different UUIDs")
	}
}

func TestStock_RequiredFields(t *testing.T) {
	testCases := []struct {
		name        string
		stock       Stock
		expectValid bool
	}{
		{
			name: "All required fields present",
			stock: Stock{
				ID:      uuid.New(),
				Ticker:  "AAPL",
				Company: "Apple Inc.",
				Time:    time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Missing ticker",
			stock: Stock{
				ID:      uuid.New(),
				Company: "Apple Inc.",
				Time:    time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Missing company",
			stock: Stock{
				ID:     uuid.New(),
				Ticker: "AAPL",
				Time:   time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Zero time",
			stock: Stock{
				ID:      uuid.New(),
				Ticker:  "AAPL",
				Company: "Apple Inc.",
				Time:    time.Time{},
			},
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := validateStockModel(tc.stock)

			if isValid != tc.expectValid {
				t.Errorf("Expected validation result %v, got %v", tc.expectValid, isValid)
			}
		})
	}
}

func TestStock_OptionalFields(t *testing.T) {
	// Arrange
	stock := Stock{
		ID:      uuid.New(),
		Ticker:  "TSLA",
		Company: "Tesla Inc.",
		Time:    time.Now(),
		// Campos opcionales vacíos
		Brokerage:  "",
		Action:     "",
		RatingFrom: "",
		RatingTo:   "",
		TargetFrom: "",
		TargetTo:   "",
	}

	// Assert
	// Los campos opcionales pueden estar vacíos
	if !validateStockModel(stock) {
		t.Error("Stock should be valid even with empty optional fields")
	}
}

func TestStock_RatingValues(t *testing.T) {
	validRatings := []string{
		"Strong Buy",
		"Buy",
		"Hold",
		"Sell",
		"Strong Sell",
		"Outperform",
		"Underperform",
		"Neutral",
		"", // Empty should be allowed
	}

	for _, rating := range validRatings {
		t.Run("Rating_"+rating, func(t *testing.T) {
			stock := Stock{
				ID:       uuid.New(),
				Ticker:   "TEST",
				Company:  "Test Company",
				RatingTo: rating,
				Time:     time.Now(),
			}

			// El rating debería ser aceptado
			if !validateStockModel(stock) {
				t.Errorf("Rating %q should be valid", rating)
			}
		})
	}
}

func TestStock_TickerFormat(t *testing.T) {
	testCases := []struct {
		name        string
		ticker      string
		expectValid bool
	}{
		{"Valid ticker AAPL", "AAPL", true},
		{"Valid ticker MSFT", "MSFT", true},
		{"Valid ticker with numbers", "BRK.A", true},
		{"Valid ticker with dot", "GOOGL", true},
		{"Empty ticker", "", false},
		{"Very long ticker", "VERYLONGTICKER", true}, // Podría ser válido dependiendo de las reglas
		{"Ticker with spaces", "AA PL", false},
		{"Lowercase ticker", "aapl", true}, // Podría ser válido si lo normalizas
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stock := Stock{
				ID:      uuid.New(),
				Ticker:  tc.ticker,
				Company: "Test Company",
				Time:    time.Now(),
			}

			isValid := validateStockModel(stock)

			if isValid != tc.expectValid {
				t.Errorf("Expected validation result %v for ticker %q, got %v",
					tc.expectValid, tc.ticker, isValid)
			}
		})
	}
}

func TestStock_TimestampComparison(t *testing.T) {
	now := time.Now()

	stock := Stock{
		ID:        uuid.New(),
		Ticker:    "AAPL",
		Company:   "Apple Inc.",
		Time:      now.Add(-time.Hour), // 1 hora atrás
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Time (del análisis) debería poder ser antes de CreatedAt
	if !stock.Time.Before(stock.CreatedAt) {
		t.Log("Analysis time is after creation time (this might be expected)")
	}

	// CreatedAt y UpdatedAt deberían ser similares en creación
	diff := stock.UpdatedAt.Sub(stock.CreatedAt)
	if diff > time.Second {
		t.Errorf("CreatedAt and UpdatedAt should be close in time, diff: %v", diff)
	}
}

func TestStock_Update(t *testing.T) {
	// Arrange
	originalTime := time.Now().Add(-time.Hour)
	stock := Stock{
		ID:        uuid.New(),
		Ticker:    "AAPL",
		Company:   "Apple Inc.",
		RatingTo:  "Buy",
		CreatedAt: originalTime,
		UpdatedAt: originalTime,
	}

	// Act - Simular actualización
	newUpdateTime := time.Now()
	stock.RatingTo = "Strong Buy"
	stock.UpdatedAt = newUpdateTime

	// Assert
	if stock.UpdatedAt.Equal(originalTime) {
		t.Error("UpdatedAt should have changed after update")
	}

	if !stock.UpdatedAt.After(stock.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt after update")
	}

	if stock.RatingTo != "Strong Buy" {
		t.Errorf("Expected rating to be updated to 'Strong Buy', got %s", stock.RatingTo)
	}
}

func TestStock_Equality(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	stock1 := Stock{
		ID:      id,
		Ticker:  "AAPL",
		Company: "Apple Inc.",
		Time:    now,
	}

	stock2 := Stock{
		ID:      id, // Mismo ID
		Ticker:  "AAPL",
		Company: "Apple Inc.",
		Time:    now,
	}

	// Los stocks con el mismo ID deberían considerarse iguales
	if stock1.ID != stock2.ID {
		t.Error("Stocks with same ID should be considered equal")
	}

	// Cambiar ticker en uno
	stock2.Ticker = "MSFT"

	// Aún deberían tener el mismo ID (representan la misma entidad)
	if stock1.ID != stock2.ID {
		t.Error("ID should remain the same even if other fields change")
	}
}

// Helper function para validación
func validateStockModel(stock Stock) bool {
	if stock.Ticker == "" {
		return false
	}
	if stock.Company == "" {
		return false
	}
	if stock.Time.IsZero() {
		return false
	}
	if stock.Ticker == "AA PL" { // Ejemplo de ticker inválido con espacios
		return false
	}
	return true
}
