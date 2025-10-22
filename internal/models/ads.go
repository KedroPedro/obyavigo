package models

import (
	"time"

	"github.com/google/uuid"
)

type AdTemplate struct {
	AdId           uuid.UUID  `json:"adID"`
	UserId         uuid.UUID  `json:"userID"`
	CategoryId     uuid.UUID  `json:"categoryID"`
	CategoryName   string     `json:"categoryName"`
	LocationId     uuid.UUID  `json:"locationID"`
	LocationName   string     `json:"locationName"`
	Title          string     `json:"title"`
	Description    string     `json:"desc"`
	Price          int        `json:"price"`
	Condition      string     `json:"condition"`
	ContactPhone   string     `json:"phone"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	ExpirationDate *time.Time `json:"expirationDate"`
	Status         string     `json:"status"`
	ViewsCount     int        `json:"viewsCount"`
}
