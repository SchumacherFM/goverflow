package goverflow

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppConfig struct {
	tickerSeconds time.Duration
	ticker        *time.Ticker
	logger        *log.Logger
}

func NewAppConfig() *AppConfig {
	inputDuration := flag.Int("seconds", 2, "Sleep duration in Seoncds")
	logFile := flag.String("logFile", "", "Log to file or if empty to os.Stderr")
	flag.Parse()
	secs := time.Second * time.Duration(*inputDuration)
	ac := &AppConfig{
		tickerSeconds: secs,
		ticker:        time.NewTicker(secs),
	}

	if "" != *logFile {
		logFilePointer, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		ac.logger = log.New(logFilePointer, "[GF] ", log.LstdFlags)
	} else {
		ac.logger = log.New(os.Stderr, "[GF] ", log.LstdFlags)
	}

	return ac
}

func (a *AppConfig) Goverflow() {
	a.catchSysCall()

	for t := range a.ticker.C {
		a.logger.Println("check SO ... Tick at", t)
	}
}

// catchSysCall ends the program correctly when receiving a sys call
func (a *AppConfig) catchSysCall() {
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
