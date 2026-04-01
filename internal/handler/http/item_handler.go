package http

import (
	itemSrv "RateNote/internal/service/item"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ItemHandler struct {
	ItemSrv itemSrv.ItemService
}

func NewItemHandler(srv itemSrv.ItemService) *ItemHandler {
	return &ItemHandler{ItemSrv: srv}
}

func (h *ItemHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/items", h.ListItem)
		r.Get("/items/{id}", h.GetItem)
		r.Post("/items", r.AddItem)
		r.Put("/items/{id}", r.UpdateItem)
		r.Delete("/items/{id}", r.DeleteItem)
	})
}

func (h *ItemHandler) ListItem(w http.ResponseWriter, r *http.Request) {
	filter := itemSrv.ItemFilter{
		Limit:  20,
		Offset: 0,
	}

	filter.Name = r.URL.Query().Get("name")

	if minRatingStr := r.URL.Query().Get("min_rating"); minRatingStr != "" {
		minRating, err := strconv.ParseFloat(minRatingStr, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid min_rating")
			return
		}
		filter.MinRating = &minRating
	}

	if maxRatingStr := r.URL.Query().Get("max_rating"); maxRatingStr != "" {
		maxRating, err := strconv.ParseFloat(maxRatingStr, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid max_rating")
			return
		}
		filter.MaxRating = &maxRating
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid offset")
			return
		}
		filter.Offset = offset
	}

	response, err := h.ItemSrv.ListItem(r.Context(), filter)
	if err != nil {
		handleServiceItemError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	item, err := h.ItemSrv.GetItem(r.Context(), id)
	if err != nil {
		handleServiceItemError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func handleServiceItemError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, itemSrv.ErrItemNotFound):
		respondError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, itemSrv.ErrItemNameRequired),
		errors.Is(err, itemSrv.ErrInvalidRating),
		errors.Is(err, itemSrv.ErrNoFields):
		respondError(w, http.StatusBadRequest, err.Error())

	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
