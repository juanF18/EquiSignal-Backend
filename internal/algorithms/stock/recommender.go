package stock

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/juanF18/EquiSignal-Backend/internal/domain/models"
)

type StockRecommendation struct {
	Ticker     string
	Company    string
	Score      int
	Reason     string
	Rating     string
	TargetFrom string
	TargetTo   string
	Time       time.Time
}

func RecommendStocks(stocks []models.Stock, limit int) []StockRecommendation {
	recommendations := []StockRecommendation{}
	now := time.Now()

	// Mapas para análisis de tendencias y consenso
	tickerAnalysis := make(map[string][]models.Stock)
	brokerageWeight := getBrokerageWeights()

	// Agrupar por ticker para análisis de consenso
	for _, st := range stocks {
		tickerAnalysis[st.Ticker] = append(tickerAnalysis[st.Ticker], st)
	}

	for _, st := range stocks {
		score := 0.0
		reason := []string{}

		// === 1. ANÁLISIS DE RATING (Peso: 35%) ===
		ratingScore, ratingReason := calculateRatingScore(st)
		score += ratingScore * 0.35
		if ratingReason != "" {
			reason = append(reason, ratingReason)
		}

		// === 2. ANÁLISIS DE PRECIO OBJETIVO (Peso: 25%) ===
		targetScore, targetReason := calculateTargetScore(st)
		score += targetScore * 0.25
		if targetReason != "" {
			reason = append(reason, targetReason)
		}

		// === 3. ANÁLISIS TEMPORAL Y MOMENTUM (Peso: 20%) ===
		timeScore, timeReason := calculateTemporalScore(st, now)
		score += timeScore * 0.20
		if timeReason != "" {
			reason = append(reason, timeReason)
		}

		// === 4. CREDIBILIDAD DEL BROKERAGE (Peso: 10%) ===
		brokerScore, brokerReason := calculateBrokerageScore(st, brokerageWeight)
		score += brokerScore * 0.10
		if brokerReason != "" {
			reason = append(reason, brokerReason)
		}

		// === 5. CONSENSO DE MERCADO (Peso: 10%) ===
		consensusScore, consensusReason := calculateConsensusScore(st, tickerAnalysis[st.Ticker])
		score += consensusScore * 0.10
		if consensusReason != "" {
			reason = append(reason, consensusReason)
		}

		// === BONIFICACIONES ESPECIALES ===
		bonusScore, bonusReason := calculateBonusScore(st, now)
		score += bonusScore
		if bonusReason != "" {
			reason = append(reason, bonusReason)
		}

		// Convertir a entero para mantener compatibilidad
		finalScore := int(score * 10) // Multiplicar por 10 para más granularidad

		recommendations = append(recommendations, StockRecommendation{
			Ticker:     st.Ticker,
			Company:    st.Company,
			Score:      finalScore,
			Reason:     joinReasons(reason),
			Rating:     st.RatingTo,
			TargetFrom: st.TargetFrom,
			TargetTo:   st.TargetTo,
			Time:       st.Time,
		})
	}

	// Ordenar por score descendente
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Score == recommendations[j].Score {
			// En caso de empate, priorizar más recientes
			return recommendations[i].Time.After(recommendations[j].Time)
		}
		return recommendations[i].Score > recommendations[j].Score
	})

	// Filtrar duplicados por ticker, manteniendo el mejor score
	uniqueRecommendations := filterUniqueByTicker(recommendations)

	// Limitar top N
	if len(uniqueRecommendations) > limit {
		uniqueRecommendations = uniqueRecommendations[:limit]
	}

	return uniqueRecommendations
}

// === FUNCIONES AUXILIARES PARA ANÁLISIS AVANZADO ===

// getBrokerageWeights retorna pesos de credibilidad para diferentes brokerages
func getBrokerageWeights() map[string]float64 {
	return map[string]float64{
		// Tier 1: Analistas premium
		"Goldman Sachs":   1.0,
		"Morgan Stanley":  1.0,
		"JPMorgan":        1.0,
		"Bank of America": 0.95,
		"Citigroup":       0.95,
		"Wells Fargo":     0.9,
		"Barclays":        0.9,
		// Tier 2: Analistas sólidos
		"Deutsche Bank": 0.85,
		"Credit Suisse": 0.85,
		"UBS":           0.85,
		"Jefferies":     0.8,
		"Cowen":         0.8,
		"Piper Sandler": 0.75,
		// Tier 3: Otros
		"Default": 0.7, // Para brokerages no listados
	}
}

