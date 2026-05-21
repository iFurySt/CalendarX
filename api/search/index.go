package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/iFurySt/CalendarX/pkg/runtimecache"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	events, meta, err := runtimecache.LoadEvents(r.Context(), runtimecache.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	limit := 25
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 && value <= 100 {
			limit = value
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"query":   r.URL.Query().Get("q"),
		"meta":    meta,
		"results": runtimecache.Search(events, r.URL.Query().Get("q"), limit),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=60")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
