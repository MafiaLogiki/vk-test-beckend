package register

import (
	"encoding/json"
	"errors"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/token"
	"net/http"
	"fmt"
)

type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"12345"`
}

type handler struct {
	logger logger.Logger
	token  *token.Service
}

func NewHandler(l logger.Logger, token *token.Service) *handler {
	return &handler{
		logger:  l,
		token:   token,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.registerNewUser)
}

// registerNewUser handles user registration.
// @Summary      Register a new user
// @Description  Registers a new user with a username and password.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User registration details"
// @Success      201 "User successfully registered"
// @Failure      400 "Invalid request payload or user already exists"
// @Failure      500 "Internal server error"
// @Router       /api/v1/users [post]
func (h *handler) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var userData RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := database.CreateNewUser(userData.Username, userData.Password)

	if err != nil {
		if errors.Is(err, database.ErrUserAlreadyExists) {
			http.Error(w, "User with this username already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, _ := h.token.GenerateToken(fmt.Sprintf("%d", id))

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully registered"))
}
