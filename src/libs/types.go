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
