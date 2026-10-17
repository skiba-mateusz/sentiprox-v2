package handlers

import (
	"net/http"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"go.uber.org/zap"
)

type HealthHandler struct {
	Base
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewHealthHandler(cfg *config.Config, logger *zap.SugaredLogger) *HealthHandler {
	return &HealthHandler{
		Base:   NewBase(logger),
		cfg:    cfg,
		logger: logger,
	}
}

func (h *HealthHandler) Health(r *http.Request) (any, error) {
	return map[string]string{"status": "ok"}, nil
}
