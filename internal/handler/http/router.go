package http

import (
	itemService "RateNote/internal/service/item"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	ItemService itemService.ItemService
}

func NewRouter(cfg RouteConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(CORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	itemPageHandler := NewItemPageHandler(cfg.ItemService)

	r.Get("/", itemPageHandler.ListItem)
	r.Get("/item/{id}", itemPageHandler.GetItemPage)

	r.Get("/add", itemPageHandler.AddItemPage)
	r.Post("/item", itemPageHandler.AddItemPage)

	r.Get("/edit/{id}", itemPageHandler.EditItemPage)
	r.Post("/item/{id}", itemPageHandler.EditItemPage)

	r.Post("/delete/{id}", itemPageHandler.DeleteItemPage)

	r.Route("/api/v1", func(r chi.Router) {
		itemHandler := NewItemHandler(cfg.ItemService)
		itemHandler.RegisterRoutes(r)
	})

	return r
}
