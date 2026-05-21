# Quality Score

Track quality by product area and architectural layer so agents can prioritize the weakest parts of the system.

## Suggested Scale

- `A`: strong coverage, stable behavior, clear docs, low operational risk.
- `B`: acceptable but still has known gaps.
- `C`: works but needs targeted hardening.
- `D`: fragile or underspecified.

## Initial Template

| Area | Score | Why | Next Step |
| --- | --- | --- | --- |
| Product surface | B | Preset GitHub Pages feeds and a Vercel custom feed builder are available. | Add IPO and event calendar families after the earnings path is stable. |
| Architecture docs | B | Go/Cobra, Pages, Vercel API, and stateless custom token flows are documented. | Document new source packages as they are added. |
| Testing | C | Unit tests cover date range, ICS output, earnings processing, token encoding, filtering, and runtime search; `make ci` is available. | Add fixture-based end-to-end generation and API handler tests. |
| Observability | C | CLI logs counts and Vercel exposes `/api/health` with cache metadata. | Add generated manifest metadata for last source fetch and per-feed cache freshness. |
| Security | C | No secrets are required; workflow actions are pinned, Vercel configs are stateless, and Nasdaq cache stays out of git. | Add dependency scanning and release provenance when binary releases exist. |
