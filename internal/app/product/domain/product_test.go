package domain_test

import (
	"CatalogService/internal/app/product/domain"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct_Lifecycle(t *testing.T) {
	now := time.Now()
	price, _ := domain.NewMoney(1000, 100) // 10.00
	id := "prod-123"

	t.Run("create product", func(t *testing.T) {
		p, err := domain.NewProduct(id, "Test Product", "Desc", "Cat", price, now)
		require.NoError(t, err)
		assert.Equal(t, domain.ProductStatusDraft, p.Status())
		assert.Equal(t, id, p.ID())

		events := p.DomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "product.created", events[0].EventType())
	})

	t.Run("activate product", func(t *testing.T) {
		p, _ := domain.NewProduct(id, "Test", "Desc", "Cat", price, now)
		p.DomainEvents() // clear creation event

		err := p.Activate(now.Add(time.Hour))
		require.NoError(t, err)
		assert.Equal(t, domain.ProductStatusActive, p.Status())
		assert.True(t, p.Changes().IsDirty(domain.FieldStatus))

		events := p.DomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "product.activated", events[0].EventType())
	})

	t.Run("apply discount to active product", func(t *testing.T) {
		p, _ := domain.NewProduct(id, "Test", "Desc", "Cat", price, now)
		_ = p.Activate(now)
		p.DomainEvents()

		discount, _ := domain.NewDiscount(big.NewRat(20, 1), now, now.Add(time.Hour))
		err := p.ApplyDiscount(discount, now)
		require.NoError(t, err)
		assert.NotNil(t, p.Discount())
		assert.True(t, p.Changes().IsDirty(domain.FieldDiscount))

		events := p.DomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "discount.applied", events[0].EventType())
	})

	t.Run("cannot apply discount to draft product", func(t *testing.T) {
		p, _ := domain.NewProduct(id, "Test", "Desc", "Cat", price, now)
		discount, _ := domain.NewDiscount(big.NewRat(20, 1), now, now.Add(time.Hour))
		err := p.ApplyDiscount(discount, now)
		assert.ErrorIs(t, err, domain.ErrProductNotActive)
	})
}

func TestProduct_Update(t *testing.T) {
	now := time.Now()
	price, _ := domain.NewMoney(100, 1)
	p, _ := domain.NewProduct("1", "Old Name", "Old Desc", "Old Cat", price, now)
	p.DomainEvents() // clear

	t.Run("update name", func(t *testing.T) {
		err := p.Update("New Name", "Old Desc", "Old Cat", now.Add(time.Hour))
		require.NoError(t, err)
		assert.Equal(t, "New Name", p.Name())
		assert.True(t, p.Changes().IsDirty(domain.FieldName))
		assert.False(t, p.Changes().IsDirty(domain.FieldDescription))

		events := p.DomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "product.updated", events[0].EventType())
	})

	t.Run("no update if same values", func(t *testing.T) {
		p, _ := domain.NewProduct("1", "Name", "Desc", "Cat", price, now)
		p.DomainEvents()
		err := p.Update("Name", "Desc", "Cat", now.Add(time.Hour))
		require.NoError(t, err)
		assert.False(t, p.Changes().HasChanges())
		assert.Len(t, p.DomainEvents(), 0)
	})

	t.Run("ignore empty fields on update", func(t *testing.T) {
		p, _ := domain.NewProduct("1", "Name", "Desc", "Cat", price, now)
		p.DomainEvents()

		// Update with empty strings
		err := p.Update("", "", "", now.Add(time.Hour))
		require.NoError(t, err)

		// Assert values remain unchanged
		assert.Equal(t, "Name", p.Name())
		assert.Equal(t, "Desc", p.Description())
		assert.Equal(t, "Cat", p.Category())

		// Assert no changes tracked
		assert.False(t, p.Changes().HasChanges())
		assert.Len(t, p.DomainEvents(), 0)
	})
}
