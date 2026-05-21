# Frontend Guide

CalendarX's current frontend is a generated static page at `docs/index.html`.

## Local Workflow

```sh
go run ./cmd/calendarx generate
open docs/index.html
```

No dev server is required. The page uses relative `ics/*.ics` links so it works both from the filesystem and from GitHub Pages.

## Style

The page follows the editorial static HTML direction from `html-effectiveness`: ivory background, serif headings, compact cards, restrained borders, and a small sticky summary toolbar.

Keep the UI focused on the actual feed list. Avoid marketing-only sections; the first screen should expose the available calendar feeds and their subscribe/download actions.

## Verification

- Confirm every public feed card has a copy action and a download link.
- Confirm text wraps cleanly on narrow mobile widths.
- Confirm `docs/index.html` uses relative links to `docs/ics/*.ics`.
