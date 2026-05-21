# Feature Release Notes

## 2026-05

| Date | Area | User Impact | Change Summary |
| --- | --- | --- | --- |
| 2026-05-21 | Custom feeds | Users can build a no-login custom earnings subscription link that directly returns `.ics` from Vercel. | Added Vercel UI and Go API routes using gzip + base64url encoded stateless configs and `/tmp` runtime caching. |
| 2026-05-21 | Earnings calendars | Added generated calendar subscriptions for Mega 7, Nasdaq-100, S&P 500, Dow 30, and all Nasdaq earnings results. | Introduced a Go/Cobra generator, Nasdaq earnings cache, `.ics` feed output, GitHub Pages site generation, and scheduled Pages deployment. |

## 2026-04

| Date | Area | User Impact | Change Summary |
| --- | --- | --- | --- |
| 2026-04-08 | Template | Introduced the base harness repository template for future services and products. | Added agent entry docs, execution-plan scaffolding, change-history templates, and docs checks. |
