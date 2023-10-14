package infra

import (
	"context"
	"time"

	"github.com/shigaichi/tranco"
)

type TrancoApiImpl struct {
	cli tranco.Client
}

func NewTrancoApiImpl() *TrancoApiImpl {
	cli := tranco.New()
	return &TrancoApiImpl{cli: *cli}
}

func (t TrancoApiImpl) GetIdByDate(date time.Time) (tranco.ListMetadata, error) {
	lists, err := t.cli.GetListMetadataByDate(context.Background(), date)
	if err != nil {
		return tranco.ListMetadata{}, err
	}

	return lists, nil
}
