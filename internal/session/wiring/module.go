package wiring

import (
	"net/http"

	"github.com/xryceu/nanami-platform-showcase/internal/session/adapters/httpapi"
	"github.com/xryceu/nanami-platform-showcase/internal/session/adapters/memory"
	"github.com/xryceu/nanami-platform-showcase/internal/session/application"
	"github.com/xryceu/nanami-platform-showcase/internal/session/domain"
)

// Module is the standalone composition root for the exported vertical slice.
// Replacing Store with the private PostgreSQL adapter does not change the
// application or HTTP layers.
type Module struct {
	Store   *memory.Store
	Service *application.Service
	Handler http.Handler
}

func New(seed ...domain.Session) Module {
	store := memory.NewStore(seed...)
	service := application.NewService(
		store,
		store,
		application.SystemClock{},
	)
	handler := httpapi.NewHandler(service).Routes()
	return Module{
		Store:   store,
		Service: service,
		Handler: handler,
	}
}
