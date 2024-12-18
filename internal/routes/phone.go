package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	controller "github.com/xinchuantw/hoki-tabloid-backend/internal/controllers/phone"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/middlewares"
)

func RegisterPhoneRoutes(root chi.Router, app *app.Registry) {
	phoneController := controller.NewPhoneController(app)

	root.Route("/phones", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(app))
			r.Post("/", phoneController.CreatePhone)
			r.Get("/{PhoneID}", phoneController.GetPhone)
			r.Patch("/{PhoneID}", phoneController.UpdatePhone)
			r.Delete("/{PhoneID}", phoneController.DeletePhone)
			r.Get("/", phoneController.GetPhones)
		})
	})
}
