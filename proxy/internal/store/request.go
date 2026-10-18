package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/db"
	"go.uber.org/zap"
)

type Request struct {
	ID          int64     `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	IP          string    `json:"ip"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Query       string    `json:"query"`
	UserAgent   string    `json:"user_agent"`
	ContentType string    `json:"content_type"`
	Features    string    `json:"-"`
	Route       string    `json:"route"`
	Source      string    `json:"source"`
	Category    string    `json:"category"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	StatusCode  int       `json:"status_code"`
	LatencyMS   int64     `json:"latency_ms"`
}

type Stats struct {
	TotalCount       int    `json:"total_count"`
	AttacksCount     int    `json:"attacks_count"`
	TopCategory      string `json:"top_category"`
	TopCategoryCount int    `json:"top_category_count"`
}

type Timeserie struct {
	Time    string `json:"time"`
	Total   int    `json:"total"`
	Attacks int    `json:"attacks"`
}

type Category struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
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

func (s *RequestStore) List(ctx context.Context, q PaginatedQuery) ([]Request, error) {
	query := `
		SELECT
			id, timestamp, ip, method, path, query, user_agent, content_type, route, source, category, score, action, status_code, latency_ms
		FROM requests_log
		OFFSET $1 LIMIT $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.pg.Pool.Query(ctx, query, q.Offset, q.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []Request{}

	for rows.Next() {
		var r Request

		err := rows.Scan(
			&r.ID,
			&r.Timestamp,
			&r.IP,
			&r.Method,
			&r.Path,
			&r.Query,
			&r.UserAgent,
			&r.ContentType,
			&r.Route,
			&r.Source,
			&r.Category,
			&r.Score,
			&r.Action,
			&r.StatusCode,
			&r.LatencyMS,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *RequestStore) Timeseries(ctx context.Context) ([]Timeserie, error) {
	query := `
		SELECT
			to_char(date_trunc('minute', timestamp), 'HH24:MI') AS time,
			count(*) AS total,
			count(*) FILTER (WHERE action NOT IN ('pass')) AS attack
		FROM requests_log
		WHERE timestamp > now() - interval '60 minutes'
		GROUP BY time ORDER BY time
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	results := []Timeserie{}

	rows, err := s.pg.Pool.Query(ctx, query)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var t Timeserie

		if err = rows.Scan(&t.Time, &t.Total, &t.Attacks); err != nil {
			return results, err
		}

		results = append(results, t)
	}

	if err = rows.Err(); err != nil {
		return results, err
	}

	return results, nil
}

func (s *RequestStore) Stats(ctx context.Context) (Stats, error) {
	totalCount, err := s.getCount(ctx)
	if err != nil {
		return Stats{}, err
	}

	attacksCount, err := s.getAttacksCount(ctx)
	if err != nil {
		return Stats{}, err
	}

	topCategory, topCount, err := s.getTopCategory(ctx)
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		TotalCount:       totalCount,
		AttacksCount:     attacksCount,
		TopCategory:      topCategory,
		TopCategoryCount: topCount,
	}, nil
}

func (s *RequestStore) Categories(ctx context.Context) ([]Category, error) {
	query := `
		SELECT category, count(id)
		FROM requests_log
		GROUP BY category
		ORDER BY count(id) DESC
	`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := []Category{}

	rows, err := s.pg.Pool.Query(ctx, query)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var c Category

		if err = rows.Scan(&c.Name, &c.Count); err != nil {
			return results, err
		}

		results = append(results, c)
	}

	return results, nil
}

func (s *RequestStore) Close() {
	s.worker.Close()
}

func (s *RequestStore) Insert(req *Request) {
	s.worker.Enqueue(req)
}

func (s *RequestStore) getCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(id) FROM requests_log
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var count int

	row := s.pg.Pool.QueryRow(ctx, query)

	if err := row.Scan(&count); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return 0, nil
		default:
			return 0, err
		}
	}

	return count, nil
}

func (s *RequestStore) getAttacksCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(id) FROM requests_log WHERE action NOT IN ('pass')
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var count int

	row := s.pg.Pool.QueryRow(ctx, query)

	if err := row.Scan(&count); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return 0, nil
		default:
			return 0, err
		}
	}

	return count, nil
}

func (s *RequestStore) getTopCategory(ctx context.Context) (string, int, error) {
	query := `
		SELECT category, COUNT(id) FROM requests_log GROUP BY category ORDER BY COUNT(id) DESC LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var top string
	var count int

	row := s.pg.Pool.QueryRow(ctx, query)

	if err := row.Scan(&top, &count); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return "", 0, nil
		default:
			return "", 0, err
		}
	}

	return top, count, nil
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
