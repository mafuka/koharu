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
	"github.com/kwaain/nakisama/lib/logger"
	"github.com/kwaain/nakisama/router"
	"go.uber.org/zap"
)

func main() {
	var err error

	err = logger.NewLogger("log/")
	if err != nil {
		panic("初始化日志记录器失败: " + err.Error())
	}

	err = conf.Load("config.yml")
	if err != nil {
		logger.Fatal("加载配置文件失败", zap.Error(err))
	}
	logger.Info("配置文件加载成功")

	router := router.SetupRouter()

	srv := &http.Server{
		Addr:    conf.Get().Server.Address,
		Handler: router,
	}

	// Gracefully Shutdown
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("监听服务错误", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务器正在关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("服务器被迫关闭", zap.Error(err))
	}

	logger.Info("服务器优雅地退出")

	// Close the logger
	err = logger.Close()
	if err != nil {
		logger.Error("关闭日志记录器失败", zap.Error(err))
	}
}
