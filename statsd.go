package main

import (
	"log"
	"net"
	"strings"
	"time"
)

const statsdReadTimeout = 250 * time.Millisecond

func statsdServer(cfg Config, shutdownCh <-chan struct{}, logger *log.Logger) {
	packetConn, err := net.ListenPacket("udp", cfg.StatsdAddr)
	if err != nil {
		logger.Fatalf("listen error: %s", err)
	}
	logger.Printf("listening on %s ...", cfg.StatsdAddr)

	// The only way to interrupt the blocking ReadFrom call below is to close
	// the connection, so we wait for the shutdown signal in a goroutine and
	// asynchronously close it.
	go func() {
		<-shutdownCh
		logger.Printf("shutting down ...")
		packetConn.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, _, err := packetConn.ReadFrom(buf)
		// Checking the error string is gross, but it's the only way I could
		// figure out to check whether the error is a result of the connection
		// being closed asynchronously by the goroutine above, which is an
		// indication that we should exit.
		if err != nil && strings.Contains(err.Error(), "closed network connection") {
			return
		}
		if cfg.Verbose {
			logger.Printf("recv: %s", string(buf[:n]))
		}
	}
}
