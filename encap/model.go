package encap

import (
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

// Snapshot is a DTO representing the Order's domain state, without persistence concerns.
type Snapshot struct {
	ID        string
	Customer  string
	Shipping  SnapshotAddress
	Items     []SnapshotLineItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SnapshotAddress struct {
	Street string
	City   string
	State  string
	Zip    string
}

type SnapshotLineItem struct {
	SKU        string
	Quantity   int
	PriceCents int64
}

// Order encapsulates its state.
type Order struct {
	id        string
	customer  string
	shipping  address
	items     []lineItem
	createdAt time.Time
	updatedAt time.Time
}

type address struct {
	street string
	city   string
	state  string
	zip    string
}

type lineItem struct {
	sku        string
	quantity   int
	priceCents int64
}

func NewOrder(id, customer string, shipping SnapshotAddress) *Order {
	return &Order{
		id:        id,
		customer:  customer,
		shipping:  address{street: shipping.Street, city: shipping.City, state: shipping.State, zip: shipping.Zip},
		items:     nil,
		createdAt: clock.Now(),
		updatedAt: clock.Now(),
	}
}

func (o *Order) AddItem(sku string, qty int, priceCents int64) {
	o.items = append(o.items, lineItem{sku: sku, quantity: qty, priceCents: priceCents})
	o.touch()
}

func (o *Order) UpdateShipping(s SnapshotAddress) {
	o.shipping = address{street: s.Street, city: s.City, state: s.State, zip: s.Zip}
	o.touch()
}

func (o *Order) touch() { o.updatedAt = clock.Now() }

// ToSnapshot exposes state for serialization elsewhere.
func (o *Order) ToSnapshot() Snapshot {
	items := make([]SnapshotLineItem, len(o.items))
	for i, it := range o.items {
		items[i] = SnapshotLineItem{SKU: it.sku, Quantity: it.quantity, PriceCents: it.priceCents}
	}
	return Snapshot{
		ID:        o.id,
		Customer:  o.customer,
		Shipping:  SnapshotAddress{Street: o.shipping.street, City: o.shipping.city, State: o.shipping.state, Zip: o.shipping.zip},
		Items:     items,
		CreatedAt: o.createdAt,
		UpdatedAt: o.updatedAt,
	}
}

// FromSnapshot constructs an Order from a snapshot.
func FromSnapshot(s Snapshot) *Order {
	items := make([]lineItem, len(s.Items))
	for i, it := range s.Items {
		items[i] = lineItem{sku: it.SKU, quantity: it.Quantity, priceCents: it.PriceCents}
	}
	return &Order{
		id:        s.ID,
		customer:  s.Customer,
		shipping:  address{street: s.Shipping.Street, city: s.Shipping.City, state: s.Shipping.State, zip: s.Shipping.Zip},
		items:     items,
		createdAt: s.CreatedAt,
		updatedAt: s.UpdatedAt,
	}
}
