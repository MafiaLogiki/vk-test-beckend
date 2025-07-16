package store

import (
	"database/sql"
	"fmt"
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
	err := s.DB.QueryRow(query, an.UserId, an.Article, an.Text, an.ImageAddress, an.CostRubles).Scan(&an.Id, &an.Date)
	return err
}

func extractAnnouncements(rows *sql.Rows) ([]model.Announcement, error) {
	defer rows.Close()

	var announcements []model.Announcement
	for rows.Next() {
		var an model.Announcement
		var isOwner sql.NullBool

		if err := rows.Scan(
			&an.OwnerUsername,
			&an.Article,
			&an.Text,
			&an.ImageAddress,
			&an.CostRubles,
			&isOwner,
		); err != nil {
			return nil, err
		}

		if isOwner.Valid {
			an.IsOwner = &isOwner.Bool
		} else {
			an.IsOwner = nil
		}

		announcements = append(announcements, an)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return announcements, nil
}

func (s *PostgresAnnouncementsStore) GetAnnouncementsByPage(page, limit, currentUserId int, sortBy string) ([]model.Announcement, error) {
	query := fmt.Sprintf(`
		SELECT users.username, title, text, image_url, price,
		CASE WHEN $3 > 0 THEN (user_id = $3) ELSE NULL END AS is_owner
		FROM announcements
		JOIN users ON announcements.user_id = users.id
		ORDER BY %s
		LIMIT $1
		OFFSET $2
	`, sortBy)
	offset := (page - 1) * limit

	rows, err := s.DB.Query(query, limit, offset, currentUserId)
	if err != nil {
		return nil, err
	}

	return extractAnnouncements(rows)
}
