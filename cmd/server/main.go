package main

import (
	"github.com/Enthreeka/auth-lab3/internal/config"
	"github.com/Enthreeka/auth-lab3/internal/server"
	"github.com/Enthreeka/auth-lab3/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to load config: %v", err)
	}

	if err := server.Run(log, cfg); err != nil {
		log.Fatal("failed to run server: %v", err)
	}
}
