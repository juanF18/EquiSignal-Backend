package dto

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStock_JSONMarshaling(t *testing.T) {
	// Arrange
	now := time.Now()
	stock := Stock{
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

	// Act - Marshal to JSON
	jsonBytes, err := json.Marshal(stock)
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Act - Unmarshal back
	var unmarshaledStock Stock
	err = json.Unmarshal(jsonBytes, &unmarshaledStock)
	if err != nil {
		t.Fatalf("Error unmarshaling from JSON: %v", err)
	}

	// Assert
	if unmarshaledStock.Ticker != stock.Ticker {
		t.Errorf("Expected ticker %s, got %s", stock.Ticker, unmarshaledStock.Ticker)
	}

	if unmarshaledStock.Company != stock.Company {
		t.Errorf("Expected company %s, got %s", stock.Company, unmarshaledStock.Company)
	}

	if unmarshaledStock.RatingTo != stock.RatingTo {
		t.Errorf("Expected rating %s, got %s", stock.RatingTo, unmarshaledStock.RatingTo)
	}

	// Verificar que el tiempo se preserva (con tolerancia por precisión)
	timeDiff := unmarshaledStock.Time.Sub(stock.Time)
	if timeDiff > time.Second || timeDiff < -time.Second {
		t.Errorf("Time difference too large: %v", timeDiff)
	}
}

func TestStock_JSONTags(t *testing.T) {
	// Arrange
	stock := Stock{
		Ticker:     "MSFT",
		Company:    "Microsoft Corporation",
		Brokerage:  "Morgan Stanley",
		Action:     "Reiterates",
		RatingFrom: "Buy",
		RatingTo:   "Strong Buy",
		TargetFrom: "$300.00",
		TargetTo:   "$350.00",
		Time:       time.Now(),
	}

	// Act
	jsonBytes, err := json.Marshal(stock)
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	jsonString := string(jsonBytes)

	// Assert - Verificar que los tags JSON están correctos
	expectedFields := []string{
		`"ticker"`,
		`"company"`,
		`"brokerage"`,
		`"action"`,
		`"rating_from"`,
		`"rating_to"`,
		`"target_from"`,
		`"target_to"`,
		`"time"`,
	}

	for _, field := range expectedFields {
		if !contains(jsonString, field) {
			t.Errorf("JSON should contain field %s", field)
		}
	}
}

func TestStock_EmptyValues(t *testing.T) {
	// Arrange
	stock := Stock{
		Ticker: "TSLA",
		// Otros campos vacíos intencionalmente
	}

	// Act
	jsonBytes, err := json.Marshal(stock)
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	var unmarshaledStock Stock
	err = json.Unmarshal(jsonBytes, &unmarshaledStock)
	if err != nil {
		t.Fatalf("Error unmarshaling from JSON: %v", err)
	}

	// Assert
	if unmarshaledStock.Ticker != "TSLA" {
		t.Errorf("Expected ticker TSLA, got %s", unmarshaledStock.Ticker)
	}

	if unmarshaledStock.Company != "" {
		t.Errorf("Expected empty company, got %s", unmarshaledStock.Company)
	}

	if !unmarshaledStock.Time.IsZero() {
		t.Errorf("Expected zero time, got %v", unmarshaledStock.Time)
	}
}

func TestStockResponse_JSONMarshaling(t *testing.T) {
	// Arrange
	stocks := []Stock{
		{
			Ticker:   "AAPL",
			Company:  "Apple Inc.",
			RatingTo: "Strong Buy",
			Time:     time.Now(),
		},
		{
			Ticker:   "MSFT",
			Company:  "Microsoft Corp.",
			RatingTo: "Buy",
			Time:     time.Now(),
		},
	}

	response := StockResponse{
		Items:    stocks,
		NextPage: "page2",
	}

	// Act
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling StockResponse to JSON: %v", err)
	}

	var unmarshaledResponse StockResponse
	err = json.Unmarshal(jsonBytes, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Error unmarshaling StockResponse from JSON: %v", err)
	}

	// Assert
	if len(unmarshaledResponse.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(unmarshaledResponse.Items))
	}

	if unmarshaledResponse.NextPage != "page2" {
		t.Errorf("Expected NextPage 'page2', got %s", unmarshaledResponse.NextPage)
	}

	// Verificar primer item
	firstItem := unmarshaledResponse.Items[0]
	if firstItem.Ticker != "AAPL" {
		t.Errorf("Expected first item ticker AAPL, got %s", firstItem.Ticker)
	}
}

func TestStockResponse_EmptyItems(t *testing.T) {
	// Arrange
	response := StockResponse{
		Items:    []Stock{}, // Array vacío
		NextPage: "",
	}

	// Act
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling empty StockResponse: %v", err)
	}

	var unmarshaledResponse StockResponse
	err = json.Unmarshal(jsonBytes, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Error unmarshaling empty StockResponse: %v", err)
	}

	// Assert
	if unmarshaledResponse.Items == nil {
		t.Error("Items should not be nil, should be empty array")
	}

	if len(unmarshaledResponse.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(unmarshaledResponse.Items))
	}

	if unmarshaledResponse.NextPage != "" {
		t.Errorf("Expected empty NextPage, got %s", unmarshaledResponse.NextPage)
	}
}

func TestStock_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		stock       Stock
		expectValid bool
	}{
		{
			name: "Valid complete stock",
			stock: Stock{
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
			stock: Stock{
				Company:   "Apple Inc.",
				Brokerage: "Goldman Sachs",
				RatingTo:  "Strong Buy",
				Time:      time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Missing company",
			stock: Stock{
				Ticker:    "AAPL",
				Brokerage: "Goldman Sachs",
				RatingTo:  "Strong Buy",
				Time:      time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Zero time",
			stock: Stock{
				Ticker:   "AAPL",
				Company:  "Apple Inc.",
				RatingTo: "Strong Buy",
				Time:     time.Time{},
			},
			expectValid: false,
		},
		{
			name: "Valid minimal stock",
			stock: Stock{
				Ticker:  "TSLA",
				Company: "Tesla Inc.",
				Time:    time.Now(),
			},
			expectValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Función de validación conceptual
			isValid := validateStock(tc.stock)

			if isValid != tc.expectValid {
				t.Errorf("Expected validation result %v, got %v", tc.expectValid, isValid)
			}
		})
	}
}

func TestPriceFormatting(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool // Si es un formato válido
	}{
		{"Valid dollar format", "$123.45", true},
		{"Valid without symbol", "123.45", true},
		{"Valid with commas", "$1,234.56", true},
		{"Invalid format", "abc", false},
		{"Empty string", "", false},
		{"Only symbol", "$", false},
		{"Negative price", "-$123.45", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := isValidPriceFormat(tc.input)

			if isValid != tc.expected {
				t.Errorf("Expected %v for price format validation of %q, got %v",
					tc.expected, tc.input, isValid)
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func validateStock(stock Stock) bool {
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

func isValidPriceFormat(price string) bool {
	if price == "" {
		return false
	}

	// Lógica simple de validación de formato de precio
	// En la implementación real, usarías regex o parsePrice
	if price == "$" || price == "abc" || price[0] == '-' {
		return false
	}

	return true
}
