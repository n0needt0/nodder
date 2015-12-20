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

package main

import (
	nodder "../libs"
	"flag"
	"fmt"
	statsd "github.com/cactus/go-statsd-client/statsd"
	logging "github.com/op/go-logging"
	"github.com/vharitonsky/iniflags"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

//some globals cant leave without them

var Appchannels = nodder.NewAppChannels()

//application config and runtime data goes here for dependency injection
var AppData = nodder.NewAppData()

var (
	Build     string //this is build version
	logFile   *os.File
	logFormat = logging.MustStringFormatter(nodder.LOGFMT)
	log       = logging.MustGetLogger("logfile")
	Statsd    statsd.Statter
)

func main() {
	//set defaults here so we can use dependency injection
	var (
		help          = flag.Bool("help", false, "Show these options")
		_             = flag.String("version", Build, "Build version")
		logfile       = flag.String("logfile", "STDOUT", "Logfile destination. STDOUT | STDERR or file path")
		loglevel      = flag.String("loglevel", "DEBUG", "Loglevel CRITICAL, ERROR, WARNING, NOTICE, INFO, DEBUG")
		localhttp     = flag.String("localhttp", "0.0.0.0:8888", "Local webserver api address")
		statsdserver  = flag.String("statsdserver", "127.0.0.1:8125", "Statsd address:port")
		statsdflushms = flag.Int("statsdflushms", 60000, "flush data to Statsd every ms")
		_             = flag.Int("workers", 10, "Number of workers")
	)

	iniflags.Parse() // use instead of flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//here we can do loggers
	var logback *logging.LogBackend

	AppData.Config["loglevel"] = *loglevel

	if loglvl, err := logging.LogLevel(*loglevel); err == nil {
		logging.SetLevel(loglvl, "")
	} else {
		logging.SetLevel(logging.DEBUG, "")
	}

	if *logfile == "STDOUT" {
		logback = logging.NewLogBackend(os.Stdout, "", 0)
	} else if *logfile == "STDERR" {
		logback = logging.NewLogBackend(os.Stderr, "", 0)
	} else {

		os.MkdirAll(filepath.Dir(*logfile), 0777)

		logFile, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Log file error: %s %s", *logfile, err)
		}

		defer func() {
			logFile.WriteString(fmt.Sprintf("closing %s", time.UnixDate))
			logFile.Close()
		}()

		logback = logging.NewLogBackend(logFile, "", 0)
	}

	logformatted := logging.NewBackendFormatter(logback, logFormat)

	logging.SetBackend(logformatted)

	//see what we have here
	log.Info("BUILD: %s\n", Build)

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name != "config" && f.Name != "dumpflags" {
			log.Info("%s = %s  # %s\n", f.Name, f.Value.String(), f.Usage)
			//publich all app data
			AppData.Config[f.Name] = f.Value.String()
		}
	})

	s := strings.Split(os.Args[0], "/")
	pidfile := fmt.Sprintf("/var/run/pid.%s", s[len(s)-1])

	err := ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
	if err != nil {
		log.Warning("Error creating pid file %s must be root!, %s", pidfile, err)
	}

	log.Info("Started %s proc: %d\n", os.Args[0], os.Getpid())

	signal.Notify(Appchannels.Kill, os.Interrupt)
	signal.Notify(Appchannels.Kill, syscall.SIGTERM)
	signal.Notify(Appchannels.Kill, syscall.SIGINT)
	signal.Notify(Appchannels.Kill, syscall.SIGQUIT)
	go func() {
		for {
			select {
			case <-Appchannels.Kill:
				if err = os.Remove(pidfile); err != nil {
					log.Error("%+v", err)
				}
				log.Fatalf("Interrupt %s", time.Now().String())
			case newlvl := <-Appchannels.Loglvl:
				loglvl := logging.DEBUG
				AppData.Config["loglevel"] = "DEBUG"

				if loglvl, err = logging.LogLevel(newlvl); err == nil {
					AppData.Config["loglevel"] = newlvl
				}

				//set to debug to output otherwise how do we know whats up
				logging.SetLevel(logging.DEBUG, "")
				res := fmt.Sprintf("\nSetting log level to: %v\n", AppData.Config["loglevel"])
				log.Debug(res)

				//now set it to what ever was specified
				logging.SetLevel(loglvl, "")
			}
		}
	}()

	// create statsd client
	Statsd, err = statsd.NewBufferedClient(*statsdserver, "", time.Duration(*statsdflushms)*time.Millisecond, 0)
	// handle any errors
	if err != nil {
		log.Fatal(err)
	}
	defer Statsd.Close()

	//start systems

	wg := &sync.WaitGroup{}
	wg.Add(1)
	//first start http interface for self stats
	//put some data
	go nodder.StartApi(localhttp, AppData, Appchannels)
	wg.Wait()
}
