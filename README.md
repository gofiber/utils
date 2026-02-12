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
Benchmark_ToLowerBytes/empty/fiber-12                            562607137   1.932  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/mut-12                        669249470   1.812  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                          370465512   3.247  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                         370317031   3.263  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/mut-12                     239448190   5.001  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                        63751792   18.22  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                    95151250   12.53  ns/op         3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/mut-12               238922935   5.000  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                  85559721   13.87  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12         45776702   25.80  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/mut-12     97528600   12.43  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12       25347337   46.83  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                       44031445   27.32  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/mut-12                   63445485   19.06  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                     17751248   67.86  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                       34313900   33.39  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/mut-12                   63220422   18.98  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                     16217996   75.03  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                       33949492   34.28  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/mut-12                   63029622   19.01  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                     16209014   73.70  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/empty/fiber-12                            636204963   1.895  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/mut-12                        665187241   1.802  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                          371806530   3.239  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                          86623316   12.62  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/mut-12                     239637328   5.070  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                        85115437   13.79  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                   254208652   4.719  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/mut-12               240077564   5.012  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                  63225559   18.84  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12         45539355   26.50  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/mut-12     95642610   12.45  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12       25478991   47.31  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                       34421308   33.81  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/mut-12                   62725418   19.12  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                     16254315   73.90  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                       27777884   42.39  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/mut-12                   63828795   19.09  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                     14958552   81.01  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                       36056834   33.74  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/mut-12                   61524001   19.17  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                     16253546   74.89  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/empty/fiber-12                                 596518689   1.959  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/mut-12                             624564645   1.912  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                               657259036   1.832  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                               90831099   12.87  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/mut-12                          229938926   5.107  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                             59003214   20.52  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                        395015073   3.055  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/mut-12                    234575518   5.135  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                      235055221   5.114  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12              37753705   32.46  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/mut-12          92610163   13.11  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12            10964001   108.7  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                            28924374   41.59  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/mut-12                        60848068   19.91  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                           6796114   177.7  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                            42915831   27.73  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/mut-12                        60532687   20.63  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                          18319328   64.96  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                            29809990   40.91  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/mut-12                        60801308   20.06  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                           5727622   211.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/empty/fiber-12                                 585763149   1.997  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/mut-12                             579111684   1.946  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                               649398088   1.886  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                              397400338   3.043  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/mut-12                          234215938   5.101  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                            282817620   4.283  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                         91238850   12.87  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/mut-12                    234980426   5.144  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                       58956351   20.68  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12              39367387   31.87  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/mut-12          91721205   13.11  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12            16687071   72.30  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                            43543207   27.72  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/mut-12                        60477774   20.04  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                          23948589   49.60  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                            29243416   40.74  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/mut-12                        60230253   19.89  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                           6784809   176.9  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                            29184238   41.25  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/mut-12                        60772956   19.83  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                           5642757   211.9  ns/op        64  B/op   1  allocs/op

# Add Trailing Slash
Benchmark_AddTrailingSlashBytes/empty-12                        1000000000  0.5877  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                   1000000000  0.8842  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                100000000   11.24  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12             1000000000  0.8823  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                  92474558   12.80  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12              1000000000  0.8845  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/empty-12                       1000000000  0.2989  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                  1000000000  0.4513  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12               100000000   11.60  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12            1000000000  0.4524  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                 83570551   14.36  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12             1000000000  0.4508  ns/op         0  B/op   0  allocs/op

# EqualFold
Benchmark_EqualFoldBytes/fiber-12                                 34969195   34.31  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                               18024433   66.58  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-12                                      35762536   33.63  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                    18098534   67.76  ns/op         0  B/op   0  allocs/op

