# CI/CD Guide

CalendarX uses GitHub Actions to validate the Go generator, refresh earnings data, generate calendar feeds, and deploy GitHub Pages.

## Workflows

- `.github/workflows/pages.yml`: runs on a twice-daily schedule, manual dispatch, and relevant pushes to `main`.

The workflow:

1. Checks out the repository.
2. Sets up Go from `go.mod`.
3. Runs `make ci`.
4. Restores the `data/earnings/` Actions cache.
5. Runs `go run ./cmd/calendarx build`.
6. Uploads `docs/` as the Pages artifact.
7. Deploys GitHub Pages.

All workflow actions are pinned to commit SHAs with the version tag noted in comments.

## Local Commands

```sh
make ci
go run ./cmd/calendarx fetch
go run ./cmd/calendarx generate
go run ./cmd/calendarx build
vercel dev --listen 4174 --yes
```

## Generated Artifacts

- `data/earnings/`: runtime cache, ignored by git and restored through Actions cache.
- `docs/ics/*.ics`: generated feed files for Pages.
- `docs/index.html`: generated feed directory page.

## Vercel

The Vercel app serves custom feeds from:

- `public/index.html`
- `api/search/index.go`
- `api/link/index.go`
- `api/ics/index.go`
- `api/health/index.go`

Deploy with:

```sh
vercel deploy --prod
```

The project is linked locally as `calendarx`; `.vercel/` is ignored and should not be committed.

## Release Posture

CalendarX currently ships static Pages artifacts and Vercel serverless functions. Add SBOM, provenance, and binary release automation once there is a packaged CLI release target.
