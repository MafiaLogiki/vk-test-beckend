package store

import "marketplace-service/internal/model"

type AnnouncementsStore interface  {
	CreateAnnouncements(announcement *model.Announcement) error
}
