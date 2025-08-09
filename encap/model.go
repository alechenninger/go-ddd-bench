package encap

import (
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

// Snapshot DTOs with deeper nesting

type SnapshotName struct{ First, Last string }

type SnapshotLoyalty struct {
	Tier   string
	Points int
}

type SnapshotCustomer struct {
	Name    SnapshotName
	Email   string
	Loyalty SnapshotLoyalty
}

type Snapshot struct {
	ID        string
	Customer  SnapshotCustomer
	Shipping  SnapshotAddress
	Billing   SnapshotAddress
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

type SnapshotMoney struct {
	Cents    int64
	Currency string
}

type SnapshotItemFlags struct{ Backorder, Digital bool }

type SnapshotLineItem struct {
	SKU      string
	Quantity int
	Price    SnapshotMoney
	Flags    SnapshotItemFlags
}

// Encapsulated model

type name struct{ first, last string }

type loyalty struct {
	tier   string
	points int
}

type customer struct {
	name    name
	email   string
	loyalty loyalty
}

type address struct {
	street string
	city   string
	state  string
	zip    string
}

type money struct {
	cents    int64
	currency string
}

type itemFlags struct{ backorder, digital bool }

type lineItem struct {
	sku      string
	quantity int
	price    money
	flags    itemFlags
}

type Order struct {
	id        string
	customer  customer
	shipping  address
	billing   address
	items     []lineItem
	createdAt time.Time
	updatedAt time.Time
}

func NewOrder(id string, cust SnapshotCustomer, shipping, billing SnapshotAddress) *Order {
	return &Order{
		id: id,
		customer: customer{
			name:    name{first: cust.Name.First, last: cust.Name.Last},
			email:   cust.Email,
			loyalty: loyalty{tier: cust.Loyalty.Tier, points: cust.Loyalty.Points},
		},
		shipping:  address{street: shipping.Street, city: shipping.City, state: shipping.State, zip: shipping.Zip},
		billing:   address{street: billing.Street, city: billing.City, state: billing.State, zip: billing.Zip},
		items:     nil,
		createdAt: clock.Now(),
		updatedAt: clock.Now(),
	}
}

func (o *Order) AddItem(sku string, qty int, priceCents int64, currency string, flags SnapshotItemFlags) {
	o.items = append(o.items, lineItem{sku: sku, quantity: qty, price: money{cents: priceCents, currency: currency}, flags: itemFlags{backorder: flags.Backorder, digital: flags.Digital}})
	o.touch()
}

func (o *Order) UpdateShipping(s SnapshotAddress) {
	o.shipping = address{street: s.Street, city: s.City, state: s.State, zip: s.Zip}
	o.touch()
}
func (o *Order) UpdateBilling(s SnapshotAddress) {
	o.billing = address{street: s.Street, city: s.City, state: s.State, zip: s.Zip}
	o.touch()
}

func (o *Order) touch() { o.updatedAt = clock.Now() }

func (o *Order) ToSnapshot() Snapshot {
	items := make([]SnapshotLineItem, len(o.items))
	for i, it := range o.items {
		items[i] = SnapshotLineItem{
			SKU:      it.sku,
			Quantity: it.quantity,
			Price:    SnapshotMoney{Cents: it.price.cents, Currency: it.price.currency},
			Flags:    SnapshotItemFlags{Backorder: it.flags.backorder, Digital: it.flags.digital},
		}
	}
	return Snapshot{
		ID: o.id,
		Customer: SnapshotCustomer{
			Name:    SnapshotName{First: o.customer.name.first, Last: o.customer.name.last},
			Email:   o.customer.email,
			Loyalty: SnapshotLoyalty{Tier: o.customer.loyalty.tier, Points: o.customer.loyalty.points},
		},
		Shipping:  SnapshotAddress{Street: o.shipping.street, City: o.shipping.city, State: o.shipping.state, Zip: o.shipping.zip},
		Billing:   SnapshotAddress{Street: o.billing.street, City: o.billing.city, State: o.billing.state, Zip: o.billing.zip},
		Items:     items,
		CreatedAt: o.createdAt,
		UpdatedAt: o.updatedAt,
	}
}

func FromSnapshot(s Snapshot) *Order {
	items := make([]lineItem, len(s.Items))
	for i, it := range s.Items {
		items[i] = lineItem{sku: it.SKU, quantity: it.Quantity, price: money{cents: it.Price.Cents, currency: it.Price.Currency}, flags: itemFlags{backorder: it.Flags.Backorder, digital: it.Flags.Digital}}
	}
	return &Order{
		id:        s.ID,
		customer:  customer{name: name{first: s.Customer.Name.First, last: s.Customer.Name.Last}, email: s.Customer.Email, loyalty: loyalty{tier: s.Customer.Loyalty.Tier, points: s.Customer.Loyalty.Points}},
		shipping:  address{street: s.Shipping.Street, city: s.Shipping.City, state: s.Shipping.State, zip: s.Shipping.Zip},
		billing:   address{street: s.Billing.Street, city: s.Billing.City, state: s.Billing.State, zip: s.Billing.Zip},
		items:     items,
		createdAt: s.CreatedAt,
		updatedAt: s.UpdatedAt,
	}
}
