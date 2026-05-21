package company

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

//go:embed watchlists/*.json
var watchlists embed.FS

var Presets = []Preset{
	{Slug: "mega7", Label: "Mega 7", File: "watchlists/mega7.json"},
	{Slug: "nasdaq100", Label: "Nasdaq-100", File: "watchlists/nasdaq100.json"},
	{Slug: "sp500", Label: "S&P 500", File: "watchlists/sp500.json"},
	{Slug: "dow30", Label: "Dow 30", File: "watchlists/dow30.json"},
}

type Preset struct {
	Slug    string   `json:"slug"`
	Label   string   `json:"label"`
	File    string   `json:"-"`
	Symbols []string `json:"symbols,omitempty"`
}

type Record struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry,omitempty"`
	Date        string `json:"date,omitempty"`
	Time        string `json:"time,omitempty"`
	MarketCap   string `json:"marketCap,omitempty"`
	HasEarnings bool   `json:"hasEarnings"`
	Source      string `json:"source,omitempty"`
}

func LoadPresets() ([]Preset, error) {
	presets := make([]Preset, 0, len(Presets))
	for _, preset := range Presets {
		entries, err := loadWatchlist(preset.File)
		if err != nil {
			return nil, fmt.Errorf("load %s preset: %w", preset.Slug, err)
		}
		symbols := make([]string, 0, len(entries))
		for _, entry := range entries {
			symbols = append(symbols, entry.Symbol)
		}
		sort.Strings(symbols)
		preset.Symbols = symbols
		presets = append(presets, preset)
	}
	return presets, nil
}

func BuildUniverse(events []calendarx.Event) ([]Record, error) {
	records := map[string]Record{}
	for _, preset := range Presets {
		entries, err := loadWatchlist(preset.File)
		if err != nil {
			return nil, fmt.Errorf("load %s preset: %w", preset.Slug, err)
		}
		for _, entry := range entries {
			mergeRecord(records, Record{
				Symbol:      entry.Symbol,
				CompanyName: entry.CompanyName,
				Industry:    entry.Industry,
				Source:      preset.Slug,
			})
		}
	}
	for _, event := range events {
		mergeRecord(records, Record{
			Symbol:      event.Symbol,
			CompanyName: event.CompanyName,
			Industry:    event.Industry,
			Date:        event.Date,
			Time:        event.Time,
			MarketCap:   event.MarketCap,
			HasEarnings: true,
			Source:      "earnings",
		})
	}

	out := make([]Record, 0, len(records))
	for _, record := range records {
		out = append(out, record)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].HasEarnings != out[j].HasEarnings {
			return out[i].HasEarnings
		}
		return out[i].Symbol < out[j].Symbol
	})
	return out, nil
}

func Search(records []Record, query string, limit int) []Record {
	query = strings.ToUpper(strings.TrimSpace(query))
	if limit <= 0 {
		limit = 25
	}
	matches := make([]Record, 0, len(records))
	for _, record := range records {
		if query != "" && !strings.Contains(record.Symbol, query) && !strings.Contains(strings.ToUpper(record.CompanyName), query) {
			continue
		}
		matches = append(matches, record)
	}
	sort.SliceStable(matches, func(i, j int) bool {
		left := matchScore(matches[i], query)
		right := matchScore(matches[j], query)
		if left != right {
			return left < right
		}
		if matches[i].HasEarnings != matches[j].HasEarnings {
			return matches[i].HasEarnings
		}
		return matches[i].Symbol < matches[j].Symbol
	})
	if len(matches) > limit {
		return matches[:limit]
	}
	return matches
}

func Default(records []Record, symbols []string, limit int) []Record {
	if limit <= 0 {
		limit = 25
	}
	wanted := map[string]struct{}{}
	for _, symbol := range symbols {
		wanted[strings.ToUpper(symbol)] = struct{}{}
	}
	out := make([]Record, 0, len(wanted))
	for _, record := range records {
		if _, ok := wanted[record.Symbol]; ok {
			out = append(out, record)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Symbol < out[j].Symbol
	})
	if len(out) > limit {
		return out[:limit]
	}
	return out
}

func mergeRecord(records map[string]Record, next Record) {
	next.Symbol = strings.ToUpper(strings.TrimSpace(next.Symbol))
	next.CompanyName = strings.TrimSpace(next.CompanyName)
	next.Industry = strings.TrimSpace(next.Industry)
	if next.Symbol == "" {
		return
	}
	current, ok := records[next.Symbol]
	if !ok {
		if next.CompanyName == "" {
			next.CompanyName = next.Symbol
		}
		records[next.Symbol] = next
		return
	}
	if current.CompanyName == "" || current.CompanyName == current.Symbol {
		current.CompanyName = next.CompanyName
	}
	if current.Industry == "" {
		current.Industry = next.Industry
	}
	if next.HasEarnings {
		current.Date = next.Date
		current.Time = next.Time
		current.MarketCap = next.MarketCap
		current.HasEarnings = true
	}
	if current.Source == "" || current.Source == "earnings" {
		current.Source = next.Source
	}
	if current.CompanyName == "" {
		current.CompanyName = current.Symbol
	}
	records[next.Symbol] = current
}

func matchScore(record Record, query string) int {
	if query == "" {
		return 0
	}
	company := strings.ToUpper(record.CompanyName)
	switch {
	case record.Symbol == query:
		return 0
	case strings.HasPrefix(record.Symbol, query):
		return 1
	case strings.Contains(record.Symbol, query):
		return 2
	case strings.HasPrefix(company, query):
		return 3
	case strings.Contains(company, query):
		return 4
	default:
		return 5
	}
}

func loadWatchlist(path string) ([]calendarx.WatchlistEntry, error) {
	raw, err := watchlists.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []calendarx.WatchlistEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}
	for i := range entries {
		entries[i].Symbol = strings.ToUpper(strings.TrimSpace(entries[i].Symbol))
		entries[i].CompanyName = strings.TrimSpace(entries[i].CompanyName)
		entries[i].Industry = strings.TrimSpace(entries[i].Industry)
		if entries[i].CompanyName == "" {
			entries[i].CompanyName = entries[i].Symbol
		}
	}
	return entries, nil
}
