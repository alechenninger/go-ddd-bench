package direct

import "time"

// OrderHeader mirrors an RDBMS-oriented header row shape.
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
	CreatedAt     int64 // unix nanos
	UpdatedAt     int64
}

// OrderItemRow mirrors an order_items row shape.
type OrderItemRow struct {
	OrderID    string
	SKU        string
	Quantity   int
	PriceCents int64
	Currency   string
	Backorder  bool
	Digital    bool
}

type persistenceRecord struct {
	Header OrderHeader
	Items  []OrderItemRow
}

func toPersistenceRecord(o *Order) persistenceRecord {
	rec := persistenceRecord{
		Header: OrderHeader{
			ID:            o.ID,
			CustomerFirst: o.Customer.Name.First,
			CustomerLast:  o.Customer.Name.Last,
			CustomerEmail: o.Customer.Email,
			LoyaltyTier:   o.Customer.Loyalty.Tier,
			LoyaltyPoints: o.Customer.Loyalty.Points,
			Street:        o.Shipping.Street,
			City:          o.Shipping.City,
			State:         o.Shipping.State,
			Zip:           o.Shipping.Zip,
			BillStreet:    o.Billing.Street,
			BillCity:      o.Billing.City,
			BillState:     o.Billing.State,
			BillZip:       o.Billing.Zip,
			CreatedAt:     o.CreatedAt.UnixNano(),
			UpdatedAt:     o.UpdatedAt.UnixNano(),
		},
	}
	rec.Items = make([]OrderItemRow, len(o.Items))
	for i, it := range o.Items {
		rec.Items[i] = OrderItemRow{OrderID: o.ID, SKU: it.SKU, Quantity: it.Quantity, PriceCents: it.Price.Cents, Currency: it.Price.Currency, Backorder: it.Flags.Backorder, Digital: it.Flags.Digital}
	}
	return rec
}

func fromPersistenceRecord(rec persistenceRecord) *Order {
	items := make([]LineItem, len(rec.Items))
	for i, it := range rec.Items {
		items[i] = LineItem{SKU: it.SKU, Quantity: it.Quantity, Price: Money{Cents: it.PriceCents, Currency: it.Currency}, Flags: ItemFlags{Backorder: it.Backorder, Digital: it.Digital}}
	}
	return &Order{
		ID: rec.Header.ID,
		Customer: Customer{
			Name:    Name{First: rec.Header.CustomerFirst, Last: rec.Header.CustomerLast},
			Email:   rec.Header.CustomerEmail,
			Loyalty: Loyalty{Tier: rec.Header.LoyaltyTier, Points: rec.Header.LoyaltyPoints},
		},
		Shipping:  Address{Street: rec.Header.Street, City: rec.Header.City, State: rec.Header.State, Zip: rec.Header.Zip},
		Billing:   Address{Street: rec.Header.BillStreet, City: rec.Header.BillCity, State: rec.Header.BillState, Zip: rec.Header.BillZip},
		Items:     items,
		CreatedAt: time.Unix(0, rec.Header.CreatedAt),
		UpdatedAt: time.Unix(0, rec.Header.UpdatedAt),
	}
}