# Trim
Benchmark_TrimRight/fiber-12                                     582591085   2.059  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                   410040874   2.919  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/fiber-12                                561094851   2.058  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                              370952655   3.241  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/fiber-12                                      581832638   2.069  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                    411236700   2.945  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/fiber-12                                 577633549   2.085  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                               404901102   2.973  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/fiber-12                                          312043243   3.825  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default-12                                        256613820   4.687  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                              239937675   5.009  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/fiber-12                                     312875478   3.829  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                   254974299   4.708  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                         221112632   5.353  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/empty-12                              1000000000  0.2952  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                             585226446   2.061  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                              397083148   2.992  ns/op   1002.76  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/spaces-12                            332733718   3.607  ns/op    831.62  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                          337526038   3.547  ns/op   2537.01  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                        272195794   4.415  ns/op   2038.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                  290465881   4.186  ns/op   5732.77  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                226599972   5.306  ns/op   4523.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                   388413254   3.087  ns/op  12309.40  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                 254461356   4.711  ns/op   8066.88  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                     340252790   3.538  ns/op   5936.01  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                   271563048   4.443  ns/op   4726.42  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                        338757052   3.536  ns/op   5656.43  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/content-type-12                      267459386   4.462  ns/op   4481.99  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                     451773350   2.653  ns/op  10553.53  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                   292073134   4.164  ns/op   6724.23  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                        337979772   3.539  ns/op   5651.46  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/query-params-12                      271706955   4.423  ns/op   4521.82  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                          272344134   4.415  ns/op  23103.26  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                        203575363   5.972  ns/op  17079.23  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                            1000000000  0.5958  ns/op   8391.80  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                           451524316   2.665  ns/op   1876.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                    289224228   4.168  ns/op   4078.33  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                  226609796   5.381  ns/op   3159.07  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/empty-12                         1000000000  0.2937  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                        511396959   2.358  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                         397237488   3.000  ns/op   1000.09  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                       309039368   3.935  ns/op    762.31  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                     338146242   3.550  ns/op   2534.88  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                   255104299   4.713  ns/op   1909.54  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12             291427654   4.132  ns/op   5808.49  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12           214732161   5.590  ns/op   4293.45  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12              387477585   3.095  ns/op  12276.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12            239973702   5.017  ns/op   7574.43  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                331869382   3.583  ns/op   5861.57  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12              254630810   4.711  ns/op   4457.86  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                   337523863   3.553  ns/op   5629.75  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                 253799433   4.754  ns/op   4207.32  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                453855430   2.654  ns/op  10551.67  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12              271331482   4.425  ns/op   6327.23  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                   338630191   3.610  ns/op   5540.21  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                 253886691   4.716  ns/op   4240.75  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                     270497842   4.415  ns/op  23101.05  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                   193911619   6.195  ns/op  16464.11  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                       1000000000  0.5940  ns/op   8417.71  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                      402750675   2.982  ns/op   1676.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12               292170618   4.179  ns/op   4067.65  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12             213880502   5.680  ns/op   2993.00  MB/s   0  B/op  0  allocs/op

# Convert
Benchmark_ConvertToBytes/fiber-12                                310905159   3.827  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int-12                                        502824020   2.401  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                       606820150   1.989  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                      487144623   2.471  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                      490471567   2.434  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                      498678242   2.403  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                       574467778   2.100  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                      612626748   1.974  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                     545585880   2.096  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                     569111924   2.100  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                     556174278   2.114  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/string-12                                     614593784   1.968  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                     97850684   12.39  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                       602868020   1.984  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                     24873472   47.78  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                     17029458   70.32  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                   12715797   92.94  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                12975626   92.55  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                    54478285   22.39  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                       52474882   22.84  ns/op         8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                      31125241   37.26  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                      8956072   134.7  ns/op       112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                              10187494   118.4  ns/op        72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                           398638814   3.007  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                            11591224   102.7  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/int-12                           1000000000  0.2751  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int8-12                          1000000000  0.2291  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int16-12                         1000000000  0.2869  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int32-12                         1000000000  0.2801  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int64-12                         1000000000  0.2759  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint-12                          1000000000  0.2368  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint8-12                         1000000000  0.2298  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint16-12                        1000000000  0.2405  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint32-12                        1000000000  0.2392  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint64-12                        1000000000  0.2387  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/string-12                        1000000000  0.2254  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/[]uint8-12                        188456032   5.449  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/bool-12                          1000000000  0.2257  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/float32-12                        189025742   6.282  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/float64-12                        136045532   8.798  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time-12                       65109538   18.11  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time#01-12                    69462532   18.51  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]string-12                       153318090   7.714  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]int-12                          230506161   5.021  ns/op         8  B/op   1  allocs/op
Benchmark_ToString_concurrency/[2]int-12                         132628204   8.854  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[][]int-12                         22405947   54.58  ns/op       112  B/op   6  allocs/op
Benchmark_ToString_concurrency/[]interface_{}-12                  34394383   37.82  ns/op        72  B/op   3  allocs/op
Benchmark_ToString_concurrency/utils.MyStringer-12              1000000000  0.3529  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/utils.CustomType-12                47802733   25.50  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeBytes/unsafe-12                                 1000000000  0.5334  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                  87169176   13.79  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                1000000000  0.3525  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                100000000   11.94  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/0-12                                          399833734   3.011  ns/op         0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                           67199912   18.03  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                         59089530   20.47  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                        59969888   19.76  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                        52984129   22.60  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                     60678589   19.52  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                     53521848   22.47  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                  62191738   19.11  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                  55012855   21.98  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                               62992538   18.86  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                               55179073   21.68  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                            64188285   18.52  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                            56490193   21.63  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                         64351789   18.50  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                         55831156   21.54  ns/op        16  B/op   1  allocs/op

