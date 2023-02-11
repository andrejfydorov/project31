package router

import (
	"github.com/go-chi/chi"
	"project31/pkg/controllers"
)

type Router struct {
	router *chi.Mux
}

func New() *Router {
	r := Router{}
	r.router = chi.NewRouter()
	return &r
}

func (r *Router) InitRoutes(c *controllers.Controllers) {
	r.router.Post("/create", c.Create)
	r.router.Post("/make_friends", c.MakeFriends)
	r.router.Post("/user", c.Delete)
	r.router.Get("/friends/{id:[0-9]+}", c.GetFriends)
	r.router.Post("/{id:[0-9]+}", c.Update)
	r.router.Get("/get", c.GetAll)
}

func (r *Router) Router() *chi.Mux {
	return r.router
}
