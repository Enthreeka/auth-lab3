package entity

import (
	"time"
)

type Token struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`

	User User `json:"user"`
}