// calculateRatingScore analiza el rating con mayor sophisticación
func calculateRatingScore(stock models.Stock) (float64, string) {
	var score float64
	var reason string

	// Análisis del rating actual
	switch strings.ToLower(strings.TrimSpace(stock.RatingTo)) {
	case "strong buy", "outperform", "overweight":
		score = 10.0
		reason = "Strong Buy rating (+10.0)"
	case "buy", "positive":
		score = 7.5
		reason = "Buy rating (+7.5)"
	case "hold", "neutral", "market perform":
		score = 3.0
		reason = "Hold rating (+3.0)"
	case "underweight", "underperform":
		score = -5.0
		reason = "Underperform rating (-5.0)"
	case "sell", "strong sell":
		score = -10.0
		reason = "Sell rating (-10.0)"
	default:
		score = 2.0
		reason = "Unknown rating (+2.0)"
	}

	// Bonificación por upgrade de rating
	if stock.RatingFrom != "" && stock.RatingTo != "" {
		fromScore := getRatingNumericValue(stock.RatingFrom)
		toScore := getRatingNumericValue(stock.RatingTo)

		if toScore > fromScore {
			upgradeBonus := (toScore - fromScore) * 2.5
			score += upgradeBonus
			reason += fmt.Sprintf("; Upgrade bonus (+%.1f)", upgradeBonus)
		} else if toScore < fromScore {
			downgradepenalty := (fromScore - toScore) * 1.5
			score -= downgradepenalty
			reason += fmt.Sprintf("; Downgrade penalty (-%.1f)", downgradepenalty)
		}
	}

	return score, reason
}

// calculateTargetScore analiza precios objetivo con más detalle
func calculateTargetScore(stock models.Stock) (float64, string) {
	if stock.TargetFrom == "" || stock.TargetTo == "" {
		return 2.0, "No target data (+2.0)"
	}

	fromPrice := parsePrice(stock.TargetFrom)
	toPrice := parsePrice(stock.TargetTo)

	if fromPrice == 0 || toPrice == 0 {
		return 1.0, "Invalid target data (+1.0)"
	}

	// Calcular porcentaje de cambio en el precio objetivo
	percentChange := ((toPrice - fromPrice) / fromPrice) * 100

	var score float64
	var reason string

	switch {
	case percentChange >= 20:
		score = 8.0
		reason = fmt.Sprintf("Major target increase (+%.1f%%, +8.0)", percentChange)
	case percentChange >= 10:
		score = 6.0
		reason = fmt.Sprintf("Strong target increase (+%.1f%%, +6.0)", percentChange)
	case percentChange >= 5:
		score = 4.0
		reason = fmt.Sprintf("Moderate target increase (+%.1f%%, +4.0)", percentChange)
	case percentChange >= 0:
		score = 2.0
		reason = fmt.Sprintf("Small target increase (+%.1f%%, +2.0)", percentChange)
	case percentChange >= -5:
		score = -2.0
		reason = fmt.Sprintf("Minor target decrease (%.1f%%, -2.0)", percentChange)
	case percentChange >= -10:
		score = -4.0
		reason = fmt.Sprintf("Moderate target decrease (%.1f%%, -4.0)", percentChange)
	case percentChange >= -20:
		score = -6.0
		reason = fmt.Sprintf("Strong target decrease (%.1f%%, -6.0)", percentChange)
	default:
		score = -8.0
		reason = fmt.Sprintf("Major target decrease (%.1f%%, -8.0)", percentChange)
	}

	// Bonificación por precio objetivo alto (indica confianza)
	if toPrice > 100 {
		score += 1.0
		reason += "; High target confidence (+1.0)"
	}

	return score, reason
}

// calculateTemporalScore evalúa timing y momentum
func calculateTemporalScore(stock models.Stock, now time.Time) (float64, string) {
	days := now.Sub(stock.Time).Hours() / 24

	var score float64
	var reason string

	// Análisis de frescura de la información
	switch {
	case days <= 1:
		score = 6.0
		reason = "Breaking news (<1 day, +6.0)"
	case days <= 3:
		score = 5.0
		reason = "Very recent (1-3 days, +5.0)"
	case days <= 7:
		score = 4.0
		reason = "Recent (3-7 days, +4.0)"
	case days <= 14:
		score = 3.0
		reason = "Current (1-2 weeks, +3.0)"
	case days <= 30:
		score = 1.5
		reason = "Relevant (2-4 weeks, +1.5)"
	case days <= 60:
		score = 0.5
		reason = "Aging (1-2 months, +0.5)"
	default:
		score = -1.0
		reason = fmt.Sprintf("Stale (%.0f days, -1.0)", days)
	}

	// Bonificación por timing de mercado (evitar weekends en análisis crítico)
	weekday := stock.Time.Weekday()
	if weekday >= 1 && weekday <= 5 { // Monday to Friday
		score += 0.5
		reason += "; Market timing (+0.5)"
	}

	return score, reason
}

