package types

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
)

const (
	DefaultRoundOnFloat     = 0.6
	DefaultRoundPlacesFloat = 2
)

var (
	ZeroPrice       = *NewPrice(0.0)
	ZeroPriceP      = NewPrice(0.0)
	ZeroPriceString = "0"
)

type Price struct {
	float64Value    float64
	precision       uint64
	cryptoPrecision uint64 // crypto precision

	needBeautyPrecision bool
}

func NewPrice(value float64, opts ...Option) *Price {
	price := &Price{
		float64Value:        value,
		precision:           DefaultRoundPlacesFloat,
		cryptoPrecision:     DefaultRoundPlacesFloat,
		needBeautyPrecision: true,
	}

	for _, opt := range opts {
		opt(price)
	}

	return price
}

func NewPriceNonePointer(value float64, opts ...Option) Price {
	return *NewPrice(value, opts...)
}

func (a Price) MarshalJSON() ([]byte, error) {
	str := a.ToString()
	return json.Marshal(str)
}

func (a *Price) UnmarshalJSON(data []byte) error {
	var stringValue string
	var floatValue float64
	err := json.Unmarshal(data, &stringValue)
	if err != nil {
		return err
	}

	// remove ',' from string value
	stringValue = strings.Replace(stringValue, ",", "", -1)
	floatValue, err = strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return err
	}

	a.float64Value = floatValue
	a.precision = getPrecisionFromStringValue(stringValue)
	a.needBeautyPrecision = a.precision == 0
	return nil
}

func (a *Price) Add(value *Price) *Price {
	if a == nil {
		return nil
	}
	newPrice := *a
	newPrice.float64Value += value.float64Value
	*a = newPrice
	return &newPrice
}

func (a *Price) Multiply(value *Price) *Price {
	if a == nil {
		return nil
	}
	newPrice := *a
	newPrice.float64Value *= value.float64Value
	*a = newPrice
	return &newPrice
}

func (a *Price) MultiplyFloat(value float64) *Price {
	if a == nil {
		return nil
	}
	newPrice := *a
	newPrice.float64Value *= value
	*a = newPrice
	return &newPrice
}

func (a *Price) Minus(value *Price) *Price {
	newPrice := *a
	newPrice.float64Value -= value.float64Value
	*a = newPrice
	return &newPrice
}

func (a *Price) Copy() *Price {
	if a == nil {
		return nil
	}
	newPrice := *a
	return &newPrice
}

func (a Price) ToStringWithCurrencySymbol(currencySymbol string) string {
	str := a.ToString()
	return fmt.Sprintf("%s%s", currencySymbol, str)
}

func (a Price) ToString() string {
	return a.RoundPriceStringFormatAtCheckout(DefaultRoundOnFloat, int(a.precision))
}

func (a Price) ToFloat64() float64 {
	return a.float64Value
}

func (a Price) ToCryptoAmount() string {
	float64ValueIsRoundedWithPrecision := a.round(a.ToFloat64(), DefaultRoundOnFloat, int(a.precision))
	float64Value := float64ValueIsRoundedWithPrecision * math.Pow10(int(a.cryptoPrecision))
	amount := a.round(float64Value, DefaultRoundOnFloat, 0)
	cryptoAmount := fmt.Sprintf("%.f", amount)
	return cryptoAmount
}

func (a Price) GetPrecision() uint64 {
	return a.precision
}

func (a Price) GetCryptoPrecision() uint64 {
	return a.cryptoPrecision
}

func (a Price) ToAdyenString(adyenInt *int) string {
	if adyenInt == nil {
		temp := a.ToAdyenInt()
		adyenInt = &temp
	}
	return fmt.Sprintf("%d", adyenInt)
}

func (a Price) ToAdyenInt() int {
	amountMulti100 := a.float64Value * 100
	amount := a.round(amountMulti100, DefaultRoundOnFloat, DefaultRoundPlacesFloat)
	adyenAmount := int(amount)
	return adyenAmount
}

func (a Price) round(val, roundOn float64, precision int) float64 {
	var round float64
	pow := math.Pow(10, float64(precision))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func (a Price) RoundPriceStringFormatAtCheckout(roundOn float64, precision int) string {
	return a.priceStringFormatAtCheckout(a.round(a.ToFloat64(), roundOn, precision), precision)
}

func (a Price) priceStringFormatAtCheckout(val float64, precision int) string {
	ac := accounting.Accounting{Precision: precision}
	price := ac.FormatMoney(val)

	if a.needBeautyPrecision {
		return a.beautyMoney(price, precision)
	}

	return price
}

// beautyMoney format 2.00 to 2
func (a Price) beautyMoney(price string, precision int) string {
	return strings.Replace(price, a.genAllZeroPrecision(precision), "", 1)
}

func (a Price) genAllZeroPrecision(precision int) string {
	return "." + strings.Repeat("0", precision)
}

// getPrecisionFromStringValue get precision from strValue 123.00
func getPrecisionFromStringValue(str string) uint64 {
	strSplit := strings.Split(str, ".")
	if len(strSplit) == 2 {
		return uint64(len(strSplit[1]))
	}

	return 0
}

func FormatNumber(n float64, precision int) string {
	ac := accounting.Accounting{Precision: precision}
	priceStr := ac.FormatMoney(n)
	return priceStr
}
