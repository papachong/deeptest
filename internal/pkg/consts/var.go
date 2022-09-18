package consts

import "sync"

var (
	Language = "zh"
	Verbose  = true

	ExecFrom  ExecFromType
	IsRelease bool
	ExecDir   string
	WorkDir   string

	EnvVar sync.Map
)
