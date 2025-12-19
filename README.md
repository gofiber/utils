# Fiber Utils

![Release](https://img.shields.io/github/release/gofiber/utils.svg)
![Test](https://github.com/gofiber/utils/workflows/Test/badge.svg)
![Codecov](https://img.shields.io/codecov/c/github/gofiber/utils?token=3Cr92CwaPQ&style=flat-square&logo=codecov&label=codecov)
![Linter](https://github.com/gofiber/utils/workflows/Linter/badge.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)

A collection of common functions for [Fiber](https://github.com/gofiber/fiber) with better performance, fewer allocations, and fewer dependencies.

## Benchmarks

Environment:
goos: darwin
goarch: amd64
pkg: github.com/gofiber/utils
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

```go
// go test -benchmem -run=^$ -bench=Benchmark_ -count=2

Benchmark_ToLowerBytes/fiber-12             29715831   36.44  ns/op    0  B/op  0  allocs/op
Benchmark_ToLowerBytes/fiber-12             33316479   36.28  ns/op    0  B/op  0  allocs/op
Benchmark_ToLowerBytes/default-12           11894427   96.98  ns/op   80  B/op  1  allocs/op
Benchmark_ToLowerBytes/default-12           12217050   97.43  ns/op   80  B/op  1  allocs/op

Benchmark_ToUpperBytes/fiber-12             22042162   46.92  ns/op    0  B/op  0  allocs/op
Benchmark_ToUpperBytes/fiber-12             25859680   46.43  ns/op    0  B/op  0  allocs/op
Benchmark_ToUpperBytes/default-12           10015346   117.2  ns/op   80  B/op  1  allocs/op
Benchmark_ToUpperBytes/default-12           10185375   117.8  ns/op   80  B/op  1  allocs/op

Benchmark_TrimRight/fiber-12               522399795   2.138  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRight/fiber-12               578245326   2.084  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRight/default-12             345155300   3.380  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRight/default-12             366972850   3.328  ns/op    0  B/op  0  allocs/op

Benchmark_TrimRightBytes/fiber-12          586471208   2.099  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRightBytes/fiber-12          576055069   2.087  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRightBytes/default-12        348849292   3.316  ns/op    0  B/op  0  allocs/op
Benchmark_TrimRightBytes/default-12        359904445   3.384  ns/op    0  B/op  0  allocs/op

Benchmark_TrimLeft/fiber-12                578044544   2.122  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeft/fiber-12                585290433   2.074  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeft/default-12              351906888   3.667  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeft/default-12              330852666   3.448  ns/op    0  B/op  0  allocs/op

Benchmark_TrimLeftBytes/fiber-12           545400109   2.239  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeftBytes/fiber-12           544800061   2.270  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeftBytes/default-12         329749006   3.521  ns/op    0  B/op  0  allocs/op
Benchmark_TrimLeftBytes/default-12         344199560   3.452  ns/op    0  B/op  0  allocs/op

Benchmark_Trim/fiber-12                    280692232   4.128  ns/op    0  B/op  0  allocs/op
Benchmark_Trim/fiber-12                    297070083   3.961  ns/op    0  B/op  0  allocs/op
Benchmark_Trim/default-12                  232522952   5.163  ns/op    0  B/op  0  allocs/op
Benchmark_Trim/default-12                  230659057   5.172  ns/op    0  B/op  0  allocs/op
Benchmark_Trim/default.trimspace-12        227328967   5.245  ns/op    0  B/op  0  allocs/op
Benchmark_Trim/default.trimspace-12        227340775   5.253  ns/op    0  B/op  0  allocs/op

Benchmark_TrimBytes/fiber-12               275612090   4.280  ns/op    0  B/op  0  allocs/op
Benchmark_TrimBytes/fiber-12               284892168   4.302  ns/op    0  B/op  0  allocs/op
Benchmark_TrimBytes/default-12             224021550   5.163  ns/op    0  B/op  0  allocs/op
Benchmark_TrimBytes/default-12             239689282   4.922  ns/op    0  B/op  0  allocs/op
Benchmark_TrimBytes/default.trimspace-12   216809300   5.514  ns/op    0  B/op  0  allocs/op
Benchmark_TrimBytes/default.trimspace-12   218177734   5.603  ns/op    0  B/op  0  allocs/op

Benchmark_EqualFoldBytes/fiber-12           22944849   47.14  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFoldBytes/fiber-12           26006342   46.82  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFoldBytes/default-12          5222006   222.5  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFoldBytes/default-12          5349398   223.2  ns/op    0  B/op  0  allocs/op

Benchmark_EqualFold/fiber-12                24761037   48.63  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFold/fiber-12                24159073   48.63  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFold/default-12               6322188   191.5  ns/op    0  B/op  0  allocs/op
Benchmark_EqualFold/default-12               6319070   193.5  ns/op    0  B/op  0  allocs/op

Benchmark_UUID/fiber-12                     22061482   49.13  ns/op   48  B/op  1  allocs/op
Benchmark_UUID/fiber-12                     24123198   48.40  ns/op   48  B/op  1  allocs/op
Benchmark_UUID/default-12                    3581961   336.9  ns/op  168  B/op  6  allocs/op
Benchmark_UUID/default-12                    3465946   344.8  ns/op  168  B/op  6  allocs/op

Benchmark_ConvertToBytes/fiber-12           53392819   23.19  ns/op    0  B/op  0  allocs/op
Benchmark_ConvertToBytes/fiber-12           51117225   23.32  ns/op    0  B/op  0  allocs/op

Benchmark_UnsafeString/unsafe-12          1000000000  0.5672  ns/op    0  B/op  0  allocs/op
Benchmark_UnsafeString/unsafe-12          1000000000  0.5683  ns/op    0  B/op  0  allocs/op
Benchmark_UnsafeString/default-12           64000897   18.45  ns/op   16  B/op  1  allocs/op
Benchmark_UnsafeString/default-12           64138909   18.13  ns/op   16  B/op  1  allocs/op

Benchmark_UnsafeBytes/unsafe-12            474777096   2.539  ns/op    0  B/op  0  allocs/op
Benchmark_UnsafeBytes/unsafe-12            469340781   2.535  ns/op    0  B/op  0  allocs/op
Benchmark_UnsafeBytes/default-12            53125656   22.15  ns/op   16  B/op  1  allocs/op
Benchmark_UnsafeBytes/default-12            52615048   22.33  ns/op   16  B/op  1  allocs/op

Benchmark_ToString-12                       22981430   51.72  ns/op   40  B/op  2  allocs/op
Benchmark_ToString-12                       22956476   52.93  ns/op   40  B/op  2  allocs/op

Benchmark_GetMIME/fiber-12                  15782622   74.99  ns/op    0  B/op  0  allocs/op
Benchmark_GetMIME/fiber-12                  13992375   93.13  ns/op    0  B/op  0  allocs/op
Benchmark_GetMIME/default-12                 6825952   147.0  ns/op    0  B/op  0  allocs/op
Benchmark_GetMIME/default-12                 9158227   132.5  ns/op    0  B/op  0  allocs/op

ParseVendorSpecificContentType
Benchmark_Parse.../vendorContentType-12     21334663   50.24  ns/op   16  B/op  1  allocs/op
Benchmark_Parse.../vendorContentType-12     23121808   51.20  ns/op   16  B/op  1  allocs/op
Benchmark_Parse.../defaultContentType-12   154423909   6.772  ns/op    0  B/op  0  allocs/op
Benchmark_Parse.../defaultContentType-12   183285117   6.662  ns/op    0  B/op  0  allocs/op

Benchmark_StatusMessage/fiber-12          1000000000  0.9796  ns/op    0  B/op  0  allocs/op
Benchmark_StatusMessage/fiber-12          1000000000  0.9706  ns/op    0  B/op  0  allocs/op
Benchmark_StatusMessage/default-12         380260562   2.989  ns/op    0  B/op  0  allocs/op
Benchmark_StatusMessage/default-12         403639642   3.124  ns/op    0  B/op  0  allocs/op

Benchmark_IsIPv4/fiber-12                   53576214   21.07  ns/op    0  B/op  0  allocs/op
Benchmark_IsIPv4/fiber-12                   62672907   22.04  ns/op    0  B/op  0  allocs/op
Benchmark_IsIPv4/default-12                 21204613   62.23  ns/op   16  B/op  1  allocs/op
Benchmark_IsIPv4/default-12                 21399847   56.61  ns/op   16  B/op  1  allocs/op

Benchmark_IsIPv6/fiber-12                   16754995   72.46  ns/op    0  B/op  0  allocs/op
Benchmark_IsIPv6/fiber-12                   17080897   74.31  ns/op    0  B/op  0  allocs/op
Benchmark_IsIPv6/default-12                  8160195   124.5  ns/op   16  B/op  1  allocs/op
Benchmark_IsIPv6/default-12                  9415326   119.8  ns/op   16  B/op  1  allocs/op

Benchmark_ToUpper/fiber-12                  13175154   81.67  ns/op   80  B/op  1  allocs/op
Benchmark_ToUpper/fiber-12                  14285533   77.27  ns/op   80  B/op  1  allocs/op
Benchmark_ToUpper/default-12                 5332206   231.8  ns/op   80  B/op  1  allocs/op
Benchmark_ToUpper/default-12                 5364650   236.0  ns/op   80  B/op  1  allocs/op

Benchmark_ToLower/fiber-12                  12996409   80.24  ns/op   80  B/op  1  allocs/op
Benchmark_ToLower/fiber-12                  16539536   69.27  ns/op   80  B/op  1  allocs/op
Benchmark_ToLower/default-12                 5132185   222.5  ns/op   80  B/op  1  allocs/op
Benchmark_ToLower/default-12                 5158561   225.3  ns/op   80  B/op  1  allocs/op

Benchmark_CalculateTimestamp/fiber-12     1000000000  0.2634  ns/op    0  B/op  0  allocs/op
Benchmark_CalculateTimestamp/fiber-12     1000000000  0.2935  ns/op    0  B/op  0  allocs/op
Benchmark_CalculateTimestamp/default-12     15740576   73.79  ns/op    0  B/op  0  allocs/op
Benchmark_CalculateTimestamp/default-12     15789036   71.12  ns/op    0  B/op  0  allocs/op

Benchmark_ParseUint/fiber-12               172685805   6.949  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint/fiber-12               172395474   6.984  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint/fiber_bytes-12         167209564   7.162  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint/fiber_bytes-12         161316883   7.248  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint/default-12              89826800   13.52  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint/default-12              89406565   13.68  ns/op    0  B/op  0  allocs/op

Benchmark_ParseInt/fiber-12                158532987   7.442  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt/fiber-12                154777971   7.710  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt/fiber_bytes-12          157400030   7.453  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt/fiber_bytes-12          148624418   7.400  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt/default-12               78927678   15.97  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt/default-12               78261080   16.26  ns/op    0  B/op  0  allocs/op

Benchmark_ParseInt32/fiber-12              164397682   7.234  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt32/fiber-12              164329150   7.319  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12        159049194   7.556  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12        155705494   7.697  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt32/default-12             71818300   16.92  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt32/default-12             78039262   15.83  ns/op    0  B/op  0  allocs/op

Benchmark_ParseInt16/fiber-12              232986352   5.165  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt16/fiber-12              232214074   5.256  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12        220412829   5.398  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12        222333234   5.409  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt16/default-12            100000000   11.07  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt16/default-12            100000000   11.08  ns/op    0  B/op  0  allocs/op

Benchmark_ParseInt8/fiber-12               260329051   4.543  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt8/fiber-12               265292354   4.541  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12         260297640   4.635  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12         260662333   4.669  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt8/default-12             134202080   8.700  ns/op    0  B/op  0  allocs/op
Benchmark_ParseInt8/default-12             137497462   8.702  ns/op    0  B/op  0  allocs/op

Benchmark_ParseUint32/fiber-12             166919528   6.991  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint32/fiber-12             172230549   7.004  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12       168104906   7.182  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12       166743417   7.189  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint32/default-12            88639659   13.70  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint32/default-12            78153198   13.74  ns/op    0  B/op  0  allocs/op

Benchmark_ParseUint16/fiber-12             265107002   4.425  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint16/fiber-12             265636831   4.517  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12       255349777   4.674  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12       250084923   4.722  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint16/default-12           133589893   9.006  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint16/default-12           136365088   8.863  ns/op    0  B/op  0  allocs/op

Benchmark_ParseUint8/fiber-12              326680580   3.719  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint8/fiber-12              313552454   3.739  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12        325318082   3.697  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12        321770954   3.699  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint8/default-12            182636923   6.678  ns/op    0  B/op  0  allocs/op
Benchmark_ParseUint8/default-12            184060842   6.756  ns/op    0  B/op  0  allocs/op

Benchmark_ParseFloat64/fiber-12            100000000   11.11  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat64/fiber-12            100000000   11.09  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12      127174759   9.525  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12      124214686   9.577  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat64/default-12           50066755   24.16  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat64/default-12           49396011   24.80  ns/op    0  B/op  0  allocs/op

Benchmark_ParseFloat32/fiber-12            100000000   11.71  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat32/fiber-12            100000000   11.85  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12      100000000   11.17  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12      100000000   10.43  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat32/default-12           46055755   25.65  ns/op    0  B/op  0  allocs/op
Benchmark_ParseFloat32/default-12           45263090   25.78  ns/op    0  B/op  0  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/>
