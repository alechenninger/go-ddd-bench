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
			Customer:  "cust",
			Shipping:  Address{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"},
			Items:     []LineItem{{SKU: "A", Quantity: 1, PriceCents: 1234}, {SKU: "B", Quantity: 2, PriceCents: 555}},
			CreatedAt: clock.Now(),
			UpdatedAt: clock.Now(),
		}
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
