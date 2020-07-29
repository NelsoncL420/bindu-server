package theme

import "github.com/go-chi/chi"

// theme request model
type theme struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Slug        string `json:"slug"`
	URL         string `json:"url" validate:"required"`
	Description string `json:"description"`
}

// Router - Group of theme router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)

	r.Route("/{theme_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r

}