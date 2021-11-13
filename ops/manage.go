package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
)

type Info struct {
	Version string `json:"version"`
}

func GetInfo(ctx context.Context) Info {
	log.InfoCtx(ctx, "Getting sunagent info...")
	return Info{
		Version: common.Version,
	}
}
