package register

import (
	"encoding/json"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/token"
	"net/http"
)

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
	var userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&userData)

	err := database.CreateNewUser(userData.Username, userData.Password)

	if err != nil {

	}
}
