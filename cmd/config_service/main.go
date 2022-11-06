package main

import (
	"distributedConfig/config"
	"distributedConfig/internal/app"
	"log"
)

func main() {
	log.Println("Starting server")
	configPath := "."
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return
	}
	log.Printf("Successfully parsed config")
	app.Run(cfg)
}
