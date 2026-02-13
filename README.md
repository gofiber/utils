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
goarch: arm64
pkg: github.com/gofiber/utils/v2
cpu: Apple M2 Pro

```text
// go test -benchmem -run=^$ -bench=Benchmark_ -count=1

# Case Conversion
Benchmark_ToLowerBytes/empty/fiber-12                              593913379   1.975  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/unsafe-12                       657748088   1.845  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                            367040010   3.281  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                           365235189   3.393  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/unsafe-12                    239092582   5.144  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                          63143354   18.46  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                      94304095   12.86  ns/op         3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/unsafe-12              239366296   5.057  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                    86047208   14.00  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12           45169880   26.27  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/unsafe-12    95687102   12.70  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12         25045154   47.56  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                         43856107   28.39  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/unsafe-12                  62907780   19.56  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                       17712895   69.18  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                         34286121   34.08  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/unsafe-12                  61959598   19.37  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                       16091281   75.42  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                         33481714   35.06  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/unsafe-12                  62478574   19.36  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                       16091731   75.18  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/empty/fiber-12                              620292880   1.946  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/unsafe-12                       642997940   1.840  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                            366694138   3.324  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                            95312475   12.75  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/unsafe-12                    236603360   5.055  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                          86667109   14.19  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                     360987979   3.425  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/unsafe-12              236142552   5.096  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                    61760695   19.24  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12           44394090   27.55  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/unsafe-12    94647172   12.75  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12         25090692   47.83  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                         35558892   34.73  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/unsafe-12                  62620129   19.31  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                       15487833   76.16  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                         43901836   27.54  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/unsafe-12                  62027790   20.57  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                       14579604   82.82  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                         33719427   35.99  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/unsafe-12                  62008428   21.01  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                       15469116   77.21  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/empty/fiber-12                                   613156160   1.955  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/unsafe-12                            649026514   1.838  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                                 650956171   1.831  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                                 92095164   12.77  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/unsafe-12                         235750344   5.036  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                               60473330   20.50  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                          398272657   3.048  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/unsafe-12                   220124466   5.078  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                        233645295   5.051  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12                36106831   33.14  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/unsafe-12         93201168   12.83  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12              11069166   107.7  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                              29007373   41.16  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/unsafe-12                       61283246   19.47  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                             6855549   174.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                              43965034   28.00  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/unsafe-12                       53715409   20.68  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                            18805570   65.60  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                              27005611   43.54  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/unsafe-12                       60924529   20.57  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                             5662490   212.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/empty/fiber-12                                   607541103   1.972  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/unsafe-12                            631531857   1.969  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                                 638706405   1.872  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                                390062401   3.080  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/unsafe-12                         231204609   5.190  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                              278482443   4.408  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                           91051068   13.16  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/unsafe-12                   233437356   5.168  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                         59555317   22.03  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12                32247734   35.17  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/unsafe-12         92616117   12.98  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12              16339488   74.21  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                              37544337   33.77  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/unsafe-12                       61856332   19.83  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                            23513985   50.17  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                              28177750   44.45  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/unsafe-12                       61319127   20.23  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                             6700268   180.1  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                              26804186   45.68  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/unsafe-12                       61659141   19.97  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                             5614172   213.1  ns/op        64  B/op   1  allocs/op

# Add Trailing Slash
Benchmark_AddTrailingSlashBytes/empty-12                          1000000000  0.6118  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                     1000000000  0.9030  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                  100000000   12.17  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12               1000000000  0.9096  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                    79386301   13.31  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12                1000000000  0.9079  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/empty-12                         1000000000  0.3001  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                    1000000000  0.5217  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12                 100000000   13.37  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12              1000000000  0.4548  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                   84763720   16.71  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12               1000000000  0.4536  ns/op         0  B/op   0  allocs/op

# EqualFold
Benchmark_EqualFoldBytes/fiber-12                                   34227160   35.89  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                                 17905046   68.31  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-12                                        34215410   36.76  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                      17103853   69.62  ns/op         0  B/op   0  allocs/op

# Trim
Benchmark_TrimRight/fiber-12                                       556198228   2.195  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                     394761454   3.080  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/fiber-12                                  552028600   2.188  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                                364481850   3.382  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/fiber-12                                        493763746   2.439  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                      394974064   3.147  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/fiber-12                                   533978953   2.129  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                                 390162962   3.090  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/fiber-12                                            303470884   4.120  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default-12                                          238874899   4.889  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                                230900814   5.395  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/fiber-12                                       302062202   3.963  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                     245799558   5.052  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                           217547649   5.520  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/empty-12                                1000000000  0.3031  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                               554941630   2.211  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                                378650127   3.139  ns/op    955.63  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/spaces-12                              313411838   3.929  ns/op    763.64  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                            325986266   3.682  ns/op   2444.19  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                          260590963   4.751  ns/op   1894.19  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                    278654619   4.279  ns/op   5608.46  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                  217345236   5.665  ns/op   4236.48  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                     374732025   3.245  ns/op  11711.17  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                   245032184   5.138  ns/op   7395.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                       322936750   3.772  ns/op   5567.24  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                     259551775   4.791  ns/op   4383.20  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                          327816716   3.700  ns/op   5404.76  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/content-type-12                        261473678   4.851  ns/op   4122.61  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                       430371238   3.136  ns/op   8929.40  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                     259204014   4.351  ns/op   6434.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                          297618063   3.686  ns/op   5425.87  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/query-params-12                        260095944   4.591  ns/op   4356.73  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                            244683963   4.594  ns/op  22203.55  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                          196174795   6.097  ns/op  16728.77  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                              1000000000  0.6452  ns/op   7750.12  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                             431454468   2.758  ns/op   1813.07  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                      279402585   4.522  ns/op   3759.80  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                    216711007   5.540  ns/op   3068.59  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/empty-12                           1000000000  0.3096  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                          486190254   2.466  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                           350271897   3.119  ns/op    961.83  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                         294683720   4.188  ns/op    716.36  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                       328983026   3.690  ns/op   2438.91  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                     242798146   5.181  ns/op   1737.10  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12               276373491   4.299  ns/op   5582.44  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12             203176566   6.066  ns/op   3956.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12                371498816   3.203  ns/op  11863.40  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12              233269664   5.123  ns/op   7417.55  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                  331386438   3.740  ns/op   5614.74  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12                248029371   5.101  ns/op   4116.44  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                     324414853   3.624  ns/op   5518.54  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                   250866906   4.764  ns/op   4198.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                  448704225   2.690  ns/op  10409.74  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12                267864716   5.395  ns/op   5189.83  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                     298958100   3.821  ns/op   5234.49  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                   247046188   4.859  ns/op   4116.00  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                       259565295   4.594  ns/op  22201.52  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                     188316823   6.343  ns/op  16079.55  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                         1000000000  0.6094  ns/op   8204.80  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                        395805396   3.032  ns/op   1649.01  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12                 281619982   4.206  ns/op   4041.49  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12               209473969   5.716  ns/op   2974.30  MB/s   0  B/op  0  allocs/op

# Convert
Benchmark_ConvertToBytes/fiber-12                                  305651658   3.974  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int-12                                          494480786   2.455  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                         601754490   2.009  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                        480638928   2.507  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                        483084573   2.505  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                        485630560   2.572  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                         542968509   2.126  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                        603426201   2.082  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                       568713753   2.111  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                       567709288   2.106  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                       568810574   2.180  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/string-12                                       606042240   2.191  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                       94285881   12.78  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                         601565197   1.995  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                       24979378   48.10  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                       17121148   70.77  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                     12894415   93.74  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                  12970261   95.06  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                      52706324   26.62  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                         51370597   24.19  ns/op         8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                        31207941   38.93  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                        8454266   142.9  ns/op       112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                                 9656547   133.4  ns/op        72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                             390543688   3.060  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                              11184109   106.9  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/int-12                             1000000000  0.3039  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int8-12                            1000000000  0.2423  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int16-12                           1000000000  0.2901  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int32-12                           1000000000  0.2935  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int64-12                           1000000000  0.2856  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint-12                            1000000000  0.3382  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint8-12                           1000000000  0.2279  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint16-12                          1000000000  0.2460  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint32-12                          1000000000  0.2469  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint64-12                          1000000000  0.2615  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/string-12                          1000000000  0.2305  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/[]uint8-12                          176655307   6.047  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/bool-12                            1000000000  0.2306  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/float32-12                          183723876   8.276  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/float64-12                          121188367   9.445  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time-12                         68860974   17.98  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time#01-12                      65536016   18.02  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]string-12                         159027694   8.212  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]int-12                            203055362   5.631  ns/op         8  B/op   1  allocs/op
Benchmark_ToString_concurrency/[2]int-12                           131550500   9.186  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[][]int-12                           21128588   55.64  ns/op       112  B/op   6  allocs/op
Benchmark_ToString_concurrency/[]interface_{}-12                    30081375   38.26  ns/op        72  B/op   3  allocs/op
Benchmark_ToString_concurrency/utils.MyStringer-12                1000000000  0.3727  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/utils.CustomType-12                  49417796   29.59  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeBytes/unsafe-12                                   1000000000  0.5450  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                    80875700   14.54  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                  1000000000  0.3579  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                   99855757   12.15  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/0-12                                            394155170   3.023  ns/op         0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                             64058362   18.38  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                           58149043   21.64  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                          58947781   20.26  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                          52060172   23.14  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                       60832776   19.66  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                       51995977   22.68  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                    61154545   19.64  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                    52380780   23.06  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                                 62655552   19.10  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                                 53917438   23.08  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                              60970962   18.98  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                              55339809   21.64  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                           63578542   18.73  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                           55347254   21.74  ns/op        16  B/op   1  allocs/op

# Format and Append
Benchmark_FormatUint/small/fiber-12                               1000000000  0.3025  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                              584345623   2.084  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                               125524198   10.04  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/strconv-12                              57687223   21.13  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                                 64259034   18.67  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/large/strconv-12                               44394020   26.98  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-12                             609931978   1.994  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                           579927732   2.046  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                             611030502   1.969  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                            82791925   14.13  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                             61593342   19.56  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                           58425434   21.66  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                             60677822   19.54  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                           59446896   20.23  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                              46969258   25.61  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                            44999648   26.57  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                              45581593   26.18  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                            45468108   26.83  ns/op        24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-12                                    126955321   9.411  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint32/strconv-12                                   58030117   21.34  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                      60967604   19.60  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                    59480410   20.12  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                    164073385   7.314  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint16/strconv-12                                   74532424   16.78  ns/op         5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                      74833182   15.94  ns/op         8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                    75867136   16.05  ns/op         8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                    1000000000  0.2999  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                    82306642   14.64  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                     1000000000  0.3015  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                     80331591   14.71  ns/op         4  B/op   1  allocs/op
Benchmark_AppendUint/fiber-12                                      135872859   8.812  ns/op         0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                    100000000   10.44  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                             345227787   3.321  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                           210340929   5.744  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                            121853198   9.739  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                          100000000   10.04  ns/op         0  B/op   0  allocs/op

# Token
Benchmark_GenerateSecureToken/16_bytes-12                            4445361   284.0  ns/op        24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                            4071276   291.0  ns/op        48  B/op   1  allocs/op
Benchmark_TokenGenerators/UUIDv4-12                                  4021219   327.0  ns/op        64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                             4073481   287.3  ns/op        48  B/op   1  allocs/op

# HTTP
Benchmark_GetMIME/fiber-12                                          29202141   41.13  ns/op         0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                        16953927   73.52  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/vendorContentType-12      154110751   7.784  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12     396583869   3.023  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/fiber-12                                  1000000000  0.3396  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                                 447403729   2.694  ns/op         0  B/op   0  allocs/op

# IP
Benchmark_IsIPv4/fiber-12                                           82770516   14.99  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                         52222729   22.87  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/fiber-12                                           27058488   44.28  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                         19518978   63.64  ns/op         0  B/op   0  allocs/op

# Parse
Benchmark_ParseUint/fiber-12                                       171968007   6.984  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                                 162498577   7.444  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                      87571795   13.70  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                        186597552   6.331  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                  180017686   6.686  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                       76678334   15.73  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber-12                                      155367661   7.717  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                                144681208   8.325  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                     76537323   15.73  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                      228710796   5.222  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                                223933951   5.322  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                    100000000   11.05  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                       257930811   4.626  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                                 240602468   4.588  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                     138314644   8.609  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                     170833845   7.138  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                               162105060   7.443  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                    85408747   13.70  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                     262682311   4.547  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                               257359222   4.658  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                   134532823   8.902  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                      350422557   3.416  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                                327501829   3.700  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                    182633575   6.816  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber-12                                    100000000   10.73  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                              100000000   10.86  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                   50053615   24.18  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                    100000000   11.62  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                              100000000   11.97  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                   46590180   25.98  ns/op         0  B/op   0  allocs/op

# Time
Benchmark_CalculateTimestamp/fiber-12                             1000000000  0.3029  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                             36542426   32.77  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                       4240807   276.8  ns/op        13  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                     3974042   305.6  ns/op         8  B/op   2  allocs/op

# XML
Benchmark_GolangXMLEncoder-12                                         591470    2025  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                        563848    2124  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                        329698    3513  ns/op      2857  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>
