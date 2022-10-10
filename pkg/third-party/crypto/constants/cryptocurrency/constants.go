package cryptocurrency

const (
	Solana   = "sol"
	Ethereum = "eth"
	USD      = "usd"
	USDT     = "usdt"
	USDC     = "usdc"
	DAI      = "dai"
)

var (
	CurrencySymbolByCurrency = map[string]string{
		Solana:   "https://s2.coinmarketcap.com/static/img/coins/64x64/5426.png",
		Ethereum: "https://s2.coinmarketcap.com/static/img/coins/64x64/1027.png",
		USD:      "https://s2.coinmarketcap.com/static/img/coins/64x64/1.png",
		USDT:     "https://s2.coinmarketcap.com/static/img/coins/64x64/825.png",
		USDC:     "https://s2.coinmarketcap.com/static/img/coins/64x64/3408.png",
		DAI:      "https://s2.coinmarketcap.com/static/img/coins/64x64/4943.png",
	}

	// CoingeckoCoinIdByCurrency get coin id by currency
	// ref: https://api.coingecko.com/api/v3/coins/list
	CoingeckoCoinIdByCurrency = map[string]string{
		Solana:   "solana",
		Ethereum: "ethereum",
		USDT:     "tether",
		USDC:     "usd-coin",
		DAI:      "dai",
	}

	RoundPlacesByCurrency = map[string]int{
		Solana:   SolonaRoundPlaces,
		Ethereum: EthereumRoundPlaces,
		USDT:     USDTRoundPlaces,
		USDC:     USDCRoundPlaces,
		DAI:      DAIRoundPlaces,
	}
)

const (
	SolonaRoundPlaces   = 9
	EthereumRoundPlaces = 18
	USDTRoundPlaces     = 6
	USDCRoundPlaces     = 6
	DAIRoundPlaces      = 18
	DefaultRoundPlaces  = 6
)
