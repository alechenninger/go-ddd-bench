## go-ddd-bench

Benchmarks for the effect of additional struct allocations and transformation layers in a DDD-style read-modify-write (RMW) cycle in Go.

### Variants

- **direct**: Domain model with public fields, repository (de)serializes the domain shape directly.
- **encap**: Encapsulated domain model; repository performs transformations domain ↔ snapshot ↔ persistence-shape before (de)serialization.
- **directflat**: Best-case coupling; the domain model matches the persistence-record shape exactly, so no transform is required before (de)serialization.

All RMW benches simulate IO via encoding/json to avoid database dependence while still exercising serialization/allocations.

### Time source

To avoid `time.Now()` syscall noise, all benchmarks use a shared fake clock (`internal/clock`) that returns a monotonically increasing timestamp. This makes allocations and transform work the dominant signal.

### How to run

- Typical:
  ```bash
  go test -bench=. -benchmem ./...
  ```
- More stable (recommended):
  ```bash
  go test -run=^$ -bench=. -benchmem -benchtime=3s -count=5 -cpu=1 ./...
  ```

### Environment

- goos: darwin
- goarch: arm64
- cpu: Apple M1 Pro
- Go: 1.22.5

### Aggregation method

- We report the median across repeated runs per benchmark. Median is preferred over mean for microbenchmarks to reduce the influence of outliers and GC jitter.
- The figures below come from runs with: `-benchtime=2s -count=3 -cpu=1`.

### Results (median of 2s x3 runs)

RMW with identical persistence shape (apples-to-apples serialization cost):

- DirectFlat_JSON_RMW: 110.495 µs/op, 30,767 B/op, 186 allocs/op
- Encap_JSON_RMW: 109.545 µs/op, 55,145 B/op, 178 allocs/op

RMW with domain JSON shapes (includes JSON-shape differences):

- Direct_RMW: 102.988 µs/op, 20,703 B/op, 104 allocs/op
- Encap_RMW: 109.346 µs/op, 55,194 B/op, 178 allocs/op

Pure transform round-trips (no serialization):

- Direct_RoundTrip_NoJSON: 243.0 ns/op, 544 B/op, 3 allocs/op
- Encap_RoundTrip_NoJSON: 437.9 ns/op, 768 B/op, 5 allocs/op

### Findings

- Persistence work dominates end-to-end cost:
  - JSON here is a stand-in for database driver encode/decode and IO. In real systems, the serialization/IO path (driver scanning, network, syscalls) dominates time and allocations much more than small in-memory transforms.
- Encapsulation overhead is small but measurable when isolated:
  - No-JSON round trips show ~+87 ns/op and +2 allocs/op for the encapsulated path due to extra representations and slice/struct copies.
- Domain-shape vs persistence-shape matters:
  - When the serialized shape matches tables (persistence shape), both variants pay similar serialization cost; the remaining difference is the transform overhead. With domain-shape JSON, encap also pays a JSON-shape cost, widening the gap (e.g. this is more like comparing using a document-style database vs relational)
- Occasional “encap faster” ns/op readings are noise:
  - Minor ns/op inversions are within run-to-run variance. Allocations consistently show encap does more work (higher B/op), even if ns/op occasionally looks lower.
- GC overhead is inferred, not isolated:
  - Higher B/op and allocs/op imply more GC pressure. Any GC time is folded into ns/op; we don’t break it out explicitly here.
- Translation to SQL/RDBMS:
  - Database drivers (e.g., pgx, database/sql) often reuse buffers and can scan into caller-provided structs, reducing absolute allocations relative to JSON approach shown here. This may make the relative transform cost a bit more visible, but it remains small in absolute terms.
  - With public fields, you can scan directly into the domain struct (akin to direct/directflat). With strict encapsulation, you typically scan into persistence DTOs and then map to the domain, matching the encap path here.
  - ORMs add mapping layers of their own; this can widen the gap in favor of persistence-aligned (directflat) models.
- Practical guidance:
  - Optimize the persistence path first (prepared statements, batching, buffer reuse, zero-copy scanning/decoding). These changes yield much larger wins than eliminating the O(1) transforms used for encapsulation.
  - In most systems, the benefits of encapsulation (enabling the model to have an independent shape, decoupled from persistence concerns) far outweigh the theoretical nanoseconds or microseconds of overhead it may incur.
