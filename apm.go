package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func apmServer(cfg Config, shutdownCh <-chan struct{}, logger *log.Logger) {
	srv := &http.Server{
		Addr: cfg.APMAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n, _ := io.Copy(ioutil.Discard, r.Body)
			if cfg.Verbose {
				logger.Printf("method=%q uri=%q bodysize=%d", r.Method, r.RequestURI, n)
			}
		}),
	}

	// exitCh will be closed when it is safe to exit, after graceful shutdown
	exitCh := make(chan struct{})

	go func() {
		<-shutdownCh
		logger.Printf("shutting down ...")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Printf("shutdown error: %s", err)
		}

		close(exitCh)
	}()

	logger.Printf("listening on %s ...", cfg.APMAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen error: %s", err)
	}

	<-exitCh
}
