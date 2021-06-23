package ops

import (
	"context"
	"expinc/sunagent/common"
)

type Info struct {
	Version string `json:"version"`
}

func GetInfo(ctx context.Context) Info {
	return Info{
		Version: common.Version,
	}
}
