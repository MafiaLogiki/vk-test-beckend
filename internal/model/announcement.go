package model

type Announcement struct {
	UserId        int64     `json:"-"`
	OwnerUsername string    `json:"owner_username"`
	Article       string    `json:"title"`
	Text          string    `json:"text"`
	CostRubles    int32     `json:"price"`
	ImageAddress  string    `json:"image_url"`
	IsOwner       *bool     `json:"is_owner,omitempty"`
}
