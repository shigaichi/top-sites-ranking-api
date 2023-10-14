package repository

import (
	"time"

	"github.com/shigaichi/tranco"
)

type TrancoApiRepository interface {
	GetIdByDate(date time.Time) (tranco.ListMetadata, error)
}
