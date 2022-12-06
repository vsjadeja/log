package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/vsjadeja/log"
)

func main() {
	ctx := context.Background()
	logger := log.NewDevelopmentLogger().Named(`main`)
	logger2 := logger.Named(`subordinate`)
	defer func() {
		logger.Info(ctx, `logger stopped`)
		_ = logger.Sync()

		logger2.Info(ctx, `logger2 stopped`)
		_ = logger2.Sync()
	}()

	sig := make(chan os.Signal, 1)
	pid := os.Getpid()
	wg := sync.WaitGroup{}

	signal.Notify(sig, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			switch <-sig {
			case syscall.SIGUSR1:
				logger.SetLevel(log.ErrorLevel)
			case syscall.SIGUSR2:
				logger.SetLevel(log.DebugLevel)
			}
		}
	}()
	fmt.Printf("Usage:\nkill -s USR1 %d -- set logging level to the ERROR.\n", pid)
	fmt.Printf("kill -s USR2 %d -- set logging level to the DEBUG.\n\n", pid)
	fmt.Println(`Press the any key to start ...`)
	_, _ = fmt.Scanln()

	logger.Infof("logger initialized with level: %s", logger.Level())
	logger2.Infof("logger2 initialized with level: %s", logger2.Level())

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			logger.Error(ctx, `error message`)
			logger2.Error(ctx, `error message`)
			time.Sleep(delay * 5)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			logger.Info(ctx, `info message`)
			logger2.Info(ctx, `info message`)
			time.Sleep(delay)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			logger.Debug(ctx, `debug message`)
			logger2.Debug(ctx, `debug message`)
			time.Sleep(delay)
		}
	}()

	wg.Wait()
}

const (
	delay = 1000 * time.Millisecond
)
