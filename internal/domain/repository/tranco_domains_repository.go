package repository

import (
	"context"
)

type TrancoDomainsRepository interface {
	GetIdByDomain(ctx context.Context, domain string) (int, error)
	Save(ctx context.Context, domain string) (int, error)
}
