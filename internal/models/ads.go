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

type AdPage struct {
	UserID      uuid.UUID
	AdStatus    string
	Title       string
	Price       float32
	CreatedAt   time.Time
	ViewsCount  int
	Condition   string
	Description string
	SellerName  string
	Online      string
	SellerPhone string
	SellerCity  string
}

type AdsPageData struct {
	Ads        []AdTemplate `json:"ads"`
	TotalCount int          `json:"totalCount"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
}

type MessagesPageData struct {
	Chats []ChatPreview `json:"chats"`
}

type ChatPreview struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	Username     string    `json:"username"`
	LastMessage  string    `json:"lastMessage"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	UnreadCount  int       `json:"unreadCount"`
	IsOnline     bool      `json:"isOnline"`
}

type LikedAdsPageData struct {
	Ads []AdTemplate `json:"ads"`
}

type AdminPanelData struct {
	Stats        AdminStats `json:"stats"`
	Ads          []AdTemplate `json:"ads"`
	Users        []UserTemplate `json:"users"`
	Complaints   []ComplaintTemplate `json:"complaints"`
	Moderation   []AdTemplate `json:"moderation"`
}

type AdminStats struct {
	TotalAds         int `json:"totalAds"`
	TotalUsers       int `json:"totalUsers"`
	PendingReports   int `json:"pendingReports"`
	PendingModeration int `json:"pendingModeration"`
}

type UserTemplate struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PhoneNumber      *string   `json:"phoneNumber"`
	RegistrationDate time.Time `json:"registrationDate"`
	Status           string    `json:"status"`
	Role             string    `json:"role"`
}

type ComplaintTemplate struct {
	ID              uuid.UUID `json:"id"`
	ListingID       *uuid.UUID `json:"listingId"`
	TargetUserID    *uuid.UUID `json:"targetUserId"`
	ComplainantID   uuid.UUID `json:"complainantId"`
	ComplaintType   string    `json:"complaintType"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	AdminID         *uuid.UUID `json:"adminId"`
	ResolutionComment *string `json:"resolutionComment"`
}
