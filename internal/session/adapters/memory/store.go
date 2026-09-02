package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/application"
	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

var (
	_ application.SessionReader  = (*Store)(nil)
	_ application.SessionRevoker = (*Store)(nil)
)

// Store is the standalone adapter used by the public showcase. The private
// application wires the same application ports to a GORM/PostgreSQL adapter.
type Store struct {
	mu       sync.RWMutex
	sessions map[domain.SessionID]domain.Session
}

func NewStore(seed ...domain.Session) *Store {
	store := &Store{
		sessions: make(map[domain.SessionID]domain.Session, len(seed)),
	}
	for _, session := range seed {
		store.sessions[session.ID] = cloneSession(session)
	}
	return store
}

func (s *Store) GetActiveByID(ctx context.Context, userID domain.UserID, sessionID domain.SessionID, now time.Time) (*domain.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok || session.UserID != userID || !session.ActiveAt(now) {
		return nil, domain.ErrNotFound
	}
	copy := cloneSession(session)
	return &copy, nil
}

func (s *Store) ListActiveByUserID(ctx context.Context, userID domain.UserID, now time.Time) ([]domain.Session, error) {
	grouped, err := s.ListActiveByUserIDs(ctx, []domain.UserID{userID}, now)
	if err != nil {
		return nil, err
	}
	return grouped[userID], nil
}

func (s *Store) ListActiveByUserIDs(ctx context.Context, userIDs []domain.UserID, now time.Time) (map[domain.UserID][]domain.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wanted := make(map[domain.UserID]struct{}, len(userIDs))
	result := make(map[domain.UserID][]domain.Session, len(userIDs))
	for _, userID := range userIDs {
		if userID == "" {
			continue
		}
		wanted[userID] = struct{}{}
		result[userID] = nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, session := range s.sessions {
		if _, ok := wanted[session.UserID]; !ok || !session.ActiveAt(now) {
			continue
		}
		result[session.UserID] = append(result[session.UserID], cloneSession(session))
	}
	for userID := range result {
		sort.SliceStable(result[userID], func(i, j int) bool {
			return result[userID][i].LastUsedAt.After(result[userID][j].LastUsedAt)
		})
	}
	return result, nil
}

func (s *Store) RevokeByID(ctx context.Context, userID domain.UserID, sessionID domain.SessionID, revokedAt time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok || session.UserID != userID || session.RevokedAt != nil {
		return domain.ErrNotFound
	}
	revoked := revokedAt.UTC()
	session.RevokedAt = &revoked
	session.LastUsedAt = revoked
	s.sessions[sessionID] = session
	return nil
}

func (s *Store) RevokeAllExcept(ctx context.Context, userID domain.UserID, exceptID domain.SessionID, revokedAt time.Time) error {
	return s.revokeMatching(ctx, userID, revokedAt, func(id domain.SessionID) bool {
		return id != exceptID
	})
}

func (s *Store) RevokeAllByUserID(ctx context.Context, userID domain.UserID, revokedAt time.Time) error {
	return s.revokeMatching(ctx, userID, revokedAt, func(domain.SessionID) bool {
		return true
	})
}

func (s *Store) revokeMatching(ctx context.Context, userID domain.UserID, revokedAt time.Time, matches func(domain.SessionID) bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	revoked := revokedAt.UTC()
	for id, session := range s.sessions {
		if session.UserID != userID || session.RevokedAt != nil || !matches(id) {
			continue
		}
		session.RevokedAt = &revoked
		session.LastUsedAt = revoked
		s.sessions[id] = session
	}
	return nil
}

func cloneSession(session domain.Session) domain.Session {
	copy := session
	if session.TenantID != nil {
		tenantID := *session.TenantID
		copy.TenantID = &tenantID
	}
	if session.RevokedAt != nil {
		revokedAt := *session.RevokedAt
		copy.RevokedAt = &revokedAt
	}
	return copy
}
