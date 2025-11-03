package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id               uuid.UUID  `db:"id" json:"id"`
	Username         string     `db:"username" json:"username"`
	Email            string     `db:"email" json:"email"`
	PasswordHash     string     `db:"password_hash" json:"password_hash"`
	Password         string     `json:"password"`
	Role             string     `db:"role" json:"role"`
	Status           string     `db:"status" json:"status"`
	PhoneNumber      *string    `db:"phone_number" json:"phone_number,omitempty"`
	RegistrationDate time.Time  `db:"registration_date" json:"registration_date"`
	LastLogin        *time.Time `db:"last_login" json:"last_login,omitempty"`
	ProfilePictureID *string    `db:"profile_picture_id" json:"profile_picture_id"`
	Bio              *string    `db:"bio" json:"bio,omitempty"`
	Settings         []byte     `db:"settings" json:"settings"`
}

type AuthUser struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type AuhtInfo struct {
	Id           uuid.UUID `db:"id" json:"id"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	Confirmed    bool      `db:"confirmed" json:"confirmed"`
	Status       string    `db:"status" json:"status"`
}

type UserData struct {
	Username    string
	Email       string
	PhoneNumber string
	Settings    json.RawMessage
}
