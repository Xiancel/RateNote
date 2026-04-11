package http

import (
	itemSrv "RateNote/internal/service/item"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ItemPageHandler struct {
	ItemSrv itemSrv.ItemService
}

func NewItemPageHandler(srv itemSrv.ItemService) *ItemPageHandler {
	return &ItemPageHandler{ItemSrv: srv}
}

func (h *ItemPageHandler) ListItem(w http.ResponseWriter, r *http.Request) {
	items, err := h.ItemSrv.ListItem(r.Context(), itemSrv.ItemFilter{})
	if err != nil {
		http.Error(w, "Failed to List items", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("ui/home.html"))
	tmpl.Execute(w, items)
}

func (h *ItemPageHandler) GetItemPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	item, err := h.ItemSrv.GetItem(r.Context(), id)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("ui/item.html"))
	tmpl.Execute(w, item)
}
