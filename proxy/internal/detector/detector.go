package detector

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/client"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/config"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/limiter"
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
	"go.uber.org/zap"
)

type Action int

const (
	PassAction Action = iota
	AttackAction
	BlockAction
	FlagAction
	ErrorAction
)

func (a Action) String() string {
	switch a {
	case PassAction:
		return "pass"
	case AttackAction:
		return "attack"
	case BlockAction:
		return "block"
	case FlagAction:
		return "flag"
	case ErrorAction:
		return "error"
	default:
		return "unknown"
	}
}

type Detector struct {
	cfg       *config.Config
	logger    *zap.SugaredLogger
	store     *store.Storage
	client    client.Client
	limiter   limiter.Limiter
	authPaths map[string]bool
}

func NewDetector(cfg *config.Config, log *zap.SugaredLogger, store *store.Storage, client client.Client, limiter limiter.Limiter) *Detector {
	authPaths := make(map[string]bool, len(cfg.AuthPaths))
	for _, path := range cfg.AuthPaths {
		authPaths[path] = true
	}

	return &Detector{
		cfg:       cfg,
		logger:    log,
		client:    client,
		store:     store,
		limiter:   limiter,
		authPaths: authPaths,
	}
}

func (d *Detector) Wrap(route config.Route, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w}
		start := time.Now()
		ip := clientIP(r)
		ctx := r.Context()
		action := PassAction
		path := strings.TrimSuffix(r.URL.Path, "/")

		req := &store.Request{
			Timestamp:   start,
			IP:          ip,
			Method:      r.Method,
			Path:        path,
			Query:       r.URL.RawQuery,
			UserAgent:   r.Header.Get("User-Agent"),
			ContentType: r.Header.Get("Content-Type"),
			Route:       route.Path,
			Category:    "benign",
		}
		defer func() {
			status := sw.statusOrDefault()
			req.StatusCode = status
			req.Action = action.String()
			req.LatencyMS = time.Since(start).Milliseconds()

			if route.BehavioralCheck && action == PassAction {
				if cat := d.limiter.RecordResponses(ctx, ip, status, d.authPaths[path]); cat != "" {
					action = FlagAction
					req.Action = action.String()
					req.Category = cat
				}
			}

			d.store.Request.Insert(req)
		}()

		reject := func(a Action, category string, code int, msg string) {
			action = a
			req.Category = category
			http.Error(sw, msg, code)
		}

		if route.BehavioralCheck {
			if ver := d.limiter.Check(ctx, ip); ver.Blocked {
				code := http.StatusForbidden
				if ver.Category == limiter.CategoryDOS {
					code = http.StatusTooManyRequests
				}
				reject(BlockAction, ver.Category, code, http.StatusText(code))
				return
			}
		}

		if route.MLCheck {
			if a, code, msg := d.mlCheck(r, req); a != PassAction {
				reject(a, req.Category, code, msg)
				return
			}
		}

		next.ServeHTTP(sw, r)
	})
}

func (d *Detector) mlCheck(r *http.Request, req *store.Request) (Action, int, string) {
	candidate := d.buildCandidate(r)

	if d.cfg.ML.SaveFeatures {
		req.Features = candidate
	}

	suspicious, _ := prescreen(candidate)
	if !suspicious {
		return PassAction, 0, ""
	}

	res, err := d.client.Classify(r.Context(), candidate)
	if err != nil {
		d.logger.Warnw("classification failed", "err", err, "fail_open", d.cfg.ML.FailOpen)
		if d.cfg.ML.FailOpen {
			return PassAction, 0, ""
		}
		req.Category = "error"
		return ErrorAction, http.StatusServiceUnavailable, "service unavailable"
	}

	req.Category, req.Score = res.Category, res.Score
	if d.isThreat(res) {
		return AttackAction, http.StatusForbidden, "request rejected"
	}

	return PassAction, 0, ""
}

func (d *Detector) isThreat(res client.Result) bool {
	return res.Category != "benign" && res.Score >= d.cfg.ML.Score
}

func (d *Detector) buildCandidate(r *http.Request) string {
	var sb strings.Builder
	sb.WriteString(r.URL.Path)
	sb.WriteString(" ")
	sb.WriteString(r.URL.RawQuery)

	if body := d.bodySample(r); body != "" {
		sb.WriteString(" ")
		sb.WriteString(body)
	}

	sb.WriteString(" ")
	sb.WriteString(r.Header.Get("User-Agent"))
	sb.WriteString(" ")
	sb.WriteString(r.Header.Get("Referer"))

	return normalize(sb.String())
}

func (d *Detector) bodySample(r *http.Request) string {
	if r.Body == nil && (r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch) {
		return ""
	}

	mt, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return ""
	}

	textual := mt == "application/json" ||
		mt == "application/x-www-form-urlencoded" ||
		mt == "application/xml" ||
		strings.HasPrefix(mt, "text/") ||
		strings.HasSuffix(mt, "+json") ||
		strings.HasSuffix(mt, "+xml")

	if !textual && mt != "multipart/form-data" {
		return ""
	}

	head, _ := io.ReadAll(io.LimitReader(r.Body, int64(d.cfg.ML.InspectBodyBytes)))
	r.Body = struct {
		io.Reader
		io.Closer
	}{io.MultiReader(bytes.NewReader(head), r.Body), r.Body}

	if mt == "multipart/form-data" {
		return multipartSummary(head, params["boundary"], 4096)
	}

	return string(head)
}

func multipartSummary(head []byte, boundary string, maxField int64) string {
	var sb strings.Builder
	mr := multipart.NewReader(bytes.NewReader(head), boundary)

	for {
		p, err := mr.NextPart()
		if err != nil {
			break
		}
		sb.WriteString(p.FormName())
		sb.WriteString(" ")
		if fn := p.FileName(); fn != "" {
			sb.WriteString(fn)
			sb.WriteString(" ")
			sb.WriteString(p.Header.Get("Content-Type"))
			sb.WriteString(" ")
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(p, maxField))
		sb.Write(b)
		sb.WriteString(" ")
	}

	return sb.String()
}

func normalize(s string) string {
	cur := s
	for i := 0; i < 3; i++ {
		dec, err := url.QueryUnescape(cur)
		if err != nil || dec == cur {
			break
		}
		cur = dec
	}
	return strings.ToLower(cur)
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) statusOrDefault() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
