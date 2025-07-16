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
	query := `INSERT INTO announcements(user_id, title, text, image_url, price) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := s.DB.QueryRow(query, an.UserId, an.Article, an.Text, an.ImageAddress, an.CostRubles).Scan(&an.ID, &an.CreatedAt)
	return err
}

func extractAnnouncements(rows *sql.Rows) ([]model.Announcement, error) {
	defer rows.Close()

	var announcements []model.Announcement
	for rows.Next() {
		var an model.Announcement
		if err := rows.Scan(
			&an.ID,
			&an.UserId,
			&an.Article,
			&an.Text,
			&an.ImageAddress,
			&an.CostRubles,
			&an.CreatedAt,
		); err != nil {
			return nil, err
		}
		announcements = append(announcements, an)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return announcements, nil
}

func (s *PostgresAnnouncementsStore) GetAnnouncementsByPage(page, limit int32) ([]model.Announcement, error) {
	query := `
		SELECT id, user_id, title, text, image_url, price, created_at 
		FROM announcements
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`
	offset := (page - 1) * limit

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	return extractAnnouncements(rows)
}
