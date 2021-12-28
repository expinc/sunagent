package util

import (
	"context"
	"expinc/sunagent/log"
)

func LogErrorIfNotNilCtx(ctx context.Context, err error) {
	if err != nil {
		log.ErrorCtx(ctx, err)
	}
}
