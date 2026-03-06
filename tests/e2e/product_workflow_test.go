package e2e_test

import (
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/app/product/usecases/activate_product"
	"CatalogService/internal/app/product/usecases/apply_discount"
	"CatalogService/internal/app/product/usecases/create_product"
	"CatalogService/internal/pkg/clock"
	"CatalogService/internal/pkg/committer"
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
func (m *MockRepo) InsertMut(p *domain.Product) committer.Mutation {
	args := m.Called(p)
	return args.Get(0)
}
func (m *MockRepo) UpdateMut(p *domain.Product) committer.Mutation {
	args := m.Called(p)
	return args.Get(0)
}

type MockOutbox struct {
	mock.Mock
}

func (m *MockOutbox) InsertEventMut(event domain.DomainEvent, now time.Time) committer.Mutation {
	args := m.Called(event, now)
	return args.Get(0)
}

type MockCommitter struct {
	mock.Mock
}

func (m *MockCommitter) Apply(ctx context.Context, plan *committer.Plan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

func TestProductWorkflowE2E(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clk := clock.FixedClock{FixedTime: now}

	repo := new(MockRepo)
	outbox := new(MockOutbox)
	committerMock := new(MockCommitter)

	createUC := create_product.New(repo, outbox, committerMock, clk)
	activateUC := activate_product.NewActivate(repo, outbox, committerMock, clk)
	applyDiscUC := apply_discount.NewApply(repo, outbox, committerMock, clk)

	var productID string

	t.Run("1. Create Product", func(t *testing.T) {
		repo.On("InsertMut", mock.Anything).Return("insert-mut")
		outbox.On("InsertEventMut", mock.Anything, now).Return("outbox-mut")
		committerMock.On("Apply", ctx, mock.MatchedBy(func(p *committer.Plan) bool {
			return len(p.Mutations()) == 2 // 1 insert + 1 outbox
		})).Return(nil)

		id, err := createUC.Execute(ctx, create_product.Request{
			Name: "Prod", Category: "Cat", PriceNum: 100, PriceDenom: 1,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, id)
		productID = id
	})

	t.Run("2. Activate Product", func(t *testing.T) {
		// Mock loading product from repo (as draft)
		price, _ := domain.NewMoney(100, 1)
		p, _ := domain.NewProduct(productID, "Prod", "D", "Cat", price, now)
		repo.On("GetByID", ctx, productID).Return(p, nil).Once()

		repo.On("UpdateMut", mock.MatchedBy(func(p *domain.Product) bool {
			return p.Status() == domain.ProductStatusActive
		})).Return("update-mut").Once()
		outbox.On("InsertEventMut", mock.Anything, now).Return("outbox-mut").Once()
		committerMock.On("Apply", ctx, mock.Anything).Return(nil).Once()

		err := activateUC.Execute(ctx, activate_product.Request{ProductID: productID})
		require.NoError(t, err)
	})

	t.Run("3. Apply Discount", func(t *testing.T) {
		// Mock loading product from repo (as active)
		price, _ := domain.NewMoney(100, 1)
		p, _ := domain.NewProduct(productID, "Prod", "D", "Cat", price, now)
		_ = p.Activate(now)
		p.DomainEvents() // clear
		repo.On("GetByID", ctx, productID).Return(p, nil).Once()

		repo.On("UpdateMut", mock.MatchedBy(func(p *domain.Product) bool {
			return p.Discount() != nil
		})).Return("update-mut").Once()
		outbox.On("InsertEventMut", mock.Anything, now).Return("outbox-mut").Once()
		committerMock.On("Apply", ctx, mock.Anything).Return(nil).Once()

		err := applyDiscUC.Execute(ctx, apply_discount.ApplyRequest{
			ProductID: productID, Percentage: big.NewRat(20, 1),
			StartDate: now, EndDate: now.Add(time.Hour),
		})
		require.NoError(t, err)
	})
}
