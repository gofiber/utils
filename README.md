# Fiber Utils

![Release](https://img.shields.io/github/release/gofiber/utils.svg)
![Test](https://github.com/gofiber/utils/workflows/Test/badge.svg)
![Codecov](https://img.shields.io/codecov/c/github/gofiber/utils?token=3Cr92CwaPQ&style=flat-square&logo=codecov&label=codecov)
![Linter](https://github.com/gofiber/utils/actions/workflows/lint.yml/badge.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)

A collection of common functions for [Fiber](https://github.com/gofiber/fiber) with better performance, fewer allocations, and fewer dependencies.

## Benchmarks

Environment:
goos: darwin
goarch: arm64
pkg: github.com/gofiber/utils/v2
cpu: Apple M2 Pro

```text
// go test ./... -benchmem -run=^$ -bench=Benchmark_ -count=1

# Case Conversion
Benchmark_ToLowerBytes/empty/fiber-12                              601228760   1.955  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/unsafe-12                       601901127   1.977  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                            361526772   3.343  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                           356280168   3.373  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/unsafe-12                    232881622   5.121  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                          64825149   18.35  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                      95134903   12.71  ns/op         3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/unsafe-12              233711541   5.175  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                    83459391   13.94  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12           46423608   25.82  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/unsafe-12    94087516   12.69  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12         25226778   47.33  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                         42942836   27.95  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/unsafe-12                  61712125   19.44  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                       17121913   69.94  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                         36060220   33.69  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/unsafe-12                  61862710   19.45  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                       16092603   74.24  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                         35001792   34.35  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/unsafe-12                  62047836   19.42  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                       16119986   74.45  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/empty/fiber-12                              621527377   1.948  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/unsafe-12                       620614620   1.903  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                            347387406   3.339  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                            93976990   12.81  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/unsafe-12                    228022617   5.177  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                          80799015   14.77  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                     353756593   3.385  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/unsafe-12              234419678   5.132  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                    63601288   18.92  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12           44992688   26.51  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/unsafe-12    95225496   12.77  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12         22269578   53.94  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                         35570046   33.42  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/unsafe-12                  61963066   19.52  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                       14771874   80.05  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                         43817876   27.94  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/unsafe-12                  61968400   19.39  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                       14641957   81.42  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                         35656184   34.25  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/unsafe-12                  56226717   19.48  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                       14955886   80.25  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/empty/fiber-12                                   539313216   2.018  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/unsafe-12                            634806298   1.889  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                                 638947015   1.882  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                                 95414156   12.67  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/unsafe-12                         234368000   5.118  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                               60732201   19.48  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                          248156660   4.837  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/unsafe-12                   234148320   5.125  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                        234811851   5.098  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12                45805240   26.20  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/unsafe-12         93765258   12.83  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12              11631184   103.2  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                              36103797   33.35  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/unsafe-12                       56551980   20.95  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                             5996400   177.5  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                              28375362   42.48  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/unsafe-12                       61935218   19.44  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                            19255944   61.54  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                              36390100   33.16  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/unsafe-12                       62494440   19.35  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                             5833207   204.4  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/empty/fiber-12                                   613171694   1.947  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/unsafe-12                            647631400   1.858  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                                 638164638   1.883  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                                394179825   3.058  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/unsafe-12                         235441498   5.085  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                              279097488   4.179  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                           98945953   12.38  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/unsafe-12                   236656982   5.051  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                         62779154   19.15  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12                47738392   25.17  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/unsafe-12         94119801   12.68  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12              18145400   65.65  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                              44083813   27.34  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/unsafe-12                       60995110   19.58  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                            24336615   48.71  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                              35895274   33.27  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/unsafe-12                       61949072   19.50  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                             7084752   168.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                              35053681   34.56  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/unsafe-12                       61658090   19.58  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                             5924719   206.8  ns/op        64  B/op   1  allocs/op

# Add Trailing Slash
Benchmark_AddTrailingSlashBytes/empty-12                          1000000000  0.6968  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                     1000000000   1.072  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                   87295214   16.10  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12               1000000000   1.055  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                    49061860   28.34  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12                1000000000   1.053  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/empty-12                         1000000000  0.2986  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                    1000000000  0.4552  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12                 100000000   11.43  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12              1000000000  0.4569  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                   84621513   14.61  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12               1000000000  0.4529  ns/op         0  B/op   0  allocs/op

# EqualFold
Benchmark_EqualFoldBytes/fiber-12                                   24560449   43.80  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                                 14990294   84.85  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-12                                        31785342   44.77  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                      12664594   87.50  ns/op         0  B/op   0  allocs/op

# Trim
Benchmark_TrimRight/fiber-12                                       485636703   2.319  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                     389052968   2.984  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/fiber-12                                  570778065   2.153  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                                354338642   3.378  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/fiber-12                                        546707935   2.161  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                      395816548   3.839  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/fiber-12                                   517814619   2.408  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                                 344487759   3.479  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/fiber-12                                            268123645   4.985  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default-12                                          216317569   5.603  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                                196718382   6.164  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/fiber-12                                       276413652   4.682  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                     206646734   6.416  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                           180688819   7.050  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/empty-12                                1000000000  0.3595  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                               530890448   2.519  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                                274174633   3.941  ns/op    761.16  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/spaces-12                              287224635   4.902  ns/op    611.98  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                            291975805   4.577  ns/op   1966.40  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                          205111491   5.559  ns/op   1618.98  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                    249115078   5.141  ns/op   4668.43  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                  181133806   6.330  ns/op   3791.28  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                     358262827   3.202  ns/op  11868.28  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                   246771914   4.861  ns/op   7817.71  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                       327337672   3.649  ns/op   5754.64  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                     263390455   4.551  ns/op   4614.63  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                          328083764   3.669  ns/op   5451.76  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/content-type-12                        263841824   4.557  ns/op   4388.98  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                       439685796   2.752  ns/op  10175.38  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                     281075425   4.245  ns/op   6595.73  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                          327527677   3.652  ns/op   5476.47  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/query-params-12                        263637592   4.546  ns/op   4399.23  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                            263431845   4.665  ns/op  21864.72  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                          196392247   6.130  ns/op  16640.29  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                              1000000000  0.6094  ns/op   8205.02  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                             438738540   2.760  ns/op   1811.78  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                      280189737   4.289  ns/op   3963.47  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                    220052497   5.464  ns/op   3111.01  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/empty-12                           1000000000  0.3036  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                          492696795   2.426  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                           390392758   3.090  ns/op    970.86  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                         294485332   4.060  ns/op    738.89  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                       326114242   3.662  ns/op   2457.47  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                     245688716   4.901  ns/op   1836.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12               281842446   4.254  ns/op   5642.35  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12             208275072   5.765  ns/op   4162.93  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12                375564568   3.186  ns/op  11925.46  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12              232178205   5.161  ns/op   7362.63  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                  327044280   3.674  ns/op   5715.36  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12                246186018   4.882  ns/op   4301.08  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                     327820933   3.673  ns/op   5445.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                   247571902   4.851  ns/op   4123.23  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                  433509044   2.763  ns/op  10132.97  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12                263725153   4.551  ns/op   6153.07  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                     327897222   3.666  ns/op   5455.18  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                   246765888   4.859  ns/op   4115.96  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                       262128422   4.558  ns/op  22377.17  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                     188152813   6.363  ns/op  16030.46  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                         1000000000  0.6211  ns/op   8050.53  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                        395639230   3.036  ns/op   1646.73  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12                 282314606   4.243  ns/op   4007.05  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12               207869355   5.766  ns/op   2948.57  MB/s   0  B/op  0  allocs/op

# Convert
Benchmark_ConvertToBytes/fiber-12                                  303731499   3.947  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int-12                                          497704682   2.406  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                         607933017   1.975  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                        484205890   2.481  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                        490413688   2.459  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                        500168806   2.405  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                         568876189   2.106  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                        609816264   1.975  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                       567892762   2.102  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                       570461389   2.123  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                       567259215   2.119  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/string-12                                       606332637   1.993  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                       93052707   12.31  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                         603432777   1.995  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                       25012396   48.04  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                       17169336   70.51  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                     12940731   92.84  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                  12885784   93.52  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                      51672447   23.11  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                         51177982   23.30  ns/op         8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                        29876820   38.88  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                        8826847   136.3  ns/op       112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                                 9783631   120.3  ns/op        72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                             398624744   3.015  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                              11562967   104.5  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/int-12                             1000000000  0.2675  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int8-12                            1000000000  0.2241  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int16-12                           1000000000  0.2761  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int32-12                           1000000000  0.2719  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int64-12                           1000000000  0.2673  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint-12                            1000000000  0.2339  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint8-12                           1000000000  0.2289  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint16-12                          1000000000  0.2354  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint32-12                          1000000000  0.2309  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint64-12                          1000000000  0.2364  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/string-12                          1000000000  0.2189  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/[]uint8-12                          205632634   5.758  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/bool-12                            1000000000  0.2346  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/float32-12                          187598121   6.133  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/float64-12                          135170138   8.719  ns/op         4  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time-12                         74157042   18.63  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time#01-12                      61104981   20.85  ns/op        24  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]string-12                         148502100   7.790  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]int-12                            230634176   5.118  ns/op         8  B/op   1  allocs/op
Benchmark_ToString_concurrency/[2]int-12                           133293684   9.173  ns/op        16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[][]int-12                           22382804   54.81  ns/op       112  B/op   6  allocs/op
Benchmark_ToString_concurrency/[]interface_{}-12                    33572145   36.32  ns/op        72  B/op   3  allocs/op
Benchmark_ToString_concurrency/utils.MyStringer-12                1000000000  0.3800  ns/op         0  B/op   0  allocs/op
Benchmark_ToString_concurrency/utils.CustomType-12                  54238104   23.27  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeBytes/unsafe-12                                   1000000000  0.5654  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                    88676498   13.68  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                  1000000000  0.3530  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                  100000000   11.92  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/0-12                                            400988882   3.004  ns/op         0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                             68500806   17.65  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                           60281438   19.95  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                          59684172   19.67  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                          53318622   26.05  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                       59837815   21.27  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                       52797439   22.84  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                    60002750   19.76  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                    52817191   22.15  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                                 64727100   19.55  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                                 54228812   21.65  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                              64025749   18.28  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                              57805672   22.45  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                           64201592   18.70  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                           58673847   21.37  ns/op        16  B/op   1  allocs/op

# Format and Append
Benchmark_FormatUint/small/fiber-12                                616120275   1.968  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                              597079038   1.990  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                                61883050   18.90  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/medium/strconv-12                              60313885   20.38  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                                 47939197   25.36  ns/op        24  B/op   1  allocs/op
Benchmark_FormatUint/large/strconv-12                               44694955   26.49  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-12                             608629350   1.963  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                           618271995   1.956  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                             626947920   1.919  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                            87661094   14.05  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                             55510044   24.13  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                           49636088   20.89  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                             63894786   19.05  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                           62000552   19.78  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                              47706525   25.50  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                            45759099   26.59  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                              47631805   25.46  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                            45561183   26.65  ns/op        24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-12                                     62901873   19.08  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint32/strconv-12                                   59485326   20.61  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                      61375447   19.29  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                    60683575   19.99  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                     80310759   15.03  ns/op         5  B/op   1  allocs/op
Benchmark_FormatUint16/strconv-12                                   74569866   15.74  ns/op         5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                      76845693   15.44  ns/op         8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                    78253000   15.80  ns/op         8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                     604374918   1.973  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                    82805967   14.56  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                      610502032   1.969  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                     82919455   14.59  ns/op         4  B/op   1  allocs/op
Benchmark_AppendUint/fiber-12                                      130580040   9.204  ns/op         0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                    100000000   10.15  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                             363979753   3.293  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                           212234758   5.689  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                            125316576   9.624  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                          122066809   9.835  ns/op         0  B/op   0  allocs/op

# Token
Benchmark_GenerateSecureToken/16_bytes-12                            4395446   268.0  ns/op        24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                            4308816   277.5  ns/op        48  B/op   1  allocs/op
Benchmark_TokenGenerators/UUIDv4-12                                  4115958   294.5  ns/op        64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                             4320008   278.1  ns/op        48  B/op   1  allocs/op

# HTTP
Benchmark_GetMIME/fiber-12                                          22006587   54.74  ns/op         0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                        17412228   69.29  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/vendorContentType-12      125887742   9.981  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12     353516896   3.504  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/fiber-12                                  1000000000  0.3850  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                                 431016290   2.751  ns/op         0  B/op   0  allocs/op

# IP
Benchmark_IsIPv4/fiber-12                                           80343021   15.43  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                         53110760   22.82  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/fiber-12                                           27125079   44.82  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                         20017320   60.60  ns/op         0  B/op   0  allocs/op

# Parse
Benchmark_ParseUint/fiber-12                                       170572611   7.031  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                                 161515281   7.414  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                      84878630   13.91  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                        188326305   6.554  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                  177401347   6.791  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                       74426298   16.02  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber-12                                      152689100   7.792  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                                143480342   8.350  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                     75226464   15.95  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                      233907660   5.214  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                                218961453   5.362  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                    100000000   11.20  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                       277644391   4.368  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                                 255406502   4.720  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                     137303502   8.716  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                     168541228   7.201  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                               160750947   7.499  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                    87199530   13.95  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                     259066075   4.625  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                               249345403   4.822  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                   132290673   9.031  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                      346359704   3.473  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                                296839494   4.125  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                    180975388   6.673  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber-12                                     84649861   14.23  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                               83656969   14.19  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                   48787262   24.61  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                     78740156   15.25  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                               78525037   15.15  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                   45575605   26.21  ns/op         0  B/op   0  allocs/op

# Time
Benchmark_CalculateTimestamp/fiber-12                             1000000000  0.3411  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                             36431808   33.04  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                       4388325   271.9  ns/op        10  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                     3753456   303.1  ns/op         8  B/op   2  allocs/op

# XML
Benchmark_GolangXMLEncoder-12                                         581364    2083  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                        591445    2080  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                        326768    3486  ns/op      2862  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>

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
