package http

import (
	itemSrv "RateNote/internal/service/item"
	"net/http"
	"strconv"
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

func (h *ItemPageHandler) EditItemPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if r.Method == http.MethodGet {
		item, _ := h.ItemSrv.GetItem(r.Context(), id)
		tmpl := template.Must(template.ParseFiles("templates/edit.html"))
		tmpl.Execute(w, item)
		return
	}
	name := r.FormValue("name")
	comm := r.FormValue("comm")
	ratingStr := r.FormValue("rating")
	imagepath := r.FormValue("image_path")

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid max_rating")
		return
	}
	r.ParseForm()
	updReq := itemSrv.UpdateItemRequest{
		Name:      &name,
		Comment:   &comm,
		Rating:    &rating,
		ImagePath: &imagepath,
	}

	_, err = h.ItemSrv.UpdateItem(r.Context(), id, updReq)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *ItemPageHandler) DeleteItemPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.ItemSrv.DeleteItem(r.Context(), id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
