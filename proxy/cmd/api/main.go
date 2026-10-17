package main

import (
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/api/handlers"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/api/server"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store/cache"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	cfg, err := config.Load("./config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	pool, err := db.NewPG(cfg.Database.Addr)
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Pool.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(logger, pool)

	rdb := cache.NewRedisClient(
		cfg.Cache.Addr,
		cfg.Cache.Pass,
		cfg.Cache.DB,
	)
	logger.Info("redis cache connection established")
	defer rdb.Close()

	cache := cache.NewStorage(rdb)

	healthHandler := handlers.NewHealthHandler(cfg, logger)
	requestsHandler := handlers.NewRequestsHandler(cfg, logger, store, cache)

	srv := server.NewAPIServer(cfg, logger, healthHandler, requestsHandler)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
