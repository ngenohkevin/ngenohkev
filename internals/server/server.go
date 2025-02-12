package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ngenohkevin/ngenohkev/components/layout"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	logger     *slog.Logger
	port       int
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(logger *slog.Logger, port int) (*Server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	return &Server{
		logger: logger,
		port:   port,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Info("Starting server", slog.Int("port", s.port))

	//define the router
	router := http.NewServeMux()

	//handle static files
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//handle the route
	router.HandleFunc("/health", s.healthCheckHandler)
	router.HandleFunc("/", s.defaultHandler)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router,
	}

	//create a channel to listen for os signals
	//A buffered channel of size 1 is created so that the signals are not missed
	stopChan := make(chan os.Signal, 1)
	//Notify the channel when an interrupt or SIGTERM signal is received
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start the HTTP server in new goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start server", slog.String("error", err.Error()))
		}
	}()
	//Block the main goroutine until a signal is received
	<-stopChan
	s.logger.Info("Shutting down server")

	//Shutdown the server gracefully
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	return nil
}

// healthCheckHandler is a handler function that returns the health status of the server
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		s.logger.Error("Failed to write response", slog.String("error", err.Error()))
	}
}

func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	homeTemplate := layout.Home()
	err := homeTemplate.Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render template", slog.String("error", err.Error()))
	}

}
