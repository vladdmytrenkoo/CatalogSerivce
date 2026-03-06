package domain

import "time"

type ProductStatus string

type Product struct {
	id          string
	name        string
	description string
	category    string
	basePrice   *Money
	discount    *Discount
	status      ProductStatus
	createdAt   time.Time
	updatedAt   time.Time
	archivedAt  *time.Time
	changes     *ChangeTracker
	events      []DomainEvent
}

func NewProduct(id, name, description, category string, basePrice *Money, now time.Time) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidProductName
	}
	if category == "" {
		return nil, ErrInvalidCategory
	}
	if basePrice == nil || basePrice.IsZero() {
		return nil, ErrInvalidPrice
	}

	p := &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		status:      ProductStatusDraft,
		createdAt:   now,
		updatedAt:   now,
		changes:     NewChangeTracker(),
	}

	p.events = append(p.events, &ProductCreatedEvent{
		ProductID: id,
		Name:      name,
		Category:  category,
		Timestamp: now,
	})

	return p, nil
}

func Hydrate(
	id, name, description, category string,
	basePrice *Money,
	discount *Discount,
	status ProductStatus,
	createdAt, updatedAt time.Time,
	archivedAt *time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		discount:    discount,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		archivedAt:  archivedAt,
		changes:     NewChangeTracker(),
	}
}

func (p *Product) ID() string              { return p.id }
func (p *Product) Name() string            { return p.name }
func (p *Product) Description() string     { return p.description }
func (p *Product) Category() string        { return p.category }
func (p *Product) BasePrice() *Money       { return p.basePrice }
func (p *Product) Discount() *Discount     { return p.discount }
func (p *Product) Status() ProductStatus   { return p.status }
func (p *Product) CreatedAt() time.Time    { return p.createdAt }
func (p *Product) UpdatedAt() time.Time    { return p.updatedAt }
func (p *Product) ArchivedAt() *time.Time  { return p.archivedAt }
func (p *Product) Changes() *ChangeTracker { return p.changes }

func (p *Product) DomainEvents() []DomainEvent {
	events := p.events
	p.events = nil
	return events
}

func (p *Product) Update(name, description, category string, now time.Time) error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}

	if name != "" && name != p.name {
		p.name = name
		p.changes.MarkDirty(FieldName)
	}

	if description != "" && description != p.description {
		p.description = description
		p.changes.MarkDirty(FieldDescription)
	}

	if category != "" && category != p.category {
		p.category = category
		p.changes.MarkDirty(FieldCategory)
	}

	if p.changes.HasChanges() {
		p.updatedAt = now
		p.events = append(p.events, &ProductUpdatedEvent{
			ProductID: p.id,
			Timestamp: now,
		})
	}

	return nil
}

func (p *Product) Activate(now time.Time) error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	if p.status == ProductStatusActive {
		return ErrProductAlreadyActive
	}

	p.status = ProductStatusActive
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)

	p.events = append(p.events, &ProductActivatedEvent{
		ProductID: p.id,
		Timestamp: now,
	})
	return nil
}

func (p *Product) Deactivate(now time.Time) error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}

	p.status = ProductStatusInactive
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)

	p.events = append(p.events, &ProductDeactivatedEvent{
		ProductID: p.id,
		Timestamp: now,
	})
	return nil
}

func (p *Product) Archive(now time.Time) error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}

	p.status = ProductStatusArchived
	p.archivedAt = &now
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldArchivedAt)
	return nil
}

func (p *Product) ApplyDiscount(discount *Discount, now time.Time) error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}

	if !discount.IsValidAt(now) {
		return ErrInvalidDiscountPeriod
	}

	if p.discount != nil && p.discount.IsValidAt(now) {
		return ErrDiscountAlreadyActive
	}

	p.discount = discount
	p.updatedAt = now
	p.changes.MarkDirty(FieldDiscount)

	p.events = append(p.events, &DiscountAppliedEvent{
		ProductID:  p.id,
		Percentage: discount.Percentage().RatString(),
		StartDate:  discount.StartDate(),
		EndDate:    discount.EndDate(),
		Timestamp:  now,
	})
	return nil
}

func (p *Product) RemoveDiscount(now time.Time) error {
	if p.discount == nil {
		return ErrNoActiveDiscount
	}
	p.discount = nil
	p.updatedAt = now
	p.changes.MarkDirty(FieldDiscount)

	p.events = append(p.events, &DiscountRemovedEvent{
		ProductID: p.id,
		Timestamp: now,
	})
	return nil
}
