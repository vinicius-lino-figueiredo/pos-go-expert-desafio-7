// Package handler TODO
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vinicius-lino-figueiredo/pos-go-expert-desafio-7/domain"
)

const (
	tokenHeader                = "API_KEY"
	internalServerErrorMessage = "Internal server error"
	tooManyReqsErrorMessage    = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

// Handler TODO
type Handler struct {
	*chi.Mux
	strategy domain.LimitStrategy
}

// NewHandler TODO
func NewHandler(strategy domain.LimitStrategy) http.Handler {
	h := &Handler{
		Mux:      chi.NewMux(),
		strategy: strategy,
	}
	h.Use(h.verifyToken)
	h.Get("/*", h.handle)
	return h
}

func (h *Handler) handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello world"))
}

func (h *Handler) verifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get(tokenHeader)
		if tkn != "" {
			allowed, err := h.strategy.GetCountByToken(r.Context(), tkn)
			if err != nil {
				http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, tooManyReqsErrorMessage, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		allowed, err := h.strategy.GetCountByIP(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, tooManyReqsErrorMessage, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
