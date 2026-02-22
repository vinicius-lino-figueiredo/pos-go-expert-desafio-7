// Package redisstrategy TODO
package redisstrategy

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/domain"
)

// RedisStrategy implements [domain.LimitStrategy].
type RedisStrategy struct {
	cl               *redis.Client
	ipLimit          int64
	ipLimitExpiry    time.Duration
	tokenLimit       int64
	tokenLimitExpiry time.Duration
}

// NewStorageStrategy TODO
func NewStorageStrategy(cl *redis.Client, ipLim, tknLim int64, ipExp, tknExp time.Duration) domain.LimitStrategy {
	return &RedisStrategy{
		cl:               cl,
		ipLimit:          ipLim,
		ipLimitExpiry:    ipExp,
		tokenLimit:       tknLim,
		tokenLimitExpiry: tknExp,
	}
}

// GetCountByIP implements [domain.LimitStrategy].
func (r *RedisStrategy) GetCountByIP(ctx context.Context, ip string) (bool, error) {
	c, err := r.increaseAndSetExpiry(ctx, "ip_res:"+ip, r.ipLimitExpiry)
	if err != nil {
		return false, err
	}
	return c <= r.ipLimit, err
}

// GetCountByToken implements [domain.LimitStrategy].
func (r *RedisStrategy) GetCountByToken(ctx context.Context, token string) (bool, error) {
	c, err := r.increaseAndSetExpiry(ctx, "token_res:"+token, r.tokenLimitExpiry)
	if err != nil {
		return false, err
	}
	return c <= r.tokenLimit, err
}

func (r *RedisStrategy) increaseAndSetExpiry(ctx context.Context, key string, expiry time.Duration) (int64, error) {
	c, err := r.cl.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if c == 1 {
		_ = r.cl.Expire(ctx, key, expiry)
	}

	return c, nil
}
