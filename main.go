package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/api/router"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/logger"
	"github.com/8tomat8/Qm9yeXMtSHVsaWk-/store"
	"github.com/go-chi/chi"
)

var (
	port            = flag.Uint("port", 8080, "Port for API listener")
	host            = flag.String("host", "0.0.0.0", "Host for API listener")
	shutdownTimeout = flag.Duration("s", 30, "Shutdown timeout")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	storage := store.NewStorage()

	addr := fmt.Sprintf("%s:%s", *host, strconv.Itoa(int(*port)))
	srv := http.Server{
		Addr:    addr,
		Handler: chi.ServerBaseContext(ctx, router.NewRouter(storage)),
	}

	done := make(chan struct{})
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop
		cancel()

		// Creating new context with cancel for Shutdown only
		ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Fatalf("could not shutdown gracefully: %v", err)
		}

		close(done)
	}()

	logger.Log.Info(fmt.Sprintf("Starting listener on %s", addr))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		close(done)
		logger.Log.Error(err)
	}

	select {
	case <-done:
		logger.Log.Info("Application stopped gracefully")
	case <-stop:
		logger.Log.Warn("Received second SIGINT. Stopping immediately")
	}
}
