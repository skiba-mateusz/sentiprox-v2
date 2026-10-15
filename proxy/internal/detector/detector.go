package detector

import (
	"bytes"
	"context"
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
		authPaths[normalizePath(path)] = true
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
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		ip := clientIP(r)
		path := normalizePath(r.URL.Path)

		req := &store.Request{
			Timestamp:   start,
			IP:          ip,
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.RawQuery,
			UserAgent:   r.Header.Get("User-Agent"),
			ContentType: r.Header.Get("Content-Type"),
			Route:       route.Path,
			Category:    CategoryBenign,
		}
		decision := pass

		defer func() {
			status := sw.statusOrDefault()

			if route.BehavioralCheck && decision.Action == ActionPass {
				if cat := d.recordBehavior(r.Context(), ip, path, status); cat != "" {
					decision = flag(SourceBehavioral, cat)
				}
			}

			d.save(req, decision, status, start)
		}()

		decision = d.inspect(r, route, req, ip)
		if decision.Blocked() {
			http.Error(sw, http.StatusText(decision.Status), decision.Status)
			return
		}

		next.ServeHTTP(sw, r)
	})
}

func (d *Detector) inspect(r *http.Request, route config.Route, req *store.Request, ip string) Decision {
	if route.BehavioralCheck {
		if ver := d.limiter.Check(r.Context(), ip); ver.Blocked {
			status := http.StatusForbidden
			if ver.Category == limiter.CategoryDOS {
				status = http.StatusTooManyRequests
			}
			return block(SourceBehavioral, ver.Category, status)
		}
	}

	if route.MLCheck {
		return d.mlCheck(r, req)
	}

	return pass
}

func (d *Detector) mlCheck(r *http.Request, req *store.Request) Decision {
	candidate := d.buildCandidate(r)

	if d.cfg.ML.SaveFeatures {
		req.Features = candidate
	}

	suspicious, _ := prescreen(candidate)
	if !suspicious {
		return pass
	}

	res, err := d.client.Classify(r.Context(), candidate)
	if err != nil {
		d.logger.Warnw("classification failed", "err", err, "fail_open", d.cfg.ML.FailOpen)
		if d.cfg.ML.FailOpen {
			return pass
		}
		return block(SourceSystem, CategoryUnknown, http.StatusServiceUnavailable)
	}

	req.Category, req.Score = res.Category, res.Score
	if d.isThreat(res) {
		return block(SourceML, res.Category, http.StatusForbidden)
	}

	return pass
}

func (d *Detector) recordBehavior(ctx context.Context, ip, path string, status int) string {
	return d.limiter.RecordResponses(ctx, ip, status, d.authPaths[path])
}

func (d *Detector) save(req *store.Request, dec Decision, status int, start time.Time) {
	req.StatusCode = status
	req.Action = string(dec.Action)
	req.Source = string(dec.Source)
	if dec.Category != "" {
		req.Category = dec.Category
	}
	req.LatencyMS = time.Since(start).Milliseconds()

	d.store.Request.Insert(req)
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

func normalizePath(p string) string {
	if len(p) > 1 {
		p = strings.TrimSuffix(p, "/")
		if p == "" {
			return "/"
		}
	}
	return p
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
