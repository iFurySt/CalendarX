package site

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/iFurySt/CalendarX/internal/calendarx"
)

type Page struct {
	GeneratedAt string
	Feeds       []calendarx.FeedSummary
}

func WriteIndex(path string, feeds []calendarx.FeedSummary, generatedAt time.Time) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return indexTemplate.Execute(file, Page{
		GeneratedAt: generatedAt.UTC().Format("2006-01-02 15:04 UTC"),
		Feeds:       feeds,
	})
}

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>CalendarX - Market calendar subscriptions</title>
<meta name="description" content="Subscribe to generated market calendar feeds for Mega 7, Nasdaq-100, S&P 500, Dow 30, and more.">
<style>
  :root {
    --ivory: #FAF9F5;
    --paper: #FFFFFF;
    --slate: #141413;
    --clay: #D97757;
    --clay-d: #B85C3E;
    --oat: #E3DACC;
    --olive: #788C5D;
    --g100: #F0EEE6;
    --g200: #E6E3DA;
    --g300: #D1CFC5;
    --g500: #87867F;
    --g700: #3D3D3A;
    --serif: ui-serif, Georgia, "Times New Roman", Times, serif;
    --sans: system-ui, -apple-system, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    --mono: ui-monospace, "SF Mono", Menlo, Monaco, Consolas, monospace;
  }
  * { box-sizing: border-box; }
  html { scroll-behavior: smooth; }
  body {
    margin: 0;
    background: var(--ivory);
    color: var(--slate);
    font-family: var(--sans);
    line-height: 1.55;
    -webkit-font-smoothing: antialiased;
  }
  .wrap { max-width: 1120px; margin: 0 auto; padding: 0 32px 96px; }
  header {
    padding: 72px 0 44px;
    border-bottom: 1.5px solid var(--g300);
    margin-bottom: 28px;
  }
  .eyebrow {
    font-family: var(--mono);
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--g500);
    margin-bottom: 18px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .eyebrow::before { content: ""; width: 24px; height: 1.5px; background: var(--clay); }
  h1 {
    font-family: var(--serif);
    font-weight: 500;
    font-size: clamp(38px, 5.2vw, 62px);
    line-height: 1.06;
    letter-spacing: 0;
    margin: 0;
    max-width: 15ch;
  }
  h1 em { color: var(--clay); font-style: italic; }
  .intro {
    color: var(--g700);
    font-size: 16.5px;
    max-width: 660px;
    margin: 22px 0 0;
  }
  .toolbar {
    position: sticky;
    top: 0;
    z-index: 2;
    display: flex;
    gap: 14px;
    align-items: center;
    justify-content: space-between;
    padding: 14px 0 12px;
    background: var(--ivory);
    border-bottom: 1.5px solid var(--g300);
  }
  .summary {
    font-family: var(--mono);
    font-size: 12px;
    color: var(--g700);
  }
  .summary b { color: var(--slate); }
  .source {
    color: var(--g500);
    font-size: 13px;
  }
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
    margin-top: 24px;
  }
  .card {
    background: var(--paper);
    border: 1.5px solid var(--g300);
    border-radius: 8px;
    padding: 18px;
    display: flex;
    flex-direction: column;
    min-height: 184px;
  }
  .card:hover { border-color: var(--slate); }
  .kicker {
    font-family: var(--mono);
    font-size: 11px;
    color: var(--clay-d);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    margin-bottom: 10px;
  }
  .card h2 {
    font-family: var(--serif);
    font-size: 24px;
    font-weight: 500;
    letter-spacing: 0;
    margin: 0 0 4px;
  }
  .count {
    font-family: var(--mono);
    font-size: 12px;
    color: var(--g500);
    margin-bottom: 10px;
  }
  .desc {
    color: var(--g700);
    font-size: 14px;
    margin: 0 0 18px;
    flex: 1;
  }
  .actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
  .btn {
    border: 1.5px solid var(--g300);
    border-radius: 999px;
    color: var(--slate);
    background: transparent;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 7px;
    font-family: var(--mono);
    font-size: 12px;
    padding: 8px 13px;
    text-decoration: none;
  }
  .btn:hover { border-color: var(--slate); }
  .btn.primary {
    background: var(--slate);
    border-color: var(--slate);
    color: var(--ivory);
  }
  .btn.primary:hover { background: var(--g700); }
  .help {
    margin-top: 40px;
    border-top: 1.5px solid var(--g300);
    padding-top: 28px;
    display: grid;
    grid-template-columns: 220px 1fr;
    gap: 32px;
  }
  .help h2 {
    font-family: var(--serif);
    font-weight: 500;
    margin: 0;
    font-size: 24px;
  }
  .help ol {
    margin: 0;
    padding-left: 20px;
    color: var(--g700);
  }
  .help li { margin-bottom: 8px; }
  code {
    font-family: var(--mono);
    background: var(--g100);
    border: 1px solid var(--g200);
    border-radius: 4px;
    padding: 1px 5px;
  }
  footer {
    margin-top: 48px;
    padding-top: 22px;
    border-top: 1.5px solid var(--g300);
    color: var(--g500);
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    justify-content: space-between;
    font-size: 13px;
  }
  footer a { color: var(--g700); text-decoration-color: var(--oat); text-underline-offset: 3px; }
  @media (max-width: 720px) {
    .wrap { padding: 0 20px 72px; }
    header { padding-top: 48px; }
    .toolbar { align-items: flex-start; flex-direction: column; }
    .help { grid-template-columns: 1fr; gap: 14px; }
  }
</style>
</head>
<body>
<main class="wrap">
  <header>
    <div class="eyebrow">CalendarX</div>
    <h1>Market dates in <em>your calendar</em></h1>
    <p class="intro">
      Daily generated calendar subscriptions for earnings dates, starting with common US stock baskets.
      The first source is Nasdaq earnings calendar data, normalized into stable .ics feeds.
    </p>
  </header>

  <div class="toolbar">
    <div class="summary"><b>{{len .Feeds}}</b> public feeds &middot; updated {{.GeneratedAt}}</div>
    <div class="source">Source: Nasdaq earnings calendar</div>
  </div>

  <section class="grid" aria-label="Calendar feeds">
    {{range .Feeds}}
    <article class="card">
      <div class="kicker">{{.Group}}</div>
      <h2>{{.Title}}</h2>
      <div class="count">{{.Count}} events in the current window</div>
      <p class="desc">{{.Description}}</p>
      <div class="actions">
        <button class="btn primary" type="button" data-copy="ics/{{.Slug}}.ics">Copy link</button>
        <a class="btn" href="ics/{{.Slug}}.ics" download>Download</a>
      </div>
    </article>
    {{end}}
  </section>

  <section class="help">
    <h2>Subscribe</h2>
    <ol>
      <li>Copy a feed link from this page.</li>
      <li>In Apple Calendar, Google Calendar, or Outlook, add a calendar by URL.</li>
      <li>Keep the subscription enabled. CalendarX refreshes the feed from GitHub Pages after each CI run.</li>
    </ol>
  </section>

  <footer>
    <span>Generated by CalendarX</span>
    <a href="https://github.com/iFurySt/CalendarX">GitHub</a>
  </footer>
</main>
<script>
  document.querySelectorAll("[data-copy]").forEach((button) => {
    button.addEventListener("click", async () => {
      const url = new URL(button.dataset.copy, window.location.href).toString();
      await navigator.clipboard.writeText(url);
      const old = button.textContent;
      button.textContent = "Copied";
      window.setTimeout(() => { button.textContent = old; }, 1200);
    });
  });
</script>
</body>
</html>
`))
