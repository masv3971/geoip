package main

import (
	"context"
	"geoip/internal/apiv1"
	"geoip/internal/httpserver"
	"geoip/internal/maxmind"
	"geoip/internal/store"
	"geoip/internal/traveler"
	"geoip/pkg/configuration"
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type service interface {
	Close(ctx context.Context) error
}

func main() {
	ctx := context.Background()
	cfg := &model.Cfg{}

	wg := &sync.WaitGroup{}

	var (
		log      *logger.Logger
		mainLog  *logger.Logger
		services = make(map[string]service)
	)

	cfg, err := configuration.Parse(logger.NewSimple("Configuration"))
	if err != nil {
		panic(err)
	}

	mainLog = logger.New("main", cfg.Production)
	log = logger.New("geoip", cfg.Production)

	storage, err := store.New(ctx, cfg, log.New("storage"))
	if err != nil {
		panic(err)
	}

	maxmind, err := maxmind.New(ctx, cfg, storage, log.New("maxmind"))
	if err != nil {
		panic(err)
	}
	services["maxmind"] = maxmind

	traveler, err := traveler.New(ctx)
	if err != nil {
		panic(err)
	}

	apiv1, err := apiv1.New(ctx, cfg, maxmind, storage, traveler, log.New("apiv1"))
	if err != nil {
		panic(err)
	}

	httpserver, err := httpserver.New(ctx, cfg, apiv1, log.New("httpserver"))
	if err != nil {
		panic(err)
	}
	services["httpserver"] = httpserver

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan // Blocks here until interrupted

	mainLog.Info("HALTING SIGNAL!")

	for serviceName, service := range services {
		if err := service.Close(ctx); err != nil {
			mainLog.Warn(serviceName)
		}
	}

	wg.Wait() // Block here until are workers are done

	mainLog.Info("Stopped")
}
