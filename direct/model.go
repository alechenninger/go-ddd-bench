package direct

import (
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

// Name nested under Customer
type Name struct {
	First string
	Last  string
}

// Loyalty nested under Customer
type Loyalty struct {
	Tier   string
	Points int
}

// Customer nested object
type Customer struct {
	Name    Name
	Email   string
	Loyalty Loyalty
}

// Address value object
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

// Money nested under LineItem
type Money struct {
	Cents    int64
	Currency string
}

// ItemFlags nested under LineItem
type ItemFlags struct {
	Backorder bool
	Digital   bool
}

// LineItem contains nested Money and Flags
type LineItem struct {
	SKU      string
	Quantity int
	Price    Money
	Flags    ItemFlags
}

// Order is the aggregate root with more nested state.
type Order struct {
	ID        string
	Customer  Customer
	Shipping  Address
	Billing   Address
	Items     []LineItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *Order) AddItem(sku string, qty int, priceCents int64, currency string, flags ItemFlags) {
	o.Items = append(o.Items, LineItem{SKU: sku, Quantity: qty, Price: Money{Cents: priceCents, Currency: currency}, Flags: flags})
	o.touch()
}

func (o *Order) UpdateShipping(addr Address) {
	o.Shipping = addr
	o.touch()
}

func (o *Order) UpdateBilling(addr Address) {
	o.Billing = addr
	o.touch()
}

func (o *Order) touch() { o.UpdatedAt = clock.Now() }
