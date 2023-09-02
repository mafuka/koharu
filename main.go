package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kwaain/nakisama/lib/conf"
	"github.com/kwaain/nakisama/router"
	"go.uber.org/zap"
)

func main() {
	err = conf.Load("config.yml")
	router := router.SetupRouter()

	srv := &http.Server{
		Addr:    conf.Get().Server.Address,
		Handler: router,
	}

	// Gracefully Shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
