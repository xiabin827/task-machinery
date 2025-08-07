package main

import (
	"log"

	"github.com/xiabin827/task-machinery/config"
	"github.com/xiabin827/task-machinery/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
