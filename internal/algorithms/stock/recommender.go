package stock

import (
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

	for _, st := range stocks {
		score := 0
		reason := []string{}

		// --- Rating ---
		switch st.RatingTo {
		case "Strong Buy":
			score += 4
			reason = append(reason, "Rating Strong Buy (+4)")
		case "Buy":
			score += 3
			reason = append(reason, "Rating Buy (+3)")
		case "Hold":
			reason = append(reason, "Rating Hold (+0)")
		case "Sell":
			score -= 3
			reason = append(reason, "Rating Sell (-3)")
		}

		// --- Price Target ---
		if st.TargetFrom != "" && st.TargetTo != "" {
			from := parsePrice(st.TargetFrom)
			to := parsePrice(st.TargetTo)

			if from < to {
				score += 2
				reason = append(reason, "Positive target increase (+2)")
			} else if from > to {
				score -= 2
				reason = append(reason, "Negative target change (-2)")
			}
		}

		// --- Recency ---
		days := now.Sub(st.Time).Hours() / 24
		if days <= 7 {
			score += 3
			reason = append(reason, "Recent update (<7 days) (+3)")
		} else if days <= 30 {
			score += 1
			reason = append(reason, "Update within 30 days (+1)")
		}

		recommendations = append(recommendations, StockRecommendation{
			Ticker:     st.Ticker,
			Company:    st.Company,
			Score:      score,
			Reason:     joinReasons(reason),
			Rating:     st.RatingTo,
			TargetFrom: st.TargetFrom,
			TargetTo:   st.TargetTo,
			Time:       st.Time,
		})
	}

	// Ordenar por score desc
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limitar top N
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

func joinReasons(rs []string) string {
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
