package client

import "context"

type Client interface {
	Classify(ctx context.Context, content string) (Result, error)
}
