package encap

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"testing"
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
		o := NewOrder(randID(), "cust", SnapshotAddress{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"})
		o.AddItem("A", 1, 1234)
		o.AddItem("B", 2, 555)
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
	orders := seedEncapOrders(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := roundTripEncap(orders[i%len(orders)])
		sinkEncap = res
		runtime.KeepAlive(sinkEncap)
	}
}
