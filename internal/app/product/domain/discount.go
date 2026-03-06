package domain

import (
	"math/big"
	"time"
)

type Discount struct {
	percentage *big.Rat
	startDate  time.Time
	endDate    time.Time
}

func NewDiscount(percentage *big.Rat, startDate, endDate time.Time) (*Discount, error) {
	zero := new(big.Rat)
	hundred := new(big.Rat).SetInt64(100)

	if percentage.Cmp(zero) <= 0 || percentage.Cmp(hundred) >= 0 {
		return nil, ErrInvalidDiscountPercent
	}

	if !endDate.After(startDate) {
		return nil, ErrInvalidDiscountPeriod
	}

	return &Discount{
		percentage: new(big.Rat).Set(percentage),
		startDate:  startDate,
		endDate:    endDate,
	}, nil
}

func (d *Discount) Percentage() *big.Rat {
	return new(big.Rat).Set(d.percentage)
}

func (d *Discount) StartDate() time.Time {
	return d.startDate
}

func (d *Discount) EndDate() time.Time {
	return d.endDate
}

func (d *Discount) IsValidAt(now time.Time) bool {
	return !now.Before(d.startDate) && now.Before(d.endDate)
}

func (d *Discount) FractionOff() *big.Rat {
	return new(big.Rat).Mul(d.percentage, new(big.Rat).SetFrac64(1, 100))
}
