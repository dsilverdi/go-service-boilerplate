package server

import (
	"go-service-boilerplate/app/core"
	"net/http"

	"github.com/go-chi/chi"
)

type AppRestService struct {
	RestService *chi.Mux
	AppCore     *core.AppCore
}

func NewService(svc *chi.Mux, app *core.AppCore) *AppRestService {
	return &AppRestService{
		RestService: svc,
		AppCore:     app,
	}
}

func (svc *AppRestService) MountServerRoute() http.Handler {
	r := svc.RestService
	r.MethodFunc(http.MethodGet, "/", svc.HelloHandler)
	return r
}
