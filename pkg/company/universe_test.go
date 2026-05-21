package company

import (
	"testing"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

func TestSearchFindsWatchlistCompanyWithoutCurrentEvent(t *testing.T) {
	records := []Record{
		{Symbol: "CMC", CompanyName: "Commercial Metals Company", Date: "2026-06-22", HasEarnings: true},
		{Symbol: "META", CompanyName: "Meta Platforms", Source: "mega7"},
		{Symbol: "NVDA", CompanyName: "NVIDIA Corporation", Date: "2026-05-20", HasEarnings: true},
	}
	matches := Search(records, "meta", 10)
	if len(matches) < 1 || matches[0].Symbol != "META" || matches[0].HasEarnings {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestBuildUniverseMergesEvents(t *testing.T) {
	records, err := BuildUniverse([]calendarx.Event{
		{Symbol: "META", CompanyName: "Meta Platforms Inc.", Date: "2026-06-01", Time: "After-hours"},
	})
	if err != nil {
		t.Fatal(err)
	}
	matches := Search(records, "META", 10)
	if len(matches) != 1 {
		t.Fatalf("expected one META match, got %+v", matches)
	}
	if !matches[0].HasEarnings || matches[0].Date != "2026-06-01" {
		t.Fatalf("event metadata was not merged: %+v", matches[0])
	}
}
