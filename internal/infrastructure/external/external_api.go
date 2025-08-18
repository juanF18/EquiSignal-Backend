package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juanF18/EquiSignal-Backend/internal/config"
	"github.com/juanF18/EquiSignal-Backend/internal/infrastructure/dto"
)

type ExternalAPI struct {
	cfg *config.Config
}

func NewExternalAPI(cfg *config.Config) *ExternalAPI {
	return &ExternalAPI{cfg: cfg}
}

func (e *ExternalAPI) FetchStocks(nextPage string) (*dto.StockResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.cfg.ExternalAPIURL, nil)
	if err != nil {
		return nil, err
	}

	// Auth header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.cfg.ExternalAPIToken))

	// query param
	if nextPage != "" {
		q := req.URL.Query()
		q.Add("next_page", nextPage)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stockResp dto.StockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResp); err != nil {
		return nil, err
	}

	return &stockResp, nil
}
