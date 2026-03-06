package domain_test

import (
	"CatalogService/internal/app/product/domain"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoney_NewMoney(t *testing.T) {
	t.Run("valid positive price", func(t *testing.T) {
		m, err := domain.NewMoney(1999, 100)
		require.NoError(t, err)
		assert.Equal(t, "19.99", m.Amount().FloatString(2))
		assert.Equal(t, int64(1999), m.Numerator())
		assert.Equal(t, int64(100), m.Denominator())
	})

	t.Run("zero price is valid", func(t *testing.T) {
		m, err := domain.NewMoney(0, 1)
		require.NoError(t, err)
		assert.True(t, m.IsZero())
	})

	t.Run("negative price is invalid", func(t *testing.T) {
		_, err := domain.NewMoney(-1, 1)
		assert.ErrorIs(t, err, domain.ErrInvalidPrice)
	})

	t.Run("zero denominator is invalid", func(t *testing.T) {
		_, err := domain.NewMoney(10, 0)
		assert.ErrorIs(t, err, domain.ErrInvalidPrice)
	})
}

func TestMoney_Arithmetic(t *testing.T) {
	m1, _ := domain.NewMoney(100, 1)
	m2, _ := domain.NewMoney(30, 1)

	t.Run("subtract", func(t *testing.T) {
		res, err := m1.Subtract(m2)
		require.NoError(t, err)
		assert.Equal(t, "70", res.Amount().RatString())
	})

	t.Run("subtract to negative is invalid", func(t *testing.T) {
		_, err := m2.Subtract(m1)
		assert.ErrorIs(t, err, domain.ErrInvalidPrice)
	})

	t.Run("multiply by rat", func(t *testing.T) {
		// 100 * 0.2 = 20
		res, err := m1.MultiplyByRat(big.NewRat(20, 100))
		require.NoError(t, err)
		assert.Equal(t, "20", res.Amount().RatString())
	})
}

func TestMoney_Equal(t *testing.T) {
	m1, _ := domain.NewMoney(100, 1)
	m2, _ := domain.NewMoney(200, 2)
	m3, _ := domain.NewMoney(101, 1)

	assert.True(t, m1.Equal(m2))
	assert.False(t, m1.Equal(m3))
	assert.False(t, m1.Equal(nil))
}
