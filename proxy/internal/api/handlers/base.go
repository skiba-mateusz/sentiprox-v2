package handlers

import (
	"errors"
	"net/http"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/api/respond"
	"go.uber.org/zap"
)

type HandlerFunc func(r *http.Request) (any, error)

type Base struct {
	logger *zap.SugaredLogger
}

func NewBase(logger *zap.SugaredLogger) Base {
	return Base{
		logger: logger,
	}
}

func (b Base) Handle(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fn(r)
		if err != nil {
			var apiErr *respond.APIError
			if !errors.As(err, &apiErr) {
				apiErr = &respond.APIError{Status: http.StatusInternalServerError, Message: "internal server error", Err: err}
			}

			if apiErr.Status >= 500 {
				b.logger.Errorw("request failed", "err", err, "method", r.Method, "path", r.URL.Path)
			}

			if err = respond.Error(w, apiErr.Status, apiErr.Message); err != nil {
				b.logger.Warnw("write error response failed", "err", err)
			}

			return
		}

		if err = respond.JSON(w, http.StatusOK, data); err != nil {
			b.logger.Warnw("write response failed", "err", err, "path", r.URL.Path)
		}
	}
}
