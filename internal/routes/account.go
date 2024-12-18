package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/controllers/public"
)

func RegisterAccountRoutes(root chi.Router, app *app.Registry) {
	root.Route("/accounts", func(r chi.Router) {
		accountController := public.NewAccountController(app)

		r.Get("/check-account-duplicate", accountController.CheckDuplicate)
	})
}
