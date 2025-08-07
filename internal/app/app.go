// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xiabin827/task-machinery/config"
	"github.com/xiabin827/task-machinery/internal/controller/http"
	"github.com/xiabin827/task-machinery/internal/repo/persistent"
	"github.com/xiabin827/task-machinery/internal/repo/webapi"
	"github.com/xiabin827/task-machinery/internal/usecase/translation"
	"github.com/xiabin827/task-machinery/pkg/httpserver"
	"github.com/xiabin827/task-machinery/pkg/logger"
	"github.com/xiabin827/task-machinery/pkg/machinery"
	"github.com/xiabin827/task-machinery/pkg/machinery/tasks/oceanengine"
	"github.com/xiabin827/task-machinery/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use-Case
	translationUseCase := translation.New(
		persistent.New(pg),
		webapi.New(),
	)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, translationUseCase, l)

	// Machinery
	machinery, err := machinery.NewMachinery(cfg, oceanengine.NewOceanEngineOpenSdk())
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - machinery.NewMachinery: %w", err))
	}

	machinery.StartWorker()

	// Start server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-machinery.Notify():
		l.Error(fmt.Errorf("app - Run - machinery.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
