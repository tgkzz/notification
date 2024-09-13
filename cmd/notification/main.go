package main

import (
	"context"
	"github.com/tgkzz/notification/internal/app"
	"github.com/tgkzz/notification/internal/config"
	"github.com/tgkzz/notification/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	ctx := context.Background()

	application, err := app.New(ctx, log)
	if err != nil {
		panic(err)
	}

	application.KafkaApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.KafkaApp.Stop()
}
