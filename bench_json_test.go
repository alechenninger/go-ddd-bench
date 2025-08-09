package bench

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/alechenninger/go-ddd-bench/directflat"
	"github.com/alechenninger/go-ddd-bench/encap"
	"github.com/alechenninger/go-ddd-bench/internal/clock"
)

func randID2() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func seedDirectFlatRepo(n int) *directflat.Repo {
	repo := directflat.NewRepo()
	for i := 0; i < n; i++ {
		rec := directflat.NewOrderRecord(randID2(), "cust")
		rec.AddItem("A", 1, 1234)
		rec.AddItem("B", 2, 555)
		_ = repo.Save(rec)
	}
	return repo
}

func seedEncapRepo2(n int) *encap.Repo {
	repo := encap.NewRepo()
	for i := 0; i < n; i++ {
		o := encap.NewOrder(randID2(), "cust", encap.SnapshotAddress{Street: "1 Main", City: "Town", State: "CA", Zip: "94000"})
		o.AddItem("A", 1, 1234)
		o.AddItem("B", 2, 555)
		_ = repo.Save(o)
	}
	return repo
}

const nSeedJSON = 1000

func BenchmarkDirectFlat_JSON_RMW(b *testing.B) {
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()
	repo := seedDirectFlatRepo(nSeedJSON)
	ids := make([]string, 0, nSeedJSON)
	for id := range repo.DataUnsafeForBench() {
		ids = append(ids, id)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		rec, err := repo.FindByID(id)
		if err != nil {
			b.Fatal(err)
		}
		rec.AddItem("C", 1, 99)
		if err := repo.Save(rec); err != nil {
			b.Fatal(err)
		}
		Blackhole = rec
	}
}

func BenchmarkEncap_JSON_RMW(b *testing.B) {
	restore := clock.UseMonotonicFake(time.Unix(0, 0), time.Nanosecond)
	defer restore()
	repo := seedEncapRepo2(nSeedJSON)
	ids := make([]string, 0, nSeedJSON)
	for id := range repo.DataUnsafeForBench() {
		ids = append(ids, id)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		o, err := repo.FindByID(id)
		if err != nil {
			b.Fatal(err)
		}
		o.AddItem("C", 1, 99)
		if err := repo.Save(o); err != nil {
			b.Fatal(err)
		}
		Blackhole = o
	}
}
