package shutdown

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
)

func Graceful(f func() error, sigChan chan os.Signal) error {

	if f == nil {
		return errors.New("nil func reference")
	}

	signal.Notify(sigChan, os.Interrupt,
		syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)

	<-sigChan

	return f()
}
