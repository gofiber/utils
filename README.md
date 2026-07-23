# Fiber Utils

![Release](https://img.shields.io/github/release/gofiber/utils.svg)
![Test](https://github.com/gofiber/utils/workflows/Test/badge.svg)
![Codecov](https://img.shields.io/codecov/c/github/gofiber/utils?token=3Cr92CwaPQ&style=flat-square&logo=codecov&label=codecov)
![Linter](https://github.com/gofiber/utils/actions/workflows/lint.yml/badge.svg)
[![Benchmarks](https://img.shields.io/badge/%F0%9F%93%8A%20benchmarks-charts-00ACD7.svg)](https://gofiber.github.io/utils/benchmarks/)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)

A collection of common functions for [Fiber](https://github.com/gofiber/fiber) with better performance, fewer allocations, and fewer dependencies.

## Benchmarks

Per-commit benchmark charts: <https://gofiber.github.io/utils/benchmarks/>

Environment:
goos: darwin
goarch: arm64
pkg: github.com/gofiber/utils/v2
cpu: Apple M2 Pro

```text
// go test ./... -benchmem -run=^$ -bench=Benchmark_ -count=1
// (swar's Benchmark_Load8_Fusion, an advisory codegen guard, is deliberately omitted)

# Case Conversion
Benchmark_ToLowerBytes/empty/fiber-12                              599341221    1.970  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/unsafe-12                       613673802    1.951  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                            363300402    3.316  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                           350302147    3.434  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/unsafe-12                    250488015    4.803  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                          65883986    18.03  ns/op     8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                      95234625    12.53  ns/op     3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/unsafe-12              242418936    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                    86405012    13.85  ns/op     8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12           52263289    22.74  ns/op    48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/unsafe-12   131711175    9.101  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12         25312339    47.13  ns/op    48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                        120148472    10.02  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/unsafe-12                 100000000    10.57  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                       16746064    71.43  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                         46270634    26.82  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/unsafe-12                 100000000    10.71  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                       15734505    76.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                         46607823    26.06  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/unsafe-12                 100000000    10.86  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                       16291084    73.76  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/empty/fiber-12                              607678140    1.976  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/unsafe-12                       611928006    1.981  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                            362632894    3.309  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                            95358867    12.55  ns/op     3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/unsafe-12                    249490471    4.905  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                          79248020    14.63  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                     365516280    3.330  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/unsafe-12              249981619    4.803  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                    64343162    19.01  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12           52506741    22.69  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/unsafe-12   130478352    9.350  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12         22415607    54.44  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                         44765679    26.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/unsafe-12                 100000000    13.68  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                       14464340    110.7  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                        100000000    11.98  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/unsafe-12                  83808488    13.94  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                       13142096    87.35  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                         33038661    32.75  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/unsafe-12                 100000000    10.93  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                       15034422    98.27  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/empty/fiber-12                                   595298628    2.309  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/unsafe-12                            595273404    2.021  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                                 636980288    1.883  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                                 97025577    12.40  ns/op     3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/unsafe-12                         244711635    4.857  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                               59790234    19.54  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                          383218539    3.056  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/unsafe-12                   247851331    4.864  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                        223044249    5.182  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12                52628502    24.13  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/unsafe-12        127726341    9.233  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12              11549816    104.0  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                              42409070    26.04  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/unsafe-12                      100000000    11.00  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                             7051911    169.0  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                             123837582    9.590  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/unsafe-12                      100000000    10.83  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                            19115356    61.51  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                              45872143    27.77  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/unsafe-12                      100000000    11.00  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                             5927710    203.1  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/empty/fiber-12                                   613625554    1.952  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/unsafe-12                            572670608    1.989  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                                 644769141    1.863  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                                395302272    3.034  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/unsafe-12                         252231990    4.813  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                              284202345    4.219  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                          100000000    12.08  ns/op     3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/unsafe-12                   252859507    4.816  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                         60415860    19.54  ns/op     8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12                53446660    22.49  ns/op    48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/unsafe-12        131152108    9.153  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12              18039788    66.79  ns/op    48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                             120227584    9.982  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/unsafe-12                      100000000    10.68  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                            24070928    48.79  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                              44699118    25.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/unsafe-12                      100000000    10.74  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                             7075200    170.0  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                              46426453    26.01  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/unsafe-12                      100000000    10.74  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                             5682778    204.3  ns/op    64  B/op   1  allocs/op

# Add Trailing Slash
Benchmark_AddTrailingSlashBytes/empty-12                          1000000000   0.5857  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                     1000000000   0.8694  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                  100000000    11.65  ns/op     4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12               1000000000   0.8867  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                    93543546    13.36  ns/op    16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12                1000000000   0.9011  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/empty-12                         1000000000   0.2926  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                    1000000000   0.4482  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12                 100000000    11.37  ns/op     4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12              1000000000   0.4502  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                   82427726    14.41  ns/op    16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12               1000000000   0.4529  ns/op     0  B/op   0  allocs/op

# EqualFold
Benchmark_EqualFoldBytes/fiber-12                                   67368261    18.48  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                                 17774803    65.91  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-12                                        80501566    14.75  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                      18314365    65.96  ns/op     0  B/op   0  allocs/op

# Search
// amd64 note: IndexAny2/IndexAny3 dispatch to package simd for 32+ byte
// inputs; these arm64 rows keep the portable SWAR path. For amd64 numbers
// see the Benchmark_Memchr2/Benchmark_Memchr3 rows in the simd section
// below (the kernels these functions dispatch to).
Benchmark_IndexAny2/7B/swar-12                                     282163779    4.297  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/7B/scalar-12                                   314553402    3.826  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/7B/stdlib-indexany-12                           50085909    23.97  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/swar-12                                     522780714    2.289  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/scalar-12                                   271494799    4.206  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/stdlib-indexany-12                           44831599    26.78  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/swar-12                                    393245086    3.047  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/scalar-12                                  168476984    7.117  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/stdlib-indexany-12                         100000000    10.03  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/swar-12                                    252199587    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/scalar-12                                   84105762    14.27  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/stdlib-indexany-12                          70291048    17.37  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/swar-12                                    144335962    8.290  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/scalar-12                                   42694686    28.60  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/stdlib-indexany-12                          37615540    32.23  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/swar-12                                    20231538    59.71  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/scalar-12                                   5071586    233.6  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/stdlib-indexany-12                          4448119    274.2  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/7B/swar-12                                     221889592    5.397  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/7B/scalar-12                                   240859849    4.977  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/8B/swar-12                                     445742876    2.689  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/8B/scalar-12                                   211981612    5.620  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/16B/swar-12                                    326475568    3.739  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/16B/scalar-12                                  100000000    10.31  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/32B/swar-12                                    196365920    6.112  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/32B/scalar-12                                   60021255    19.81  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/64B/swar-12                                    100000000    10.98  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/64B/scalar-12                                   31180405    38.97  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/512B/swar-12                                    15466316    79.22  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/512B/scalar-12                                   3911382    307.1  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-8B/swar-12                                332497509    3.634  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-8B/scalar-12                              491493494    2.408  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-32B/swar-12                               149769112    7.928  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-32B/scalar-12                              44041342    27.13  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-32B/swar-12                                138732252    8.758  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-32B/scalar-12                               33871633    30.73  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-tail/swar-12                               100000000    12.34  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-tail/scalar-12                              16138168    75.39  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-64B/swar-12                               100000000    11.37  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-64B/scalar-12                              17902330    67.78  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-512B/swar-12                               17640939    67.65  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-512B/scalar-12                              2211582    544.4  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-hit-64B/swar-12                           100000000    11.58  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-hit-64B/scalar-12                          23906326    51.33  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-miss-64B/swar-12                          100000000    11.33  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-miss-64B/scalar-12                         23553600    50.27  ns/op     0  B/op   0  allocs/op

# ASCII
// amd64 note: IsASCII dispatches to package simd for 32+ byte inputs;
// these arm64 rows keep the portable SWAR path. See the simd section
// below for amd64 numbers.
Benchmark_IsASCII/7B/swar-12                                       238699563    5.001  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/7B/scalar-12                                     292669958    4.064  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/swar-12                                       439961722    2.687  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/scalar-12                                     274732192    4.364  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/16B/swar-12                                      367179319    3.255  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/16B/scalar-12                                    179331987    6.706  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/swar-12                                      368345228    3.271  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/scalar-12                                    100000000    11.62  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/swar-12                                      265505815    4.486  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/scalar-12                                     57864440    21.33  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/swar-12                                      47664915    25.14  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/scalar-12                                     7081790    169.5  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/7B/swar-12                              142234436    8.424  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/7B/scalar-12                            142394628    8.174  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/8B/swar-12                              401721262    3.099  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/8B/scalar-12                            138047568    8.654  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/16B/swar-12                             244592613    5.059  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/16B/scalar-12                            75413396    15.94  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/32B/swar-12                             138857455    8.875  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/32B/scalar-12                            39539425    30.18  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/64B/swar-12                              75193269    16.46  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/64B/scalar-12                            20898807    57.81  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/512B/swar-12                              9634008    124.9  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/512B/scalar-12                            2619062    460.6  ns/op     0  B/op   0  allocs/op

# Trim
Benchmark_TrimRight/fiber-12                                       587360972    2.043  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                     417142952    2.887  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRightBytes/fiber-12                                  590828253    2.067  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                                374548395    3.236  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeft/fiber-12                                        590822072    2.030  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                      408184784    2.967  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/fiber-12                                   592875732    2.033  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                                 405921434    2.967  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/fiber-12                                            317625432    3.758  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/default-12                                          250173384    4.735  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                                241883282    5.022  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/fiber-12                                       308885211    3.858  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                     254064188    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                           224231778    5.330  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/empty-12                                1000000000   0.3017  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                               583064172    2.074  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                                387354229    3.059  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/spaces-12                              320059388    3.651  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                            345078337    3.535  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                          270603267    4.465  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                    291494109    4.141  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                  224087944    5.329  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                     390765345    3.047  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                   258943531    4.718  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                       344369745    3.490  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                     277203638    4.347  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                          344263705    3.506  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/content-type-12                        276213367    4.343  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                       460119685    2.653  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                     296468840    4.058  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                          330053276    3.520  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/query-params-12                        275688495    4.352  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                            275511898    4.342  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                          206143018    5.807  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                              1000000000   0.5914  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                             451842741    2.634  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                      296083767    4.054  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                    224007750    5.294  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/empty-12                           1000000000   0.2914  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                          519488766    2.352  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                           410959255    2.966  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                         315602949    3.827  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                       341217429    3.519  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                     247799318    4.667  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12               295130076    4.182  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12             209162235    5.685  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12                393366796    3.053  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12              240882009    4.930  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                  343673244    3.512  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12                250962042    4.642  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                     344382758    3.516  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                   250622116    4.693  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                  459435145    2.621  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12                276595178    4.359  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                     341144192    3.526  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                   255843174    4.873  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                       270183530    4.445  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                     189740618    6.335  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                         1000000000   0.5969  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                        406072189    2.941  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12                 289382116    4.057  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12               218559462    5.493  ns/op     0  B/op   0  allocs/op

# Convert
Benchmark_ConvertToBytes/fiber-12                                  320332576    3.787  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int-12                                          499429904    2.382  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                         614406423    1.965  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                        486644316    2.463  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                        495227163    2.422  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                        502362236    2.429  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                         561693333    2.101  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                        611239028    1.965  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                       573383422    2.127  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                       575723191    2.093  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                       575901750    2.089  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/string-12                                       613909396    1.959  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                       96832432    12.11  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                         608058547    1.990  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                       23195016    47.11  ns/op     4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                       17204126    70.24  ns/op     4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                     13103935    91.85  ns/op    24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                  13198113    91.22  ns/op    24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                      65187865    18.35  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                         54417859    21.53  ns/op     8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                        32691942    36.57  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                        9382754    128.8  ns/op   112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                                10584483    114.5  ns/op    72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                             365979882    3.277  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                              11649330    102.5  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/int-12                             1000000000   0.2619  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int8-12                            1000000000   0.2202  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int16-12                           1000000000   0.2708  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int32-12                           1000000000   0.2658  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int64-12                           1000000000   0.2658  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint-12                            1000000000   0.2281  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint8-12                           1000000000   0.2223  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint16-12                          1000000000   0.2396  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint32-12                          1000000000   0.2273  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint64-12                          1000000000   0.2271  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/string-12                          1000000000   0.2117  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/[]uint8-12                          213842515    5.718  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/bool-12                            1000000000   0.2118  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/float32-12                          201081819    6.011  ns/op     4  B/op   1  allocs/op
Benchmark_ToString_concurrency/float64-12                          142000966    8.799  ns/op     4  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time-12                         71451783    18.04  ns/op    24  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time#01-12                      69313578    17.74  ns/op    24  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]string-12                         187564670    6.455  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]int-12                            262707805    4.498  ns/op     8  B/op   1  allocs/op
Benchmark_ToString_concurrency/[2]int-12                           129274200    8.777  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[][]int-12                           22638411    52.65  ns/op   112  B/op   6  allocs/op
Benchmark_ToString_concurrency/[]interface_{}-12                    34120837    35.36  ns/op    72  B/op   3  allocs/op
Benchmark_ToString_concurrency/utils.MyStringer-12                1000000000   0.4259  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/utils.CustomType-12                  54806094    25.95  ns/op    16  B/op   1  allocs/op
Benchmark_UnsafeBytes/unsafe-12                                   1000000000   0.5658  ns/op     0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                    88407558    13.50  ns/op    16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                  1000000000   0.3545  ns/op     0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                  100000000    12.09  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/0-12                                            395954607    3.023  ns/op     0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                             65295756    18.15  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                           58264801    20.92  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                          60886531    19.80  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                          53409986    22.69  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                       60969927    19.41  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                       53647869    22.12  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                    63792180    19.05  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                    54560540    22.02  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                                 64080169    18.85  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                                 52957434    22.01  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                              65629659    18.33  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                              56535770    21.73  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                           64569356    18.56  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                           56481996    21.10  ns/op    16  B/op   1  allocs/op

# Format and Append
Benchmark_FormatUint/small/fiber-12                                607527520    1.966  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                              607521366    1.973  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                                64427509    19.04  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/medium/strconv-12                              59566404    20.42  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                                 45395149    26.16  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint/large/strconv-12                               45620366    26.38  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-12                             608898928    1.963  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                           598721480    1.976  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                             616958124    1.946  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                            87181050    13.77  ns/op     3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                             63105722    19.12  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                           57454753    20.70  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                             63631222    19.10  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                           61296942    20.17  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                              48604558    25.08  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                            45605844    26.25  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                              47388217    25.09  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                            45871998    26.10  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-12                                     63452196    18.63  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint32/strconv-12                                   59485942    20.08  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                      63030724    18.98  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                    61887302    19.69  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                     80149611    14.48  ns/op     5  B/op   1  allocs/op
Benchmark_FormatUint16/strconv-12                                   77022440    15.43  ns/op     5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                      79013646    15.23  ns/op     8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                    79695165    15.24  ns/op     8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                     599811685    1.947  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                    84144096    14.18  ns/op     3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                      630717244    1.985  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                     83826201    14.24  ns/op     4  B/op   1  allocs/op
Benchmark_AppendUint/fiber-12                                      130817422    9.139  ns/op     0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                    100000000    10.08  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                             369512403    3.266  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                           213341739    5.627  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                            125937435    9.531  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                          122490632    9.791  ns/op     0  B/op   0  allocs/op

# Token
Benchmark_GenerateSecureToken/16_bytes-12                            4565527    258.6  ns/op    24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                            4306404    280.8  ns/op    48  B/op   1  allocs/op
Benchmark_TokenGenerators/UUIDv4-12                                  4019139    288.2  ns/op    64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                             4360036    276.9  ns/op    48  B/op   1  allocs/op

# HTTP
Benchmark_GetMIME/fiber-12                                          22321272    53.90  ns/op     0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                        17492530    68.46  ns/op     0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/vendorContentType-12      125467813    9.619  ns/op     0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12     404169172    2.974  ns/op     0  B/op   0  allocs/op
Benchmark_StatusMessage/fiber-12                                  1000000000   0.3349  ns/op     0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                                 460456902    2.610  ns/op     0  B/op   0  allocs/op

# IP
Benchmark_IsIPv4/fiber-12                                           80408070    14.58  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                         54025154    21.99  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv6/fiber-12                                           28283431    42.67  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                         20747964    58.18  ns/op     0  B/op   0  allocs/op

# Parse
Benchmark_ParseUint/fiber-12                                       236748204    5.049  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                                 224650773    5.348  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                      88588122    13.75  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/6digit-12                            207540178    5.769  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/8digit-12                            222360716    5.403  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/12digit-12                           156132950    7.699  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/19digit-12                           135049720    8.880  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/1digit-12                            338642533    3.553  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/4digit-12                            212370570    5.672  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                        299791800    3.981  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                  288714712    4.165  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                       76370292    15.48  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber-12                                      200451391    5.988  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                                201629558    5.956  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                     78384832    15.71  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                      220670613    5.462  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                                220175068    5.430  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                    100000000    10.71  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                       277948545    4.279  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                                 264784368    4.566  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                     142184748    8.373  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                     238916116    5.025  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                               222286592    5.409  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                    90395479    13.33  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                     251398184    4.828  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                               248772012    4.838  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                   137015901    8.788  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                      360430984    3.427  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                                327526336    3.659  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                    181904056    6.438  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber-12                                     82912768    14.76  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                               95692504    12.68  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                   50109265    25.05  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                     76648135    15.76  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                               86422119    13.63  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                   47198646    25.41  ns/op     0  B/op   0  allocs/op

# Time
Benchmark_CalculateTimestamp/fiber-12                             1000000000   0.3366  ns/op     0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                             36163686    32.69  ns/op     0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                       4471299    267.8  ns/op    12  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                     3991147    298.7  ns/op     8  B/op   2  allocs/op

# XML
Benchmark_GolangXMLEncoder-12                                         574852     2040  ns/op  4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                        593738     2012  ns/op  4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                        309031     3876  ns/op  2892  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>

## Case-insensitive matching helpers

`EqualFold`, `IndexFold`, `ContainsFold`, `HasPrefixFold`, and
`HasSuffixFold` compare ASCII case-insensitively using the SWAR word
loops: only `A`..`Z`/`a`..`z` fold, every other byte (including bytes
>= 0x80) must match exactly. `HasPrefixFold` and `HasSuffixFold` cover
the header checks Fiber otherwise spells as an allocation-prone
`ToLower` + `HasPrefix`/`HasSuffix` pair (authorization schemes such as
`Bearer`, content-type prefixes such as `multipart/form-data`, suffixes
such as `+json`), and like the other Fold helpers they take the needle
as a plain string because call sites pass constant tokens. Their
`Benchmark_HasPrefixFold`/`Benchmark_HasSuffixFold` numbers are tracked
on the [benchmark charts](https://gofiber.github.io/utils/benchmarks/)
and join the catalog above on its next regeneration.

## HTTP dates

`AppendHTTPDate` and `FormatHTTPDate` write a time in the RFC 9110
preferred HTTP date format (`Mon, 02 Jan 2006 15:04:05 GMT`,
`net/http.TimeFormat`), byte-identical to `time.Format` with that layout
but without walking a layout string: the fixed-width template is copied
once and only the fields are overwritten. `ParseHTTPDate` is the reverse:
canonical preferred-format input takes a strict scalar fast path, and
everything else — the obsolete RFC 850 and asctime forms, unusual casing,
non-GMT zone names, padding — falls back to `time.Parse` with
`net/http.ParseTime` semantics, including its errors. `Date`,
`Last-Modified`, and `If-Modified-Since` handling sit on every
request/response, which makes these the highest-leverage helpers in this
group for Fiber.

## URL escaping

`AppendQueryEscape`/`AppendPathEscape` and
`AppendQueryUnescape`/`AppendPathUnescape` produce byte-identical results
(and, for the unescape pair, identical `url.EscapeError` values) to their
`net/url` counterparts, as append-style, allocation-free single passes.
The escape tables are pinned to `net/url.shouldEscape` by exhaustive
per-byte tests. Unescaping jumps between escape sites with the vectorized
scans (`IndexAny2` when `+` needs rewriting, `bytes.IndexByte` otherwise)
and copies clean spans wholesale, so route parameters and query values
without escapes — the common case — cost one scan and one copy. Decoding
never grows the input, so `dst` may be `s[:0]` on a common backing array
to unescape in place; escaping can grow the input, so there `dst` must
not alias `s`.

## JSON string escaping

`AppendJSONString` appends a value as a double-quoted JSON string,
byte-identical to `encoding/json.Marshal` of the same string — including
its default HTML escaping (`<`, `>`, `&`), the `\ufffd` replacement of
invalid UTF-8 bytes, and the U+2028/U+2029 escapes — without any of
`Marshal`'s reflection or allocation. Clean spans are located with a SWAR
scan (the same first-match-mask technique as `IndexNonQuotable`) and
copied wholesale. This is the building block for hand-rolled JSON hot
paths such as access-log lines and error bodies.

## IP address parsing

`ParseIPv4` and `ParseIPv6` parse addresses into `netip.Addr`, accepting
exactly the strings `netip.ParseAddr` accepts for the respective family
(pinned by fuzzing) while reporting failure with a `bool` instead of a
constructed error. Both take strings or byte slices, so fasthttp-style
callers skip the `string` conversion entirely; remote-address and
trusted-proxy checks parse an IP on every request in Fiber.

## Header key canonicalization

`CanonicalHeaderKey` Title-Cases an HTTP header key exactly like
`net/http.CanonicalHeaderKey`, including its return-unchanged guard for
keys containing non-token bytes. Already-canonical keys — the common case
on receive paths — are validated in a single table pass and returned
as-is with zero allocations for strings and byte slices alike. Keys that
do need rewriting cost one allocation; note that for the ~40 header names
in the stdlib's interning table the stdlib returns a cached string
without allocating, so this helper's edge there is time, not allocations.

These helpers were added on an amd64 machine, so like the `simd` numbers
their benchmarks are recorded separately from the arm64 catalog above and
join it on its next regeneration:

Environment:
goos: linux
goarch: amd64
pkg: github.com/gofiber/utils/v2
cpu: Intel(R) Xeon(R) Processor @ 2.80GHz

```text
// go test -benchmem -run=^$ -bench='Benchmark_(Append|Parse)HTTPDate|Benchmark_Append(Query|Path)|Benchmark_AppendJSONString|Benchmark_ParseIPv[46]|Benchmark_CanonicalHeaderKey' -count=1 .
Benchmark_AppendHTTPDate/fiber-4                                    24732448    50.32  ns/op     0  B/op   0  allocs/op
Benchmark_AppendHTTPDate/default-4                                   7009405    171.2  ns/op     0  B/op   0  allocs/op
Benchmark_ParseHTTPDate/fiber-4                                     35094189    32.26  ns/op     0  B/op   0  allocs/op
Benchmark_ParseHTTPDate/default-4                                    4741988    251.7  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/clean-64B/fiber-4                       18618994    64.21  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/clean-64B/default-4                      4653057    254.4  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/mixed-64B/fiber-4                        4637103    244.2  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/mixed-64B/default-4                      1600225    756.7  ns/op   192  B/op   2  allocs/op
Benchmark_AppendQueryUnescape/plain-64B/fiber-4                     56108641    20.74  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/plain-64B/default-4                    7959500    136.8  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/escaped-64B/fiber-4                    4737538    250.3  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/escaped-64B/default-4                  3816861    312.7  ns/op    48  B/op   1  allocs/op
Benchmark_AppendPathUnescape/plain-25B/fiber-4                      74870614    16.38  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/plain-25B/default-4                    11201378    90.13  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/escaped-31B/fiber-4                    19748607    62.69  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/escaped-31B/default-4                   6227439    173.6  ns/op    24  B/op   1  allocs/op
Benchmark_AppendJSONString/clean-16B/fiber-4                        47903053    22.13  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/clean-16B/default-4                       7762933    151.9  ns/op    40  B/op   2  allocs/op
Benchmark_AppendJSONString/clean-64B/fiber-4                        18578377    64.12  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/clean-64B/default-4                       3901436    271.8  ns/op    96  B/op   2  allocs/op
Benchmark_AppendJSONString/escaped-64B/fiber-4                       7056486    165.8  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/escaped-64B/default-4                     3842581    319.8  ns/op    96  B/op   2  allocs/op
Benchmark_AppendJSONString/unicode-64B/fiber-4                       4654760    248.2  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/unicode-64B/default-4                     3786951    309.2  ns/op   112  B/op   2  allocs/op
Benchmark_ParseIPv4/fiber-4                                         52309202    20.02  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv4/default-4                                       42021915    28.74  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/compressed/fiber-4                              22675394    51.88  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/compressed/default-4                            18522648    64.65  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/full/fiber-4                                    20229412    61.40  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/full/default-4                                  14830760    83.42  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical/fiber-4                      54643375    22.69  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical/default-4                    36882357    28.46  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/common-lower/fiber-4                   22508378    52.04  ns/op    16  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/common-lower/default-4                 21160863    56.48  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/custom-lower/fiber-4                   17017753    73.05  ns/op    24  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/custom-lower/default-4                 12678414    102.6  ns/op    24  B/op   1  allocs/op
```

## SWAR primitives

The [`swar`](swar/) package exports the SWAR (SIMD within a register)
building blocks the helpers above are composed from: `Load8`/`Store8`,
`Broadcast`, `ZeroLanes`, `MatchByteMask`/`MatchRangeMask`,
`ToLowerWord`/`ToUpperWord`, `FirstLane`/`LastLane`, and the `WordLen`,
`Ones`, `HighBits`, and `LowSeven` constants. They are exported so
downstream packages (Fiber itself, middleware) can fuse their own byte
scans — for example, finding a delimiter while classifying the bytes
before it — without re-deriving the bit tricks. The contracts (unchecked
bounds preconditions, `ZeroLanes`' approximate mask, the little-endian
lane order) and runnable examples live in the package documentation.

The package benchmarks its primitives against their stdlib counterparts
the same way the helpers above do, measuring the canonical loops composed
from them: `Load8`/`Store8` vs `encoding/binary`'s little-endian
`Uint64`/`PutUint64`, a `ZeroLanes` first-match scan vs
`strings.IndexByte`, a `MatchByteMask`+`LastLane` reverse scan vs
`bytes.LastIndexByte`, a `MatchRangeMask` digit scan vs
`strings.IndexAny`, and an in-place `ToLowerWord` loop vs
`bytes.ToLower`. The composed loops are pinned to the stdlib results by a
dedicated test, and the numbers are tracked per commit on the
[benchmark charts](https://gofiber.github.io/utils/benchmarks/). Note
that `strings.IndexByte` is hand-written SIMD assembly and overtakes the
portable SWAR scan on large inputs; the SWAR primitives earn their keep
on short HTTP-sized inputs and on fused scans the stdlib has no single
function for. The package's remaining benchmark, `Benchmark_Load8_Fusion`,
is an advisory codegen guard rather than an API performance promise, so
it is deliberately not part of the catalog above.

## SIMD-accelerated search (package simd)

The [`simd`](simd/) package is the vector-width counterpart to `swar`: byte
searching and validation whose scan kernels dispatch to AVX2 assembly (32
bytes per iteration) on amd64 CPUs for inputs of `simd.MinLen` (32) bytes
or more, and fall back to portable loops built on the `swar` primitives
everywhere else, so it builds and behaves identically on every platform.
`simd.Accelerated()` reports which mode is active. CPU capability detection
uses [`golang.org/x/sys/cpu`](https://pkg.go.dev/golang.org/x/sys/cpu),
which also verifies OS support for saving the vector register state.

The AVX2-backed kernels are the multi-needle scans (`Memchr2`, `Memchr3`),
the paired-byte scan (`MemchrPair`), the class scans
(`MemchrDigit`/`MemchrDigitAt`, `MemchrWord`/`MemchrNotWord`), ASCII
validation (`IsASCII`), and the substring search (`Memmem`), which
prefilters on the needle's rarest bytes (`SelectRareBytes`, `ByteRank`)
before verifying candidates — 2-6 byte needles containing two distinct
values are scanned for their two rarest bytes at the exact relative
distance, longer or single-valued needles for the single rarest byte.
Two smaller tiers are documented explicitly:
`FirstNonASCII` and `CountNonASCII` are SWAR-only (8 bytes per iteration,
no assembly kernel), and `MemchrInTable`/`MemchrNotInTable` are plain
scalar loops, since an arbitrary 256-entry membership test has no cheap
vector form. There is deliberately no single-needle `Memchr`:
`bytes.IndexByte` is already vector-accelerated by the Go runtime. `Memmem`
delegates to `bytes.Index` below 128 bytes — a conservative routing point
above the paired scan's measured break-even (`Benchmark_Memmem_Prefilter`
forces the prefilter helpers at every size to keep it measurable) — and
bails to `bytes.Index` after a bounded number of failed candidate
verifications, so even adversarial inputs — a rare byte everywhere plus a
long almost-matching needle — stay within a constant of the stdlib's
O(n+m) (see `Benchmark_Memmem_Adversarial`). The
package operates on `[]byte` only; the generic top-level helpers `IsASCII`,
`IndexAny2`, and `IndexAny3` dispatch into it automatically for inputs of
`simd.MinLen`+ bytes on amd64, which is where the AVX2 kernels overtake the
SWAR word loops (-35% to -88% ns/op at 32-512B in `benchstat -count=10`
runs; shorter inputs keep the existing SWAR paths).

The kernels are adapted from the [coregex](https://github.com/coregx/coregex)
project's `simd` package (MIT License, Copyright (c) 2025 Andrey Kolkov and
contributors; see [`simd/LICENSE`](simd/LICENSE)), with the legacy-SSE
register moves in the kernels replaced by VEX-encoded ones to avoid AVX-SSE
transition stalls, the scalar tail loops replaced by overlapping vector
rescans at the buffer end (a 63-byte scan previously cost ~3.7x a 64-byte
one), and
the fallbacks reworked around the stdlib and the `swar` package as
described above.

Because the AVX2 kernels only engage on amd64, their benchmark numbers are
recorded separately from the arm64 catalog above:

Environment:
goos: linux
goarch: amd64
pkg: github.com/gofiber/utils/v2/simd
cpu: Intel(R) Xeon(R) Processor @ 2.80GHz (AVX2)

```text
// go test ./simd/ -benchmem -run=^$ -bench=Benchmark_ -count=1
Benchmark_Memchr2/8B/simd-4                                        160932234    7.510  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/8B/default-4                                      31407770    39.39  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/32B/simd-4                                       180379407    7.463  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/32B/default-4                                     28221582    42.47  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/64B/simd-4                                       170949894    6.883  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/64B/default-4                                     21863884    57.19  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/512B/simd-4                                       59448526    21.04  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/512B/default-4                                     3130293    391.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/4096B/simd-4                                      10117922    117.0  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/4096B/default-4                                     412353     3055  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/8B/simd-4                                        141613051    7.994  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/8B/default-4                                      29630468    40.43  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/32B/simd-4                                       171311121    6.895  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/32B/default-4                                     30218668    40.41  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/64B/simd-4                                       150505892    7.856  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/64B/default-4                                     19079786    64.43  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/512B/simd-4                                       52560612    23.74  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/512B/default-4                                     3086830    388.6  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/4096B/simd-4                                       7735116    150.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/4096B/default-4                                     400952     2975  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/8B/simd-4                                      89391540    12.82  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/8B/scalar-4                                   216545959    5.584  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/32B/simd-4                                     54730342    21.36  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/32B/scalar-4                                   56326855    21.45  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/64B/simd-4                                    100000000    10.74  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/64B/scalar-4                                   27203602    44.19  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/512B/simd-4                                    55834921    21.79  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/512B/scalar-4                                   3446038    342.1  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/4096B/simd-4                                    9817900    122.1  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/4096B/scalar-4                                   445594     2675  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/8B/simd-4                                    169659172    7.221  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/8B/scalar-4                                  234911584    5.033  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/32B/simd-4                                   174639595    7.164  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/32B/scalar-4                                  89061945    12.65  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/64B/simd-4                                   151388893    7.846  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/64B/scalar-4                                  52860561    23.28  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/512B/simd-4                                   47377659    25.40  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/512B/scalar-4                                  7123378    168.5  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/4096B/simd-4                                   6884252    176.9  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/4096B/scalar-4                                  884161     1317  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/8B/simd-4                                  100000000    10.07  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/8B/scalar-4                                134836845    8.672  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/32B/simd-4                                 149667877    8.193  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/32B/scalar-4                                44375330    26.63  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/64B/simd-4                                 100000000    10.22  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/64B/scalar-4                                23879234    50.68  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/512B/simd-4                                 22524218    46.90  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/512B/scalar-4                                2931174    408.6  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/4096B/simd-4                                 3752906    335.2  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/4096B/scalar-4                                345715     3136  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/8B/simd-4                                         100000000    11.61  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/8B/default-4                                      134337361    9.246  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/32B/simd-4                                         73584829    15.61  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/32B/default-4                                      88738117    13.46  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/64B/simd-4                                         61874611    19.30  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/64B/default-4                                      74323999    16.76  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/96B/simd-4                                         22135388    63.57  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/96B/default-4                                      23687415    50.17  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/128B/simd-4                                        47676351    24.76  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/128B/default-4                                     19612291    60.68  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/192B/simd-4                                        47072698    26.98  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/192B/default-4                                     11572932    100.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/512B/simd-4                                        34756200    36.05  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/512B/default-4                                      4604108    261.4  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/4096B/simd-4                                        8868234    135.8  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/4096B/default-4                                      613701     1962  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/8B/paired-4                              86125052    14.52  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/8B/paired-stdlib-4                      136569471    9.120  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/paired-4                             46744827    25.76  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/paired-stdlib-4                      95989609    12.81  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/single-4                             69678790    17.61  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/single-stdlib-4                      68305318    18.38  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/paired-4                             64893210    17.20  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/paired-stdlib-4                      70806460    16.34  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/single-4                             34043565    36.62  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/single-stdlib-4                      24664317    47.43  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/paired-4                             67596938    19.75  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/paired-stdlib-4                      21772609    50.35  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/single-4                             25962633    47.83  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/single-stdlib-4                      31263175    38.29  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/paired-4                            65862615    18.80  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/paired-stdlib-4                     18729648    64.33  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/single-4                            18243868    65.45  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/single-stdlib-4                     23761232    51.21  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/paired-4                            60481297    20.08  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/paired-stdlib-4                     13275300    93.19  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/single-4                            11936652    102.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/single-stdlib-4                     14467984    83.21  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/paired-4                            42353164    28.11  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/paired-stdlib-4                      4483054    258.0  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/single-4                             4454024    282.8  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/single-stdlib-4                      4753104    251.3  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/paired-4                            9504126    124.9  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/paired-stdlib-4                      613930     1985  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/single-4                             563440     1980  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/single-stdlib-4                      591222     1924  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Adversarial/simd-4                                      715  1698023  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Adversarial/default-4                                   722  1722337  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/simd-4                                        211383338    5.542  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/swar-4                                        316762524    4.007  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/simd-4                                       240245251    5.077  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/swar-4                                       164716737    7.257  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/simd-4                                       191159047    6.444  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/swar-4                                        88207845    12.81  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/simd-4                                       75486206    14.96  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/swar-4                                       14351298    90.64  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/4096B/simd-4                                      11831648    101.2  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/4096B/swar-4                                       1512061    685.9  ns/op     0  B/op   0  allocs/op
```

## ☕ Supporters

Fiber is an open-source project that runs on donations to pay the bills, e.g., our domain name, hosting, and serverless infrastructure. If you want to support Fiber, please become a [GitHub Sponsor](https://github.com/sponsors/gofiber).

<!-- sponsors -->

### 📅 Monthly Sponsors

<table>
<tr><td valign="top"><strong>🔥 Fiber Guardian</strong></td><td><a href="https://www.coderabbit.ai/?utm_source=cr_org&amp;utm_medium=github" title="@coderabbitai"><img src="https://github.com/coderabbitai.png" width="50" alt="@coderabbitai" /></a></td></tr>
<tr><td valign="top"><strong>☕ Fiber Supporter</strong></td><td><a href="https://ndole.studio" title="@NdoleStudio"><img src="https://github.com/NdoleStudio.png" width="34" alt="@NdoleStudio" /></a>&nbsp;<a href="https://cyberapper.ai" title="@petercool"><img src="https://github.com/petercool.png" width="34" alt="@petercool" /></a></td></tr>
<tr><td valign="top"><strong>🪴 Fiber Friend</strong></td><td><a href="https://github.com/simonheisstpeter" title="@simonheisstpeter"><img src="https://github.com/simonheisstpeter.png" width="32" alt="@simonheisstpeter" /></a></td></tr>
</table>

### 🎁 One-time Sponsors

<table>
<tr><td valign="top"><strong>🚀 Fiber Hero</strong></td><td><a href="https://www.thanks.dev" title="@thnxdev"><img src="https://github.com/thnxdev.png" width="40" alt="@thnxdev" /></a></td></tr>
</table>
<!-- sponsors -->
