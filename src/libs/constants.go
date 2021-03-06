/**
Copyright 2015 andrew@yasinsky.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
