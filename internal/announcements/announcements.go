package announcements

import (
	"encoding/json"
	"marketplace-service/internal/auth"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/model"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"
	"net/http"
	"strconv"
)



type handler struct {
	db     store.AnnouncementsStore
	logger logger.Logger
	token  *token.Service
}

type AnnouncementsPostRequest struct {
	Article  string
	Text     string
	ImageURL string
	Cost     int32
}

func NewHandler(db store.AnnouncementsStore, logger logger.Logger, token *token.Service) *handler {
	return &handler{
		db: db,
		logger: logger,
		token: token,
	}
}

func (h *handler) RegisterService(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/announcements", auth.AuthMiddleware(h.token, h.createAnnouncement))
}

// Announcement creation
// @Summary      Create an announcement
// @Description  Create an announcement for authorized users.
// @Tags         Announcements 
// @Accept       json
// @Produce      json
// @Param        request body AnnouncementsPostRequest true "Announcement details"
// @Success      201 "Announcement successfully created"
// @Failure      400 "Invalid request payload"
// @Failure      500 "Internal server error"
// @Router       /api/v1/announcements [post]
// @Security     Bearer
func (h *handler) createAnnouncement (w http.ResponseWriter, r *http.Request) {
	var apr AnnouncementsPostRequest

	err := json.NewDecoder(r.Body).Decode(&apr)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token := token.ExtractToken(r)
	userId, _ := h.token.ValidateToken(token)
	
	var an model.Announcement
	
	an.Article = apr.Article
	an.Text = apr.Text
	an.CostRubles = apr.Cost
	id, err := strconv.Atoi(userId)
	h.logger.Debug("id: ", id, " err: ", err)
	an.UserId = int64(id)

	err = h.db.CreateAnnouncement(&an)
	if err != nil {
		h.logger.Info(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Announcement successfully created"))
}
