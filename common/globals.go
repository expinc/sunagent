package common

import "os"

const (
	Version           = "1.0.0"
	TraceIdContextKey = "traceId"
	ProcName          = "sunagentd"
)

var (
	Pid = os.Getpid()
)
