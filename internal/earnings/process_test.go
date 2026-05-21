package earnings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/iFurySt/CalendarX/internal/calendarx"
)

func TestProcessWindowFiltersDedupesAndSorts(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "earnings")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeDay := func(date string, rows []calendarx.NasdaqRow) {
		t.Helper()
		raw, err := json.Marshal(calendarx.NasdaqDay{Date: date, Data: rows})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cacheDir, date+".json"), raw, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeDay("2026-05-21", []calendarx.NasdaqRow{
		{Symbol: "AAPL", Name: "Apple", Time: "time-pre-market"},
		{Symbol: "MSFT", Name: "Microsoft"},
	})
	writeDay("2026-05-22", []calendarx.NasdaqRow{
		{Symbol: "AAPL", Name: "Apple Updated", Time: "time-after-hours"},
	})

	events, err := ProcessWindow(ProcessOptions{
		DataDir: dir,
		Anchor:  "2026-05-21",
		After:   1,
		Watchlist: []calendarx.WatchlistEntry{
			{Symbol: "AAPL", CompanyName: "Apple Inc.", Industry: "Technology"},
		},
		UseFilter: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one filtered/deduped event, got %d", len(events))
	}
	if events[0].Symbol != "AAPL" || events[0].Date != "2026-05-22" || events[0].Time != "After-hours" {
		t.Fatalf("unexpected event: %+v", events[0])
	}
	if events[0].CompanyName != "Apple Inc." || events[0].Industry != "Technology" {
		t.Fatalf("watchlist enrichment missing: %+v", events[0])
	}
}
