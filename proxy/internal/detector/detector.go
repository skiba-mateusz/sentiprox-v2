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
	"github.com/skiba-mateusz/sentiprox-v2/proxy/internal/store"
	"go.uber.org/zap"
)

type Action int

const (
	PassAction Action = iota
	BlockAction
	ErrorAction
)

func (a Action) String() string {
	switch a {
	case PassAction:
		return "pass"
	case BlockAction:
		return "block"
	case ErrorAction:
		return "error"
	default:
		return "unknown"
	}
}

type Detector struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
	store  *store.Storage
	client client.Client
}

func NewDetector(cfg *config.Config, l *zap.SugaredLogger, c client.Client, s *store.Storage) *Detector {
	return &Detector{
		cfg:    cfg,
		logger: l,
		client: c,
		store:  s,
	}
}

func (d *Detector) Wrap(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		candidate := d.buildCandidate(r)
		action := PassAction

		req := &store.Request{
			Timestamp:   start,
			IP:          clientIP(r),
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.RawQuery,
			UserAgent:   r.Header.Get("User-Agent"),
			ContentType: r.Header.Get("Content-Type"),
			Features:    candidate,
			Route:       route,
			Category:    "benign",
		}
		defer func() {
			req.Action = action.String()
			req.StatusCode = sw.statusOrDefault()
			req.LatencyMS = time.Since(start).Milliseconds()
			d.store.Request.Insert(req)
		}()

		var res client.Result
		var err error

		if suspicious, _ := prescreen(candidate); suspicious {
			res, err = d.client.Classify(r.Context(), candidate)
			if err != nil {
				d.logger.Warnw("classification failed", "err", err, "fail_open", d.cfg.ML.FailOpen)
				if !d.cfg.ML.FailOpen {
					action = ErrorAction
					http.Error(w, "service unavailable", http.StatusServiceUnavailable)
					return
				}
				next.ServeHTTP(sw, r)
				return
			}
			req.Category, req.Score = res.Category, res.Score
		}

		if d.isThreat(res) {
			action = BlockAction
			http.Error(w, "request rejected", http.StatusForbidden)
			return
		}

		next.ServeHTTP(sw, r)
	})
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
