package merge

import (
	"sync"

	"google.golang.org/protobuf/types/known/structpb"
)

var (
	merger Merger
	mutex  sync.RWMutex
)

// Merger defines an interface for merging multiple protobuf Structs into one
type Merger interface {
	// Merge combines multiple structpb.Struct values into a single struct
	Merge(values ...*structpb.Struct) *structpb.Struct
}

// SetMerger sets the global merger instance with thread-safe protection.
// This function uses a mutex to ensure atomic write operations when updating
// the shared 'merger' variable.
//
// Parameters:
//
//	m - The Merger implementation to be set as the global merger.
//	    Passing nil will effectively clear the current merger.
func SetMerger(m Merger) {
	// Synchronize access to shared 'merger' variable
	mutex.Lock()
	merger = m
	mutex.Unlock()
}

// GetMerger returns the current Merger instance in a thread-safe manner.
//
// This function uses a mutex to ensure safe concurrent access to the shared
// merger variable. It locks the mutex before reading the value and unlocks
// immediately after to minimize contention.
//
// Returns:
//
//	Merger - The current merger instance being used
func GetMerger() Merger {
	// Lock mutex to ensure thread-safe access to shared merger variable
	mutex.RLock()
	m := merger
	mutex.RUnlock()
	return m
}
