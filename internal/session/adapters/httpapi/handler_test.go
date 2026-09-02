package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/adapters/httpapi"
	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
	"github.com/xryceu/nanami-platform-showcase/internal/session/wiring"
)

func TestListRequiresAuthenticatedActorContext(t *testing.T) {
	t.Parallel()

	module := wiring.New()
	request := httptest.NewRequest(http.MethodGet, "/me/sessions", nil)
	response := httptest.NewRecorder()

	module.Handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestListMapsDomainSessionsAndMarksCurrent(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	module := wiring.New(
		domain.Session{
			ID:         "current",
			UserID:     "user-1",
			DeviceName: "workstation",
			LastUsedAt: now.Add(-time.Minute),
			ExpiresAt:  now.Add(time.Hour),
		},
		domain.Session{
			ID:         "tablet",
			UserID:     "user-1",
			DeviceName: "tablet",
			LastUsedAt: now.Add(-2 * time.Minute),
			ExpiresAt:  now.Add(time.Hour),
		},
		domain.Session{
			ID:         "other-user",
			UserID:     "user-2",
			DeviceName: "other",
			LastUsedAt: now,
			ExpiresAt:  now.Add(time.Hour),
		},
	)
	request := httptest.NewRequest(http.MethodGet, "/me/sessions", nil)
	request = request.WithContext(httpapi.WithActor(request.Context(), httpapi.Actor{
		UserID:           "user-1",
		CurrentSessionID: "current",
	}))
	response := httptest.NewRecorder()

	module.Handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	var payload []struct {
		ID      domain.SessionID `json:"id"`
		Current bool             `json:"current"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload) != 2 {
		t.Fatalf("expected tenant actor's two sessions, got %#v", payload)
	}
	if payload[0].ID != "current" || !payload[0].Current {
		t.Fatalf("current session mapping is wrong: %#v", payload)
	}
}
