package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type stats struct {
	ns     []float64
	bytes  []float64
	allocs []float64
}

func (s *stats) add(ns, bytes, allocs float64) {
	s.ns = append(s.ns, ns)
	s.bytes = append(s.bytes, bytes)
	s.allocs = append(s.allocs, allocs)
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	cp := make([]float64, len(values))
	copy(cp, values)
	sort.Float64s(cp)
	n := len(cp)
	if n%2 == 1 {
		return cp[n/2]
	}
	return (cp[n/2-1] + cp[n/2]) / 2
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

var (
	benchLinePrefix = "Benchmark"
	valueBeforeUnit = regexp.MustCompile(`([0-9]+\.?[0-9]*)\s+(ns/op|B/op|allocs/op)`) // captures value and unit
)

func parseBenchLine(line string) (name string, ns, bytes, allocs float64, ok bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, benchLinePrefix) {
		return "", 0, 0, 0, false
	}
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return "", 0, 0, 0, false
	}
	name = fields[0]
	// Find value preceding units using regex to be robust to spacing
	matches := valueBeforeUnit.FindAllStringSubmatch(line, -1)
	// Expect three metrics per line
	for _, m := range matches {
		if len(m) != 3 {
			continue
		}
		valStr := m[1]
		unit := m[2]
		v, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}
		switch unit {
		case "ns/op":
			ns = v
		case "B/op":
			bytes = v
		case "allocs/op":
			allocs = v
		}
	}
	if ns == 0 && bytes == 0 && allocs == 0 {
		return "", 0, 0, 0, false
	}
	return name, ns, bytes, allocs, true
}

func main() {
	file := flag.String("file", "bench_results_stable.txt", "path to benchmark results file")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	byName := make(map[string]*stats)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		name, ns, bytes, allocs, ok := parseBenchLine(line)
		if !ok {
			continue
		}
		st := byName[name]
		if st == nil {
			st = &stats{}
			byName[name] = st
		}
		st.add(ns, bytes, allocs)
	}
	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(1)
	}

	// Deterministic output: sort names
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Printf("%-34s  %12s  %12s  %12s  |  %12s  %12s  %12s  |  %s\n", "BENCHMARK", "med ns/op", "med B/op", "med allocs", "mean ns/op", "mean B/op", "mean allocs", "n")
	for _, name := range names {
		st := byName[name]
		fmt.Printf("%-34s  %12.3f  %12.0f  %12.0f  |  %12.3f  %12.0f  %12.2f  |  %d\n",
			name,
			median(st.ns), median(st.bytes), median(st.allocs),
			mean(st.ns), mean(st.bytes), mean(st.allocs),
			len(st.ns),
		)
	}
}
