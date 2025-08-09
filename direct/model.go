package direct

import (
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

// Address is embedded as a public struct to allow direct (de)serialization.
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

// LineItem is a public struct to allow direct (de)serialization.
type LineItem struct {
	SKU        string
	Quantity   int
	PriceCents int64
}

// Order is the aggregate root with public fields.
type Order struct {
	ID        string
	Customer  string
	Shipping  Address
	Items     []LineItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *Order) AddItem(sku string, qty int, priceCents int64) {
	o.Items = append(o.Items, LineItem{SKU: sku, Quantity: qty, PriceCents: priceCents})
	o.touch()
}

func (o *Order) UpdateShipping(addr Address) {
	o.Shipping = addr
	o.touch()
}

func (o *Order) touch() { o.UpdatedAt = clock.Now() }
