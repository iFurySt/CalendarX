# Vercel Custom Feeds

## Request

Build a Vercel-hosted custom CalendarX feed generator using gzip + base64url encoded URL configs, without user login or a database.

## Changes

- Added a Vercel static UI for searching companies, selecting symbols, and generating subscription/download links.
- Added Go Vercel API routes for `/api/search`, `/api/link`, `/api/ics`, and `/api/health`.
- Added clean rewritten URLs for `/api/ics/c/<token>.ics` and `/api/download/c/<token>.ics`.
- Added compressed config token support in `pkg/custom`.
- Added runtime Nasdaq cache support in `pkg/runtimecache`, using `/tmp` as ephemeral storage.
- Moved reusable Go packages from `internal/` to `pkg/` so Vercel functions can import them.
- Added tests for config token handling, filtering, and runtime search.

## Design Intent

The custom feed link is stateless: the selected symbols are encoded inside the URL token. Vercel does not need a database or KV store, and function restarts do not break existing subscription links. The `/tmp` cache only accelerates Nasdaq data refreshes and can be rebuilt at any time.

## Important Files

- `public/index.html`
- `api/search/index.go`
- `api/link/index.go`
- `api/ics/index.go`
- `api/health/index.go`
- `pkg/custom/`
- `pkg/runtimecache/`
- `vercel.json`
