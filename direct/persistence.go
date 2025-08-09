package direct

import "time"

// OrderHeader mirrors an RDBMS-oriented header row shape.
type OrderHeader struct {
	ID        string
	Customer  string
	Street    string
	City      string
	State     string
	Zip       string
	CreatedAt int64 // unix nanos
	UpdatedAt int64
}

// OrderItemRow mirrors an order_items row shape.
type OrderItemRow struct {
	OrderID    string
	SKU        string
	Quantity   int
	PriceCents int64
}

type persistenceRecord struct {
	Header OrderHeader
	Items  []OrderItemRow
}

func toPersistenceRecord(o *Order) persistenceRecord {
	rec := persistenceRecord{
		Header: OrderHeader{
			ID:        o.ID,
			Customer:  o.Customer,
			Street:    o.Shipping.Street,
			City:      o.Shipping.City,
			State:     o.Shipping.State,
			Zip:       o.Shipping.Zip,
			CreatedAt: o.CreatedAt.UnixNano(),
			UpdatedAt: o.UpdatedAt.UnixNano(),
		},
	}
	rec.Items = make([]OrderItemRow, len(o.Items))
	for i, it := range o.Items {
		rec.Items[i] = OrderItemRow{OrderID: o.ID, SKU: it.SKU, Quantity: it.Quantity, PriceCents: it.PriceCents}
	}
	return rec
}

func fromPersistenceRecord(rec persistenceRecord) *Order {
	items := make([]LineItem, len(rec.Items))
	for i, it := range rec.Items {
		items[i] = LineItem{SKU: it.SKU, Quantity: it.Quantity, PriceCents: it.PriceCents}
	}
	return &Order{
		ID:        rec.Header.ID,
		Customer:  rec.Header.Customer,
		Shipping:  Address{Street: rec.Header.Street, City: rec.Header.City, State: rec.Header.State, Zip: rec.Header.Zip},
		Items:     items,
		CreatedAt: time.Unix(0, rec.Header.CreatedAt),
		UpdatedAt: time.Unix(0, rec.Header.UpdatedAt),
	}
}
