package calendarx

import "path/filepath"

func EarningsCacheFile(dataDir string, isoDate string) string {
	return filepath.Join(dataDir, "earnings", isoDate+".json")
}

func WatchlistFile(dataDir string, slug string) string {
	return filepath.Join(dataDir, "watchlists", slug+".json")
}

func ICSFile(outDir string, slug string) string {
	return filepath.Join(outDir, "ics", slug+".ics")
}

func SiteIndexFile(outDir string) string {
	return filepath.Join(outDir, "index.html")
}
