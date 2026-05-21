package runtimecache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
	"github.com/iFurySt/CalendarX/pkg/earnings"
)

const (
	DefaultBefore = 1
	DefaultAfter  = 45
	DefaultTTL    = 4 * time.Hour
)

var mu sync.Mutex

type Options struct {
	DataDir string
	Anchor  string
	Before  int
	After   int
	TTL     time.Duration
}

type Meta struct {
	Anchor    string    `json:"anchor"`
	Before    int       `json:"before"`
	After     int       `json:"after"`
	UpdatedAt time.Time `json:"updatedAt"`
	Saved     int       `json:"saved"`
	Kept      int       `json:"kept"`
}

type SearchResult struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	Date        string `json:"date"`
	Time        string `json:"time,omitempty"`
	MarketCap   string `json:"marketCap,omitempty"`
}

func DefaultOptions() Options {
	return Options{
		DataDir: filepath.Join(os.TempDir(), "calendarx"),
		Anchor:  calendarx.TodayUTC(),
		Before:  DefaultBefore,
		After:   DefaultAfter,
		TTL:     DefaultTTL,
	}
}

func LoadEvents(ctx context.Context, opts Options) ([]calendarx.Event, Meta, error) {
	opts = normalizeOptions(opts)
	mu.Lock()
	defer mu.Unlock()

	meta, stale := readMeta(opts)
	if stale {
		var err error
		meta, err = refresh(ctx, opts)
		if err != nil && !hasAnyCache(opts) {
			return nil, meta, err
		}
	}

	events, err := earnings.ProcessWindow(earnings.ProcessOptions{
		DataDir: opts.DataDir,
		Anchor:  opts.Anchor,
		Before:  opts.Before,
		After:   opts.After,
	})
	if err != nil {
		return nil, meta, err
	}
	return events, meta, nil
}

func Search(events []calendarx.Event, query string, limit int) []SearchResult {
	query = strings.ToUpper(strings.TrimSpace(query))
	if limit <= 0 {
		limit = 25
	}
	results := make([]SearchResult, 0, limit)
	for _, event := range events {
		if query != "" && !strings.Contains(event.Symbol, query) && !strings.Contains(strings.ToUpper(event.CompanyName), query) {
			continue
		}
		results = append(results, SearchResult{
			Symbol:      event.Symbol,
			CompanyName: event.CompanyName,
			Date:        event.Date,
			Time:        event.Time,
			MarketCap:   event.MarketCap,
		})
		if len(results) >= limit {
			break
		}
	}
	return results
}

func normalizeOptions(opts Options) Options {
	defaults := DefaultOptions()
	if opts.DataDir == "" {
		opts.DataDir = defaults.DataDir
	}
	if opts.Anchor == "" {
		opts.Anchor = defaults.Anchor
	}
	if opts.Before == 0 {
		opts.Before = defaults.Before
	}
	if opts.After == 0 {
		opts.After = defaults.After
	}
	if opts.TTL == 0 {
		opts.TTL = defaults.TTL
	}
	return opts
}

func readMeta(opts Options) (Meta, bool) {
	raw, err := os.ReadFile(metaFile(opts))
	if err != nil {
		return Meta{}, true
	}
	var meta Meta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return Meta{}, true
	}
	if meta.Anchor != opts.Anchor || meta.Before != opts.Before || meta.After != opts.After {
		return meta, true
	}
	if time.Since(meta.UpdatedAt) > opts.TTL {
		return meta, true
	}
	return meta, !windowExists(opts)
}

func refresh(ctx context.Context, opts Options) (Meta, error) {
	dates, err := calendarx.DateRange(opts.Anchor, opts.Before, opts.After)
	if err != nil {
		return Meta{}, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	type result struct {
		date string
		day  calendarx.NasdaqDay
		err  error
	}
	jobs := make(chan string)
	results := make(chan result)
	workers := 6
	if len(dates) < workers {
		workers = len(dates)
	}
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for date := range jobs {
				day, err := earnings.FetchOneDay(ctx, client, date)
				results <- result{date: date, day: day, err: err}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, date := range dates {
			select {
			case <-ctx.Done():
				return
			case jobs <- date:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	meta := Meta{Anchor: opts.Anchor, Before: opts.Before, After: opts.After, UpdatedAt: time.Now().UTC()}
	var firstErr error
	for result := range results {
		if result.err == nil {
			if err := earnings.WriteDay(opts.DataDir, result.date, result.day); err != nil {
				return meta, err
			}
			meta.Saved++
			continue
		}
		if fileExists(calendarx.EarningsCacheFile(opts.DataDir, result.date)) {
			meta.Kept++
			continue
		}
		if firstErr == nil {
			firstErr = fmt.Errorf("%s: %w", result.date, result.err)
		}
	}
	if firstErr != nil {
		return meta, firstErr
	}
	if err := writeMeta(opts, meta); err != nil {
		return meta, err
	}
	return meta, nil
}

func writeMeta(opts Options, meta Meta) error {
	if err := os.MkdirAll(filepath.Dir(metaFile(opts)), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(metaFile(opts), raw, 0o644)
}

func windowExists(opts Options) bool {
	dates, err := calendarx.DateRange(opts.Anchor, opts.Before, opts.After)
	if err != nil {
		return false
	}
	for _, date := range dates {
		if !fileExists(calendarx.EarningsCacheFile(opts.DataDir, date)) {
			return false
		}
	}
	return true
}

func hasAnyCache(opts Options) bool {
	files, err := filepath.Glob(filepath.Join(opts.DataDir, "earnings", "*.json"))
	return err == nil && len(files) > 0
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func metaFile(opts Options) string {
	return filepath.Join(opts.DataDir, "meta", opts.Anchor+".json")
}

func SortEvents(events []calendarx.Event) {
	sort.Slice(events, func(i, j int) bool {
		if events[i].Date == events[j].Date {
			return events[i].Symbol < events[j].Symbol
		}
		return events[i].Date < events[j].Date
	})
}
