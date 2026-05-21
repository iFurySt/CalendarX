package custom

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/iFurySt/CalendarX/pkg/calendarx"
)

const (
	CurrentVersion = 1
	MaxSymbols     = 250
	MaxTokenBytes  = 64 * 1024
)

var symbolPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9.-]{0,14}$`)

type Config struct {
	Version int      `json:"v"`
	Symbols []string `json:"symbols"`
	Title   string   `json:"title,omitempty"`
}

func NewConfig(symbols []string) (Config, error) {
	normalized, err := NormalizeSymbols(symbols)
	if err != nil {
		return Config{}, err
	}
	return Config{Version: CurrentVersion, Symbols: normalized}, nil
}

func NormalizeSymbols(symbols []string) ([]string, error) {
	seen := map[string]struct{}{}
	normalized := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.ToUpper(strings.TrimSpace(symbol))
		if symbol == "" {
			continue
		}
		if !symbolPattern.MatchString(symbol) {
			return nil, fmt.Errorf("invalid symbol %q", symbol)
		}
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		normalized = append(normalized, symbol)
	}
	sort.Strings(normalized)
	if len(normalized) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}
	if len(normalized) > MaxSymbols {
		return nil, fmt.Errorf("too many symbols: got %d, max %d", len(normalized), MaxSymbols)
	}
	return normalized, nil
}

func Encode(config Config) (string, error) {
	normalized, err := NormalizeSymbols(config.Symbols)
	if err != nil {
		return "", err
	}
	config.Version = CurrentVersion
	config.Symbols = normalized
	config.Title = strings.TrimSpace(config.Title)

	raw, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	var compressed bytes.Buffer
	zipper, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := zipper.Write(raw); err != nil {
		_ = zipper.Close()
		return "", err
	}
	if err := zipper.Close(); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(compressed.Bytes()), nil
}

func Decode(token string) (Config, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return Config{}, fmt.Errorf("missing config token")
	}
	compressed, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return Config{}, fmt.Errorf("invalid config token: %w", err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return Config{}, fmt.Errorf("invalid compressed config: %w", err)
	}
	defer reader.Close()
	raw, err := io.ReadAll(io.LimitReader(reader, MaxTokenBytes+1))
	if err != nil {
		return Config{}, err
	}
	if len(raw) > MaxTokenBytes {
		return Config{}, fmt.Errorf("config token is too large")
	}
	var config Config
	if err := json.Unmarshal(raw, &config); err != nil {
		return Config{}, fmt.Errorf("invalid config json: %w", err)
	}
	if config.Version != CurrentVersion {
		return Config{}, fmt.Errorf("unsupported config version %d", config.Version)
	}
	config.Symbols, err = NormalizeSymbols(config.Symbols)
	if err != nil {
		return Config{}, err
	}
	config.Title = strings.TrimSpace(config.Title)
	return config, nil
}

func FilterEvents(events []calendarx.Event, symbols []string) []calendarx.Event {
	normalized, err := NormalizeSymbols(symbols)
	if err != nil {
		return nil
	}
	wanted := make(map[string]struct{}, len(normalized))
	for _, symbol := range normalized {
		wanted[symbol] = struct{}{}
	}
	filtered := make([]calendarx.Event, 0, len(wanted))
	for _, event := range events {
		if _, ok := wanted[event.Symbol]; ok {
			filtered = append(filtered, event)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Date == filtered[j].Date {
			return filtered[i].Symbol < filtered[j].Symbol
		}
		return filtered[i].Date < filtered[j].Date
	})
	return filtered
}

func Feed(config Config, events []calendarx.Event) calendarx.Feed {
	title := config.Title
	if title == "" {
		title = "Custom"
	}
	return calendarx.Feed{
		Slug:        "custom",
		Title:       title,
		Description: "Custom CalendarX earnings feed for " + strings.Join(config.Symbols, ", "),
		Group:       "Custom",
		Events:      FilterEvents(events, config.Symbols),
	}
}
