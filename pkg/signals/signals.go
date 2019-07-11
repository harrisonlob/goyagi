package signals

import (
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})

func Setup() <-chan struct{} {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c 
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
