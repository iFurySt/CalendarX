package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/iFurySt/CalendarX/pkg/runtimecache"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	events, meta, err := runtimecache.LoadEvents(r.Context(), runtimecache.Options{})
	status := http.StatusOK
	payload := map[string]any{
		"ok":        err == nil,
		"checkedAt": time.Now().UTC().Format(time.RFC3339),
		"events":    len(events),
		"meta":      meta,
	}
	if err != nil {
		status = http.StatusBadGateway
		payload["error"] = err.Error()
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
