package domain

const (
	FieldName        = "name"
	FieldDescription = "description"
	FieldCategory    = "category"
	FieldBasePrice   = "base_price"
	FieldDiscount    = "discount"
	FieldStatus      = "status"
	FieldArchivedAt  = "archived_at"
	FieldUpdatedAt   = "updated_at"
)

type ChangeTracker struct {
	dirty map[string]bool
}

func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{dirty: make(map[string]bool)}
}

func (c *ChangeTracker) MarkDirty(field string) {
	if c == nil {
		return
	}
	c.dirty[field] = true
}

func (c *ChangeTracker) Dirty(field string) bool {
	if c == nil {
		return false
	}
	return c.dirty[field]
}
