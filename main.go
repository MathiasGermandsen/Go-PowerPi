package main

import (
	"fmt"
	"net/http"

	"Power-Pi/apis"
	"Power-Pi/config"
	"Power-Pi/database"
	"Power-Pi/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel)

	database.Connect(cfg)

	router := apis.NewRouter()

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Log.Info().Str("addr", addr).Msg("Server listening")

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Log.Fatal().Err(err).Msg("server error")
	}
}
