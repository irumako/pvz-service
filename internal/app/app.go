package app

import (
	"fmt"
	"os"
	"os/signal"
	"pvz-service/config"
	v1 "pvz-service/internal/controller/http"
	"pvz-service/pkg/httpserver"
	"pvz-service/pkg/logger"
	"pvz-service/pkg/postgres"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// HTTP Server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Duration(cfg.HTTP.ReadTimeout)*time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.HTTP.WriteTimeout)*time.Second),
		httpserver.ShutdownTimeout(time.Duration(cfg.App.ShutdownTimeout)*time.Second),
	)

	uc := v1.Usecases{}
	v1.SetRouters(httpServer.Router, l, uc)

	httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
