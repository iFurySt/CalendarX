package calendarx

type NasdaqDay struct {
	Date string      `json:"date"`
	Data []NasdaqRow `json:"data"`
}

type NasdaqRow struct {
	Symbol              string `json:"symbol"`
	Name                string `json:"name"`
	MarketCap           string `json:"marketCap"`
	FiscalQuarterEnding string `json:"fiscalQuarterEnding"`
	Time                string `json:"time"`
	EPSForecast         string `json:"epsForecast"`
	NoOfEsts            string `json:"noOfEsts"`
}

type WatchlistEntry struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry"`
}

type Event struct {
	Symbol              string
	CompanyName         string
	Industry            string
	Date                string
	MarketCap           string
	FiscalQuarterEnding string
	Time                string
	EPSForecast         string
	NoOfEsts            string
}

type Feed struct {
	Slug        string
	Title       string
	Description string
	Group       string
	Events      []Event
}

type FeedSummary struct {
	Slug        string
	Title       string
	Description string
	Group       string
	Count       int
}
