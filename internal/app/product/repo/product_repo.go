package repo

import (
	"context"
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
	committer "github.com/vladdmytrenkoo/commiter"
	"google.golang.org/api/iterator"

	"CatalogService/internal/app/product/contracts"
	"CatalogService/internal/app/product/domain"
	"CatalogService/internal/models/m_product"
)

var _ contracts.ProductRepository = (*ProductRepo)(nil)

type ProductRepo struct {
	client *spanner.Client
}

func NewProductRepo(client *spanner.Client) *ProductRepo {
	return &ProductRepo{client: client}
}

func (r *ProductRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	row, err := r.client.Single().ReadRow(ctx, m_product.Table, spanner.Key{id}, m_product.AllColumns())
	if err != nil {
		if spanner.ErrCode(err) == 5 {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return r.rowToProduct(row)
}

func (r *ProductRepo) InsertMut(p *domain.Product) committer.Mutation {
	data := r.toData(p)
	return m_product.InsertMut(data)
}

func (r *ProductRepo) UpdateMut(p *domain.Product) committer.Mutation {
	changes := p.Changes()
	if !changes.HasChanges() {
		return nil
	}

	cols := make(map[string]interface{})

	if changes.IsDirty(domain.FieldName) {
		cols[m_product.Name] = p.Name()
	}
	if changes.IsDirty(domain.FieldDescription) {
		cols[m_product.Description] = p.Description()
	}
	if changes.IsDirty(domain.FieldCategory) {
		cols[m_product.Category] = p.Category()
	}
	if changes.IsDirty(domain.FieldStatus) {
		cols[m_product.Status] = string(p.Status())
	}
	if changes.IsDirty(domain.FieldDiscount) {
		if d := p.Discount(); d != nil {
			cols[m_product.DiscountPercent] = d.Percentage()
			sd := d.StartDate()
			ed := d.EndDate()
			cols[m_product.DiscountStartDate] = &sd
			cols[m_product.DiscountEndDate] = &ed
		} else {
			cols[m_product.DiscountPercent] = (*big.Rat)(nil)
			cols[m_product.DiscountStartDate] = (*time.Time)(nil)
			cols[m_product.DiscountEndDate] = (*time.Time)(nil)
		}
	}
	if changes.IsDirty(domain.FieldArchivedAt) {
		cols[m_product.ArchivedAt] = p.ArchivedAt()
	}

	if len(cols) == 0 {
		return nil
	}

	cols[m_product.UpdatedAt] = p.UpdatedAt()
	return m_product.UpdateMut(p.ID(), cols)
}

func (r *ProductRepo) toData(p *domain.Product) *m_product.Data {
	d := &m_product.Data{
		ProductID:            p.ID(),
		Name:                 p.Name(),
		Category:             p.Category(),
		BasePriceNumerator:   p.BasePrice().Numerator(),
		BasePriceDenominator: p.BasePrice().Denominator(),
		Status:               string(p.Status()),
		CreatedAt:            p.CreatedAt(),
		UpdatedAt:            p.UpdatedAt(),
		ArchivedAt:           p.ArchivedAt(),
	}

	desc := p.Description()
	if desc != "" {
		d.Description = &desc
	}

	if disc := p.Discount(); disc != nil {
		pct := disc.Percentage()
		sd := disc.StartDate()
		ed := disc.EndDate()
		d.DiscountPercent = pct
		d.DiscountStartDate = &sd
		d.DiscountEndDate = &ed
	}

	return d
}

func (r *ProductRepo) rowToProduct(row *spanner.Row) (*domain.Product, error) {
	var (
		id          string
		name        string
		description spanner.NullString
		category    string
		priceNum    int64
		priceDenom  int64
		discPct     spanner.NullNumeric
		discStart   spanner.NullTime
		discEnd     spanner.NullTime
		status      string
		createdAt   time.Time
		updatedAt   time.Time
		archivedAt  spanner.NullTime
	)

	if err := row.Columns(
		&id, &name, &description, &category,
		&priceNum, &priceDenom,
		&discPct, &discStart, &discEnd,
		&status, &createdAt, &updatedAt, &archivedAt,
	); err != nil {
		return nil, err
	}

	basePrice, err := domain.NewMoney(priceNum, priceDenom)
	if err != nil {
		return nil, err
	}

	var discount *domain.Discount
	if discPct.Valid && discStart.Valid && discEnd.Valid {
		discount, err = domain.NewDiscount(&discPct.Numeric, discStart.Time, discEnd.Time)
		if err != nil {
			return nil, err
		}
	}

	desc := ""
	if description.Valid {
		desc = description.StringVal
	}

	var archived *time.Time
	if archivedAt.Valid {
		archived = &archivedAt.Time
	}

	return domain.Hydrate(
		id, name, desc, category,
		basePrice, discount,
		domain.ProductStatus(status),
		createdAt, updatedAt, archived,
	), nil
}

var _ contracts.ProductReadModel = (*ProductReadRepo)(nil)

type ProductReadRepo struct {
	client *spanner.Client
}

func NewProductReadRepo(client *spanner.Client) *ProductReadRepo {
	return &ProductReadRepo{client: client}
}

func (r *ProductReadRepo) GetByID(ctx context.Context, id string) (*contracts.ProductView, error) {
	row, err := r.client.Single().ReadRow(ctx, m_product.Table, spanner.Key{id}, m_product.AllColumns())
	if err != nil {
		if spanner.ErrCode(err) == 5 {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return r.rowToView(row)
}

func (r *ProductReadRepo) ListActive(ctx context.Context, filter contracts.ListFilter) (*contracts.ProductListResult, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	query := "SELECT " + colList() + " FROM " + m_product.Table +
		" WHERE status = @status"
	params := map[string]interface{}{
		"status": string(domain.ProductStatusActive),
	}

	if filter.Category != "" {
		query += " AND category = @category"
		params["category"] = filter.Category
	}

	query += " ORDER BY created_at DESC"
	query += " LIMIT @limit"
	params["limit"] = int64(pageSize + 1)

	if filter.PageToken != "" {
		query += " OFFSET @offset"
		params["offset"] = pageTokenToOffset(filter.PageToken)
	}

	stmt := spanner.Statement{SQL: query, Params: params}
	iter := r.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var views []*contracts.ProductView
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		v, err := r.rowToView(row)
		if err != nil {
			return nil, err
		}
		views = append(views, v)
	}

	result := &contracts.ProductListResult{}
	if len(views) > int(pageSize) {
		result.Products = views[:pageSize]
		result.NextPageToken = offsetToPageToken(currentOffset(filter.PageToken) + int64(pageSize))
	} else {
		result.Products = views
	}

	return result, nil
}

func (r *ProductReadRepo) rowToView(row *spanner.Row) (*contracts.ProductView, error) {
	var (
		id         string
		name       string
		desc       spanner.NullString
		category   string
		priceNum   int64
		priceDenom int64
		discPct    spanner.NullNumeric
		discStart  spanner.NullTime
		discEnd    spanner.NullTime
		status     string
		createdAt  time.Time
		updatedAt  time.Time
		ntime      spanner.NullTime
	)

	if err := row.Columns(
		&id, &name, &desc, &category,
		&priceNum, &priceDenom,
		&discPct, &discStart, &discEnd,
		&status, &createdAt, &updatedAt, &ntime,
	); err != nil {
		return nil, err
	}

	base := new(big.Rat).SetFrac64(priceNum, priceDenom)
	effective := new(big.Rat).Set(base)

	var discountPct *big.Rat
	if discPct.Valid && discStart.Valid && discEnd.Valid {
		now := time.Now().UTC()
		if !now.Before(discStart.Time) && now.Before(discEnd.Time) {
			discountPct = &discPct.Numeric
			fraction := new(big.Rat).Mul(&discPct.Numeric, new(big.Rat).SetFrac64(1, 100))
			discountAmt := new(big.Rat).Mul(base, fraction)
			effective = new(big.Rat).Sub(base, discountAmt)
		}
	}

	description := ""
	if desc.Valid {
		description = desc.StringVal
	}

	return &contracts.ProductView{
		ID:              id,
		Name:            name,
		Description:     description,
		Category:        category,
		BasePrice:       base,
		EffectivePrice:  effective,
		DiscountPercent: discountPct,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func colList() string {
	cols := m_product.AllColumns()
	s := cols[0]
	for _, c := range cols[1:] {
		s += ", " + c
	}
	return s
}

func pageTokenToOffset(token string) int64 {
	return currentOffset(token)
}

func currentOffset(token string) int64 {
	if token == "" {
		return 0
	}
	var n int64
	for _, c := range token {
		n = n*10 + int64(c-'0')
	}
	return n
}

func offsetToPageToken(offset int64) string {
	if offset == 0 {
		return ""
	}
	return big.NewInt(offset).String()
}
