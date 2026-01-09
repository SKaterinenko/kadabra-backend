package main

import (
	"context"
	"errors"
	"fmt"
	"kadabra/internal/core/app"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, router, cfg, cleanup := app.App()
	defer cleanup()

	fmt.Println("Config", cfg)
	fmt.Println("Server is listening on port", cfg.SERVER_PORT)

	server := &http.Server{
		Addr:    cfg.SERVER_PORT,
		Handler: router,
	}

	fmt.Println("PID:", os.Getpid())

	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Got shutdown signal")
	case err := <-serverErr:
		fmt.Printf("Server error: %v\n", err)
		return
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	fmt.Println("Gracefully shut down the server")
}
