package repository

import (
	"net/url"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoCsvRepository interface {
	Get(url url.URL) ([]model.SiteRanking, error)
}
