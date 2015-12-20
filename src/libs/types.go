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
	"os"
)

type AppData struct {
	Config map[string]string
}

func NewAppData() *AppData {
	return &AppData{Config: make(map[string]string)}
}

type AppChannels struct {
	Kill   chan os.Signal
	Loglvl chan string
}

func NewAppChannels() *AppChannels {
	return &AppChannels{
		Kill:   make(chan os.Signal, 1),
		Loglvl: make(chan string, 1),
	}
}
