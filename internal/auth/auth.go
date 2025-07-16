package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"
	"net/http"
	"strings"
)

type handler struct {
	db     store.UserStore
	logger logger.Logger
	token  *token.Service
}

type AuthRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"12345"`
}

func NewHandler(db store.UserStore, l logger.Logger, token *token.Service) *handler {
	return &handler{
		db:     db,
		logger: l,
		token:  token,
	}
}

func (h *handler) RegisterService(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth", h.authHandler)
}



// Authorization of user
// @Summary      Auth user
// @Description  Auth user by username and password.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body AuthRequest true "User authorization details"
// @Success      201 "User successfully authorized"
// @Failure      400 "Invalid request payload or invalid username or invalid password"
// @Failure      500 "Internal server error"
// @Router       /api/v1/auth [post]
func (h *handler) authHandler(w http.ResponseWriter, r *http.Request) {
	var userData AuthRequest
	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	user, err := h.db.GetUserByCredentials(userData.Username, userData.Password)
	if err != nil {
		if errors.Is(err, database.ErrInvalidUsernameOrPassword) {
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, _ := h.token.GenerateToken(fmt.Sprintf("%d", user.ID))

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully authorized"))
}


func isAuthorized(tok *token.Service, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}
	_, err := tok.ValidateToken(strings.Split(token, " ")[1])

	if err != nil {
		return false
	}

	return true
}

func authMiddleware(tok *token.Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuthorized(tok, r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func optionalAuthMiddleware(tok *token.Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !isAuthorized(tok, r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

}
