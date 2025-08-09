package encap

import (
	"encoding/json"
	"errors"
	"sync"
)

// RDBMS-oriented DTOs (tables) — flat structures intended for persistence.
// OrderHeader corresponds to an orders table.
type OrderHeader struct {
	ID            string
	CustomerFirst string
	CustomerLast  string
	CustomerEmail string
	LoyaltyTier   string
	LoyaltyPoints int
	Street        string
	City          string
	State         string
	Zip           string
	BillStreet    string
	BillCity      string
	BillState     string
	BillZip       string
	CreatedAt     int64 // unix nanos for storage-form convenience
	UpdatedAt     int64
}

// OrderItemRow corresponds to an order_items table.
type OrderItemRow struct {
	OrderID    string
	SKU        string
	Quantity   int
	PriceCents int64
	Currency   string
	Backorder  bool
	Digital    bool
}

// persistenceRecord simulates multiple tables grouped together for a single aggregate.
// We simulate IO by JSON round-tripping these records.
type persistenceRecord struct {
	Header OrderHeader
	Items  []OrderItemRow
}

// Repo simulates a repository with multiple transformations:
// domain <-> snapshot <-> persistence DTOs <-> bytes
// We store JSON blobs to emulate IO and avoid in-memory aliasing.
type Repo struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewRepo() *Repo { return &Repo{data: make(map[string][]byte)} }

func (r *Repo) Save(o *Order) error {
	s := o.ToSnapshot()
	rec := toPersistenceRecord(s)
	blob, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.data[s.ID] = blob
	r.mu.Unlock()
	return nil
}

func (r *Repo) FindByID(id string) (*Order, error) {
	r.mu.RLock()
	blob, ok := r.data[id]
	r.mu.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	var rec persistenceRecord
	if err := json.Unmarshal(blob, &rec); err != nil {
		return nil, err
	}
	s := fromPersistenceRecord(rec)
	return FromSnapshot(s), nil
}

// DataUnsafeForBench returns a copy of the keys to iterate in benchmarks.
func (r *Repo) DataUnsafeForBench() map[string]struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make(map[string]struct{}, len(r.data))
	for k := range r.data {
		ids[k] = struct{}{}
	}
	return ids
}

func toPersistenceRecord(s Snapshot) persistenceRecord {
	rec := persistenceRecord{
		Header: OrderHeader{
			ID:            s.ID,
			CustomerFirst: s.Customer.Name.First,
			CustomerLast:  s.Customer.Name.Last,
			CustomerEmail: s.Customer.Email,
			LoyaltyTier:   s.Customer.Loyalty.Tier,
			LoyaltyPoints: s.Customer.Loyalty.Points,
			Street:        s.Shipping.Street,
			City:          s.Shipping.City,
			State:         s.Shipping.State,
			Zip:           s.Shipping.Zip,
			BillStreet:    s.Billing.Street,
			BillCity:      s.Billing.City,
			BillState:     s.Billing.State,
			BillZip:       s.Billing.Zip,
			CreatedAt:     s.CreatedAt.UnixNano(),
			UpdatedAt:     s.UpdatedAt.UnixNano(),
		},
	}
	rec.Items = make([]OrderItemRow, len(s.Items))
	for i, it := range s.Items {
		rec.Items[i] = OrderItemRow{OrderID: s.ID, SKU: it.SKU, Quantity: it.Quantity, PriceCents: it.Price.Cents, Currency: it.Price.Currency, Backorder: it.Flags.Backorder, Digital: it.Flags.Digital}
	}
	return rec
}

func fromPersistenceRecord(rec persistenceRecord) Snapshot {
	s := Snapshot{
		ID: rec.Header.ID,
		Customer: SnapshotCustomer{
			Name:    SnapshotName{First: rec.Header.CustomerFirst, Last: rec.Header.CustomerLast},
			Email:   rec.Header.CustomerEmail,
			Loyalty: SnapshotLoyalty{Tier: rec.Header.LoyaltyTier, Points: rec.Header.LoyaltyPoints},
		},
		Shipping:  SnapshotAddress{Street: rec.Header.Street, City: rec.Header.City, State: rec.Header.State, Zip: rec.Header.Zip},
		Billing:   SnapshotAddress{Street: rec.Header.BillStreet, City: rec.Header.BillCity, State: rec.Header.BillState, Zip: rec.Header.BillZip},
		CreatedAt: unixToTime(rec.Header.CreatedAt),
		UpdatedAt: unixToTime(rec.Header.UpdatedAt),
	}
	s.Items = make([]SnapshotLineItem, len(rec.Items))
	for i, row := range rec.Items {
		s.Items[i] = SnapshotLineItem{SKU: row.SKU, Quantity: row.Quantity, Price: SnapshotMoney{Cents: row.PriceCents, Currency: row.Currency}, Flags: SnapshotItemFlags{Backorder: row.Backorder, Digital: row.Digital}}
	}
	return s
}
