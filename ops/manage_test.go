package ops

import (
	"context"
	"expinc/sunagent/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManage_GetInfo(t *testing.T) {
	assert.Equal(t, common.Version, GetInfo(context.Background()).Version)
}
