package calculate

import "github.com/shopspring/decimal"


var oneHundred decimal.Decimal = decimal.NewFromInt(100)
var onethousand decimal.Decimal = decimal.NewFromInt(1000)

func Fen2Yuan(fen int64) float64 {
	y, _ := decimal.NewFromInt(fen).Div(oneHundred).Truncate(2).Float64()
	return y
}


func Yuan2Fen(yuan float64) int64 {
	f, _ := decimal.NewFromFloat(yuan).Mul(oneHundred).Float64()
	return int64(f)
}


func  Ge2Qian(Ge int64) float64 {
	y, _ := decimal.NewFromInt(Ge).Div(onethousand).Truncate(2).Float64()
	return y
}


func Qian2Ge(Qian float64) int64 {
	f, _ := decimal.NewFromFloat(Qian).Mul(onethousand).Float64()
	return int64(f)
}
