# Fiber Utils

![Release](https://img.shields.io/github/release/gofiber/utils.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/gofiber/utils/workflows/Test/badge.svg)
![Security](https://github.com/gofiber/utils/workflows/Security/badge.svg)
![Linter](https://github.com/gofiber/utils/workflows/Linter/badge.svg)

A collection of common functions but with better performance, less allocations and less dependencies created for [Fiber](https://github.com/gofiber/fiber).

## Benchmarks

Environment:
goos: darwin
goarch: amd64
pkg: github.com/gofiber/utils
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

```go
// go test -benchmem -run=^$ -bench=Benchmark_ -count=2

Benchmark_ToLowerBytes/fiber-12                29715831        36.44 ns/op       0 B/op       0 allocs/op
Benchmark_ToLowerBytes/fiber-12                33316479        36.28 ns/op       0 B/op       0 allocs/op
Benchmark_ToLowerBytes/default-12              11894427        96.98 ns/op      80 B/op       1 allocs/op
Benchmark_ToLowerBytes/default-12              12217050        97.43 ns/op      80 B/op       1 allocs/op

Benchmark_ToUpperBytes/fiber-12                22042162        46.92 ns/op       0 B/op       0 allocs/op
Benchmark_ToUpperBytes/fiber-12                25859680        46.43 ns/op       0 B/op       0 allocs/op
Benchmark_ToUpperBytes/default-12              10015346        117.2 ns/op      80 B/op       1 allocs/op
Benchmark_ToUpperBytes/default-12              10185375        117.8 ns/op      80 B/op       1 allocs/op

Benchmark_EqualFoldBytes/fiber-12              22944849        47.14 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFoldBytes/fiber-12              26006342        46.82 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFoldBytes/default-12             5222006        222.5 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFoldBytes/default-12             5349398        223.2 ns/op       0 B/op       0 allocs/op

Benchmark_EqualFold/fiber-12                   24761037        48.63 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFold/fiber-12                   24159073        48.63 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFold/default-12                  6322188        191.5 ns/op       0 B/op       0 allocs/op
Benchmark_EqualFold/default-12                  6319070        193.5 ns/op       0 B/op       0 allocs/op

Benchmark_UUID/fiber-12                        22061482        49.13 ns/op      48 B/op       1 allocs/op
Benchmark_UUID/fiber-12                        24123198        48.40 ns/op      48 B/op       1 allocs/op
Benchmark_UUID/default-12                       3581961        336.9 ns/op     168 B/op       6 allocs/op
Benchmark_UUID/default-12                       3465946        344.8 ns/op     168 B/op       6 allocs/op

Benchmark_ConvertToBytes/fiber-12              53392819        23.19 ns/op       0 B/op       0 allocs/op
Benchmark_ConvertToBytes/fiber-12              51117225        23.32 ns/op       0 B/op       0 allocs/op

Benchmark_UnsafeString/unsafe-12               1000000000     0.5672 ns/op       0 B/op       0 allocs/op
Benchmark_UnsafeString/unsafe-12               1000000000     0.5683 ns/op       0 B/op       0 allocs/op
Benchmark_UnsafeString/default-12              64000897        18.45 ns/op      16 B/op       1 allocs/op
Benchmark_UnsafeString/default-12              64138909        18.13 ns/op      16 B/op       1 allocs/op

Benchmark_UnsafeBytes/unsafe-12                474777096       2.539 ns/op       0 B/op       0 allocs/op
Benchmark_UnsafeBytes/unsafe-12                469340781       2.535 ns/op       0 B/op       0 allocs/op
Benchmark_UnsafeBytes/default-12               53125656        22.15 ns/op      16 B/op       1 allocs/op
Benchmark_UnsafeBytes/default-12               52615048        22.33 ns/op      16 B/op       1 allocs/op

Benchmark_ToString-12                          22981430        51.72 ns/op      40 B/op       2 allocs/op
Benchmark_ToString-12                          22956476        52.93 ns/op      40 B/op       2 allocs/op

Benchmark_GetMIME/fiber-12                     15782622        74.99 ns/op       0 B/op       0 allocs/op
Benchmark_GetMIME/fiber-12                     13992375        93.13 ns/op       0 B/op       0 allocs/op
Benchmark_GetMIME/default-12                    6825952        147.0 ns/op       0 B/op       0 allocs/op
Benchmark_GetMIME/default-12                    9158227        132.5 ns/op       0 B/op       0 allocs/op

ParseVendorSpecificContentType
Benchmark_Parse.../vendorContentType-12        21334663        50.24 ns/op      16 B/op       1 allocs/op
Benchmark_Parse.../vendorContentType-12        23121808        51.20 ns/op      16 B/op       1 allocs/op
Benchmark_Parse.../defaultContentType-12       154423909       6.772 ns/op       0 B/op       0 allocs/op
Benchmark_Parse.../defaultContentType-12       183285117       6.662 ns/op       0 B/op       0 allocs/op

Benchmark_StatusMessage/fiber-12             1000000000       0.9796 ns/op       0 B/op       0 allocs/op
Benchmark_StatusMessage/fiber-12             1000000000       0.9706 ns/op       0 B/op       0 allocs/op
Benchmark_StatusMessage/default-12            380260562        2.989 ns/op       0 B/op       0 allocs/op
Benchmark_StatusMessage/default-12            403639642        3.124 ns/op       0 B/op       0 allocs/op

Benchmark_IsIPv4/fiber-12                      53576214        21.07 ns/op       0 B/op       0 allocs/op
Benchmark_IsIPv4/fiber-12                      62672907        22.04 ns/op       0 B/op       0 allocs/op
Benchmark_IsIPv4/default-12                    21204613        62.23 ns/op      16 B/op       1 allocs/op
Benchmark_IsIPv4/default-12                    21399847        56.61 ns/op      16 B/op       1 allocs/op

Benchmark_IsIPv6/fiber-12                      16754995        72.46 ns/op       0 B/op       0 allocs/op
Benchmark_IsIPv6/fiber-12                      17080897        74.31 ns/op       0 B/op       0 allocs/op
Benchmark_IsIPv6/default-12                     8160195        124.5 ns/op      16 B/op       1 allocs/op
Benchmark_IsIPv6/default-12                     9415326        119.8 ns/op      16 B/op       1 allocs/op

Benchmark_ToUpper/fiber-12                     13175154        81.67 ns/op      80 B/op       1 allocs/op
Benchmark_ToUpper/fiber-12                     14285533        77.27 ns/op      80 B/op       1 allocs/op
Benchmark_ToUpper/default-12                    5332206        231.8 ns/op      80 B/op       1 allocs/op
Benchmark_ToUpper/default-12                    5364650        236.0 ns/op      80 B/op       1 allocs/op

Benchmark_ToLower/fiber-12                     12996409        80.24 ns/op      80 B/op       1 allocs/op
Benchmark_ToLower/fiber-12                     16539536        69.27 ns/op      80 B/op       1 allocs/op
Benchmark_ToLower/default-12                    5132185        222.5 ns/op      80 B/op       1 allocs/op
Benchmark_ToLower/default-12                    5158561        225.3 ns/op      80 B/op       1 allocs/op

Benchmark_CalculateTimestamp/fiber-12        1000000000       0.2634 ns/op       0 B/op       0 allocs/op
Benchmark_CalculateTimestamp/fiber-12        1000000000       0.2935 ns/op       0 B/op       0 allocs/op
Benchmark_CalculateTimestamp/default-12        15740576        73.79 ns/op       0 B/op       0 allocs/op
Benchmark_CalculateTimestamp/default-12        15789036        71.12 ns/op       0 B/op       0 allocs/op
```

See all the benchmarks under https://gofiber.github.io/utils/