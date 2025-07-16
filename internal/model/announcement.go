package model

import "time"

type Announcement struct {
	ID           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	Article      string    `json:"title"`
	Text         string    `json:"text"`
	CostRubles   int32     `json:"price"`
	ImageAddress string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
}
