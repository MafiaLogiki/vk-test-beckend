package store

import "marketplace-service/internal/model"

var ValidSortColumns = map[string]string {
	"price_asc":  "price ASC",
	"price_desc": "price DESC",
	"date_asc":   "created_at ASC",
	"date_desc":  "created_at DESC",
}

const DefaultSorting = "date_desc"

type AnnouncementsStore interface  {
	CreateAnnouncement(an *model.Announcement) error
	GetAnnouncementsByPage(page, limit, currentUserId int, sortBy string, maxPrice, minPrice int) ([]model.Announcement, error)
}
