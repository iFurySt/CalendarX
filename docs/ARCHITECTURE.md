# Architecture

CalendarX is a static calendar-feed generator. The current runtime is a Go CLI that fetches external calendar data, normalizes it, writes `.ics` feeds, and builds a GitHub Pages site from generated files.

## Repository Shape

- `cmd/calendarx/`: CLI executable entry point.
- `internal/app/`: Cobra command wiring and top-level orchestration.
- `internal/earnings/`: Nasdaq earnings fetch and cached-data processing.
- `internal/calendarx/`: shared domain types, feed metadata, date helpers, paths, and ICS rendering.
- `internal/watchlist/`: committed watchlist loading and normalization.
- `internal/site/`: static HTML page generation for GitHub Pages.
- `data/watchlists/`: committed basket membership for public feeds.
- `data/earnings/`: local and CI cache for fetched Nasdaq data; ignored by git.
- `docs/`: repository knowledge base and GitHub Pages artifact root.
- `.github/workflows/`: CI and Pages deployment.

## Data Flow

1. `calendarx fetch` requests Nasdaq earnings data for each date in a rolling window.
2. Each day is normalized into `data/earnings/YYYY-MM-DD.json`.
3. `calendarx generate` reads the cache, filters events through committed watchlists, and deduplicates by symbol.
4. The generator writes one `.ics` feed per public slug under `docs/ics/`.
5. The generator writes `docs/index.html` with subscribe and download links.
6. GitHub Actions uploads `docs/` as the Pages artifact.

## Package Boundaries

- Source-specific behavior belongs in packages named after the source or calendar family, such as `internal/earnings`.
- Calendar-family-agnostic rendering belongs in `internal/calendarx` or a focused renderer package.
- Static site presentation belongs in `internal/site`; it should consume feed summaries rather than raw source payloads.
- Future IPO, conference, launch, or product-event feeds should add new source/process packages and feed metadata without changing the public `.ics` renderer contract.

## Operational Model

CalendarX has no long-running server. GitHub Actions is the scheduler and deployment runtime. The only persisted runtime state is the Actions cache for `data/earnings/`; generated Pages artifacts are rebuilt on every workflow run.
