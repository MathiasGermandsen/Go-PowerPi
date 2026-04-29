package apis

import (
	"encoding/json"
	"net/http"
	"time"

	"Power-Pi/auth"
	"Power-Pi/config"
	"Power-Pi/database"

	"github.com/gorilla/mux"
)

type IssueTokenRequest struct {
	Service  string   `json:"service"`
	Audience string   `json:"audience"`
	Scopes   []string `json:"scopes"`
	OldJTI string `json:"old_jti,omitempty"`
}

type TokenResponse struct {
	Token     string   `json:"token"`
	JTI       string   `json:"jti"`
	ExpiresAt string   `json:"expires_at"`
	Service   string   `json:"service"`
	Audience  string   `json:"audience"`
	Scopes    []string `json:"scopes"`
}

// @Summary      Issue a service token
// @Description  Issues a new HS256 JWT for a named service with a specific audience and scopes.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        X-Admin-Key  header    string            true  "Admin API key"
// @Param        body         body      IssueTokenRequest true  "Token request"
// @Success      201          {object}  TokenResponse
// @Failure      400          {string}  string  "service, audience, and scopes are required"
// @Failure      403          {string}  string  "Forbidden"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /admin/tokens [post]
func IssueServiceToken(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Service == "" || req.Audience == "" || len(req.Scopes) == 0 {
			http.Error(w, "service, audience, and scopes are required", http.StatusBadRequest)
			return
		}

		tokenStr, jti, err := auth.GenerateToken(cfg.JWTSecret, req.Service, req.Audience, req.Scopes, cfg.HoursToExp)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().UTC().Add(time.Duration(cfg.HoursToExp) * time.Hour)

		scopesJSON, _ := json.Marshal(req.Scopes)
		reg := database.ServiceRegistration{
			ServiceName:   req.Service,
			Audience:      req.Audience,
			Scopes:        string(scopesJSON),
			JTI:           jti,
			ExpiresAt:     expiresAt,
			LastRotatedAt: time.Now().UTC(),
		}
		database.DB.Create(&reg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TokenResponse{
			Token:     tokenStr,
			JTI:       jti,
			ExpiresAt: expiresAt.Format(time.RFC3339),
			Service:   req.Service,
			Audience:  req.Audience,
			Scopes:    req.Scopes,
		})
	}
}

// @Summary      Revoke a service token
// @Description  Permanently revokes the token identified by its JTI. The token cannot be used after this call even if it has not yet expired.
// @Tags         admin
// @Param        X-Admin-Key  header  string  true  "Admin API key"
// @Param        jti          path    string  true  "JWT ID (jti claim) of the token to revoke"
// @Success      204
// @Failure      403          {string}  string  "Forbidden"
// @Failure      404          {string}  string  "Token not found"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /admin/tokens/{jti} [delete]
func RevokeToken(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jti := mux.Vars(r)["jti"]
		if jti == "" {
			http.Error(w, "jti is required", http.StatusBadRequest)
			return
		}

		var reg database.ServiceRegistration
		if result := database.DB.Where("jti = ?", jti).First(&reg); result.Error != nil {
			http.Error(w, "token not found", http.StatusNotFound)
			return
		}

		revoked := database.RevokedToken{
			JTI:       jti,
			ExpiresAt: reg.ExpiresAt,
		}
		if result := database.DB.Create(&revoked); result.Error != nil {
			http.Error(w, "failed to revoke token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// @Summary      Rotate a service token
// @Description  Revokes the token identified by old_jti and issues a new token with the provided service, audience, and scopes.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        X-Admin-Key  header    string            true  "Admin API key"
// @Param        body         body      IssueTokenRequest true  "Rotation request (old_jti is required)"
// @Success      200          {object}  TokenResponse
// @Failure      400          {string}  string  "service, audience, scopes, and old_jti are required"
// @Failure      403          {string}  string  "Forbidden"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /admin/tokens/rotate [post]
func RotateServiceToken(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Service == "" || req.Audience == "" || len(req.Scopes) == 0 || req.OldJTI == "" {
			http.Error(w, "service, audience, scopes, and old_jti are required", http.StatusBadRequest)
			return
		}

		var oldReg database.ServiceRegistration
		if result := database.DB.Where("jti = ?", req.OldJTI).First(&oldReg); result.Error == nil {
			database.DB.Create(&database.RevokedToken{
				JTI:       req.OldJTI,
				ExpiresAt: oldReg.ExpiresAt,
			})
		}

		tokenStr, jti, err := auth.GenerateToken(cfg.JWTSecret, req.Service, req.Audience, req.Scopes, cfg.HoursToExp)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().UTC().Add(time.Duration(cfg.HoursToExp) * time.Hour)

		scopesJSON, _ := json.Marshal(req.Scopes)
		database.DB.Create(&database.ServiceRegistration{
			ServiceName:   req.Service,
			Audience:      req.Audience,
			Scopes:        string(scopesJSON),
			JTI:           jti,
			ExpiresAt:     expiresAt,
			LastRotatedAt: time.Now().UTC(),
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{
			Token:     tokenStr,
			JTI:       jti,
			ExpiresAt: expiresAt.Format(time.RFC3339),
			Service:   req.Service,
			Audience:  req.Audience,
			Scopes:    req.Scopes,
		})
	}
}

// @Summary      List revoked tokens
// @Description  Returns the full revocation list. Useful for auditing and verifying that a token has been properly revoked.
// @Tags         admin
// @Produce      json
// @Param        X-Admin-Key  header  string  true  "Admin API key"
// @Success      200          {array}   RevokedTokenDoc
// @Failure      403          {string}  string  "Forbidden"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /admin/tokens/revoked [get]
func ListRevokedTokens(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokens []database.RevokedToken
		if result := database.DB.Find(&tokens); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	}
}
