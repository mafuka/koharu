package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kwaain/nakisama/lib/conf"
	"github.com/kwaain/nakisama/lib/event"
	"github.com/kwaain/nakisama/lib/logger"
	"github.com/kwaain/nakisama/lib/msgchain"
	"github.com/kwaain/nakisama/router"
	"go.uber.org/zap"
)

func main() {
	var err error

	// Logger
	logger.MustInit("log/")
	logger.Info("(1/5) Initialized logger")

	// Config
	err = conf.Load("config.yml")
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		panic(err.Error())
	}
	logger.Info("(2/5) Configuration loaded")

	// Router
	router := router.SetupRouter()
	logger.Info("(3/5) Router loaded")

	// Message types
	msgchain.Register()
	logger.Info("(4/5) Message types registered")

	// Events
	event.Register()
	logger.Info("(5/5) Events registered")

	// Server
	srv := &http.Server{
		Addr:    conf.Get().Server.Address,
		Handler: router,
	}
	logger.Info(`
	
██╗  ██╗██╗       ███╗   ██╗ █████╗ ██╗  ██╗██╗██╗
██║  ██║██║       ████╗  ██║██╔══██╗██║ ██╔╝██║██║
███████║██║       ██╔██╗ ██║███████║█████╔╝ ██║██║
██╔══██║██║       ██║╚██╗██║██╔══██║██╔═██╗ ██║╚═╝
██║  ██║██║▄█╗    ██║ ╚████║██║  ██║██║  ██╗██║██╗
╚═╝  ╚═╝╚═╝╚═╝    ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝                          
`)
	logger.Info("Server is ready!")

	// Gracefully shutdown
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server startup failed", zap.Error(err))
			panic(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Server forced shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown gracefully")

	// Close the logger
	err = logger.Close()
	if err != nil {
		panic("Failed to close logger" + err.Error())
	}
}
