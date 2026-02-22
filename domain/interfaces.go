// Package domain TODO
package domain

import (
	"context"
)

// LimitStrategy TODO
type LimitStrategy interface {
	GetCountByIP(ctx context.Context, ip string) (bool, error)
	GetCountByToken(ctx context.Context, token string) (bool, error)
}
