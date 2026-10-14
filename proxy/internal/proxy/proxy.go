package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
	"go.uber.org/zap"
)

type Middleware func(route config.Route, next http.Handler) http.Handler

type Server struct {
	cfg       *config.Config
	logger    *zap.SugaredLogger
	srv       *http.Server
	transport *http.Transport
	store     *store.Storage
	detector  Middleware
}

func NewServer(cfg *config.Config, l *zap.SugaredLogger, s *store.Storage, mw Middleware) *Server {
	t := &http.Transport{
		DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
		ResponseHeaderTimeout: 20 * time.Second,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
	}

	return &Server{
		cfg:       cfg,
		logger:    l,
		transport: t,
		store:     s,
		detector:  mw,
	}
}

func (s *Server) Run() error {
	r, err := s.mount()
	if err != nil {
		return fmt.Errorf("failed to mount handlers: %w", err)
	}

	s.srv = &http.Server{
		Addr:         s.cfg.Proxy.Addr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	s.logger.Infow("server is listening", "addr", s.cfg.Proxy.Addr)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) mount() (http.Handler, error) {
	r := chi.NewRouter()

	for _, route := range s.cfg.Routes {
		handler, err := s.makeHandler(route)
		if err != nil {
			return nil, fmt.Errorf("route %s: %w", route.Path, err)
		}

		chiPath := chiPattern(route.Path)

		handler = s.detector(route, handler)

		r.Handle(chiPath, handler)
	}

	return r, nil
}

func (s *Server) makeHandler(route config.Route) (http.Handler, error) {
	switch route.Type {
	case config.StaticRoute:
		return s.makeStaticHandler(route), nil
	case config.PassRoute:
		return s.makePassHandler(route)
	default:
		return nil, fmt.Errorf("invalid route type")
	}
}

func (s *Server) makeStaticHandler(route config.Route) http.Handler {
	prefix := strings.TrimSuffix(route.Path, "/")
	return http.StripPrefix(prefix, http.FileServer(http.Dir(route.Target)))
}

func (s *Server) makePassHandler(route config.Route) (http.Handler, error) {
	target, err := url.Parse(route.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %w", route.Target, err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.SetXForwarded()
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			s.logger.Errorw("proxy error", "method", r.Method, "url", r.URL.Path, "err", err)
			http.Error(w, "bad gateway", http.StatusBadGateway)
		},
		Transport: s.transport,
	}

	return http.StripPrefix(route.Path, proxy), nil
}

func chiPattern(path string) string {
	if strings.HasSuffix(path, "/*") {
		return path
	}
	if path == "/" {
		return "/*"
	}
	return strings.TrimSuffix(path, "/") + "/*"
}
