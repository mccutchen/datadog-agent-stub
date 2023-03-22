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
			statusCode := http.StatusOK

			// Respond with a 404 to agent health checks, to indicate that
			// remote configuration is not enabled or supported.
			//
			// See, e.g., the relevant code from the dd-trace-py client impl:
			// https://github.com/DataDog/dd-trace-py/blob/v1.9.4/ddtrace/internal/agent.py#L151-L169
			if r.URL.Path == "/info" {
				statusCode = http.StatusNotFound
			}

			// Otherwise, assume we received trace data via PUT /v0.5/traces or
			// POST /v0.4/traces
			n, _ := io.Copy(ioutil.Discard, r.Body)

			w.WriteHeader(statusCode)
			if cfg.Verbose {
				logger.Printf("status=%d method=%q uri=%q bodysize=%d", statusCode, r.Method, r.RequestURI, n)
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
