package api

import (
	"../controller"
	"github.com/go-chi/chi"
)

type Router interface {
	Route(chi.Router)
}

type NewParams struct {
	ControllerSvc controller.Controller
}

func New(p NewParams) Router {
	return &router{
		controllerSvc: p.ControllerSvc,
	}
}

type router struct {
	controllerSvc controller.Controller
}

// Route routes all requests
func (s *router) Route(r chi.Router) {
	// r.Use() if we want to setup a middleware to authenticate the api caller
	r.Route("/auth", s.authRouter)
}

func (s *router) authRouter(r chi.Router) {
	r.Post("/ip", s.controllerSvc.Handle)
}
