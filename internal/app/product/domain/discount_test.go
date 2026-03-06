package domain_test

import (
	"CatalogService/internal/app/product/domain"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscount_NewDiscount(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	t.Run("valid discount", func(t *testing.T) {
		d, err := domain.NewDiscount(big.NewRat(20, 1), now, later)
		require.NoError(t, err)
		assert.Equal(t, "20", d.Percentage().RatString())
		assert.True(t, d.IsValidAt(now.Add(1*time.Hour)))
	})

	t.Run("invalid percentage", func(t *testing.T) {
		_, err := domain.NewDiscount(big.NewRat(101, 1), now, later)
		assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)

		_, err = domain.NewDiscount(big.NewRat(-1, 1), now, later)
		assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)

		_, err = domain.NewDiscount(big.NewRat(0, 1), now, later)
		assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)
	})

	t.Run("invalid period", func(t *testing.T) {
		_, err := domain.NewDiscount(big.NewRat(20, 1), later, now)
		assert.ErrorIs(t, err, domain.ErrInvalidDiscountPeriod)
	})
}

func TestDiscount_IsValidAt(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	d, _ := domain.NewDiscount(big.NewRat(10, 1), start, end)

	assert.False(t, d.IsValidAt(start.Add(-time.Second)))
	assert.True(t, d.IsValidAt(start))
	assert.True(t, d.IsValidAt(start.Add(time.Hour)))
	assert.False(t, d.IsValidAt(end))
}

func TestDiscount_FractionOff(t *testing.T) {
	d, _ := domain.NewDiscount(big.NewRat(25, 1), time.Now(), time.Now().Add(time.Hour))
	assert.Equal(t, "1/4", d.FractionOff().String())
}
