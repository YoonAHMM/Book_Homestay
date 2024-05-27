package calculate

import "github.com/shopspring/decimal"


var oneHundred decimal.Decimal = decimal.NewFromInt(100)


func Fen2Yuan(fen int64) float64 {
	y, _ := decimal.NewFromInt(fen).Div(oneHundred).Truncate(2).Float64()
	return y
}


func Yuan2Fen(yuan float64) int64 {
	f, _ := decimal.NewFromFloat(yuan).Mul(oneHundred).Float64()
	return int64(f)

}