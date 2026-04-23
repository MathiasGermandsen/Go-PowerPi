// @title           Power-Pi API
// @version         1.0
// @description     Service-to-service API for Power Table management. All /power-table endpoints require a JWT Bearer token. Admin endpoints require an X-Admin-Key header.
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT service token. Format: "Bearer {token}". Obtain a token via POST /admin/tokens.

package main

import (
	"fmt"
	"net/http"
	"strings"

	_ "Power-Pi/docs"

	"Power-Pi/apis"
	"Power-Pi/config"
	"Power-Pi/database"
	"Power-Pi/logger"

	"github.com/rs/cors"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel)

	database.Connect(cfg)

	router := apis.NewRouter(cfg)

	origins := strings.Split(cfg.CORSOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Admin-Key"},
		AllowCredentials: false,
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Log.Info().Str("addr", addr).Msg("Server listening")

	if err := http.ListenAndServe(addr, c.Handler(router)); err != nil {
		logger.Log.Fatal().Err(err).Msg("server error")
	}
}
