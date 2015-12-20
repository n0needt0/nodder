package nodder

import (
	logging "github.com/op/go-logging"
)

const (
	STATE_UP   = 1
	STATE_DOWN = 0
	USER_AGENT = "Ironpoke/1.0 (https://www.ironpoke.com/)"
	LOGFMT     = "%{color}%{time:15:04:05.000000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"
)

var log = logging.MustGetLogger("logfile")
