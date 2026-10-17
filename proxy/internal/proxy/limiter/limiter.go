package limiter

import "context"

type Limiter interface {
	Check(ctx context.Context, ip string) Verdict
	RecordResponses(ctx context.Context, ip string, status int, isAuth bool) string
}
