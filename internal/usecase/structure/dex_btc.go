package structure

type PSBTData struct {
}

type GMDashBoardPercent struct {
	PastUSDT           float64 `json:"past_usdt"`
	PastContributors   int64   `json:"past_contributors"`
	USDT               float64 `json:"usdt"`
	Contributor        int64   `json:"contributor"`
	PercentUSDT        float64 `json:"percent_usdt"`
	PercentContributor float64 `json:"percent_contributor"`
}
