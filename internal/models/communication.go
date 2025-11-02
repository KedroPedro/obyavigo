package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatId     *uuid.UUID
	SellerId   *uuid.UUID
	CustomerId *uuid.UUID
	ListingId  *uuid.UUID
	CreatedAt  *time.Time
}

type Message struct {
	Id        *uuid.UUID
	SenderId  *uuid.UUID
	Text      string
	CreatedAt *time.Time
	ChatId    *uuid.UUID
}
