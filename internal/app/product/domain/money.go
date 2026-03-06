package domain

import "math/big"

type Money struct {
	amount *big.Rat
}

func NewMoney(numerator, denominator int64) (*Money, error) {
	if denominator == 0 {
		return nil, ErrInvalidPrice
	}

	amount := new(big.Rat).SetFrac64(numerator, denominator)
	if amount.Sign() < 0 {
		return nil, ErrInvalidPrice
	}
	return &Money{amount: amount}, nil
}

func NewMoneyFromRat(amount *big.Rat) (*Money, error) {
	if amount == nil || amount.Sign() < 0 {
		return nil, ErrInvalidPrice
	}
	return &Money{amount: new(big.Rat).Set(amount)}, nil
}

func (m *Money) Amount() *big.Rat {
	return new(big.Rat).Set(m.amount)
}

func (m *Money) Numerator() int64 {
	return m.amount.Num().Int64()
}

func (m *Money) Denominator() int64 {
	return m.amount.Denom().Int64()
}

func (m *Money) Subtract(other *Money) (*Money, error) {
	result := new(big.Rat).Sub(m.amount, other.amount)
	return NewMoneyFromRat(result)
}

func (m *Money) MultiplyByRat(r *big.Rat) (*Money, error) {
	result := new(big.Rat).Mul(m.amount, r)
	return NewMoneyFromRat(result)
}

func (m *Money) IsZero() bool {
	return m.amount.Sign() == 0
}

func (m *Money) Equal(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount.Cmp(other.amount) == 0
}
