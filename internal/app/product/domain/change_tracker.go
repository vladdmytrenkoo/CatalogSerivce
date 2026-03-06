package domain

type ChangeTracker struct {
	dirtyFields map[string]bool
}

func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{dirtyFields: make(map[string]bool)}
}

func (t *ChangeTracker) MarkDirty(field string) {
	if field == "" {
		return
	}
	t.dirtyFields[field] = true
}

func (t *ChangeTracker) Clear(field string) {
	if field == "" {
		return
	}
	delete(t.dirtyFields, field)
}

func (t *ChangeTracker) IsDirty(field string) bool {
	return t.dirtyFields[field]
}

func (t *ChangeTracker) HasChanges() bool {
	return len(t.dirtyFields) > 0
}
