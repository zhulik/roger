package storage

import "sync"

var instance *Storage
var once sync.Once

// GetInstance - returns shared instance of storage
func GetInstance() *Storage {
	once.Do(func() {
		instance = &Storage{}
	})
	return instance
}
