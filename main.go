// main TODO
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/adapter/handler"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/adapter/redisstrategy"
)

type configs struct {
	IPLimit          int64         `env:"IP_REQUEST_LIMIT,required"`
	IPLimitExpiry    time.Duration `env:"IP_REQUEST_LIMIT_EXPIRY"`
	TokenLimit       int64         `env:"TOKEN_REQUEST_LIMIT"`
	TokenLimitExpiry time.Duration `env:"TOKEN_REQUEST_LIMIT_EXPIRY"`
	RedisAddr        string        `env:"REDIS_ADDRESS"`
}

func main() {
	if err := godotenv.Overload(); err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := env.ParseAs[configs]()
	if err != nil {
		log.Fatal(err.Error())
	}

	cl := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	strat := redisstrategy.NewStorageStrategy(
		cl,
		cfg.IPLimit,
		cfg.TokenLimit,
		cfg.IPLimitExpiry,
		cfg.TokenLimitExpiry,
	)

	h := handler.NewHandler(strat)

	server := http.Server{Addr: ":8080", Handler: h}

	ch := make(chan os.Signal, 1)
	srvErr := make(chan error, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		srvErr <- server.ListenAndServe()
	}()

	select {
	case <-ch:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	case err := <-srvErr:
		log.Fatal(err.Error())
	}

}
