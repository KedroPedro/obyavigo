package models

import (
	"time"

	"github.com/google/uuid"
)

type AdTemplate struct {
	AdId            uuid.UUID  `json:"adID"`
	UserId          uuid.UUID  `json:"userID"`
	CategoryId      uuid.UUID  `json:"categoryID"`
	CategoryName    string     `json:"categoryName"`
	SubcategoryName string     `json:"subcategoryName"`
	LocationId      uuid.UUID  `json:"locationID"`
	LocationName    string     `json:"locationName"`
	Title           string     `json:"title"`
	Description     string     `json:"desc"`
	Price           int        `json:"price"`
	Condition       string     `json:"condition"`
	ContactPhone    string     `json:"phone"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	ExpirationDate  *time.Time `json:"expirationDate"`
	Status          string     `json:"status"`
	ViewsCount      int        `json:"viewsCount"`
	ImageID         *string    `json:"imageID,omitempty"`
}

type AdPage struct {
	AdId        uuid.UUID
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
	Images      []string
}

type AdsPageData struct {
	Ads        []AdTemplate `json:"ads"`
	TotalCount int          `json:"totalCount"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	Category   string       `json:"category"`
}

type AdSearchFilters struct {
	Category    string
	Subcategory string
	Region      string
	Location    string
	MinPrice    *int
	MaxPrice    *int
	Condition   string
	SearchQuery string
	SortBy      string
	Page        int
	Limit       int
}

type LikedAdsPageData struct {
	Ads []AdTemplate `json:"ads"`
}

type AdminPanelData struct {
	Stats      AdminStats          `json:"stats"`
	Ads        []AdTemplate        `json:"ads"`
	Users      []UserTemplate      `json:"users"`
	Complaints []ComplaintTemplate `json:"complaints"`
	Moderation []AdTemplate        `json:"moderation"`
}

type AdminStats struct {
	TotalAds          int `json:"totalAds"`
	TotalUsers        int `json:"totalUsers"`
	PendingReports    int `json:"pendingReports"`
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
	ID                uuid.UUID  `json:"id"`
	ListingID         *uuid.UUID `json:"listing_id"`
	TargetUserID      *uuid.UUID `json:"target_user_id"`
	ComplainantID     uuid.UUID  `json:"complainant_id"`
	ComplainantEmail  string     `json:"complainant_email"`
	ComplaintType     string     `json:"complaint_type"`
	Description       string     `json:"description"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	AdminID           *uuid.UUID `json:"admin_id"`
	ResolutionComment *string    `json:"resolution_comment"`
	AdOwnerID         *uuid.UUID `json:"ad_owner_id"`
}
