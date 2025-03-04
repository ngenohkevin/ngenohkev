package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ngenohkevin/ngenohkev/components/layout"
	"github.com/ngenohkevin/ngenohkev/components/pages"
	"github.com/ngenohkevin/ngenohkev/internals/blog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	router.HandleFunc("/about", s.about)
	router.HandleFunc("/projects", s.projects)

	// Blog routes
	router.HandleFunc("/posts", s.postsListHandler)
	router.HandleFunc("/posts/", func(writer http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/posts/" {
			s.postsListHandler(writer, r)
			return
		}
		s.postHandler(writer, r)
	})

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
	err := layout.Layout(homeTemplate, "Kevin's Blog", "/").Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render template", slog.String("error", err.Error()))
	}

}

func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	//Extract slug from URL
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")

	//Load the post
	post, err := blog.LoadPost(slug)
	if err != nil {
		s.logger.Error("Failed to load post", slog.String("slug", slug), slog.String("error", err.Error()))
		http.NotFound(w, r)
		return
	}

	//Render post template
	postTemplate := pages.Post(post)
	err = layout.Layout(postTemplate, post.Title, "/posts/"+slug).Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render post template", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) postsListHandler(w http.ResponseWriter, r *http.Request) {
	// List all posts
	posts, err := blog.ListPosts()
	if err != nil {
		s.logger.Error("Failed to list posts", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Render posts list template
	postsTemplate := pages.Posts(posts)
	err = layout.Layout(postsTemplate, "Blog Posts", "/posts").Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render posts template", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) about(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	aboutTemplate := pages.About()
	err := layout.Layout(aboutTemplate, "About Kevin", "/about").Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render template", slog.String("error", err.Error()))
	}
}

func (s *Server) projects(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	projectsTemplate := pages.Projects()
	err := layout.Layout(projectsTemplate, "Projects", "/projects").Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render template", slog.String("error", err.Error()))
	}
}
