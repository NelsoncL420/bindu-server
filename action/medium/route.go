package medium

import (
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// medium model
type medium struct {
	Name        string         `json:"name" validate:"required"`
	Slug        string         `json:"slug" validate:"required"`
	Type        string         `json:"type" validate:"required"`
	Title       string         `json:"title" validate:"required"`
	Description string         `json:"description"`
	Caption     string         `json:"caption"`
	AltText     string         `json:"alt_text"`
	FileSize    int64          `json:"file_size" validate:"required"`
	URL         postgres.Jsonb `json:"url"`
	Dimensions  string         `json:"dimensions" validate:"required"`
}

// Router - Group of medium router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)

	r.Route("/{medium_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r

}
