package domain

import "time"

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

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
	archivedAt  time.Time
	changes     *ChangeTracker
	events      []DomainEvent
}

func NewProduct(id, name, description, category string, basePrice *Money, now time.Time) *Product {
	if basePrice == nil {
		return nil
	}
	changes := NewChangeTracker()
	changes.MarkDirty(FieldName)
	changes.MarkDirty(FieldDescription)
	changes.MarkDirty(FieldCategory)
	changes.MarkDirty(FieldBasePrice)
	changes.MarkDirty(FieldStatus)
	changes.MarkDirty(FieldUpdatedAt)
	p := &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		status:      ProductStatusActive,
		createdAt:   now,
		updatedAt:   now,
		changes:     changes,
		events:      make([]DomainEvent, 0, 1),
	}
	p.events = append(p.events, &ProductCreatedEvent{
		ProductID:   id,
		Name:        name,
		Description: description,
		Category:    category,
		BasePrice:   basePrice,
		At:          now,
	})
	return p
}

func RestoreProduct(id, name, description, category string, basePrice *Money, discount *Discount, status ProductStatus, createdAt, updatedAt, archivedAt time.Time) *Product {
	if basePrice == nil {
		return nil
	}
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
		events:      nil,
	}
}

func (p *Product) ID() string {
	if p == nil {
		return ""
	}
	return p.id
}

func (p *Product) Name() string {
	if p == nil {
		return ""
	}
	return p.name
}

func (p *Product) Description() string {
	if p == nil {
		return ""
	}
	return p.description
}

func (p *Product) Category() string {
	if p == nil {
		return ""
	}
	return p.category
}

func (p *Product) BasePrice() *Money {
	if p == nil {
		return nil
	}
	return p.basePrice
}

func (p *Product) Discount() *Discount {
	if p == nil {
		return nil
	}
	return p.discount
}

func (p *Product) Status() ProductStatus {
	if p == nil {
		return ""
	}
	return p.status
}

func (p *Product) CreatedAt() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.createdAt
}

func (p *Product) UpdatedAt() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.updatedAt
}

func (p *Product) ArchivedAt() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.archivedAt
}

func (p *Product) Changes() *ChangeTracker {
	if p == nil {
		return nil
	}
	return p.changes
}

func (p *Product) DomainEvents() []DomainEvent {
	if p == nil {
		return nil
	}
	return p.events
}

func (p *Product) UpdateDetails(name, description, category string, now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	p.name = name
	p.description = description
	p.category = category
	p.updatedAt = now
	p.changes.MarkDirty(FieldName)
	p.changes.MarkDirty(FieldDescription)
	p.changes.MarkDirty(FieldCategory)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &ProductUpdatedEvent{
		ProductID:   p.id,
		Name:        name,
		Description: description,
		Category:    category,
		At:          now,
	})
	return nil
}

func (p *Product) Activate(now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	if p.status == ProductStatusActive {
		return nil
	}
	p.status = ProductStatusActive
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &ProductActivatedEvent{ProductID: p.id, At: now})
	return nil
}

func (p *Product) Deactivate(now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	if p.status == ProductStatusInactive {
		return nil
	}
	p.status = ProductStatusInactive
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &ProductDeactivatedEvent{ProductID: p.id, At: now})
	return nil
}

func (p *Product) Archive(now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status == ProductStatusArchived {
		return nil
	}
	p.status = ProductStatusArchived
	p.archivedAt = now
	p.updatedAt = now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldArchivedAt)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &ProductArchivedEvent{ProductID: p.id, At: now})
	return nil
}

func (p *Product) ApplyDiscount(discount *Discount, now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}
	if discount == nil || !discount.IsValidAt(now) {
		return ErrInvalidDiscountPeriod
	}
	p.discount = discount
	p.updatedAt = now
	p.changes.MarkDirty(FieldDiscount)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &DiscountAppliedEvent{
		ProductID: p.id,
		Percent:   discount.Percentage(),
		StartDate: discount.StartDate(),
		EndDate:   discount.EndDate(),
		At:        now,
	})
	return nil
}

func (p *Product) RemoveDiscount(now time.Time) error {
	if p == nil {
		return nil
	}
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	if p.discount == nil {
		return nil
	}
	p.discount = nil
	p.updatedAt = now
	p.changes.MarkDirty(FieldDiscount)
	p.changes.MarkDirty(FieldUpdatedAt)
	p.events = append(p.events, &DiscountRemovedEvent{ProductID: p.id, At: now})
	return nil
}
