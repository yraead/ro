# HyperLogLog Plugin

The HyperLogLog plugin provides operators for approximate distinct counting using the HyperLogLog algorithm.

## Installation

```bash
go get github.com/samber/ro/plugins/hyperloglog
```

## Usage

```go
package main

import (
    rohyperloglog "github.com/samber/ro/plugins/hyperloglog"
    "github.com/samber/ro"
)

func main() {
    // Count distinct strings
    observable := ro.Just("a", "b", "a", "c", "b", "a")
    count := observable.Pipe(
        rohyperloglog.CountDistinct[string](
            8,
            true,
            rohyperloglog.StringHash.FNV64a()
        ),
    )
    
    fmt.Printf("Distinct count: %d\n", count) // Output: 3
}
```


## Hash Algorithm Benchmarks

The plugin includes comprehensive benchmarks for all hash algorithms. Results show performance and memory usage for both string and bytes hashing.

### Running Benchmarks

To run all benchmarks:

```bash
cd plugins/hyperloglog
go test -bench=. -benchmem
```

To run specific benchmark suites:

```bash
# String hashers
go test -bench=BenchmarkStringHashers -benchmem

# Bytes hashers  
go test -bench=BenchmarkBytesHashers -benchmem
```

### Sample Results

Recent benchmark results (Apple M3, Go 1.24):

**String Hashers (1000 items per iteration):**
- MapHash: ~122,622 ops/sec, 0 B/op
- Adler32: ~110,823 ops/sec, 0 B/op
- Loselose: ~89,503 ops/sec, 0 B/op  
- CRC32: ~65,228 ops/sec, 48KB/op
- DJB2: ~57,729 ops/sec, 0 B/op
- FNV64a: ~20,853 ops/sec, 56KB/op

**Bytes Hashers (1000 items per iteration):**
- CRC32: ~209,274 ops/sec, 0 B/op
- DJB2: ~190,988 ops/sec, 0 B/op
- Adler32: ~186,105 ops/sec, 0 B/op
- Loselose: ~174,127 ops/sec, 0 B/op
- Jenkins: ~147,613 ops/sec, 0 B/op
- SDBM: ~139,095 ops/sec, 0 B/op

### Collision Analysis

Theoretical collision probabilities based on hash size and birthday paradox:

**64-bit Hashes (FNV64a, FNV64, CRC64, MapHash):**
- 50% collision probability: ~5.1 billion items (2^32.5)
- 1% collision probability: ~671 million items
- 0.1% collision probability: ~67 million items

**32-bit Hashes (FNV32a, FNV32, CRC32, Adler32, Jenkins, DJB2, SDBM):**
- 50% collision probability: ~77,000 items (2^16.5)
- 1% collision probability: ~9,300 items
- 0.1% collision probability: ~930 items

**Cryptographic Hashes (SHA256, SHA512, SHA1, MD5):**
- SHA256 (256-bit): 50% collision probability at ~2^128 items
- SHA512 (512-bit): 50% collision probability at ~2^256 items
- SHA1 (160-bit): 50% collision probability at ~2^80 items
- MD5 (128-bit): 50% collision probability at ~2^64 items

**Loselose Hash:**
- **Extremely poor collision resistance** - known to have massive collision rates
- In practice: 56% collision rate on 1000 unique strings, 33% on 1000 unique bytes
- **Avoid for any serious application**

### Recommendations

- **Fastest**: MapHash and Loselose for strings, CRC32 and Loselose for bytes
- **Best Balance**: FNV64a offers good speed and collision resistance (64-bit space)
- **Cryptographic Security**: SHA256/SHA512 for security-critical applications
- **Avoid**: Loselose has extremely poor collision resistance
- **General Purpose**: FNV64a or MapHash for most use cases
