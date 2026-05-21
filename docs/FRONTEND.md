# Frontend Guide

CalendarX has two frontend surfaces:

- `docs/index.html`: generated GitHub Pages directory for preset feeds.
- `public/index.html`: Vercel custom feed builder.

## Local Workflow

```sh
go run ./cmd/calendarx generate
open docs/index.html
vercel dev --listen 4174 --yes
```

No dev server is required for the GitHub Pages output. The Vercel custom UI should be checked through `vercel dev` so the Go API routes and rewrites are active.

## Style

Both pages follow the editorial static HTML direction from `html-effectiveness`: ivory background, serif headings, compact cards, restrained borders, and small operational status surfaces.

Keep the UI focused on the actual feed workflow. Avoid marketing-only sections; the first screen should expose available feeds or the custom feed builder.

## Verification

- Confirm every public feed card has a copy action and a download link.
- Confirm text wraps cleanly on narrow mobile widths.
- Confirm `docs/index.html` uses relative links to `docs/ics/*.ics`.
- Confirm the Vercel UI can search watchlist companies even when they do not have a current-window event.
- Confirm the Vercel UI starts with a Mega 7 available list, supports preset buttons, persists selected companies in `localStorage`, generates a compressed subscription URL, and downloads `.ics`.
- Confirm `/api/ics/c/<token>.ics` returns `text/calendar`.
