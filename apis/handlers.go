package apis

import (
	"encoding/json"
	"net/http"

	"Power-Pi/database"

	"gorm.io/gorm/clause"
)

func GetPowerTable(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")

	var rows []database.PowerTable

	query := database.DB
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if result := query.Find(&rows); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func CreatePowerTable(w http.ResponseWriter, r *http.Request) {
	var req database.PowerTable

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "missing userId", http.StatusBadRequest)
		return
	}

	result := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"company", "price"}),
	}).Create(&req)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
