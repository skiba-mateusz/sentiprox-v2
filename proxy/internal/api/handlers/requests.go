package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/api/respond"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store/cache"
	"go.uber.org/zap"
)

type RequestsHandler struct {
	Base
	cfg    *config.Config
	logger *zap.SugaredLogger
	store  *store.Storage
	cache  *cache.Storage
}

func NewRequestsHandler(cfg *config.Config, logger *zap.SugaredLogger, store *store.Storage, cache *cache.Storage) *RequestsHandler {
	return &RequestsHandler{
		Base:   NewBase(logger),
		cfg:    cfg,
		logger: logger,
		store:  store,
		cache:  cache,
	}
}

func (h *RequestsHandler) List(r *http.Request) (any, error) {
	query := store.PaginatedQuery{
		Limit:  20,
		Offset: 0,
	}

	query, err := query.Parse(r)
	if err != nil {
		return nil, respond.BadRequest("invalid query parameters", err)
	}

	return h.store.Request.List(r.Context(), query)
}

type Summary struct {
	store.Stats
	TotalBlocks int `json:"total_blocks"`
}

func (h *RequestsHandler) Summary(r *http.Request) (any, error) {
	stats, err := h.store.Request.Stats(r.Context())
	if err != nil {
		return stats, err
	}

	entries, err := h.cache.Blocklist.List(r.Context())
	if err != nil {
		return stats, err
	}

	s := Summary{
		Stats:       stats,
		TotalBlocks: len(entries),
	}

	return s, nil
}

func (h *RequestsHandler) Timeseries(r *http.Request) (any, error) {
	return h.store.Request.Timeseries(r.Context())
}

func (h *RequestsHandler) Categories(r *http.Request) (any, error) {
	return h.store.Request.Categories(r.Context())
}

func (h *RequestsHandler) Blocklist(r *http.Request) (any, error) {
	return h.cache.Blocklist.List(r.Context())
}

func (h *RequestsHandler) BlocklistDelete(r *http.Request) (any, error) {
	ip := chi.URLParam(r, "ip")

	return nil, h.cache.Blocklist.Delete(r.Context(), ip)
}
