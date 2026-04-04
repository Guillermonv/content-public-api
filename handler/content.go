package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	gocache "github.com/patrickmn/go-cache"

	"content-public-api/model"
	"content-public-api/store"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	cacheTTL        = 15 * time.Second
	queryTimeout    = 5 * time.Second
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
	page, pageSize, err := parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("content:page=%d:page_size=%d", page, pageSize)

	if cached, found := h.cache.Get(cacheKey); found {
		writeJSON(w, cached)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), queryTimeout)
	defer cancel()

	rows, err := h.store.GetDoneContent(ctx, page, pageSize)
	if err != nil {
		log.Printf("GetContent error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result := model.ContentPage{
		Data:     rows,
		Page:     page,
		PageSize: pageSize,
		HasNext:  len(rows) == pageSize,
	}

	h.cache.Set(cacheKey, result, gocache.DefaultExpiration)
	writeJSON(w, result)
}

func parseParams(r *http.Request) (page, pageSize int, err error) {
	page = defaultPage
	if p := r.URL.Query().Get("page"); p != "" {
		v, e := strconv.Atoi(p)
		if e != nil || v <= 0 {
			return 0, 0, fmt.Errorf("invalid page")
		}
		page = v
	}

	pageSize = defaultPageSize
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		v, e := strconv.Atoi(ps)
		if e != nil || v <= 0 {
			return 0, 0, fmt.Errorf("invalid page_size")
		}
		if v > maxPageSize {
			v = maxPageSize
		}
		pageSize = v
	}

	return page, pageSize, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
