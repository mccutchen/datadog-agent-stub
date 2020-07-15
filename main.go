package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config contains application configuration
type Config struct {
	APMAddr    string
	StatsdAddr string
	Verbose    bool
}

// loadConfig loads application configuration from the environment.
func loadConfig(getenv func(string) string) Config {
	cfg := Config{
		APMAddr:    ":8126",
		StatsdAddr: ":8125",
		Verbose:    false,
	}
	if apmAddr := getenv("APM_ADDR"); apmAddr != "" {
		cfg.APMAddr = apmAddr
	}
	if statsdAddr := getenv("STATSD_ADDR"); statsdAddr != "" {
		cfg.StatsdAddr = statsdAddr
	}
	if verbose := getenv("VERBOSE"); verbose != "" {
		cfg.Verbose = verbose == "true" || verbose == "t" || verbose == "1"
	}
	return cfg
}

func main() {
	cfg := loadConfig(os.Getenv)

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
