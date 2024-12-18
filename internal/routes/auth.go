package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/controllers/auth"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/middlewares"
)

func RegisterAuthRoutes(root chi.Router, app *app.Registry) {
	root.Route("/auth", func(r chi.Router) {
		r.Mount("/admin", AdminAuthRoutes(app))
	})
}

func AdminAuthRoutes(app *app.Registry) chi.Router {
	controller := auth.NewAuthAdminController(app)
	r := chi.NewRouter()

	r.Post("/", controller.LoginByXinchuanAuth)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AdminAuthMiddleware(app))
		r.Get("/", controller.Me)
		r.Delete("/", controller.Logout)
	})

	return r
}
