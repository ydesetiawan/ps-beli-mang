package helper

import "github.com/shopspring/decimal"

func PointerBool(data bool) *bool {
	return &data
}

func PointerString(data string) *string {
	return &data
}

func PointerFloat64(data float64) *float64 {
	return &data
}

func PointerDecimal(data decimal.Decimal) *decimal.Decimal {
	return &data
}
