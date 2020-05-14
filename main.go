package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APMAddr            string        `default:":8216" split_words:"true"`
	APMShutdownTimeout time.Duration `default:"500ms" split_words:"true"`
	APMReadTimeout     time.Duration `default:"500ms" split_words:"true"`
	APMWriteTimeout    time.Duration `default:"250ms" split_words:"true"`
	StatsdAddr         string        `default:":8215" split_words:"true"`
	StatsdReadTimeout  time.Duration `default:"250ms" split_words:"true"`
	Verbose            bool
}

func main() {
	var cfg Config
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("invalid configuration: %s", err)
	}

	var (
		logPrefix = "[datadog-agent-stub]"
		logFlags  = log.LstdFlags
		apmLog    = log.New(os.Stdout, logPrefix+"[apm] ", logFlags)
		statsdLog = log.New(os.Stdout, logPrefix+"[statsd] ", logFlags)
	)

	// sigCh triggers graceful shutdown on SIGINT or SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// shutdownCh will be closed when it is time for the listeners to shut down
	shutdownCh := make(chan struct{})

	// wg waits for goroutines to exit
	var wg sync.WaitGroup

	// start apm listener
	wg.Add(1)
	go func() {
		defer wg.Done()
		apmServer(cfg, shutdownCh, apmLog)
	}()

	// start statsd listener
	wg.Add(1)
	go func() {
		defer wg.Done()
		statsdServer(cfg, shutdownCh, statsdLog)
	}()

	// wait for term signal
	<-sigCh

	// tell listeners to shut down
	close(shutdownCh)

	// wait for listeners to exit
	wg.Wait()
}
