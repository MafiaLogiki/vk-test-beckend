package store

import (
	"database/sql"
	"marketplace-service/internal/model"
)

type PostgresAnnouncementsStore struct {
	DB *sql.DB
}

func NewPostgresAnnouncementsStore(db *sql.DB) *PostgresAnnouncementsStore {
	return &PostgresAnnouncementsStore{DB: db}
}

func (s *PostgresAnnouncementsStore) CreateAnnouncement(an *model.Announcement) error {
	query := `INSERT INTO announcements(user_id, title, text, image_url, price) VALUES($1, $2, $3, $4, $5) RETURNING ID`
	var id int64
	err := s.DB.QueryRow(query, an.UserId, an.Article, an.Text, an.ImageAddress, an.CostRubles).Scan(&id)
	return err
}
