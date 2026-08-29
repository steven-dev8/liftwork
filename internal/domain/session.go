package domain

import "time"

type Session struct {
	ID               int64
	UserID           int64
	CreatedAt        time.Time
	ExpiresAt        time.Time
	RevokedAt        *time.Time
}

func (s Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s Session) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

func (s Session) IsActive(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}
