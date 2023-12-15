package core

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*gin.Engine
	httpServer *http.Server
}

func NewServer(cfg ServerConfig) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(Recovery())
	// TODO: replace with impl core.Recovery

	s := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	return &Server{r, s}
}

func DefaultServer() *Server {
	return NewServer(DefaultServerConfig())
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

type Middleware func(*Context)

func handlerFunc(handler Middleware) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		handler(&Context{gContext: c})
	}
}

func (s *Server) POST(relativePath string, middlewares ...Middleware) {
	gHandles := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		gHandles = append(gHandles, handlerFunc(middleware))
	}

	s.Engine.POST(relativePath, gHandles...)
}

func (s *Server) Use(middlewares ...Middleware) {
	gMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		gMiddlewares = append(gMiddlewares, handlerFunc(middleware))
	}

	s.RouterGroup.Use(gMiddlewares...)
}
