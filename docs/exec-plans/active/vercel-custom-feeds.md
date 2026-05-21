# Vercel Custom Feeds

## Goal

Add a Vercel-hosted CalendarX surface that lets users search companies, select symbols, and generate self-contained compressed subscription links that return live `.ics` calendars without requiring login or a database.

## Scope

- In scope:
  - Static Vercel UI under `public/`.
  - Go Vercel API routes for search, link generation, dynamic `.ics`, and health.
  - Gzip + base64url encoded config tokens.
  - Runtime Nasdaq cache under Vercel `/tmp`.
  - Documentation, tests, and deployment verification.
- Out of scope:
  - Random short-code links backed by KV or a database.
  - User accounts.
  - Persisting custom configurations outside the URL.

## Context

- Relevant docs:
  - `docs/ARCHITECTURE.md`
  - `docs/CICD.md`
  - `docs/FRONTEND.md`
  - `docs/HISTORY_GUIDE.md`
- Relevant code paths:
  - `public/index.html`
  - `api/*/index.go`
  - `pkg/custom/`
  - `pkg/runtimecache/`
  - `pkg/earnings/`
- Constraints:
  - Calendar apps must receive a URL that directly returns `text/calendar`.
  - No DB means the custom configuration must be encoded in the URL.
  - Vercel functions can use `/tmp` only as an ephemeral cache.

## Risks

- Risk: Long symbol lists can still create long URLs.
  - Mitigation: gzip + base64url encode normalized JSON and cap custom feeds at 250 symbols.
- Risk: Cold functions must rebuild the earnings cache.
  - Mitigation: fetch the rolling Nasdaq window concurrently into `/tmp` and reuse it for several hours per warm instance.
- Risk: Vercel Go functions cannot import Go `internal/` packages after runtime wrapping.
  - Mitigation: move reusable generator libraries to `pkg/`, keeping only CLI orchestration in `internal/app`.

## Milestones

1. Confirm Vercel Go route and `/tmp` constraints.
2. Implement compressed config tokens and runtime cache.
3. Implement Vercel UI and API routes.
4. Validate locally with `vercel dev`.
5. Deploy and verify production URLs.

## Validation

- Commands:
  - `go test ./...`
  - `vercel dev --listen 4174 --yes`
  - `curl /api/health`
  - `curl /api/link`
  - `curl /api/ics/c/<token>.ics`
- Manual checks:
  - Search for a company in the UI.
  - Add a symbol and generate a link.
  - Confirm the subscription URL returns `text/calendar`.
  - Confirm mobile width has no horizontal overflow.

## Progress Log

- [x] Confirmed Vercel CLI login and Go runtime route format.
- [x] Implemented compressed config token round-trip tests.
- [x] Implemented runtime cache and API routes.
- [x] Verified local `vercel dev` API and UI.
- [x] Deploy and verify production Vercel URL.

## Decision Log

- 2026-05-21: Use self-contained gzip + base64url config tokens instead of DB-backed short codes so subscription links survive stateless Vercel restarts.
- 2026-05-21: Move shared code from `internal/` to `pkg/` because Vercel's Go runtime wrapper cannot import repository `internal` packages from generated handler paths.
