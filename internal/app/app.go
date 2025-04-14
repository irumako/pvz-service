package app

import (
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"os"
	"os/signal"
	"pvz-service/config"
	v1 "pvz-service/internal/controller/http"
	repo "pvz-service/internal/repo/postgres"
	"pvz-service/internal/usecase/product"
	"pvz-service/internal/usecase/pvz"
	"pvz-service/internal/usecase/reception"
	"pvz-service/pkg/httpserver"
	"pvz-service/pkg/jwt"
	"pvz-service/pkg/logger"
	"pvz-service/pkg/postgres"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	jwtUtils := jwt.New(cfg.Jwt.Secret, time.Duration(cfg.Expiration)*time.Hour)

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

	// transaction manager
	trManager := manager.Must(trmpgx.NewDefaultFactory(pg.Pool))

	// Repository
	pvzRepo := repo.NewPvzRepo(pg, trmpgx.DefaultCtxGetter)
	receptionRepo := repo.NewReceptionRepo(pg, trmpgx.DefaultCtxGetter)
	productRepo := repo.NewProductRepo(pg, trmpgx.DefaultCtxGetter)

	//Use-cases
	uc := v1.Usecases{
		PvzUC:       pvz.New(pvzRepo),
		ReceptionUC: reception.New(receptionRepo, trManager),
		ProductUC:   product.New(productRepo, receptionRepo, trManager),
	}
	v1.SetRouters(httpServer.Router, l, jwtUtils, uc)

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
