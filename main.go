package main

import (
	"context"
	"errors"
	"fmt"
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

const LOGO = `
	
██╗  ██╗██╗       ███╗   ██╗ █████╗ ██╗  ██╗██╗██╗
██║  ██║██║       ████╗  ██║██╔══██╗██║ ██╔╝██║██║
███████║██║       ██╔██╗ ██║███████║█████╔╝ ██║██║
██╔══██║██║       ██║╚██╗██║██╔══██║██╔═██╗ ██║╚═╝
██║  ██║██║▄█╗    ██║ ╚████║██║  ██║██║  ██╗██║██╗
╚═╝  ╚═╝╚═╝╚═╝    ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝                          
`

func main() {
	var err error

	// Logger
	logger.MustInit("log/")
	logger.Info("(1/6) Initialized logger")

	// Config
	err = conf.Load("config.yml")
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		panic(err.Error())
	}
	logger.Info("(2/6) Configuration loaded")

	// Message types
	msgchain.Register()
	logger.Info("(3/6) Message types registered")

	// Events
	event.Register()
	logger.Info("(4/6) Events registered")

	// Router
	router := router.SetupRouter()
	logger.Info("(5/6) Router loaded")

	// Server
	srv := &http.Server{
		Addr:    conf.Get().Server.Address,
		Handler: router,
	}

	started := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(started)

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server startup failed", zap.Error(err))
			panic(err.Error())
		}

	}()

	<-started
	logger.Info("(6/6) Gin engine loaded")

	logger.Info(LOGO)
	logger.Info("Server is ready!")

	// Gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Print("\n") // Separate "^C" and next logs
	logger.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Server forced shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown gracefully")
}
