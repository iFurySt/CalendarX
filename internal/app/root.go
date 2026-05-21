package app

import (
	"context"
	"fmt"
	"time"

	"github.com/iFurySt/CalendarX/internal/calendarx"
	"github.com/iFurySt/CalendarX/internal/earnings"
	"github.com/iFurySt/CalendarX/internal/site"
	"github.com/iFurySt/CalendarX/internal/watchlist"
	"github.com/spf13/cobra"
)

type options struct {
	dataDir string
	outDir  string
	anchor  string
	before  int
	after   int
}

func NewRootCommand() *cobra.Command {
	opts := options{
		dataDir: "data",
		outDir:  "docs",
		anchor:  calendarx.TodayUTC(),
		before:  1,
		after:   45,
	}

	root := &cobra.Command{
		Use:   "calendarx",
		Short: "Generate market calendar subscription feeds",
		Long:  "CalendarX fetches market calendar data and publishes stable .ics feeds plus a static index page.",
	}
	root.PersistentFlags().StringVar(&opts.dataDir, "data-dir", opts.dataDir, "data directory containing cache and watchlists")
	root.PersistentFlags().StringVar(&opts.outDir, "out-dir", opts.outDir, "output directory for generated Pages assets")
	root.PersistentFlags().StringVar(&opts.anchor, "anchor", opts.anchor, "anchor date in YYYY-MM-DD")
	root.PersistentFlags().IntVar(&opts.before, "before", opts.before, "number of days before the anchor date")
	root.PersistentFlags().IntVar(&opts.after, "after", opts.after, "number of days after the anchor date")

	root.AddCommand(fetchCommand(&opts))
	root.AddCommand(generateCommand(&opts))
	root.AddCommand(buildCommand(&opts))
	return root
}

func fetchCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Fetch the Nasdaq earnings window into the local cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := earnings.FetchWindow(cmd.Context(), earnings.FetchOptions{
				DataDir: opts.dataDir,
				Anchor:  opts.anchor,
				Before:  opts.before,
				After:   opts.after,
			})
			fmt.Printf("[fetch] saved=%d kept=%d failed=%d\n", result.Saved, result.Kept, result.Failed)
			return err
		},
	}
}

func generateCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate .ics feeds and the GitHub Pages index from cached data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Generate(cmd.Context(), *opts)
		},
	}
}

func buildCommand(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Fetch data, then generate feeds and the static site",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := earnings.FetchWindow(cmd.Context(), earnings.FetchOptions{
				DataDir: opts.dataDir,
				Anchor:  opts.anchor,
				Before:  opts.before,
				After:   opts.after,
			})
			fmt.Printf("[fetch] saved=%d kept=%d failed=%d\n", result.Saved, result.Kept, result.Failed)
			if err != nil {
				return err
			}
			return Generate(cmd.Context(), *opts)
		},
	}
}

func Generate(ctx context.Context, opts options) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	generatedAt := time.Now().UTC()
	feeds := make([]calendarx.FeedSummary, 0, len(calendarx.BasketFeeds)+1)
	for _, meta := range calendarx.BasketFeeds {
		entries, err := watchlist.Load(calendarx.WatchlistFile(opts.dataDir, meta.Slug))
		if err != nil {
			return err
		}
		events, err := earnings.ProcessWindow(earnings.ProcessOptions{
			DataDir:   opts.dataDir,
			Anchor:    opts.anchor,
			Before:    opts.before,
			After:     opts.after,
			Watchlist: entries,
			UseFilter: true,
		})
		if err != nil {
			return err
		}
		feed := calendarx.Feed{
			Slug:        meta.Slug,
			Title:       meta.Title,
			Description: meta.Description,
			Group:       meta.Group,
			Events:      events,
		}
		if err := calendarx.WriteICS(calendarx.ICSFile(opts.outDir, meta.Slug), feed, generatedAt); err != nil {
			return err
		}
		meta.Count = len(events)
		feeds = append(feeds, meta)
		fmt.Printf("[generate] %s: %d events\n", meta.Slug, len(events))
	}

	allMeta := calendarx.AllEarningsFeed
	allEvents, err := earnings.ProcessWindow(earnings.ProcessOptions{
		DataDir: opts.dataDir,
		Anchor:  opts.anchor,
		Before:  opts.before,
		After:   opts.after,
	})
	if err != nil {
		return err
	}
	if err := calendarx.WriteICS(calendarx.ICSFile(opts.outDir, allMeta.Slug), calendarx.Feed{
		Slug:        allMeta.Slug,
		Title:       allMeta.Title,
		Description: allMeta.Description,
		Group:       allMeta.Group,
		Events:      allEvents,
	}, generatedAt); err != nil {
		return err
	}
	allMeta.Count = len(allEvents)
	feeds = append(feeds, allMeta)
	fmt.Printf("[generate] %s: %d events\n", allMeta.Slug, len(allEvents))

	if err := site.WriteIndex(calendarx.SiteIndexFile(opts.outDir), feeds, generatedAt); err != nil {
		return err
	}
	fmt.Printf("[generate] wrote %s\n", calendarx.SiteIndexFile(opts.outDir))
	return nil
}
