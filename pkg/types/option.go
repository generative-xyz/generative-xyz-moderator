package types

type Option func(p *Price)

func WithPrecision(precision int) Option {
	return func(p *Price) {
		p.precision = uint64(precision)
	}
}

func WithCryptoPrecision(precision int) Option {
	return func(p *Price) {
		p.cryptoPrecision = uint64(precision)
	}
}

func WithNeedBeautyZeroPrecision(b bool) Option {
	return func(p *Price) {
		p.needBeautyPrecision = b
	}
}
