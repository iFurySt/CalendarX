package watchlist

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

func Load(path string) ([]calendarx.WatchlistEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []calendarx.WatchlistEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("decode watchlist %s: %w", path, err)
	}

	for i := range entries {
		entries[i].Symbol = strings.ToUpper(strings.TrimSpace(entries[i].Symbol))
		entries[i].CompanyName = strings.TrimSpace(entries[i].CompanyName)
		entries[i].Industry = strings.TrimSpace(entries[i].Industry)
		if entries[i].Symbol == "" {
			return nil, fmt.Errorf("watchlist %s has blank symbol at index %d", path, i)
		}
		if entries[i].CompanyName == "" {
			entries[i].CompanyName = entries[i].Symbol
		}
	}

	return entries, nil
}

func Map(entries []calendarx.WatchlistEntry) map[string]calendarx.WatchlistEntry {
	out := make(map[string]calendarx.WatchlistEntry, len(entries))
	for _, entry := range entries {
		out[strings.ToUpper(entry.Symbol)] = entry
	}
	return out
}
