package directflat

import (
	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

type OrderHeader struct {
	ID        string
	Customer  string
	Street    string
	City      string
	State     string
	Zip       string
	CreatedAt int64
	UpdatedAt int64
}

type OrderItemRow struct {
	OrderID    string
	SKU        string
	Quantity   int
	PriceCents int64
}

type OrderRecord struct {
	Header OrderHeader
	Items  []OrderItemRow
}

func NewOrderRecord(id, customer string) *OrderRecord {
	now := clock.Now().UnixNano()
	return &OrderRecord{
		Header: OrderHeader{ID: id, Customer: customer, Street: "1 Main", City: "Town", State: "CA", Zip: "94000", CreatedAt: now, UpdatedAt: now},
		Items:  nil,
	}
}

func (r *OrderRecord) AddItem(sku string, qty int, priceCents int64) {
	r.Items = append(r.Items, OrderItemRow{OrderID: r.Header.ID, SKU: sku, Quantity: qty, PriceCents: priceCents})
	r.Header.UpdatedAt = clock.Now().UnixNano()
}
