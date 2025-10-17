package application

import "context"

type LeaderLookup interface {
	FindLeaderIDByDomain(ctx context.Context, domain string) (int64, error)
}
