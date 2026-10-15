package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"go.uber.org/zap"
)

type Request struct {
	ID          int64
	Timestamp   time.Time
	IP          string
	Method      string
	Path        string
	Query       string
	UserAgent   string
	ContentType string
	Features    string
	Route       string
	Source      string
	Category    string
	Score       float64
	Action      string
	StatusCode  int
	LatencyMS   int64
}

type RequestStore struct {
	logger *zap.SugaredLogger
	pg     *db.PG
	worker *db.Worker[*Request]
}

func newRequestStore(l *zap.SugaredLogger, pg *db.PG) *RequestStore {
	s := &RequestStore{
		pg: pg,
	}
	s.worker = db.NewWorker(l, s.flush, 200, time.Duration(time.Second*2))
	return s
}

func (s *RequestStore) Close() {
	s.worker.Close()
}

func (s *RequestStore) Insert(req *Request) {
	s.worker.Enqueue(req)
}

func (s *RequestStore) flush(batch []*Request) error {
	rows := make([][]any, len(batch))
	for i, r := range batch {
		rows[i] = []any{r.Timestamp, r.IP, r.Method, r.Path, r.Query, r.UserAgent, r.ContentType, r.Features, r.Route, r.Source, r.Category, r.Score, r.Action, r.StatusCode, r.LatencyMS}
	}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	_, err := s.pg.Pool.CopyFrom(ctx,
		pgx.Identifier{"requests_log"}, []string{"timestamp", "ip", "method", "path", "query", "user_agent", "content_type", "features", "route", "source", "category", "score", "action", "status_code", "latency_ms"},
		pgx.CopyFromRows(rows),
	)
	return err
}
