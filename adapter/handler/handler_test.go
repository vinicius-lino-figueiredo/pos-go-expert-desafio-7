package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/adapter/handler"
)

type mockStrategy struct {
	mock.Mock
}

func (m *mockStrategy) GetCountByIP(ctx context.Context, ip string) (bool, error) {
	args := m.Called(ctx, ip)
	return args.Bool(0), args.Error(1)
}

func (m *mockStrategy) GetCountByToken(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

type HandlerTestSuite struct {
	suite.Suite
	strategy *mockStrategy
	server   *httptest.Server
}

func (s *HandlerTestSuite) SetupTest() {
	s.strategy = new(mockStrategy)
	s.server = httptest.NewServer(handler.NewHandler(s.strategy))
}

func (s *HandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *HandlerTestSuite) TestAllowedIPRequest() {
	s.strategy.On("GetCountByIP", mock.Anything, mock.Anything).Return(true, nil)
	resp, err := http.Get(s.server.URL)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestBlockedIPRequest() {
	s.strategy.On("GetCountByIP", mock.Anything, mock.Anything).Return(false, nil)
	resp, err := http.Get(s.server.URL)
	s.Require().NoError(err)
	s.Equal(http.StatusTooManyRequests, resp.StatusCode)
}

func (s *HandlerTestSuite) TestAllowedTokenRequest() {
	s.strategy.On("GetCountByToken", mock.Anything, "my-token").Return(true, nil)
	req, _ := http.NewRequest(http.MethodGet, s.server.URL, nil)
	req.Header.Set("API_KEY", "my-token")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *HandlerTestSuite) TestBlockedTokenRequest() {
	s.strategy.On("GetCountByToken", mock.Anything, "my-token").Return(false, nil)
	req, _ := http.NewRequest(http.MethodGet, s.server.URL, nil)
	req.Header.Set("API_KEY", "my-token")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusTooManyRequests, resp.StatusCode)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
