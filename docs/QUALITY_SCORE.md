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
| Product surface | C | First earnings calendar feeds and static Pages directory are defined. | Add IPO and event calendar families after the earnings path is stable. |
| Architecture docs | B | Go/Cobra generator boundaries and Pages data flow are documented. | Document new source packages as they are added. |
| Testing | C | Unit tests cover date range, ICS output, and earnings processing; `make ci` is available. | Add fixture-based end-to-end generation tests. |
| Observability | C | CLI logs fetch and generation counts in local and CI runs. | Add generated manifest metadata for last source fetch and per-feed cache freshness. |
| Security | C | No secrets are required; workflow actions are pinned and Nasdaq cache stays out of git. | Add dependency scanning and release provenance when binary releases exist. |
