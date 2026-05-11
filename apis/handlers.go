package apis

import (
	"encoding/json"
	"net/http"

	"Power-Pi/database"

	"gorm.io/gorm/clause"
)

// @Summary      List power table entries
// @Description  Returns all power table entries. Optionally filter by userId query parameter.
// @Tags         power-table
// @Produce      json
// @Param        userId  query     string  false  "Filter by user ID"
// @Success      200     {array}   PowerTableDoc
// @Failure      500     {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /power-table [get]
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

// @Summary      Create or update a power table entry
// @Description  Upserts a power table entry by userId. Updates company, price, selectionMode, and numberOfHours on conflict.
// @Tags         power-table
// @Accept       json
// @Produce      json
// @Param        body  body      PowerTableDoc  true  "Power table entry"
// @Success      201   {object}  PowerTableDoc
// @Failure      400   {string}  string  "Bad request"
// @Failure      500   {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /power-table [post]
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

	updateCols := []string{"company", "selection_mode"}
	switch req.SelectionMode {
	case "price-based":
		updateCols = append(updateCols, "price")
	case "hour-based":
		updateCols = append(updateCols, "number_of_hours")
	default:
		updateCols = append(updateCols, "price", "number_of_hours")
	}

	result := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns(updateCols),
	}).Create(&req)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// @Summary      Get charging status
// @Description  Returns the charging status for a given userId.
// @Tags         charging
// @Produce      json
// @Param        userId  query     string  true  "User ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {string}  string  "Bad request"
// @Failure      404     {string}  string  "Not found"
// @Security     BearerAuth
// @Router       /charging [get]
func GetCharging(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "missing userId query parameter", http.StatusBadRequest)
		return
	}

	var row database.PowerTable
	if result := database.DB.Where("user_id = ?", userID).First(&row); result.Error != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":   row.UserID,
		"charging": row.Charging,
	})
}

// @Summary      Set charging status
// @Description  Updates the charging status for a given userId.
// @Tags         charging
// @Accept       json
// @Produce      json
// @Param        userId  query     string                  true  "User ID"
// @Param        body    body      object{charging=bool}   true  "Charging status"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {string}  string  "Bad request"
// @Failure      404     {string}  string  "Not found"
// @Failure      500     {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /charging [post]
func SetCharging(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "missing userId query parameter", http.StatusBadRequest)
		return
	}

	var req struct {
		Charging bool `json:"charging"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var row database.PowerTable
	if result := database.DB.Where("user_id = ?", userID).First(&row); result.Error != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if result := database.DB.Model(&row).Update("charging", req.Charging); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":   userID,
		"charging": req.Charging,
	})
}
