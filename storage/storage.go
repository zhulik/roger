package storage

// ProgressItem - represents progress of one item
type ProgressItem struct {
	Total    int64
	Progress int64
}

// Storage represent current synchronization state
type Storage struct {
	Total *ProgressItem
	Files map[string]ProgressItem
}

// Reset - resets storage state
func (s *Storage) Reset() {
	s.Total = nil
	s.Files = map[string]ProgressItem{}
}

// Initialize - reinitializes storage with new items
func (s *Storage) Initialize(total int64, fileSizes map[string]int64) {
	s.Reset()
	s.Total = &ProgressItem{Total: total, Progress: 0}
	for name, size := range fileSizes {
		s.Files[name] = ProgressItem{Total: size, Progress: 0}
	}
}
