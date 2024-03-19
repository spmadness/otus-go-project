package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spmadness/otus-go-project/internal/app"
	"github.com/spmadness/otus-go-project/internal/scraper"
	"github.com/spmadness/otus-go-project/internal/server"
)

var (
	configFile string
	port       int
)

func init() {
	flag.StringVar(&configFile, "config", "./config/config_daemon.toml", "Path to configuration file")
	flag.IntVar(&port, "port", 50051, "grpc server port")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)

	scrapers := scraper.NewCollection(
		config.Metrics.LoadAverageSystem,
		config.Metrics.LoadAverageCPU,
		config.Metrics.LoadDisks,
	)

	if len(scrapers) == 0 {
		fmt.Println("no scrapers are specified in config")
		os.Exit(1)
	}

	fmt.Println("starting monitoring daemon...")

	service := app.New(scrapers)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go service.GatherMetrics(ctx)

	serverGRPC := server.NewServer(service, port)

	go func() {
		<-ctx.Done()

		serverGRPC.Stop()
	}()
	serverGRPC.Start()
}
