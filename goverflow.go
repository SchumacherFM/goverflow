package goverflow

import (
	"github.com/SchumacherFM/goverflow/poster"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type goverflowApp struct {
	tickerSeconds  time.Duration
	ticker         *time.Ticker
	logger         *log.Logger
	configFileName *string
}

func NewGoverflowApp() *goverflowApp {
	return &goverflowApp{}
}

func (a *goverflowApp) SetInterval(interval *int) {
	secs := time.Second * time.Duration(*interval)
	a.tickerSeconds = secs
	a.ticker = time.NewTicker(secs)

}

func (a *goverflowApp) SetLogFile(logFile *string) {
	if "" != *logFile && nil != logFile {
		logFilePointer, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		a.logger = log.New(logFilePointer, "[GF] ", log.LstdFlags)
	} else {
		a.logger = log.New(os.Stderr, "[GF] ", log.LstdFlags)
	}

}

func (a *goverflowApp) SetConfigFileName(f *string) {
	a.configFileName = f
}

func (a *goverflowApp) GetLogger() *log.Logger {
	return a.logger
}

// Goverflow is the main method
func (a *goverflowApp) Goverflow() {
	a.catchSysCall()

	thePoster := poster.NewPoster(a.configFileName)
	thePoster.SetLogger(a.GetLogger())
	for _ = range a.ticker.C {
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
			a.logger.Printf("Received signal: %s\n", sig.String())
			a.ticker.Stop()
			a.logger.Println("Ticker stopped and good bye!")
			os.Exit(0)
		}
	}()
}
