package application

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

const PresenceFreshnessWindow = 5 * time.Minute

// Service owns session authorization, lifecycle, and presence rules. It does
// not know whether sessions are stored in PostgreSQL, memory, or another
// adapter.
type Service struct {
	reader  SessionReader
	revoker SessionRevoker
	clock   Clock
}

func NewService(
	reader SessionReader,
	revoker SessionRevoker,
	clock Clock,
) *Service {
	if clock == nil {
		clock = SystemClock{}
	}
	return &Service{
		reader:  reader,
		revoker: revoker,
		clock:   clock,
	}
}

// ValidateAccessSession keeps access tokens fail-closed after revocation,
// database restoration, user recreation, or tenant mismatch.
func (s *Service) ValidateAccessSession(
	ctx context.Context,
	user *domain.User,
	sessionID domain.SessionID,
) error {
	if s == nil || s.reader == nil || user == nil || !user.Active() {
		return domain.ErrUnauthorized
	}
	if sessionID == "" {
		return domain.ErrInvalidSession
	}

	session, err := s.reader.GetActiveByID(
		ctx,
		user.ID,
		sessionID,
		s.clock.Now(),
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrUnauthorized
		}
		return err
	}
	if session == nil || session.UserID != user.ID {
		return domain.ErrUnauthorized
	}
	if !sameTenant(user.TenantID, session.TenantID) {
		return domain.ErrUnauthorized
	}
	return nil
}

func (s *Service) ListActive(
	ctx context.Context,
	userID domain.UserID,
) ([]domain.Session, error) {
	if s == nil || s.reader == nil || userID == "" {
		return nil, domain.ErrUnauthorized
	}
	sessions, err := s.reader.ListActiveByUserID(ctx, userID, s.clock.Now())
	if err != nil {
		return nil, err
	}
	sort.SliceStable(sessions, func(i, j int) bool {
		return sessions[i].LastUsedAt.After(sessions[j].LastUsedAt)
	})
	return sessions, nil
}

func (s *Service) ResolvePresenceBatch(
	ctx context.Context,
	users []domain.User,
) (map[domain.UserID]domain.Presence, error) {
	if s == nil || s.reader == nil {
		return nil, domain.ErrInvalidSession
	}

	now := s.clock.Now()
	result := make(map[domain.UserID]domain.Presence, len(users))
	userIDs := make([]domain.UserID, 0, len(users))
	for i := range users {
		user := &users[i]
		if user.ID == "" {
			continue
		}
		result[user.ID] = PresenceFromAccount(user)
		userIDs = append(userIDs, user.ID)
	}
	if len(userIDs) == 0 {
		return result, nil
	}

	sessions, err := s.reader.ListActiveByUserIDs(ctx, userIDs, now)
	if err != nil {
		return nil, err
	}
	for i := range users {
		user := &users[i]
		if user.ID == "" {
			continue
		}
		result[user.ID] = PresenceFromSessions(user, sessions[user.ID], now)
	}
	return result, nil
}

func (s *Service) Revoke(
	ctx context.Context,
	userID domain.UserID,
	sessionID domain.SessionID,
) error {
	if s == nil || s.revoker == nil || userID == "" || sessionID == "" {
		return domain.ErrInvalidSession
	}
	return s.revoker.RevokeByID(ctx, userID, sessionID, s.clock.Now())
}

func (s *Service) RevokeOthers(
	ctx context.Context,
	userID domain.UserID,
	currentSessionID domain.SessionID,
) error {
	if s == nil || s.revoker == nil || userID == "" || currentSessionID == "" {
		return domain.ErrInvalidSession
	}
	return s.revoker.RevokeAllExcept(
		ctx,
		userID,
		currentSessionID,
		s.clock.Now(),
	)
}

func (s *Service) RevokeAll(
	ctx context.Context,
	userID domain.UserID,
) error {
	if s == nil || s.revoker == nil || userID == "" {
		return domain.ErrInvalidSession
	}
	return s.revoker.RevokeAllByUserID(ctx, userID, s.clock.Now())
}

func PresenceFromAccount(user *domain.User) domain.Presence {
	if user == nil {
		return domain.Presence{
			State:  domain.PresenceUnknown,
			Source: domain.PresenceUnavailable,
		}
	}

	presence := domain.Presence{
		State:  domain.PresenceUnknown,
		Source: domain.PresenceUnavailable,
	}
	if user.LastLoginAt != nil && !user.LastLoginAt.IsZero() {
		lastSeen := user.LastLoginAt.UTC()
		presence.State = domain.PresenceOffline
		presence.Source = domain.PresenceFromLastLogin
		presence.LastSeenAt = &lastSeen
	}
	if user.Status == domain.UserStatusDisabled {
		presence.State = domain.PresenceOffline
		presence.Source = domain.PresenceFromAccountStatus
	}
	return presence
}

func PresenceFromSessions(
	user *domain.User,
	sessions []domain.Session,
	now time.Time,
) domain.Presence {
	presence := PresenceFromAccount(user)
	if user == nil {
		return presence
	}

	var latestSeen *time.Time
	for i := range sessions {
		session := sessions[i]
		if session.UserID != user.ID {
			continue
		}
		seen := session.LastUsedAt.UTC()
		if latestSeen == nil || seen.After(*latestSeen) {
			latestSeen = &seen
		}
	}
	if latestSeen == nil {
		return presence
	}

	freshUntil := latestSeen.Add(PresenceFreshnessWindow)
	presence.LastSeenAt = latestSeen
	presence.FreshUntil = &freshUntil
	presence.Source = domain.PresenceFromSession
	if user.Status == domain.UserStatusDisabled {
		presence.State = domain.PresenceOffline
		presence.Realtime = false
		return presence
	}
	presence.Realtime = true
	if freshUntil.Before(now.UTC()) {
		presence.State = domain.PresenceOffline
	} else {
		presence.State = domain.PresenceOnline
	}
	return presence
}

func sameTenant(left, right *domain.TenantID) bool {
	switch {
	case left == nil && right == nil:
		return true
	case left == nil || right == nil:
		return false
	default:
		return *left == *right
	}
}
