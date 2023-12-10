package infra

import (
	"context"
	"time"

	"github.com/shigaichi/tranco"
)

type TrancoAPIImpl struct {
	cli tranco.Client
}

func NewTrancoAPIImpl() *TrancoAPIImpl {
	cli := tranco.New()
	return &TrancoAPIImpl{cli: *cli}
}

func (t TrancoAPIImpl) GetIDByDate(date time.Time) (tranco.ListMetadata, error) {
	lists, err := t.cli.GetListMetadataByDate(context.Background(), date)
	if err != nil {
		return tranco.ListMetadata{}, err
	}

	return lists, nil
}
