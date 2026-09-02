package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/adapters/memory"
	"github.com/xryceu/nanami-platform-showcase/internal/session/application"
	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

type fixedClock struct {
	now time.Time
}

func (clock fixedClock) Now() time.Time {
	return clock.now
}

func TestValidateAccessSessionFailsClosedAcrossTenantBoundary(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 1, 8, 0, 0, 0, time.UTC)
	tenantA := domain.TenantID("tenant-a")
	tenantB := domain.TenantID("tenant-b")
	store := memory.NewStore(domain.Session{
		ID:        "session-1",
		UserID:    "user-1",
		TenantID:  &tenantB,
		ExpiresAt: now.Add(time.Hour),
	})
	service := application.NewService(store, store, fixedClock{now: now})

	err := service.ValidateAccessSession(context.Background(), &domain.User{
		ID:       "user-1",
		TenantID: &tenantA,
		Status:   domain.UserStatusActive,
	}, "session-1")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected tenant mismatch to fail closed, got %v", err)
	}
}

func TestRevokeOthersPreservesCurrentSession(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 1, 8, 0, 0, 0, time.UTC)
	store := memory.NewStore(
		domain.Session{
			ID:        "current",
			UserID:    "user-1",
			ExpiresAt: now.Add(time.Hour),
		},
		domain.Session{
			ID:        "other",
			UserID:    "user-1",
			ExpiresAt: now.Add(time.Hour),
		},
		domain.Session{
			ID:        "another-user",
			UserID:    "user-2",
			ExpiresAt: now.Add(time.Hour),
		},
	)
	service := application.NewService(store, store, fixedClock{now: now})

	if err := service.RevokeOthers(
		context.Background(),
		"user-1",
		"current",
	); err != nil {
		t.Fatalf("revoke others: %v", err)
	}

	activeUserOne, err := service.ListActive(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("list user one: %v", err)
	}
	if len(activeUserOne) != 1 || activeUserOne[0].ID != "current" {
		t.Fatalf("current session was not preserved: %#v", activeUserOne)
	}

	activeUserTwo, err := service.ListActive(context.Background(), "user-2")
	if err != nil {
		t.Fatalf("list user two: %v", err)
	}
	if len(activeUserTwo) != 1 || activeUserTwo[0].ID != "another-user" {
		t.Fatalf("cross-user session was changed: %#v", activeUserTwo)
	}
}

func TestResolvePresenceBatchUsesLatestFreshSession(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 1, 8, 0, 0, 0, time.UTC)
	store := memory.NewStore(
		domain.Session{
			ID:         "older",
			UserID:     "user-1",
			LastUsedAt: now.Add(-10 * time.Minute),
			ExpiresAt:  now.Add(time.Hour),
		},
		domain.Session{
			ID:         "latest",
			UserID:     "user-1",
			LastUsedAt: now.Add(-2 * time.Minute),
			ExpiresAt:  now.Add(time.Hour),
		},
	)
	service := application.NewService(store, store, fixedClock{now: now})

	result, err := service.ResolvePresenceBatch(
		context.Background(),
		[]domain.User{{
			ID:     "user-1",
			Status: domain.UserStatusActive,
		}},
	)
	if err != nil {
		t.Fatalf("resolve presence: %v", err)
	}

	presence := result["user-1"]
	if presence.State != domain.PresenceOnline {
		t.Fatalf("expected online presence, got %q", presence.State)
	}
	if presence.LastSeenAt == nil || !presence.LastSeenAt.Equal(
		now.Add(-2*time.Minute),
	) {
		t.Fatalf("expected latest session timestamp, got %v", presence.LastSeenAt)
	}
	if !presence.Realtime {
		t.Fatal("fresh session must produce realtime presence")
	}
}
