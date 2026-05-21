# CalendarX

CalendarX publishes generated calendar subscriptions for market events.

The first product slice is a Nasdaq-backed earnings calendar generator for common US stock baskets:

- Mega 7
- Nasdaq-100
- S&P 500
- Dow 30
- All companies returned by the Nasdaq earnings calendar window

Generated output lives under `docs/` so GitHub Pages can publish it directly.

## Data Source

Earnings dates come from Nasdaq's earnings calendar endpoint:

```text
https://api.nasdaq.com/api/calendar/earnings?date=YYYY-MM-DD
```

`calendarx fetch` writes one normalized cache file per day into `data/earnings/`. The cache is ignored by git and preserved in CI with GitHub Actions cache.

Watchlist membership for the first stock baskets is committed under `data/watchlists/`.

## Local Usage

Install Go 1.26 or newer, then run:

```sh
go test ./...
go run ./cmd/calendarx fetch
go run ./cmd/calendarx generate
```

Or run the full local build:

```sh
go run ./cmd/calendarx build
```

The generated page is `docs/index.html`; generated feeds are under `docs/ics/`.

Useful flags:

```sh
go run ./cmd/calendarx build --anchor 2026-05-21 --before 1 --after 45
go run ./cmd/calendarx generate --data-dir data --out-dir docs
```

## CI And Pages

`.github/workflows/pages.yml` runs tests, refreshes the Nasdaq earnings cache, generates `.ics` files and `docs/index.html`, then deploys `docs/` through GitHub Pages.

The workflow runs twice daily at 04:34 and 16:34 UTC, on manual dispatch, and on relevant pushes to `main`.

## Repository Guide

Start with `AGENTS.md` and the docs under `docs/` before changing behavior. If code or workflow changes make docs stale, update them in the same task.
