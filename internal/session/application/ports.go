package application

import (
	"context"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

// SessionReader is the smallest persistence contract needed by session reads.
type SessionReader interface {
	GetActiveByID(
		ctx context.Context,
		userID domain.UserID,
		sessionID domain.SessionID,
		now time.Time,
	) (*domain.Session, error)
	ListActiveByUserID(
		ctx context.Context,
		userID domain.UserID,
		now time.Time,
	) ([]domain.Session, error)
	ListActiveByUserIDs(
		ctx context.Context,
		userIDs []domain.UserID,
		now time.Time,
	) (map[domain.UserID][]domain.Session, error)
}

// SessionRevoker is kept separate so read-only consumers do not depend on
// mutation capabilities.
type SessionRevoker interface {
	RevokeByID(
		ctx context.Context,
		userID domain.UserID,
		sessionID domain.SessionID,
		revokedAt time.Time,
	) error
	RevokeAllExcept(
		ctx context.Context,
		userID domain.UserID,
		exceptID domain.SessionID,
		revokedAt time.Time,
	) error
	RevokeAllByUserID(
		ctx context.Context,
		userID domain.UserID,
		revokedAt time.Time,
	) error
}

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
