package apis

import (
	"net/http"

	"Power-Pi/auth"
	"Power-Pi/config"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)

	// Swagger UI (public, no auth required)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Power Table routes — JWT required on both; scope enforced per method.
	jwtMw := auth.JWTMiddleware(cfg.JWTSecret)

	r.Handle("/power-table",
		jwtMw(auth.RequireScope("power-table:read")(http.HandlerFunc(GetPowerTable))),
	).Methods(http.MethodGet)

	r.Handle("/power-table",
		jwtMw(auth.RequireScope("power-table:write")(http.HandlerFunc(CreatePowerTable))),
	).Methods(http.MethodPost)

	// Admin routes — protected by X-Admin-Key header.
	// Routes are registered most-specific first so /tokens/rotate is not swallowed by /tokens/{jti}.
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(adminKeyMiddleware(cfg.AdminAPIKey))

	admin.HandleFunc("/tokens/rotate", RotateServiceToken(cfg)).Methods(http.MethodPost)
	admin.HandleFunc("/tokens/revoked", ListRevokedTokens(cfg)).Methods(http.MethodGet)
	admin.HandleFunc("/tokens", IssueServiceToken(cfg)).Methods(http.MethodPost)
	admin.HandleFunc("/tokens/{jti}", RevokeToken(cfg)).Methods(http.MethodDelete)

	return r
}
