package infra

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoCsvImpl struct {
}

func NewTrancoCsvImpl() *TrancoCsvImpl {
	return &TrancoCsvImpl{}
}

func (t TrancoCsvImpl) Get(url url.URL) ([]model.SiteRanking, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to download csv from %s. error: %w", url.String(), err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to download csv from %s. response status: %d", url.String(), resp.StatusCode)
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error while parsing CSV: %w", err)
	}

	var rankings []model.SiteRanking
	for _, record := range records {
		rank, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting rank to int: %w", err)
		}
		rankings = append(rankings, model.SiteRanking{
			Rank:   rank,
			Domain: record[1],
		})
	}

	return rankings, nil
}
