package apis

import (
	"net/http"

	"Power-Pi/auth"
	"Power-Pi/config"
	"Power-Pi/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	jwtMw := auth.JWTMiddleware(cfg.JWTSecret)

	r.Handle("/power-table",
		jwtMw(auth.RequireScope("power-table:read")(http.HandlerFunc(GetPowerTable))),
	).Methods(http.MethodGet)

	r.Handle("/power-table",
		jwtMw(auth.RequireScope("power-table:write")(http.HandlerFunc(CreatePowerTable))),
	).Methods(http.MethodPost)

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminKeyMiddleware(cfg.AdminAPIKey))

	admin.HandleFunc("/tokens/rotate", RotateServiceToken(cfg)).Methods(http.MethodPost)
	admin.HandleFunc("/tokens/revoked", ListRevokedTokens(cfg)).Methods(http.MethodGet)
	admin.HandleFunc("/tokens", IssueServiceToken(cfg)).Methods(http.MethodPost)
	admin.HandleFunc("/tokens/{jti}", RevokeToken(cfg)).Methods(http.MethodDelete)

	return r
}
