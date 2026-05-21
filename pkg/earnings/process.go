package earnings

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
	"github.com/iFurySt/CalendarX/pkg/watchlist"
)

type ProcessOptions struct {
	DataDir   string
	Anchor    string
	Before    int
	After     int
	Watchlist []calendarx.WatchlistEntry
	UseFilter bool
}

func ProcessWindow(opts ProcessOptions) ([]calendarx.Event, error) {
	if opts.Anchor == "" {
		opts.Anchor = calendarx.TodayUTC()
	}
	if opts.DataDir == "" {
		opts.DataDir = "data"
	}
	dates, err := calendarx.DateRange(opts.Anchor, opts.Before, opts.After)
	if err != nil {
		return nil, err
	}

	filter := map[string]calendarx.WatchlistEntry(nil)
	if opts.UseFilter {
		filter = watchlist.Map(opts.Watchlist)
	}

	eventsBySymbol := map[string]calendarx.Event{}
	for _, date := range dates {
		day, err := readDay(opts.DataDir, date)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, row := range day.Data {
			symbol := strings.ToUpper(strings.TrimSpace(row.Symbol))
			if symbol == "" {
				continue
			}
			entry, ok := filter[symbol]
			if opts.UseFilter && !ok {
				continue
			}
			eventsBySymbol[symbol] = normalizeEvent(row, firstNonBlank(day.Date, date), entry, ok)
		}
	}

	events := make([]calendarx.Event, 0, len(eventsBySymbol))
	for _, event := range eventsBySymbol {
		events = append(events, event)
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].Date == events[j].Date {
			return events[i].Symbol < events[j].Symbol
		}
		return events[i].Date < events[j].Date
	})
	return events, nil
}

func readDay(dataDir string, isoDate string) (calendarx.NasdaqDay, error) {
	raw, err := os.ReadFile(calendarx.EarningsCacheFile(dataDir, isoDate))
	if err != nil {
		return calendarx.NasdaqDay{}, err
	}
	var day calendarx.NasdaqDay
	if err := json.Unmarshal(raw, &day); err != nil {
		return calendarx.NasdaqDay{}, err
	}
	if day.Data == nil {
		day.Data = []calendarx.NasdaqRow{}
	}
	return day, nil
}

func normalizeEvent(row calendarx.NasdaqRow, date string, entry calendarx.WatchlistEntry, hasEntry bool) calendarx.Event {
	symbol := strings.ToUpper(strings.TrimSpace(row.Symbol))
	name := strings.TrimSpace(row.Name)
	industry := ""
	if hasEntry {
		name = firstNonBlank(entry.CompanyName, name, symbol)
		industry = entry.Industry
	}
	return calendarx.Event{
		Symbol:              symbol,
		CompanyName:         firstNonBlank(name, symbol),
		Industry:            industry,
		Date:                date,
		MarketCap:           strings.TrimSpace(row.MarketCap),
		FiscalQuarterEnding: strings.TrimSpace(row.FiscalQuarterEnding),
		Time:                normalizeTime(row.Time),
		EPSForecast:         strings.TrimSpace(row.EPSForecast),
		NoOfEsts:            strings.TrimSpace(row.NoOfEsts),
	}
}

func normalizeTime(raw string) string {
	switch strings.TrimSpace(raw) {
	case "time-pre-market":
		return "Pre-market"
	case "time-after-hours", "time-after-market":
		return "After-hours"
	default:
		return ""
	}
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
