package common

import (
	"os"
	"path/filepath"
)

const (
	Version                = "1.3.0"
	TraceIdContextKey      = "traceId"
	ProcName               = "sunagentd"
	DefaultRegularFileMode = 0644
)

var (
	Pid                = os.Getpid()
	CurrentDir, _      = os.Getwd()
	TmpFolder          = filepath.Join(CurrentDir, "tmp")
	FileUploadMaxBytes = 104857600 // 100 MB
)
