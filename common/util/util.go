package util

import (
	"context"
)

func CheckCtx(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	default:
	}
	return true
}
