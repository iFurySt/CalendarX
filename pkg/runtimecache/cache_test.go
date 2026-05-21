package runtimecache

import (
	"testing"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

func TestSearchMatchesSymbolAndName(t *testing.T) {
	events := []calendarx.Event{
		{Symbol: "AAPL", CompanyName: "Apple", Date: "2026-05-21"},
		{Symbol: "MSFT", CompanyName: "Microsoft", Date: "2026-05-22"},
	}
	results := Search(events, "micro", 10)
	if len(results) != 1 || results[0].Symbol != "MSFT" {
		t.Fatalf("unexpected results: %+v", results)
	}
}
