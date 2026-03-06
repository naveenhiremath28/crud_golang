package main

import (
	"practise/go_fiber/internal/containers"

	"go.uber.org/zap"
)

func main() {
	// Bootstrap logger for pre-container errors
	bootstrapLogger, _ := zap.NewDevelopment()
	sugar := bootstrapLogger.Sugar()
	defer sugar.Sync()

	container, err := containers.NewContainer()
	if err != nil {
		sugar.Fatalw("Failed to initialize dependency container", "error", err)
	}

	// Start server (blocking)
	if err := container.Invoke(containers.StartServer); err != nil {
		sugar.Fatalw("Failed to start server", "error", err)
	}
}
