package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/model"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type RegisterRequest struct {
	Username string `json:"username" example:"testUser123" validate:"required,min=1,max=32"`
	Password string `json:"password" example:"StrongP@ssw0rd!" validate:"required,min=8,max=64"`
}

type RegisterResponse struct {
	UserId   int    `json:"user_id" example:"10"`
	Username string `json:"username" example:"CoolUsername"`
	Password string `json:"password" example:"CoolPassword"`
}

type handler struct {
	db     store.UserStore
	logger logger.Logger
	token  *token.Service
}

func NewHandler(db store.UserStore, l logger.Logger, token *token.Service) *handler {
	return &handler{
		db:      db,
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
// @Success      201 {object} RegisterResponse "User successfully registered"
// @Failure      400 "Invalid request payload or user already exists"
// @Failure      500 "Internal server error"
// @Router       /api/v1/users [post]
func (h *handler) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var requestData RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &model.User{Username: requestData.Username, Password: requestData.Password}
	id, err := h.db.CreateUser(user)

	if err != nil {
		if errors.Is(err, database.ErrUserAlreadyExists) {
			http.Error(w, "User with this username already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	response := RegisterResponse{UserId: int(id), Username: user.Username, Password: user.Password} 

	token, _ := h.token.GenerateToken(fmt.Sprintf("%d", id))

	w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusCreated)
}
