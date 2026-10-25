package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/api/handlers"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"go.uber.org/zap"
)

type APIServer struct {
	cfg             *config.Config
	logger          *zap.SugaredLogger
	healthHandler   *handlers.HealthHandler
	requestsHandler *handlers.RequestsHandler
}

func NewAPIServer(cfg *config.Config, log *zap.SugaredLogger, healthHandler *handlers.HealthHandler, requestHandler *handlers.RequestsHandler) *APIServer {
	return &APIServer{
		cfg:             cfg,
		logger:          log,
		healthHandler:   healthHandler,
		requestsHandler: requestHandler,
	}
}

func (s *APIServer) Run() error {
	r := s.mount()

	srv := &http.Server{
		Addr:         s.cfg.Dashboard.Addr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	s.logger.Infow("server is listening", "addr", s.cfg.Dashboard.Addr)

	return srv.ListenAndServe()
}

func (s *APIServer) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.healthHandler.Handle(s.healthHandler.Health))

		r.Route("/requests", func(r chi.Router) {
			r.Get("/", s.requestsHandler.Handle(s.requestsHandler.List))
			r.Get("/summary", s.requestsHandler.Handle(s.requestsHandler.Summary))
			r.Get("/timeseries", s.requestsHandler.Handle(s.requestsHandler.Timeseries))
			r.Get("/categories", s.requestsHandler.Handle(s.requestsHandler.Categories))

			r.Route("/blocklist", func(r chi.Router) {
				r.Get("/", s.requestsHandler.Handle(s.requestsHandler.Blocklist))
				r.Delete("/{ip}", s.requestsHandler.Handle(s.requestsHandler.BlocklistDelete))
			})
		})
	})

	r.Get("/*", s.spaHandler("./web/dist"))

	return r
}

func (s *APIServer) spaHandler(dir string) http.HandlerFunc {
	root := http.Dir(dir)
	fs := http.FileServer(root)

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		f, err := root.Open(path)
		if err != nil {
			http.ServeFile(w, r, filepath.Join(dir, "index.html"))
			return
		}
		f.Close()

		fs.ServeHTTP(w, r)
	}
}
