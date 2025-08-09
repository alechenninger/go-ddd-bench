package bench

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/alechenninger/go-ddd-bench/direct"
	"github.com/alechenninger/go-ddd-bench/encap"
)

func randID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func seedDirectRepo(n int) *direct.DirectRepo {
	repo := direct.NewDirectRepo()
	for i := 0; i < n; i++ {
		order := &direct.Order{
			ID:        randID(),
			Customer:  "cust",
			Shipping:  direct.Address{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"},
			Items:     []direct.LineItem{{SKU: "A", Quantity: 1, PriceCents: 1234}, {SKU: "B", Quantity: 2, PriceCents: 555}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = repo.Save(order)
	}
	return repo
}

func seedEncapRepo(n int) *encap.Repo {
	repo := encap.NewRepo()
	for i := 0; i < n; i++ {
		order := encap.NewOrder(randID(), "cust", encap.SnapshotAddress{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"})
		order.AddItem("A", 1, 1234)
		order.AddItem("B", 2, 555)
		_ = repo.Save(order)
	}
	return repo
}

// Benchmark parameters
const (
	nSeed = 1000 // number of seeded orders
)

// BenchmarkDirect_RMW simulates read-modify-write via direct (de)serialization.
func BenchmarkDirect_RMW(b *testing.B) {
	repo := seedDirectRepo(nSeed)
	ids := make([]string, 0, nSeed)
	for id := range repo.DataUnsafeForBench() { // helper method returns map copy
		ids = append(ids, id)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		order, err := repo.FindByID(id)
		if err != nil {
			b.Fatal(err)
		}
		order.AddItem("C", 1, 99)
		if err := repo.Save(order); err != nil {
			b.Fatal(err)
		}
		Blackhole = order
	}
}

// BenchmarkEncap_RMW simulates read-modify-write via snapshot + persistence DTO transforms.
func BenchmarkEncap_RMW(b *testing.B) {
	repo := seedEncapRepo(nSeed)
	ids := make([]string, 0, nSeed)
	for id := range repo.DataUnsafeForBench() { // helper method returns map copy
		ids = append(ids, id)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		order, err := repo.FindByID(id)
		if err != nil {
			b.Fatal(err)
		}
		order.AddItem("C", 1, 99)
		if err := repo.Save(order); err != nil {
			b.Fatal(err)
		}
		Blackhole = order
	}
}
