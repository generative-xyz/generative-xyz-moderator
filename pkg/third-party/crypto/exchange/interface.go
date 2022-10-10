package exchange

type Client interface {
	GetRate(currency string) (float64, error)
}
