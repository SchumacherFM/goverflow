/*
	Copyright (C) 2014  Cyrill AT Schumacher dot fm

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

    Contribute @ https://github.com/SchumacherFM/goverflow
*/

package goverflow

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SchumacherFM/goverflow/poster"
	log "github.com/segmentio/go-log"
)

type goverflowApp struct {
	tickerSeconds  time.Duration
	ticker         *time.Ticker
	logger         *log.Logger
	configFileName string
}

func NewGoverflowApp() *goverflowApp {
	return &goverflowApp{}
}

func (a *goverflowApp) SetInterval(interval int) {
	secs := time.Second * time.Duration(interval)
	a.tickerSeconds = secs
	a.ticker = time.NewTicker(secs)

}

func (a *goverflowApp) SetLogFile(logFile string, logLevel int) {
	var logMap = map[int]log.Level{
		0: log.DEBUG,
		1: log.INFO,
		2: log.NOTICE,
		3: log.WARNING,
		4: log.ERROR,
		5: log.CRITICAL,
		6: log.ALERT,
		7: log.EMERGENCY,
	}

	validLogLevel, isSetLevel := logMap[logLevel]
	if false == isSetLevel {
		validLogLevel = log.DEBUG
	}

	if "" != logFile {
		logFilePointer, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		a.logger = log.New(logFilePointer, validLogLevel, "[GF] ")
	} else {
		a.logger = log.New(os.Stderr, validLogLevel, "[GF] ")
	}

}

func (a *goverflowApp) SetConfigFileName(f string) {
	a.configFileName = f
}

func (a *goverflowApp) GetLogger() *log.Logger {
	return a.logger
}

// Goverflow is the main method
func (a *goverflowApp) Goverflow() {
	a.catchSysCall()

	thePoster := poster.NewPoster(&a.configFileName, a.GetLogger())
	go thePoster.RoutinePoster()
	for range a.ticker.C {
		go thePoster.RoutinePoster()
	}
}

// catchSysCall ends the program correctly when receiving a sys call
func (a *goverflowApp) catchSysCall() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(
		signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for sig := range signalChannel {
			a.logger.Debug("Received signal: %s\n", sig.String())
			a.ticker.Stop()
			// here we can now save the DB to a file ... if we would use a non memory version.
			a.logger.Debug("Ticker stopped and good bye!")
			os.Exit(0)
		}
	}()
}
