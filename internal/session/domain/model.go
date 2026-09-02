package domain

import (
	"errors"
	"time"
)

var (
	ErrUnauthorized   = errors.New("session actor is not authorized")
	ErrInvalidSession = errors.New("session is invalid")
	ErrNotFound       = errors.New("session not found")
)

type UserID string
type TenantID string
type SessionID string

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID          UserID
	TenantID    *TenantID
	Status      UserStatus
	LastLoginAt *time.Time
}

func (u User) Active() bool {
	return u.ID != "" && u.Status == UserStatusActive
}

type Session struct {
	ID         SessionID
	UserID     UserID
	TenantID   *TenantID
	DeviceName string
	IPAddress  string
	LastUsedAt time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
}

func (s Session) ActiveAt(now time.Time) bool {
	return s.RevokedAt == nil && s.ExpiresAt.After(now.UTC())
}

type PresenceState string

const (
	PresenceOnline  PresenceState = "online"
	PresenceOffline PresenceState = "offline"
	PresenceUnknown PresenceState = "unknown"
)

type PresenceSource string

const (
	PresenceFromSession       PresenceSource = "session"
	PresenceFromLastLogin     PresenceSource = "last_login"
	PresenceFromAccountStatus PresenceSource = "account_status"
	PresenceUnavailable       PresenceSource = "not_available"
)

type Presence struct {
	State      PresenceState
	Source     PresenceSource
	Realtime   bool
	LastSeenAt *time.Time
	FreshUntil *time.Time
}
