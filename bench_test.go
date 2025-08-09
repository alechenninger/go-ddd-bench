package bench

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/alechenninger/go-ddd-bench/direct"
	"github.com/alechenninger/go-ddd-bench/encap"
	"github.com/alechenninger/go-ddd-bench/internal/clock"
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
			ID: randID(),
			Customer: direct.Customer{
				Name:    direct.Name{First: "Ada", Last: "Lovelace"},
				Email:   "ada@example.com",
				Loyalty: direct.Loyalty{Tier: "gold", Points: 100},
			},
			Shipping:  direct.Address{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"},
			Billing:   direct.Address{Street: "2 Main", City: "Town", State: "CA", Zip: "94000"},
			Items:     nil,
			CreatedAt: clock.Now(),
			UpdatedAt: clock.Now(),
		}
		order.AddItem("A", 1, 1234, "USD", direct.ItemFlags{})
		order.AddItem("B", 2, 555, "USD", direct.ItemFlags{Backorder: true})
		_ = repo.Save(order)
	}
	return repo
}

func seedEncapRepo(n int) *encap.Repo {
	repo := encap.NewRepo()
	for i := 0; i < n; i++ {
		cust := encap.SnapshotCustomer{Name: encap.SnapshotName{First: "Ada", Last: "Lovelace"}, Email: "ada@example.com", Loyalty: encap.SnapshotLoyalty{Tier: "gold", Points: 100}}
		ship := encap.SnapshotAddress{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"}
		bill := encap.SnapshotAddress{Street: "2 Main", City: "Town", State: "CA", Zip: "94000"}
		order := encap.NewOrder(randID(), cust, ship, bill)
		order.AddItem("A", 1, 1234, "USD", encap.SnapshotItemFlags{})
		order.AddItem("B", 2, 555, "USD", encap.SnapshotItemFlags{Backorder: true})
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
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()

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
		order.AddItem("C", 1, 99, "USD", direct.ItemFlags{Digital: true})
		if err := repo.Save(order); err != nil {
			b.Fatal(err)
		}
		Blackhole = order
	}
}

// BenchmarkEncap_RMW simulates read-modify-write via snapshot + persistence DTO transforms.
func BenchmarkEncap_RMW(b *testing.B) {
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()

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
		order.AddItem("C", 1, 99, "USD", encap.SnapshotItemFlags{Digital: true})
		if err := repo.Save(order); err != nil {
			b.Fatal(err)
		}
		Blackhole = order
	}
}
