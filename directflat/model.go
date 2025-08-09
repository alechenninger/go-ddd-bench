package directflat

import (
	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

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
	CreatedAt     int64
	UpdatedAt     int64
}

type OrderItemRow struct {
	OrderID    string
	SKU        string
	Quantity   int
	PriceCents int64
	Currency   string
	Backorder  bool
	Digital    bool
}

type OrderRecord struct {
	Header OrderHeader
	Items  []OrderItemRow
}

func NewOrderRecord(id, first, last, email string, loyaltyTier string, loyaltyPts int) *OrderRecord {
	now := clock.Now().UnixNano()
	return &OrderRecord{
		Header: OrderHeader{ID: id, CustomerFirst: first, CustomerLast: last, CustomerEmail: email, LoyaltyTier: loyaltyTier, LoyaltyPoints: loyaltyPts, Street: "1 Main", City: "Town", State: "CA", Zip: "94000", BillStreet: "2 Main", BillCity: "Town", BillState: "CA", BillZip: "94000", CreatedAt: now, UpdatedAt: now},
		Items:  nil,
	}
}

func (r *OrderRecord) AddItem(sku string, qty int, priceCents int64, currency string, backorder, digital bool) {
	r.Items = append(r.Items, OrderItemRow{OrderID: r.Header.ID, SKU: sku, Quantity: qty, PriceCents: priceCents, Currency: currency, Backorder: backorder, Digital: digital})
	r.Header.UpdatedAt = clock.Now().UnixNano()
}
