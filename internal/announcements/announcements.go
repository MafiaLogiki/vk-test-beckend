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
	mux.HandleFunc("GET /api/v1/announcements", auth.OptionalAuthMiddleware(h.token, h.getAnnouncements))
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


// GetAnnouncements gets a paginated list of announcements
// @Summary      Get announcements list
// @Description  Get a paginated list of announcements. This endpoint is public.
// @Tags         Announcements
// @Produce      json
// @Param        page   query      int  true  "Page number for pagination (starts from 1)"
// @Param        limit  query      int  true  "Number of items per page"
// @Success      200    {array}    model.Announcement
// @Failure      400    {string}   string "Invalid page or limit parameter"
// @Failure      500    {string}   string "Internal server error"
// @Router       /api/v1/announcements [get]
func (h *handler) getAnnouncements(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		http.Error(w, "Choose a page to display", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageString)

	if err != nil || page < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		http.Error(w, "Choose an items count on page limit", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitString)
	
	if err != nil || limit < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	announcements, err := h.db.GetAnnouncementsByPage(int32(page), int32(limit))
	h.logger.Debug(announcements)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(announcements); err != nil {
		h.logger.Info(err)
	}
}
