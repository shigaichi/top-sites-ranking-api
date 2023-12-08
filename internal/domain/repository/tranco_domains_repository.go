package repository

import (
	"context"
)

type TrancoDomainsRepository interface {
	GetIDByDomain(ctx context.Context, domain string) (int, error)
	Save(ctx context.Context, domain string) (int, error)
}
