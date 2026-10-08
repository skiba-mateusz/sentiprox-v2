package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/client"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/detector"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/proxy"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
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
	logger.Info("db connection pool established")

	store := store.NewStorage(logger, pool)
	defer store.Request.Close()

	mlClient := client.NewMLClient(cfg.ML.Addr, time.Duration(cfg.ML.TimeoutMS)*time.Millisecond)
	detector := detector.NewDetector(cfg, logger, mlClient, store)

	srv := proxy.NewServer(cfg, logger, store, detector.Wrap)

	go func() {
		if err = srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("server stopped")
}
