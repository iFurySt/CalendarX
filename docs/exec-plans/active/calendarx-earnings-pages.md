# CalendarX Earnings Pages

## Goal

Build the first CalendarX product slice: a Go + Cobra CLI that fetches Nasdaq earnings dates, generates `.ics` calendar feeds for common stock baskets, and publishes a GitHub Pages site with subscription and download links. The shape should leave room for future calendar families such as IPOs and major technology events.

## Scope

- In scope:
  - Go module and Cobra command entry point.
  - Nasdaq earnings fetch cache.
  - Watchlist-backed feeds for Mega 7, Nasdaq-100, S&P 500, Dow 30, and an all-earnings feed.
  - Static GitHub Pages output under `docs/`.
  - Daily GitHub Actions generation and Pages deployment.
  - Repository docs and history updates.
- Out of scope:
  - IPO and conference/event feeds.
  - Automated refresh of index constituent lists.
  - Custom user watchlists.

## Context

- Relevant docs:
  - `docs/REPO_COLLAB_GUIDE.md`
  - `docs/ARCHITECTURE.md`
  - `docs/CICD.md`
  - `docs/FRONTEND.md`
- Relevant code paths:
  - `cmd/calendarx/`
  - `internal/`
  - `data/watchlists/`
  - `docs/`
  - `.github/workflows/`
- Constraints:
  - Use Go and Cobra, not Node.js.
  - Keep generated Pages assets usable from GitHub Pages without a server.
  - Keep raw earnings cache out of git and preserve it through CI cache.

## Risks

- Risk: Nasdaq can return empty days, blocked requests, or changed JSON fields.
  - Mitigation: cache one normalized JSON file per date, keep previous cached days on transient failures, and tolerate missing rows during generation.
- Risk: `.ics` feeds become unstable if filenames change.
  - Mitigation: treat feed slugs as public API and keep them documented.
- Risk: first implementation narrows the architecture around earnings only.
  - Mitigation: keep package boundaries as source, feed, ICS, and site generation so new calendar families can be added without replacing the CLI.

## Milestones

1. Discovery and design alignment.
2. Implement fetch, process, ICS, static site, and watchlist slices.
3. Add CI/Pages workflow and docs.
4. Validate with tests, local generation, and rendered page inspection.
5. Push to GitHub and verify Pages deployment if remote access is available.

## Validation

- Commands:
  - `go test ./...`
  - `go run ./cmd/calendarx generate`
  - `go run ./cmd/calendarx build --after 7`
  - `make ci`
- Manual checks:
  - Open `docs/index.html` and verify feed cards show subscribe and download links.
  - Inspect one generated `.ics` file for valid calendar headers and all-day events.
- Observability checks:
  - CI logs should show fetch counts, generated feed counts, and Pages artifact upload.

## Progress Log

- [x] Confirmed CalendarX started from the base template without Go code.
- [x] Reviewed the prior Nasdaq-backed implementation and the HTML reference style.
- [x] Implement the first Go/Cobra generator slice.
- [x] Add CI/Pages deployment and docs.
- [x] Run validation and inspect the rendered page.
- [ ] Push and verify GitHub Pages when the remote is available.

## Decision Log

- 2026-05-21: Use Nasdaq earnings calendar as the first data source because it matches the proven prior implementation and covers the requested stock baskets.
- 2026-05-21: Commit initial watchlists for Mega 7, Nasdaq-100, S&P 500, and Dow 30 so feed generation can run from cached earnings data without scraping constituent pages in the critical path.
