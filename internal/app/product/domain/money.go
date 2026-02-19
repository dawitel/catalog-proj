package domain

import "math/big"

// Money is a value object for precise decimal amounts using *big.Rat (numerator/denominator).
type Money struct {
	rat *big.Rat
}

func NewMoney(numerator, denominator int64) *Money {
	if denominator == 0 {
		return nil
	}
	return &Money{rat: big.NewRat(numerator, denominator)}
}

func NewMoneyFromRat(r *big.Rat) *Money {
	if r == nil {
		return nil
	}
	return &Money{rat: new(big.Rat).Set(r)}
}

func (m *Money) Numerator() int64 {
	if m == nil || m.rat == nil {
		return 0
	}
	return m.rat.Num().Int64()
}

func (m *Money) Denominator() int64 {
	if m == nil || m.rat == nil {
		return 0
	}
	return m.rat.Denom().Int64()
}

func (m *Money) Rat() *big.Rat {
	if m == nil || m.rat == nil {
		return new(big.Rat)
	}
	return new(big.Rat).Set(m.rat)
}
