package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Server wraps http.Server with graceful shutdown support.
type Server struct {
	srv *http.Server
}

// NewServer creates a Server listening on the given port with safe timeouts.
func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// ListenAndServe starts the server and blocks until SIGTERM/SIGINT is received,
// then drains connections with a 30-second timeout.
func (s *Server) ListenAndServe() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("gateway listening", slog.String("addr", s.srv.Addr))
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received, draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	slog.Info("gateway stopped cleanly")
	return nil
}
