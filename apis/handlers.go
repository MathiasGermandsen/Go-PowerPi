package apis

import (
	"encoding/json"
	"net/http"

	"Power-Pi/database"
)

func GetPowerTable(w http.ResponseWriter, r *http.Request) {
	var rows []database.PowerTable

	if result := database.DB.Find(&rows); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func CreatePowerTable(w http.ResponseWriter, r *http.Request) {
	var row database.PowerTable

	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := database.DB.Create(&row); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(row)
}
