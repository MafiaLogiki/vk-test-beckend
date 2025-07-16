package store

import "marketplace-service/internal/model"

type AnnouncementsStore interface  {
	CreateAnnouncement(an *model.Announcement) error
	GetAnnouncementsByPage(page, limit int32) ([]model.Announcement, error)
}
