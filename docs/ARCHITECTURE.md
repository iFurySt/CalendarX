# Architecture

CalendarX has two surfaces: a static preset-feed generator for GitHub Pages and a Vercel-hosted custom feed builder. Both use Go to fetch external calendar data, normalize it, and render `.ics` calendars.

## Repository Shape

- `cmd/calendarx/`: CLI executable entry point.
- `internal/app/`: Cobra command wiring and top-level CLI orchestration.
- `pkg/earnings/`: Nasdaq earnings fetch and cached-data processing.
- `pkg/calendarx/`: shared domain types, feed metadata, date helpers, paths, and ICS rendering.
- `pkg/watchlist/`: committed watchlist loading and normalization.
- `pkg/site/`: static HTML page generation for GitHub Pages.
- `pkg/custom/`: stateless compressed custom-feed config tokens.
- `pkg/runtimecache/`: Vercel runtime cache orchestration using `/tmp`.
- `api/`: Vercel Go functions for search, link generation, dynamic `.ics`, and health.
- `public/`: Vercel custom-feed UI.
- `data/watchlists/`: committed basket membership for public feeds.
- `data/earnings/`: local and CI cache for fetched Nasdaq data; ignored by git.
- `docs/`: repository knowledge base and GitHub Pages artifact root.
- `.github/workflows/`: CI and Pages deployment.
- `vercel.json`: Vercel route rewrites for clean dynamic `.ics` URLs.

## Data Flow

1. `calendarx fetch` requests Nasdaq earnings data for each date in a rolling window.
2. Each day is normalized into `data/earnings/YYYY-MM-DD.json`.
3. `calendarx generate` reads the cache, filters events through committed watchlists, and deduplicates by symbol.
4. The generator writes one `.ics` feed per public slug under `docs/ics/`.
5. The generator writes `docs/index.html` with subscribe and download links.
6. GitHub Actions uploads `docs/` as the Pages artifact.

## Vercel Custom Data Flow

1. The user searches companies in `public/index.html`.
2. `/api/search` loads the current rolling Nasdaq window from `pkg/runtimecache`.
3. Selected symbols are posted to `/api/link`.
4. `/api/link` normalizes the symbols and returns `/api/ics/c/<token>.ics`, where `<token>` is gzip-compressed JSON encoded with raw base64url.
5. Calendar apps request the `.ics` URL directly.
6. `/api/ics` decodes the token, refreshes or reads the `/tmp` Nasdaq cache, filters events by symbol, and returns `text/calendar`.

## Package Boundaries

- Source-specific behavior belongs in packages named after the source or calendar family, such as `pkg/earnings`.
- Calendar-family-agnostic rendering belongs in `pkg/calendarx` or a focused renderer package.
- Static site presentation belongs in `pkg/site`; it should consume feed summaries rather than raw source payloads.
- Vercel-imported packages must live under `pkg/`; Vercel's generated handler build path cannot import repository `internal/` packages.
- Future IPO, conference, launch, or product-event feeds should add new source/process packages and feed metadata without changing the public `.ics` renderer contract.

## Operational Model

CalendarX has no long-running server. GitHub Actions is the scheduler and deployment runtime for preset feeds. Vercel functions are stateless request handlers for custom feeds.

The only persisted preset-feed runtime state is the Actions cache for `data/earnings/`; generated Pages artifacts are rebuilt on every workflow run. Vercel stores custom selections in the URL, not in a database. Its `/tmp` cache only accelerates Nasdaq fetches and may disappear at any time.
