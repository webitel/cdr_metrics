package main

import (
	"github.com/webitel/cdr_metrics/app"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	interruptChan := make(chan os.Signal, 1)

	a, err := app.NewApp()
	if err != nil {
		time.Sleep(time.Second)
		panic(err.Error())
	}

	defer a.Shutdown()

	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}
