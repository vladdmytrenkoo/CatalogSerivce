package services_test

import (
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/app/product/domain/services"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPricingCalculator_EffectivePrice(t *testing.T) {
	pc := services.NewPricingCalculator()
	now := time.Now()
	price, _ := domain.NewMoney(1000, 100) // 10.00
	id := "p1"
	p, _ := domain.NewProduct(id, "P1", "D", "C", price, now)
	_ = p.Activate(now)

	t.Run("no discount", func(t *testing.T) {
		eff := pc.EffectivePrice(p, now)
		assert.Equal(t, "10", eff.RatString())
	})

	t.Run("active discount 20%", func(t *testing.T) {
		discount, _ := domain.NewDiscount(big.NewRat(20, 1), now, now.Add(time.Hour))
		_ = p.ApplyDiscount(discount, now)

		eff := pc.EffectivePrice(p, now)
		// 10 - 20% = 8
		assert.Equal(t, "8", eff.RatString())
	})

	t.Run("discount not yet started", func(t *testing.T) {
		p, _ := domain.NewProduct(id, "P1", "D", "C", price, now)
		_ = p.Activate(now)
		discount, _ := domain.NewDiscount(big.NewRat(50, 1), now.Add(time.Hour), now.Add(2*time.Hour))
		_ = p.ApplyDiscount(discount, now)

		eff := pc.EffectivePrice(p, now)
		assert.Equal(t, "10", eff.RatString())
	})

	t.Run("discount expired", func(t *testing.T) {
		p, _ := domain.NewProduct(id, "P1", "D", "C", price, now)
		_ = p.Activate(now)
		discount, _ := domain.NewDiscount(big.NewRat(50, 1), now.Add(-2*time.Hour), now.Add(-time.Hour))
		_ = p.ApplyDiscount(discount, now)

		eff := pc.EffectivePrice(p, now)
		assert.Equal(t, "10", eff.RatString())
	})
}
