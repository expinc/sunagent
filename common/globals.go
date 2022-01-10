package common

import (
	"os"
	"path/filepath"
)

const (
	Version           = "1.2.0"
	TraceIdContextKey = "traceId"
	ProcName          = "sunagentd"
)

var (
	Pid           = os.Getpid()
	CurrentDir, _ = os.Getwd()
	TmpFolder     = filepath.Join(CurrentDir, "tmp")
)
