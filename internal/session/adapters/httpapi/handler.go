package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

type SessionUseCases interface {
	ListActive(context.Context, domain.UserID) ([]domain.Session, error)
	Revoke(context.Context, domain.UserID, domain.SessionID) error
	RevokeOthers(context.Context, domain.UserID, domain.SessionID) error
	RevokeAll(context.Context, domain.UserID) error
}

type Handler struct {
	sessions SessionUseCases
}

func NewHandler(sessions SessionUseCases) *Handler {
	return &Handler{sessions: sessions}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /me/sessions", h.list)
	mux.HandleFunc("DELETE /me/sessions/{id}", h.revoke)
	mux.HandleFunc("POST /me/sessions/revoke-others", h.revokeOthers)
	mux.HandleFunc("POST /me/sessions/revoke-all", h.revokeAll)
	return mux
}

type actorContextKey struct{}

type Actor struct {
	UserID           domain.UserID
	CurrentSessionID domain.SessionID
}

// WithActor represents the upstream authentication middleware boundary. The
// HTTP adapter consumes verified identity from context and never accepts actor
// or tenant authority from request headers or payloads.
func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func actorFromRequest(request *http.Request) (Actor, bool) {
	actor, ok := request.Context().Value(actorContextKey{}).(Actor)
	return actor, ok && actor.UserID != ""
}

type sessionResponse struct {
	ID         domain.SessionID `json:"id"`
	DeviceName string           `json:"deviceName,omitempty"`
	IPAddress  string           `json:"ipAddress,omitempty"`
	LastUsedAt time.Time        `json:"lastUsedAt"`
	ExpiresAt  time.Time        `json:"expiresAt"`
	Current    bool             `json:"current"`
}

func (h *Handler) list(response http.ResponseWriter, request *http.Request) {
	actor, ok := actorFromRequest(request)
	if !ok {
		writeError(response, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessions, err := h.sessions.ListActive(request.Context(), actor.UserID)
	if err != nil {
		writeServiceError(response, err)
		return
	}

	payload := make([]sessionResponse, 0, len(sessions))
	for _, session := range sessions {
		payload = append(payload, sessionResponse{
			ID:         session.ID,
			DeviceName: session.DeviceName,
			IPAddress:  session.IPAddress,
			LastUsedAt: session.LastUsedAt,
			ExpiresAt:  session.ExpiresAt,
			Current:    session.ID == actor.CurrentSessionID,
		})
	}
	writeJSON(response, http.StatusOK, payload)
}

func (h *Handler) revoke(response http.ResponseWriter, request *http.Request) {
	actor, ok := actorFromRequest(request)
	if !ok {
		writeError(response, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessionID := domain.SessionID(strings.TrimSpace(request.PathValue("id")))
	if sessionID == "" {
		writeError(response, http.StatusBadRequest, "session id is required")
		return
	}
	if err := h.sessions.Revoke(request.Context(), actor.UserID, sessionID); err != nil {
		writeServiceError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) revokeOthers(response http.ResponseWriter, request *http.Request) {
	actor, ok := actorFromRequest(request)
	if !ok {
		writeError(response, http.StatusUnauthorized, "unauthorized")
		return
	}
	if actor.CurrentSessionID == "" {
		writeError(response, http.StatusBadRequest, "current session is missing")
		return
	}
	if err := h.sessions.RevokeOthers(request.Context(), actor.UserID, actor.CurrentSessionID); err != nil {
		writeServiceError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) revokeAll(response http.ResponseWriter, request *http.Request) {
	actor, ok := actorFromRequest(request)
	if !ok {
		writeError(response, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.sessions.RevokeAll(request.Context(), actor.UserID); err != nil {
		writeServiceError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func writeServiceError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(response, http.StatusNotFound, "session not found")
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(response, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrInvalidSession):
		writeError(response, http.StatusBadRequest, "invalid session")
	default:
		writeError(response, http.StatusInternalServerError, "internal error")
	}
}

func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, map[string]string{"error": message})
}

func writeJSON(response http.ResponseWriter, status int, payload any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(payload)
}
