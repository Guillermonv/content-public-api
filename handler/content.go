package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	gocache "github.com/patrickmn/go-cache"

	"content-public-api/model"
	"content-public-api/store"
)

const (
	defaultLimit = 20
	maxLimit     = 100
	cacheTTL     = 15 * time.Second
	queryTimeout = 5 * time.Second
)

type ContentHandler struct {
	store *store.ContentStore
	cache *gocache.Cache
}

func NewContentHandler(s *store.ContentStore) *ContentHandler {
	return &ContentHandler{
		store: s,
		cache: gocache.New(cacheTTL, 2*cacheTTL),
	}
}

func (h *ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	limit, cursor, err := parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("content:cursor=%v:limit=%d", cursor, limit)

	if cached, found := h.cache.Get(cacheKey); found {
		writeJSON(w, cached)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	rows, err := h.store.GetDoneContent(ctx, cursor, limit)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	page := model.ContentPage{
		Data:  rows,
		Limit: limit,
	}
	if len(rows) == limit {
		last := rows[len(rows)-1].ID
		page.NextCursor = &last
	}

	h.cache.Set(cacheKey, page, gocache.DefaultExpiration)
	writeJSON(w, page)
}

func parseParams(r *http.Request) (int, *int64, error) {
	limit := defaultLimit
	if l := r.URL.Query().Get("limit"); l != "" {
		v, err := strconv.Atoi(l)
		if err != nil || v <= 0 {
			return 0, nil, fmt.Errorf("invalid limit")
		}
		if v > maxLimit {
			v = maxLimit
		}
		limit = v
	}

	var cursor *int64
	if c := r.URL.Query().Get("cursor"); c != "" {
		v, err := strconv.ParseInt(c, 10, 64)
		if err != nil || v <= 0 {
			return 0, nil, fmt.Errorf("invalid cursor")
		}
		cursor = &v
	}

	return limit, cursor, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
