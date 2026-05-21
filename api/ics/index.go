package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
	"github.com/iFurySt/CalendarX/pkg/custom"
	"github.com/iFurySt/CalendarX/pkg/runtimecache"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.URL.Query().Get("c")
	if token == "" {
		token = tokenFromPath(r.URL.Path)
	}
	config, err := custom.Decode(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, _, err := runtimecache.LoadEvents(r.Context(), runtimecache.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	content := calendarx.BuildICS(custom.Feed(config, events), time.Now().UTC())
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300, s-maxage=600")
	w.Header().Set("X-CalendarX-Symbols", strings.Join(config.Symbols, ","))
	if r.URL.Query().Get("download") == "1" || strings.Contains(r.URL.Path, "/download/") {
		w.Header().Set("Content-Disposition", `attachment; filename="calendarx-custom.ics"`)
	}
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write([]byte(content))
}

func tokenFromPath(path string) string {
	path = strings.TrimSuffix(path, ".ics")
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "c" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
