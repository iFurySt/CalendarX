package custom

import (
	"strings"
	"testing"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

func TestEncodeDecodeRoundTripNormalizesSymbols(t *testing.T) {
	config, err := NewConfig([]string{"msft", "AAPL", "MSFT"})
	if err != nil {
		t.Fatal(err)
	}
	token, err := Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(token, "+/=") {
		t.Fatalf("token is not raw base64url safe: %s", token)
	}
	decoded, err := Decode(token)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(decoded.Symbols, ",")
	if got != "AAPL,MSFT" {
		t.Fatalf("symbols mismatch: %s", got)
	}
}

func TestNormalizeRejectsInvalidSymbol(t *testing.T) {
	if _, err := NormalizeSymbols([]string{"AAPL<script>"}); err == nil {
		t.Fatal("expected invalid symbol error")
	}
}

func TestFeedFiltersEvents(t *testing.T) {
	config, err := NewConfig([]string{"NVDA"})
	if err != nil {
		t.Fatal(err)
	}
	feed := Feed(config, []calendarx.Event{
		{Symbol: "AAPL", Date: "2026-05-21"},
		{Symbol: "NVDA", Date: "2026-05-20"},
	})
	if len(feed.Events) != 1 || feed.Events[0].Symbol != "NVDA" {
		t.Fatalf("unexpected feed events: %+v", feed.Events)
	}
}
