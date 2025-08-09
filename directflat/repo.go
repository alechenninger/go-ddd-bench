package directflat

import (
	"encoding/json"
	"errors"
	"sync"
)

type Repo struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewRepo() *Repo { return &Repo{data: make(map[string][]byte)} }

func (r *Repo) Save(rec *OrderRecord) error {
	blob, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.data[rec.Header.ID] = blob
	r.mu.Unlock()
	return nil
}

func (r *Repo) FindByID(id string) (*OrderRecord, error) {
	r.mu.RLock()
	blob, ok := r.data[id]
	r.mu.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	var rec OrderRecord
	if err := json.Unmarshal(blob, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Repo) DataUnsafeForBench() map[string]struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make(map[string]struct{}, len(r.data))
	for k := range r.data {
		ids[k] = struct{}{}
	}
	return ids
}
