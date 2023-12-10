package repository

import (
	"time"

	"github.com/shigaichi/tranco"
)

type TrancoAPIRepository interface {
	GetIDByDate(date time.Time) (tranco.ListMetadata, error)
}
