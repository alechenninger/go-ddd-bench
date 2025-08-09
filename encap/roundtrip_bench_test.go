package encap

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"testing"
	"time"

	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

var sinkEncap any

func randID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func seedEncapOrders(n int) []*Order {
	orders := make([]*Order, 0, n)
	for i := 0; i < n; i++ {
		cust := SnapshotCustomer{Name: SnapshotName{First: "Ada", Last: "Lovelace"}, Email: "ada@example.com", Loyalty: SnapshotLoyalty{Tier: "gold", Points: 100}}
		ship := SnapshotAddress{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"}
		bill := SnapshotAddress{Street: "2 Main", City: "Town", State: "CA", Zip: "94000"}
		o := NewOrder(randID(), cust, ship, bill)
		o.AddItem("A", 1, 1234, "USD", SnapshotItemFlags{})
		o.AddItem("B", 2, 555, "USD", SnapshotItemFlags{Backorder: true})
		orders = append(orders, o)
	}
	return orders
}

//go:noinline
func roundTripEncap(o *Order) *Order {
	s := o.ToSnapshot()
	rec := toPersistenceRecord(s)
	s2 := fromPersistenceRecord(rec)
	return FromSnapshot(s2)
}

func BenchmarkEncap_RoundTrip_NoJSON(b *testing.B) {
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()
	orders := seedEncapOrders(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := roundTripEncap(orders[i%len(orders)])
		sinkEncap = res
		runtime.KeepAlive(sinkEncap)
	}
}
