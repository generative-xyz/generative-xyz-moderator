package types

type Tax struct {
	// RateBeauty := taxRate + 1
	// example: if taxRate = 0.03 then RateBeauty := 1.03
	RateBeauty   *Price `json:"rateBeauty"`
	IsIncludeTax bool   `json:"isIncludeTax"` // IncludeTax at EU
}

func NewTax(taxRate float64) *Tax {
	result := &Tax{}
	result.setRateBeauty(taxRate)
	return result
}

// setRateBeauty set rate beauty := taxRate + 1
func (s *Tax) setRateBeauty(taxRate float64) {
	if taxRate <= 0.0 {
		return
	}

	// if hasTaxRate then include tax to final price
	s.IsIncludeTax = true

	// taxRateBeauty := taxRate + 1
	one := NewPrice(1)
	taxRateBeauty := NewPrice(taxRate)
	taxRateBeauty.Add(one)
	s.RateBeauty = taxRateBeauty
}

// PriceWithTax get price with tax
func (s *Tax) PriceWithTax(price *Price) *Price {
	if s != nil && price != nil {
		// price := price * RateBeauty
		if s.IsIncludeTax {
			return price.Copy().Multiply(s.RateBeauty)
		}

		return price.Copy()
	}
	return ZeroPriceP
}
