package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatId     *uuid.UUID `json:"chat_id"`
	SellerId   *uuid.UUID `json:"seller_id"`
	CustomerId *uuid.UUID `json:"customer_id"`
	ListingId  *uuid.UUID `json:"listing_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

type ChatPreview struct {
	ChatId            uuid.UUID  `json:"chat_id"`
	ListingId         uuid.UUID  `json:"listing_id"`
	ListingTitle      string     `json:"listing_title"`
	CompanionId       uuid.UUID  `json:"companion_id"`
	CompanionName     string     `json:"companion_name"`
	CompanionAvatarId *string    `json:"companion_avatar_id"`
	LastMessage       string     `json:"last_message"`
	LastMessageTime   *time.Time `json:"last_message_time"`
}

type Message struct {
	Id        *uuid.UUID `json:"id"`
	SenderId  *uuid.UUID `json:"sender_id"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
	ChatId    *uuid.UUID `json:"chat_id"`
}
