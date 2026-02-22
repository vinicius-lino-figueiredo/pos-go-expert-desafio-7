package redisstrategy_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/adapter/redisstrategy"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/domain"
)

type StrategyTestSuite struct {
	suite.Suite
	mr       *miniredis.Miniredis
	strategy domain.LimitStrategy
}

func (s *StrategyTestSuite) SetupTest() {
	var err error
	s.mr, err = miniredis.Run()
	s.Require().NoError(err)
	cl := redis.NewClient(&redis.Options{Addr: s.mr.Addr()})
	s.strategy = redisstrategy.NewStorageStrategy(cl, 2, 2, time.Minute, time.Minute)
}

func (s *StrategyTestSuite) TearDownTest() {
	s.mr.Close()
}

func (s *StrategyTestSuite) TestWithinIPLimit() {
	allowed, err := s.strategy.GetCountByIP(context.Background(), "192.168.1.1")
	s.Require().NoError(err)
	s.True(allowed)
}

func (s *StrategyTestSuite) TestExceedIPLimit() {
	for range 2 {
		_, _ = s.strategy.GetCountByIP(context.Background(), "192.168.1.1")
	}
	allowed, err := s.strategy.GetCountByIP(context.Background(), "192.168.1.1")
	s.Require().NoError(err)
	s.False(allowed)
}

func (s *StrategyTestSuite) TestWithinTokenLimit() {
	allowed, err := s.strategy.GetCountByToken(context.Background(), "abc123")
	s.Require().NoError(err)
	s.True(allowed)
}

func (s *StrategyTestSuite) TestExceedTokenLimit() {
	for range 2 {
		_, _ = s.strategy.GetCountByToken(context.Background(), "abc123")
	}
	allowed, err := s.strategy.GetCountByToken(context.Background(), "abc123")
	s.Require().NoError(err)
	s.False(allowed)
}

func TestStrategyTestSuite(t *testing.T) {
	suite.Run(t, new(StrategyTestSuite))
}
