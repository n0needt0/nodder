package nodder

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	logging "github.com/op/go-logging"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ServiceAPI struct {
	AppData     *AppData
	AppChannels *AppChannels
}

func StartApi(localhttp *string, appdata *AppData, appchannels *AppChannels) {
	a := &ServiceAPI{AppData: appdata, AppChannels: appchannels}
	r := pat.New()
	r.Get("/health", http.HandlerFunc(a.GetHealth))
	r.Get("/loglevel/{loglevel}", http.HandlerFunc(a.SetLogLevel))
	http.Handle("/", r)
	log.Notice("HTTP Listening %s", *localhttp)
	err := http.ListenAndServe(*localhttp, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: ", err)
	}
}

func (api *ServiceAPI) GetHealth(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]string)
	res["ts"] = time.Now().Format("2006-01-02 15:04:01 PST")
	res["config"] = fmt.Sprintf("%+v", api.AppData.Config)
	who := fmt.Sprintf("%s", os.Args[0])
	res["who"] = who
	b, err := json.Marshal(res)
	if err != nil {
		log.Error("error: %s", err)
	}
	io.WriteString(w, string(b[:])+"\n")
}

func (api *ServiceAPI) SetLogLevel(w http.ResponseWriter, r *http.Request) {
	newlvl := strings.ToUpper(r.URL.Query().Get(":loglevel"))
	res := fmt.Sprintf("\nInvalid log level: %s \n Valid log levels are CRITICAL, ERROR,  WARNING, NOTICE, INFO, DEBUG\n", newlvl)

	if _, err := logging.LogLevel(newlvl); err == nil {
		api.AppChannels.Loglvl <- newlvl
		res = fmt.Sprintf("\nSetting log level to: %s \n", newlvl)
	}
	io.WriteString(w, res)
}
