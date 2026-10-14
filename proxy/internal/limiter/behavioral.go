package limiter

import (
	"context"
	"net/http"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store/cache"
	"go.uber.org/zap"
)

const (
	CategoryDOS        = "dos"
	CategoryBruteForce = "brute-force"
	CategoryScan       = "scan"
	CategoryError      = "error"
)

type Verdict struct {
	Blocked  bool
	Category string
}

type BehavioralLimiter struct {
	cfg    *config.Cache
	logger *zap.SugaredLogger
	cache  *cache.Storage
}

func NewBehavioralLimiter(cfg *config.Cache, log *zap.SugaredLogger, cache *cache.Storage) *BehavioralLimiter {
	return &BehavioralLimiter{
		cfg:    cfg,
		logger: log,
		cache:  cache,
	}
}

func (l *BehavioralLimiter) Check(ctx context.Context, ip string) Verdict {
	reason, blocked, err := l.cache.Blocklist.IsBlocked(ctx, ip)
	if err != nil {
		l.logger.Warnw("limiter: blocklist lookup failed", "err", err, "fail_open", l.cfg.FailOpen)
		if l.cfg.FailOpen {
			return Verdict{}
		}
		return Verdict{Blocked: true, Category: CategoryError}
	}

	if blocked {
		return Verdict{Blocked: true, Category: reason}
	}

	if cat := l.hit(ctx, ip, "dos:", l.cfg.DosWindowSec, l.cfg.DosMax, CategoryDOS); cat != "" {
		return Verdict{Blocked: true, Category: cat}
	}

	return Verdict{}
}

func (l *BehavioralLimiter) RecordResponses(ctx context.Context, ip string, status int, isAuth bool) string {
	switch {
	case isAuth && (status == http.StatusUnauthorized || status == http.StatusForbidden):
		return l.hit(ctx, ip, "bf:", l.cfg.BruteforceWindowSec, l.cfg.BruteforceMax, CategoryBruteForce)
	case status == http.StatusNotFound:
		return l.hit(ctx, ip, "scan:", l.cfg.ScanWindowSec, l.cfg.ScanMax, CategoryScan)
	}
	return ""
}

func (l *BehavioralLimiter) hit(ctx context.Context, ip, prefix string, windowSec, max int, category string) string {
	n, err := l.cache.Activity.IncrWindow(ctx, prefix+ip, windowSec)
	if err != nil {
		l.logger.Warnw("limiter: incr failed", "category", category, "ip", ip, "err", err)
		return ""
	}

	if int(n) <= max {
		return ""
	}

	l.logger.Infow("limiter: threshold exceeded - blocking", "category", category, "ip", ip, "count", n)
	if err = l.cache.Blocklist.Add(ctx, ip, category, l.cfg.BlockTTLSec); err != nil {
		l.logger.Warnw("limiter: blocklist add failed", "category", category, "ip", ip, "err", err)
	}

	return category
}
