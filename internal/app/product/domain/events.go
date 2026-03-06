package domain

import "time"

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type ProductCreatedEvent struct {
	ProductID string
	Name      string
	Category  string
	Timestamp time.Time
}

func (e *ProductCreatedEvent) EventType() string     { return "product.created" }
func (e *ProductCreatedEvent) AggregateID() string   { return e.ProductID }
func (e *ProductCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type ProductUpdatedEvent struct {
	ProductID string
	Timestamp time.Time
}

func (e *ProductUpdatedEvent) EventType() string     { return "product.updated" }
func (e *ProductUpdatedEvent) AggregateID() string   { return e.ProductID }
func (e *ProductUpdatedEvent) OccurredAt() time.Time { return e.Timestamp }

type ProductActivatedEvent struct {
	ProductID string
	Timestamp time.Time
}

func (e *ProductActivatedEvent) EventType() string     { return "product.activated" }
func (e *ProductActivatedEvent) AggregateID() string   { return e.ProductID }
func (e *ProductActivatedEvent) OccurredAt() time.Time { return e.Timestamp }

type ProductDeactivatedEvent struct {
	ProductID string
	Timestamp time.Time
}

func (e *ProductDeactivatedEvent) EventType() string     { return "product.deactivated" }
func (e *ProductDeactivatedEvent) AggregateID() string   { return e.ProductID }
func (e *ProductDeactivatedEvent) OccurredAt() time.Time { return e.Timestamp }

type DiscountAppliedEvent struct {
	ProductID  string
	Percentage string
	StartDate  time.Time
	EndDate    time.Time
	Timestamp  time.Time
}

func (e *DiscountAppliedEvent) EventType() string     { return "discount.applied" }
func (e *DiscountAppliedEvent) AggregateID() string   { return e.ProductID }
func (e *DiscountAppliedEvent) OccurredAt() time.Time { return e.Timestamp }

type DiscountRemovedEvent struct {
	ProductID string
	Timestamp time.Time
}

func (e *DiscountRemovedEvent) EventType() string     { return "discount.removed" }
func (e *DiscountRemovedEvent) AggregateID() string   { return e.ProductID }
func (e *DiscountRemovedEvent) OccurredAt() time.Time { return e.Timestamp }
