package announcements

import (
	"encoding/json"
	"marketplace-service/internal/auth"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/middleware"
	"marketplace-service/internal/model"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type handler struct {
	db     store.AnnouncementsStore
	logger logger.Logger
	token  *token.Service
}

type AnnouncementsPostRequest struct {
	Article  string `json:"article" example:"Продам старый диван" validate:"required,min=5,max=200"`
	Text     string `json:"text" example:"Продается диван б/у, в хорошем состоянии, самовывоз. Торг уместен." validate:"required,min=10,max=2000"`
	ImageURL string `json:"image_url" example:"http://example.com/images/sofa.jpg" validate:"omitempty,url,max=255"`
	Cost     int32  `json:"cost" example:"5000" validate:"required,min=0"`
}

type AnnouncementsPostResponse struct {
	Id            int64     `json:"id" example:"11"`
	UserId        int64     `json:"user_id" example:"3"`
	Article       string    `json:"title" example:"Продам машину"`
	Text          string    `json:"text" example:"Продается машина, 120000км пробег"`
	CostRubles    int32     `json:"price" example:"700000"`
	ImageAddress  string    `json:"image_url" example:"http://example.com/images/car"`
	Date          time.Time `json:"created_at" example:"2025-07-16T22:39:54.789179Z"`
}

type AnnouncementsGetResponse struct {
	OwnerUsername string    `json:"owner_username" example:"CoolUsername"`
	Article       string    `json:"title" example:"Продам машину"`
	Text          string    `json:"text" example:"Продам машину, 120000км пробег"`
	CostRubles    int32     `json:"price" example:"700000"`
	ImageAddress  string    `json:"image_url" example:"http://example.com/images/car"`
	IsOwner       *bool     `json:"is_owner,omitempty" example:"false"`
}

func NewHandler(db store.AnnouncementsStore, logger logger.Logger, token *token.Service) *handler {
	return &handler{
		db: db,
		logger: logger,
		token: token,
	}
}

func (h *handler) RegisterService(mux *http.ServeMux) {
	
	ValidationRules := []middleware.ValidationRule {
		{
			ParamName: "page",
			DefaultValue: "1",
			Validator: middleware.ValidatePositiveInt,
			ContextKey: middleware.PageKey,
		},
		{
			ParamName: "limit",
			DefaultValue: "10",
			Validator: middleware.ValidatePositiveInt,
			ContextKey: middleware.LimitKey,
		},
		{
			ParamName: "sort_by",
			DefaultValue: store.DefaultSorting,
			Validator: middleware.ValidateByMap(store.ValidSortColumns),
			ContextKey: middleware.SortKey,
		},
		{
			ParamName: "min_price",
			DefaultValue: "1",
			Validator: middleware.ValidatePositiveInt,
			ContextKey: middleware.MinPrice,
		},
		{
			ParamName: "max_price",
			DefaultValue: "2147483647",
			Validator: middleware.ValidatePositiveInt,
			ContextKey: middleware.MaxPrice,
		},
	}

	getAnnouncementsHandler := http.HandlerFunc(h.getAnnouncements)
	validationMiddleware := middleware.ValidateQueryParams(ValidationRules...)
	authMiddleware := auth.OptionalAuthMiddleware(h.token)

	finalHandler := middleware.Chain(
		getAnnouncementsHandler,
		validationMiddleware,
		authMiddleware,
	)

	mux.HandleFunc("POST /api/v1/announcements", auth.AuthMiddleware(h.token, h.createAnnouncement))
	mux.Handle("GET /api/v1/announcements", finalHandler)
}

// Announcement creation
// @Summary      Create an announcement
// @Description  Create an announcement for authorized users.
// @Tags         Announcements 
// @Accept       json
// @Produce      json
// @Param        request body AnnouncementsPostRequest true "Announcement details"
// @Success      201 {object} AnnouncementsPostResponse "Announcement successfully created"
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

	if err := validate.Struct(apr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := token.ExtractToken(r)
	userId, _ := h.token.ValidateToken(token)
	
	var an model.Announcement
	
	an.Article = apr.Article
	an.Text = apr.Text
	an.CostRubles = apr.Cost

	id, err := strconv.Atoi(userId)
	an.UserId = int64(id)

	err = h.db.CreateAnnouncement(&an)
	if err != nil {
		h.logger.Info(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response AnnouncementsPostResponse

	response.Article = an.Article
	response.Text = an.Text
	response.CostRubles = an.CostRubles
	response.ImageAddress = an.ImageAddress
	response.Date = an.Date
	response.Id = an.Id
	response.UserId = an.UserId
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


// GetAnnouncements gets a paginated list of announcements
// @Summary      Get announcements list
// @Description  Get a paginated list of announcements. This endpoint is public.
// @Tags         Announcements
// @Produce      json
// @Param        page     query      int    false  "Page number for pagination (starts from 1). Defaults to 1."
// @Param        limit    query      int    false  "Number of items per page. Defaults to 10."
// @Param        sort_by  query      string false "Sort order for announcements." Enums(price_asc, price_desc, date_asc, date_desc)
// @Param        min_price query     int false "All announcements with price more than min_price. Default is 0"
// @Param        max_price query     int false "All announcements with price less than max_price. Default is (1 << 31) - 1"
// @Success      200      {array}    AnnouncementsGetResponse
// @Failure      400      "Invalid page or limit parameter"
// @Failure      500      "Internal server error"
// @Router       /api/v1/announcements [get]
// @Security     Bearer
func (h *handler) getAnnouncements(w http.ResponseWriter, r *http.Request) {

	page := r.Context().Value(middleware.PageKey).(int)
	limit := r.Context().Value(middleware.LimitKey).(int)
	sortBy := r.Context().Value(middleware.SortKey).(string)

	minPrice := r.Context().Value(middleware.MinPrice).(int)
	maxPrice := r.Context().Value(middleware.MaxPrice).(int)

	currentUserIdString, _ := h.token.ValidateToken(token.ExtractToken(r))
	currentUserId, err := strconv.Atoi(currentUserIdString)
	h.logger.Debug(currentUserId, err)

	announcements, err := h.db.GetAnnouncementsByPage(page, limit, currentUserId, sortBy, minPrice, maxPrice)
	h.logger.Debug(announcements)

	if err != nil {
		h.logger.Info(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	
	response := make([]AnnouncementsGetResponse, len(announcements))

	for i, an := range announcements {
		response[i].Article       = an.Article
		response[i].IsOwner       = an.IsOwner
		response[i].ImageAddress  = an.ImageAddress
		response[i].CostRubles    = an.CostRubles
		response[i].Text          = an.Text
		response[i].OwnerUsername = an.OwnerUsername
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Info(err)
	}
}