# Format and Append
Benchmark_FormatUint/small/fiber-12                             1000000000  0.3018  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                            589275553   2.041  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                             127757719   9.410  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/strconv-12                            58343175   20.70  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                               64500378   18.51  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/large/strconv-12                             44637529   26.84  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-12                           613004246   1.988  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                         567523132   2.067  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                           612006158   2.070  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                          84422815   14.00  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                           61412750   19.66  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                         57247695   20.84  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                           60452638   19.57  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                         59671557   20.39  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                            46695860   25.44  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                          44990439   27.77  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                            45807645   26.03  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                          44509147   26.67  ns/op        24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-12                                  126939831   9.451  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint32/strconv-12                                 57669547   21.30  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                    61047212   19.63  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                  59661667   20.17  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                  150809730   7.958  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint16/strconv-12                                 74377861   16.20  ns/op         5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                    74190662   15.93  ns/op         8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                  74672527   16.26  ns/op         8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                  1000000000  0.2983  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                  81523582   14.70  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                   1000000000  0.2988  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                   81463402   14.69  ns/op         4  B/op   1  allocs/op
Benchmark_AppendUint/fiber-12                                    135975609   8.784  ns/op         0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                  100000000   10.24  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                           362662759   3.334  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                         210232480   5.720  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                          124164187   9.592  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                        120759157   9.969  ns/op         0  B/op   0  allocs/op

# Token
Benchmark_GenerateSecureToken/16_bytes-12                          4622641   260.9  ns/op        24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                          4403814   272.6  ns/op        48  B/op   1  allocs/op
Benchmark_TokenGenerators/UUIDv4-12                                4195326   287.9  ns/op        64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                           4402108   273.8  ns/op        48  B/op   1  allocs/op

# HTTP
Benchmark_GetMIME/fiber-12                                        29484692   41.69  ns/op         0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                      17397680   69.14  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/vendorContentType-12    131862534   9.228  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12   397674488   3.033  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/fiber-12                                1000000000  0.3375  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                               450421849   2.649  ns/op         0  B/op   0  allocs/op

# IP
Benchmark_IsIPv4/fiber-12                                         79822177   14.85  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                       53036230   23.02  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/fiber-12                                         27212709   51.37  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                       19845904   60.93  ns/op         0  B/op   0  allocs/op

# Parse
Benchmark_ParseUint/fiber-12                                     170000696   7.154  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                               158467486   7.395  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                    87813716   13.73  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                      189251251   7.887  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                161616590   6.842  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                     74706812   16.44  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber-12                                    154112598   7.862  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                              142281115   8.354  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                   73937535   16.18  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                    218293975   5.513  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                              224941040   5.378  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                  100000000   11.25  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                     276938360   4.313  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                               257654323   4.699  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                   137141400   8.721  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                   166424791   7.125  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                             161026850   7.482  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                  88484428   13.76  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                   259818124   4.600  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                             254667759   4.730  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                 133563345   9.041  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                    343606695   3.430  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                              321277720   3.691  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                  180413733   6.689  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber-12                                  100000000   12.28  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                            100000000   15.33  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                 42162712   30.44  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                   83766469   14.12  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                            100000000   13.69  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                 42878646   30.11  ns/op         0  B/op   0  allocs/op

# Time
Benchmark_CalculateTimestamp/fiber-12                           1000000000  0.3015  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                           36791053   32.82  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                     4432272   272.0  ns/op        14  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                   3979948   300.7  ns/op         8  B/op   2  allocs/op

# XML
Benchmark_GolangXMLEncoder-12                                       581358    2059  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                      578254    2067  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                      325248    3488  ns/op      2864  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>
