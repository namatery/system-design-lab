package model

import "time"

type Link struct {
	Alias     string     `json:"alias"`
	LongURL   string     `json:"long_url"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (l Link) Expired(now time.Time) bool {
	return l.ExpiresAt != nil && !l.ExpiresAt.After(now)
}
