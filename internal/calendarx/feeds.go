package calendarx

var BasketFeeds = []FeedSummary{
	{
		Slug:        "mega7",
		Title:       "Mega 7",
		Description: "Apple, Microsoft, Nvidia, Amazon, Alphabet, Meta, and Tesla.",
		Group:       "Stock baskets",
	},
	{
		Slug:        "nasdaq100",
		Title:       "Nasdaq-100",
		Description: "Earnings dates for Nasdaq-100 constituents.",
		Group:       "Stock baskets",
	},
	{
		Slug:        "sp500",
		Title:       "S&P 500",
		Description: "Earnings dates for S&P 500 constituents.",
		Group:       "Stock baskets",
	},
	{
		Slug:        "dow30",
		Title:       "Dow 30",
		Description: "Earnings dates for Dow Jones Industrial Average constituents.",
		Group:       "Stock baskets",
	},
}

var AllEarningsFeed = FeedSummary{
	Slug:        "all",
	Title:       "All Earnings",
	Description: "Every company returned by the Nasdaq earnings calendar window.",
	Group:       "Full feed",
}
