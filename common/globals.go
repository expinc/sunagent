package common

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/units"
)

const (
	Version                = "1.3.0"
	TraceIdContextKey      = "traceId"
	ProcName               = "sunagentd"
	DefaultRegularFileMode = 0644
)

var (
	Pid                   = os.Getpid()
	CurrentDir, _         = os.Getwd()
	TmpFolder             = filepath.Join(CurrentDir, "tmp")
	FileTransferSizeLimit = 100 * units.MiB // 100 MB
)
