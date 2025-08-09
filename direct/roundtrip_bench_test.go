package direct

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"testing"
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

var sinkDirect any

func randID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func seedDirectOrders(n int) []*Order {
	orders := make([]*Order, 0, n)
	for i := 0; i < n; i++ {
		o := &Order{
			ID:        randID(),
			Customer:  Customer{Name: Name{First: "Ada", Last: "Lovelace"}, Email: "ada@example.com", Loyalty: Loyalty{Tier: "gold", Points: 100}},
			Shipping:  Address{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"},
			Billing:   Address{Street: "2 Main", City: "Town", State: "CA", Zip: "94000"},
			Items:     nil,
			CreatedAt: clock.Now(),
			UpdatedAt: clock.Now(),
		}
		o.AddItem("A", 1, 1234, "USD", ItemFlags{})
		o.AddItem("B", 2, 555, "USD", ItemFlags{Backorder: true})
		orders = append(orders, o)
	}
	return orders
}

//go:noinline
func roundTripDirect(o *Order) *Order {
	rec := toPersistenceRecord(o)
	return fromPersistenceRecord(rec)
}

func BenchmarkDirect_RoundTrip_NoJSON(b *testing.B) {
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()
	orders := seedDirectOrders(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := roundTripDirect(orders[i%len(orders)])
		sinkDirect = res
		runtime.KeepAlive(sinkDirect)
	}
}
