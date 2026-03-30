package apis

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(LoggingMiddleware)

	r.HandleFunc("/power-table", GetPowerTable).Methods("GET")
	r.HandleFunc("/power-table", CreatePowerTable).Methods("POST")

	return r
}
