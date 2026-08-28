package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	core_http_middleware "github.com/loundxr/expense-tracker/internal/core/transport/http/middleware"
)

type HTTPServer struct {
	router     *chi.Mux
	cfg        Config
	logger     *slog.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	cfg Config,
	l *slog.Logger,
	mdw ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		router:     chi.NewRouter(),
		cfg:        cfg,
		logger:     l,
		middleware: mdw,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, r := range routers {
		prefix := "/api/" + string(r.version)
		s.router.Mount(prefix, r.Router())
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	const op = "core_http_server.Run"

	for _, m := range s.middleware {
		s.router.Use(m)
	}

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.router,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting http server", slog.String("addr", s.cfg.Addr))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("%s: listen and serve: %w", op, err)
	case <-ctx.Done():
		s.logger.Info("shutting down http server", slog.String("addr", s.cfg.Addr))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("%s: server shutdown: %v", op, err)
		}
		s.logger.Info("HTTP server gracefully stopped")
		return nil
	}
}
