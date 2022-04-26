package metric

import "context"

type logger interface {
	Error(ctx context.Context, err error)
}
