package earnings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

const NasdaqURL = "https://api.nasdaq.com/api/calendar/earnings"

type FetchOptions struct {
	DataDir string
	Anchor  string
	Before  int
	After   int
	Sleep   time.Duration
	Client  *http.Client
}

type FetchResult struct {
	Saved  int
	Kept   int
	Failed int
}

type nasdaqResponse struct {
	Data struct {
		AsOf string                `json:"asOf"`
		Rows []calendarx.NasdaqRow `json:"rows"`
	} `json:"data"`
}

func FetchWindow(ctx context.Context, opts FetchOptions) (FetchResult, error) {
	if opts.Anchor == "" {
		opts.Anchor = calendarx.TodayUTC()
	}
	if opts.DataDir == "" {
		opts.DataDir = "data"
	}
	if opts.Sleep == 0 {
		opts.Sleep = 120 * time.Millisecond
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: 25 * time.Second}
	}

	dates, err := calendarx.DateRange(opts.Anchor, opts.Before, opts.After)
	if err != nil {
		return FetchResult{}, err
	}

	var result FetchResult
	var firstErr error
	for _, date := range dates {
		day, err := FetchOneDay(ctx, client, date)
		if err == nil {
			if err := WriteDay(opts.DataDir, date, day); err != nil {
				return result, err
			}
			result.Saved++
			fmt.Printf("[earnings] saved %s (%d rows)\n", date, len(day.Data))
		} else if cacheExists(opts.DataDir, date) {
			result.Kept++
			fmt.Printf("[earnings] kept cached %s after fetch error: %v\n", date, err)
		} else {
			result.Failed++
			fmt.Printf("[earnings] failed %s without cache: %v\n", date, err)
			if firstErr == nil {
				firstErr = err
			}
		}

		if opts.Sleep > 0 {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(opts.Sleep):
			}
		}
	}

	if result.Failed > 0 {
		return result, fmt.Errorf("%d earnings day(s) failed; first error: %w", result.Failed, firstErr)
	}
	return result, nil
}

func FetchOneDay(ctx context.Context, client *http.Client, isoDate string) (calendarx.NasdaqDay, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		day, err := requestDay(ctx, client, isoDate)
		if err == nil {
			return day, nil
		}
		lastErr = err
		timer := time.NewTimer(time.Duration(attempt*attempt) * 300 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return calendarx.NasdaqDay{}, ctx.Err()
		case <-timer.C:
		}
	}
	return calendarx.NasdaqDay{}, lastErr
}

func requestDay(ctx context.Context, client *http.Client, isoDate string) (calendarx.NasdaqDay, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, NasdaqURL+"?date="+isoDate, nil)
	if err != nil {
		return calendarx.NasdaqDay{}, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nasdaq.com")
	req.Header.Set("Referer", "https://www.nasdaq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 CalendarX/0.1")

	res, err := client.Do(req)
	if err != nil {
		return calendarx.NasdaqDay{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return calendarx.NasdaqDay{}, fmt.Errorf("nasdaq status %d: %s", res.StatusCode, string(body))
	}

	var payload nasdaqResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return calendarx.NasdaqDay{}, err
	}
	day := calendarx.NasdaqDay{
		Date: isoDate,
		Data: payload.Data.Rows,
	}
	if day.Data == nil {
		day.Data = []calendarx.NasdaqRow{}
	}
	return day, nil
}

func WriteDay(dataDir string, isoDate string, day calendarx.NasdaqDay) error {
	path := calendarx.EarningsCacheFile(dataDir, isoDate)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(day, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

func cacheExists(dataDir string, isoDate string) bool {
	_, err := os.Stat(calendarx.EarningsCacheFile(dataDir, isoDate))
	return err == nil
}
