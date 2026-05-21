# CalendarX Earnings Pages

## Request

Build a Go + Cobra version of the Nasdaq-backed earnings calendar workflow, starting with common feed baskets such as Mega 7, Nasdaq-100, S&P 500, and Dow 30, and publish generated links through GitHub Pages on a daily CI schedule.

## Changes

- Added a Go module with a `calendarx` Cobra CLI.
- Added Nasdaq earnings fetch caching under `data/earnings/`.
- Added watchlist-backed feed generation for Mega 7, Nasdaq-100, S&P 500, Dow 30, and all earnings rows.
- Added `.ics` generation and a static `docs/index.html` feed directory.
- Added a scheduled GitHub Pages workflow with pinned Actions.
- Updated repository docs, CI docs, frontend notes, quality scoring, release notes, and an execution plan.

## Design Intent

The implementation keeps source fetching, event normalization, ICS rendering, and HTML generation separate so future calendar families can reuse the feed and Pages path without rewriting the earnings source. Watchlists are committed for the first launch while raw earnings data remains a CI cache.

## Important Files

- `cmd/calendarx/main.go`
- `internal/app/root.go`
- `internal/earnings/`
- `internal/calendarx/`
- `internal/site/`
- `data/watchlists/`
- `.github/workflows/pages.yml`
- `docs/exec-plans/active/calendarx-earnings-pages.md`
