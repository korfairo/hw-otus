package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/korfairo/hw-otus/hw12_13_14_15_calendar/internal/logger"

	"github.com/pkg/errors"
)

type Server struct {
	srv *http.Server
	app Application
	log Logger
}

type Logger interface {
	WithField(key string, value interface{}) *logger.Logger
	WithError(err error) *logger.Logger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}

type Application interface{}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "root triggered")
}

func NewServer(host, port string, app Application, logger Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(handleRoot)))

	return &Server{
		srv: &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: mux,
		},
		log: logger,
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to stop server")
	}
	return nil
}
