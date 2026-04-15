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

	render(w, "ui/home.html", items)
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

	render(w, "ui/item.html", item)
}

func (h *ItemPageHandler) EditItemPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if r.Method == http.MethodGet {
		item, err := h.ItemSrv.GetItem(r.Context(), id)
		if err != nil {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		render(w, "ui/edit.html", item)
		return
	}
	r.ParseForm()

	name := r.FormValue("name")
	comm := r.FormValue("comment")
	ratingStr := r.FormValue("rating")
	imagepath := r.FormValue("image_path")

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid max_rating")
		return
	}

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

func (h *ItemPageHandler) AddItemPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "ui/add.html", nil)
		return
	}
	name := r.FormValue("name")
	comm := r.FormValue("comment")
	ratingStr := r.FormValue("rating")
	imagepath := r.FormValue("image_path")

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid max_rating")
		return
	}

	createReq := itemSrv.CreateItemRequest{
		Name:      name,
		Comment:   comm,
		Rating:    rating,
		ImagePath: imagepath,
	}

	_, err = h.ItemSrv.AddItem(r.Context(), createReq)
	if err != nil {
		http.Error(w, "Create Item failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func render(w http.ResponseWriter, file string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
