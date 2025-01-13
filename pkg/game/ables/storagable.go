package ables

// Storagable can store stuff.
type Storagable []StorageItem

// AddItem adds an item to storage the storage, incrementing count if existing.
func (s *Storagable) AddItem(name string, tag string) {
	for i, item := range *s {
		if item.Tag == tag {
			(*s)[i].Count++
			return
		}
	}
	*s = append(*s, StorageItem{Name: name, Tag: tag, Count: 1})
}

// RemoveItem removes an item from storage, decrementing count if existing.
func (s *Storagable) RemoveItem(tag string) {
	for i, item := range *s {
		if item.Tag == tag {
			(*s)[i].Count--
			if (*s)[i].Count == 0 {
				*s = append((*s)[:i], (*s)[i+1:]...)
			}
			return
		}
	}
}

// HasItem returns true if the storage has the item.
func (s Storagable) HasItem(name string) bool {
	for _, item := range s {
		if item.Name == name {
			return true
		}
	}
	return false
}

// StorageItem is an item in storage.
type StorageItem struct {
	Name  string
	Tag   string
	Count int
}
