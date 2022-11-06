package db

import (
	"context"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type DomainBlacklist struct {
	ID     uint64 `db:"id"`
	Domain string `db:"domain"`
}

func (q *Queries) toTypesDomainBlacklist(d *DomainBlacklist) (*types.DomainBlacklist, error) {
	return &types.DomainBlacklist{
		Domain: d.Domain,
	}, nil
}

func (q *Queries) ListBlacklistedDomains(ctx context.Context) ([]*types.DomainBlacklist, error) {
	stmt := "SELECT * FROM domain_blocklist"
	rs, err := q.QueryxContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	domains := make([]*types.DomainBlacklist, 0)
	for rs.Next() {
		var domainBlacklist DomainBlacklist
		if err := rs.StructScan(&domainBlacklist); err != nil {
			return nil, err
		}

		lt, err := q.toTypesDomainBlacklist(&domainBlacklist)
		if err != nil {
			return nil, err
		}

		domains = append(domains, lt)
	}

	return domains, nil
}

func (q *Queries) IsBlacklistedDomain(ctx context.Context, domain string) (bool, error) {
	stmt := "SELECT count(*) FROM domain_blocklist WHERE domain = ?"
	var count int
	if err := q.GetContext(ctx, &count, stmt, domain); err != nil {
		return false, err
	}

	return count == 1, nil
}
