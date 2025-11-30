# TrimSpace Performance Analysis

## String Benchmarks (Fiber vs Standard Library)

| Scenario | Fiber Time | Std Time | Speedup | Throughput Gain |
|----------|------------|----------|---------|-----------------|
| **empty** | 0.63 ns | 2.88 ns | **4.6x faster** ⚡ | - |
| **no-trim** | 2.83 ns | 3.80 ns | **34% faster** | 1764 vs 1314 MB/s |
| **no-trim-long** | 2.81 ns | 3.75 ns | **33% faster** | 12454 vs 9335 MB/s |
| **leading-space** | 3.76 ns | 4.15 ns | **10% faster** | 1598 vs 1444 MB/s |
| **trailing-space** | 4.87 ns | 4.91 ns | **Similar** | ~1230 MB/s |
| **both-spaces** | 5.16 ns | 5.63 ns | **8% faster** | 1355 vs 1242 MB/s |
| **both-many-spaces** | 9.63 ns | 9.27 ns | **Similar** | ~2000 MB/s |
| **all-spaces** | 5.46 ns | 10.12 ns | **46% faster** ⚡ | 1831 vs 988 MB/s |
| **medium-no-trim** | 2.92 ns | 3.79 ns | **30% faster** | 21953 vs 16895 MB/s |
| **large-no-trim** | 2.92 ns | 3.75 ns | **28% faster** | 85609 vs 66758 MB/s |

## Byte Slice Benchmarks (Fiber vs Standard Library)

| Scenario | Fiber Time | Std Time | Speedup | Throughput Gain |
|----------|------------|----------|---------|-----------------|
| **empty** | 0.47 ns | 3.04 ns | **6.4x faster** ⚡⚡ | - |
| **no-trim** | 2.83 ns | 4.35 ns | **54% faster** | 1766 vs 1149 MB/s |
| **no-trim-long** | 2.78 ns | 4.25 ns | **53% faster** | 12571 vs 8227 MB/s |
| **trailing-space** | 3.50 ns | 5.20 ns | **49% faster** | 1716 vs 1153 MB/s |
| **both-spaces** | 5.55 ns | 6.36 ns | **13% faster** | 1261 vs 1100 MB/s |
| **both-many-spaces** | 9.64 ns | 11.40 ns | **18% faster** | 1971 vs 1666 MB/s |
| **all-spaces** | 8.09 ns | 10.02 ns | **24% faster** | 1236 vs 997 MB/s |
| **medium-no-trim** | 2.91 ns | 4.29 ns | **47% faster** | 21982 vs 14923 MB/s |
| **large-no-trim** | 2.83 ns | 4.37 ns | **54% faster** | 88404 vs 57178 MB/s |
| **large-with-trim** | 9.59 ns | 11.22 ns | **17% faster** | 26905 vs 22986 MB/s |

## Key Findings

### 🚀 Massive Improvements for Common Cases
- **Empty strings**: 4.6x faster (strings), 6.4x faster (bytes)
- **No trimming needed**: 28-34% faster (strings), 47-54% faster (bytes)
- **Large inputs without trim**: Throughput > 85 GB/s for strings, > 88 GB/s for bytes

### 📊 Performance Characteristics
1. **Best case**: Empty inputs and no-trim scenarios (most common in web frameworks)
2. **Good performance**: All-whitespace removal scenarios (24-46% faster)
3. **Competitive**: Mixed whitespace scenarios (similar or slightly better)

### 💡 Why It's Faster
1. **Lookup table** for O(1) whitespace detection
2. **Zero allocations** when no trimming needed
3. **Early returns** for edge cases
4. **Efficient byte-level operations**
5. **Generic implementation** works for both string and []byte

### 🎯 Use Cases
Perfect for:
- HTTP header parsing
- Query parameter processing
- Path normalization
- Cookie value handling
- Form data processing

All common web framework operations benefit from these optimizations!
