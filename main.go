package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"todo/internal/server"
)

type Config struct {
	Port int
}

func loadConfig(env string) *Config {
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	return &Config{
		Port: port,
	}
}

func main() {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	config := loadConfig(env)

	slog.Info("Starting server on port", "port", config.Port, "env", env)

	server := server.NewServer(config.Port)
	err := server.Start()
	if err != nil {
		slog.Error("Error starting server:", "error", err)
	}
}
