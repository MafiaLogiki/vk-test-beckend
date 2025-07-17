package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/middleware"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type handler struct {
	db     store.UserStore
	logger logger.Logger
	token  *token.Service
}

type AuthRequest struct {
	Username string `json:"username" example:"testUser123" validate:"required,min=1,max=32"`
	Password string `json:"password" example:"StrongP@ssw0rd!" validate:"required,min=8,max=64"`
}

func NewHandler(db store.UserStore, l logger.Logger, token *token.Service) *handler {
	return &handler{
		db:     db,
		logger: l,
		token:  token,
	}
}

func (h *handler) RegisterService(mux *http.ServeMux) {
	loggerMiddleware := middleware.LoggingMiddleware(h.logger)

	mux.Handle("POST /api/v1/auth", loggerMiddleware(http.HandlerFunc(h.authHandler)))
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

	// --- Валидация с помощью go-playground/validator ---
	if err := validate.Struct(userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// --- Конец валидации ---

	user, err := h.db.GetUserByCredentials(userData.Username, userData.Password)
	if err != nil {
		if errors.Is(err, store.ErrInvalidUsernameOrPassword) {
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



