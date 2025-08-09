package direct

import (
	"encoding/json"
	"errors"
	"sync"
)

// DirectRepo simulates a repository that (de)serializes the model directly.
type DirectRepo struct {
	mu   sync.RWMutex
	data map[string][]byte // stores JSON blobs
}

func NewDirectRepo() *DirectRepo { return &DirectRepo{data: make(map[string][]byte)} }

func (r *DirectRepo) Save(o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	blob, err := json.Marshal(o)
	if err != nil {
		return err
	}
	r.data[o.ID] = blob
	return nil
}

func (r *DirectRepo) FindByID(id string) (*Order, error) {
	r.mu.RLock()
	blob, ok := r.data[id]
	r.mu.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	var o Order
	if err := json.Unmarshal(blob, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// DataUnsafeForBench returns a copy of the keys to iterate in benchmarks.
func (r *DirectRepo) DataUnsafeForBench() map[string]struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make(map[string]struct{}, len(r.data))
	for k := range r.data {
		ids[k] = struct{}{}
	}
	return ids
}