// calculateBrokerageScore evalúa credibilidad del brokerage
func calculateBrokerageScore(stock models.Stock, weights map[string]float64) (float64, string) {
	weight, exists := weights[stock.Brokerage]
	if !exists {
		weight = weights["Default"]
	}

	score := weight * 5.0 // Base score multiplied by weight
	reason := fmt.Sprintf("%s credibility (%.1fx, +%.1f)", stock.Brokerage, weight, score)

	return score, reason
}

// calculateConsensusScore analiza consenso de múltiples analistas
func calculateConsensusScore(stock models.Stock, allAnalysis []models.Stock) (float64, string) {
	if len(allAnalysis) <= 1 {
		return 1.0, "Single analysis (+1.0)"
	}

	// Contar ratings positivos vs negativos
	var positive, negative, neutral int
	for _, analysis := range allAnalysis {
		switch strings.ToLower(strings.TrimSpace(analysis.RatingTo)) {
		case "strong buy", "outperform", "overweight", "buy", "positive":
			positive++
		case "sell", "strong sell", "underweight", "underperform":
			negative++
		default:
			neutral++
		}
	}

	total := len(allAnalysis)
	positiveRatio := float64(positive) / float64(total)

	var score float64
	var reason string

	switch {
	case positiveRatio >= 0.8:
		score = 5.0
		reason = fmt.Sprintf("Strong consensus (%d/%d positive, +5.0)", positive, total)
	case positiveRatio >= 0.6:
		score = 3.0
		reason = fmt.Sprintf("Good consensus (%d/%d positive, +3.0)", positive, total)
	case positiveRatio >= 0.4:
		score = 1.0
		reason = fmt.Sprintf("Mixed consensus (%d/%d positive, +1.0)", positive, total)
	case positiveRatio >= 0.2:
		score = -1.0
		reason = fmt.Sprintf("Negative consensus (%d/%d positive, -1.0)", positive, total)
	default:
		score = -3.0
		reason = fmt.Sprintf("Very negative consensus (%d/%d positive, -3.0)", positive, total)
	}

	return score, reason
}

// calculateBonusScore aplica bonificaciones especiales
func calculateBonusScore(stock models.Stock, now time.Time) (float64, string) {
	var totalBonus float64
	var reasons []string

	// Bonificación por acción específica
	switch strings.ToLower(strings.TrimSpace(stock.Action)) {
	case "initiates", "initiated":
		totalBonus += 2.0
		reasons = append(reasons, "New coverage (+2.0)")
	case "reiterates", "reiterated":
		totalBonus += 1.0
		reasons = append(reasons, "Reaffirmed position (+1.0)")
	case "raises", "raised":
		totalBonus += 1.5
		reasons = append(reasons, "Raised expectations (+1.5)")
	case "lowers", "lowered":
		totalBonus -= 1.5
		reasons = append(reasons, "Lowered expectations (-1.5)")
	}

	// Bonificación por empresa de alto perfil (basado en nombre)
	companyName := strings.ToLower(stock.Company)
	highProfileKeywords := []string{"apple", "microsoft", "google", "amazon", "tesla", "nvidia", "meta"}
	for _, keyword := range highProfileKeywords {
		if strings.Contains(companyName, keyword) {
			totalBonus += 1.0
			reasons = append(reasons, "High-profile company (+1.0)")
			break
		}
	}

	// Penalización por volatilidad excesiva (múltiples cambios recientes)
	// Esta lógica se podría expandir con datos históricos

	reason := joinReasons(reasons)
	return totalBonus, reason
}

// filterUniqueByTicker elimina duplicados manteniendo el mejor score por ticker
func filterUniqueByTicker(recommendations []StockRecommendation) []StockRecommendation {
	seen := make(map[string]bool)
	unique := []StockRecommendation{}

	for _, rec := range recommendations {
		if !seen[rec.Ticker] {
			seen[rec.Ticker] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

// getRatingNumericValue convierte rating a valor numérico para comparaciones
func getRatingNumericValue(rating string) float64 {
	switch strings.ToLower(strings.TrimSpace(rating)) {
	case "strong buy", "outperform", "overweight":
		return 5.0
	case "buy", "positive":
		return 4.0
	case "hold", "neutral", "market perform":
		return 3.0
	case "underweight", "underperform":
		return 2.0
	case "sell", "strong sell":
		return 1.0
	default:
		return 3.0 // Default to neutral
	}
}

func joinReasons(rs []string) string {
	if len(rs) == 0 {
		return ""
	}
	out := ""
	for i, r := range rs {
		if i > 0 {
			out += "; "
		}
		out += r
	}
	return out
}

func parsePrice(s string) float64 {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", ".")
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return price
}
