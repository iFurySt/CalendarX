package calendarx

import (
	"strings"
	"testing"
	"time"
)

func TestBuildICSEscapesTextAndWritesAllDayEvent(t *testing.T) {
	content := BuildICS(Feed{
		Slug:        "test",
		Title:       "Test Feed",
		Description: "A feed with punctuation",
		Events: []Event{
			{
				Symbol:      "ACME",
				CompanyName: "Acme, Inc.; Research",
				Date:        "2026-05-21",
				Time:        "Pre-market",
			},
		},
	}, time.Date(2026, 5, 21, 1, 2, 3, 0, time.UTC))

	for _, want := range []string{
		"BEGIN:VCALENDAR\r\n",
		"DTSTART;VALUE=DATE:20260521\r\n",
		"DTEND;VALUE=DATE:20260522\r\n",
		"SUMMARY:ACME Acme\\, Inc.\\; Research (Pre-market) earnings\r\n",
		"DTSTAMP:20260521T010203Z\r\n",
		"END:VCALENDAR\r\n",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected ICS to contain %q\n%s", want, content)
		}
	}
}

func TestDateRange(t *testing.T) {
	got, err := DateRange("2026-05-21", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2026-05-20", "2026-05-21", "2026-05-22", "2026-05-23"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("DateRange mismatch: got %v want %v", got, want)
	}
}
