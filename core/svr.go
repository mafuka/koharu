package core

import (
	"context"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	ginEngine  *gin.Engine
	httpServer *http.Server
}

type ServerConfig struct {
	Address string `yaml:"address"` // HTTP listener address
	Secret  string `yaml:"secret"`  // Authorization header
}

// NewServer creates a new Server.
func NewServer(cfg ServerConfig) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(Recovery())
	// TODO: replace with impl core.Recovery

	return &Server{
		r,
		&http.Server{Addr: cfg.Address, Handler: r},
	}
}

func (s *Server) Run() error {
	srv := s.httpServer

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			Log().Fatal("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	Log().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return xerrors.Errorf("server shutdown failed: %w", err)
	}

	Log().Info("Server gracefully stopped")
	return nil
}

func (s *Server) PProf() {
	pprof.Register(s.ginEngine)
}

type Middleware func(*Context)

func ginMiddleware(m Middleware) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		m(&Context{context: c})
	}
}

func (s *Server) POST(relativePath string, middlewares ...Middleware) {
	ginHandles := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		ginHandles = append(ginHandles, ginMiddleware(middleware))
	}

	s.ginEngine.POST(relativePath, ginHandles...)
}

func (s *Server) Use(middlewares ...Middleware) {
	ginMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		ginMiddlewares = append(ginMiddlewares, ginMiddleware(middleware))
	}

	s.ginEngine.RouterGroup.Use(ginMiddlewares...)
}
