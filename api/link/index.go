package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iFurySt/CalendarX/pkg/custom"
)

type linkRequest struct {
	Symbols []string `json:"symbols"`
	Title   string   `json:"title"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req linkRequest
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32*1024)).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
	} else {
		req.Symbols = strings.Split(r.URL.Query().Get("symbols"), ",")
		req.Title = r.URL.Query().Get("title")
	}

	config, err := custom.NewConfig(req.Symbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config.Title = strings.TrimSpace(req.Title)
	token, err := custom.Encode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	base := origin(r)
	subscribePath := "/api/ics/c/" + token + ".ics"
	downloadPath := "/api/download/c/" + token + ".ics"
	writeJSON(w, http.StatusOK, map[string]any{
		"token":         token,
		"symbols":       config.Symbols,
		"subscribeUrl":  base + subscribePath,
		"downloadUrl":   base + downloadPath,
		"subscribePath": subscribePath,
		"downloadPath":  downloadPath,
	})
}

func origin(r *http.Request) string {
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "https"
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		if r.Header.Get("X-Forwarded-Proto") == "" {
			proto = "http"
		}
	}
	return proto + "://" + host
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
