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

```go
// go test -benchmem -run=^$ -bench=Benchmark_ -count=1

Benchmark_ToLowerBytes/fiber-12                                   52217712   22.65  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/default-12                                 13652452   88.28  ns/op        80  B/op   1  allocs/op

Benchmark_ToUpperBytes/fiber-12                                   53019828   22.65  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/default-12                                 12197802   100.6  ns/op        80  B/op   1  allocs/op

Benchmark_AddTrailingSlashBytes/empty-12                        1000000000  0.6182  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                   1000000000  0.9144  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                100000000   11.64  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12             1000000000   1.008  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                  85574719   13.53  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12              1000000000  0.9137  ns/op         0  B/op   0  allocs/op

Benchmark_EqualFoldBytes/fiber-12                                 33551104   43.94  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                               17184109   69.09  ns/op         0  B/op   0  allocs/op

Benchmark_EqualFold/fiber-12                                      34499653   34.97  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                    17140081   70.01  ns/op         0  B/op   0  allocs/op

Benchmark_TrimRight/fiber-12                                     566421138   2.128  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                   400385982   3.009  ns/op         0  B/op   0  allocs/op

Benchmark_TrimRightBytes/fiber-12                                563772230   2.148  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                              353540287   3.408  ns/op         0  B/op   0  allocs/op

Benchmark_TrimLeft/fiber-12                                      552920756   2.139  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                    396947145   3.104  ns/op         0  B/op   0  allocs/op

Benchmark_TrimLeftBytes/fiber-12                                 566912395   2.150  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                               365988206   3.205  ns/op         0  B/op   0  allocs/op

Benchmark_Trim/fiber-12                                          293546023   4.058  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default-12                                        244428660   4.873  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                              226544163   5.206  ns/op         0  B/op   0  allocs/op

Benchmark_TrimBytes/fiber-12                                     305831440   3.948  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                   241251229   4.906  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                         219920371   5.527  ns/op         0  B/op   0  allocs/op

Benchmark_TrimSpace/fiber/empty-12                              1000000000  0.3343  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                             539124730   2.200  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                              380309074   3.178  ns/op    943.93  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/spaces-12                            323342026   4.017  ns/op    746.85  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                          305200390   3.672  ns/op   2450.94  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                        264772220   4.646  ns/op   1936.96  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                  281984287   4.258  ns/op   5635.84  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                219716269   5.472  ns/op   4385.89  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                   376020445   3.179  ns/op  11954.53  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                 248763286   4.826  ns/op   7873.39  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                     329946304   3.694  ns/op   5685.60  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                   265207191   4.637  ns/op   4529.26  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                        328451793   3.642  ns/op   5491.69  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/content-type-12                      264308883   4.539  ns/op   4406.72  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                     439010670   2.728  ns/op  10262.42  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                   283859259   4.246  ns/op   6594.82  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                        329072942   3.651  ns/op   5477.44  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/query-params-12                      263042356   4.541  ns/op   4404.12  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                          265053544   4.533  ns/op  22500.56  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                        197880016   6.055  ns/op  16844.72  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                            1000000000  0.6066  ns/op   8242.41  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                           440016036   2.721  ns/op   1837.51  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                    283569202   4.241  ns/op   4008.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                  220461590   5.459  ns/op   3114.31  MB/s   0  B/op  0  allocs/op

Benchmark_TrimSpaceBytes/fiber/empty-12                         1000000000  0.3063  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                        488930062   2.420  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                         384257913   3.142  ns/op    954.77  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                       296069066   4.080  ns/op    735.22  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                     327912043   3.641  ns/op   2471.75  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                   248219776   4.847  ns/op   1857.01  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12             282536949   4.265  ns/op   5627.43  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12           207584353   5.822  ns/op   4122.41  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12              377500597   3.180  ns/op  11950.47  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12            233271138   5.132  ns/op   7404.96  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                325874539   3.654  ns/op   5747.30  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12              248342956   4.833  ns/op   4345.28  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                   331166112   3.621  ns/op   5524.04  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                 248150503   4.854  ns/op   4120.32  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                427332102   2.750  ns/op  10180.00  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12              264074143   4.530  ns/op   6181.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                   329397716   3.647  ns/op   5484.14  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                 247608130   4.853  ns/op   4121.02  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                     263368270   4.538  ns/op  22474.68  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                   189202016   6.345  ns/op  16076.90  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                       1000000000  0.6119  ns/op   8171.03  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                      382449446   3.019  ns/op   1656.19  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12               283825689   4.237  ns/op   4012.62  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12             208716372   5.750  ns/op   2956.38  MB/s   0  B/op  0  allocs/op

Benchmark_ConvertToBytes/fiber-12                                305627947   3.929  ns/op         0  B/op   0  allocs/op

Benchmark_GenerateSecureToken/16_bytes-12                          4514146   266.7  ns/op        24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                          4297660   278.9  ns/op        48  B/op   1  allocs/op

Benchmark_TokenGenerators/UUIDv4-12                                3983023   295.1  ns/op        64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                           4314350   279.4  ns/op        48  B/op   1  allocs/op

Benchmark_ToString/int-12                                        498315158   2.417  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                       603348963   1.984  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                      481211367   2.503  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                      484490580   2.467  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                      496510891   2.428  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                       566930694   2.117  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                      605057905   1.983  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                     565606969   2.139  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                     565128951   2.119  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                     566257761   2.140  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/string-12                                     606878588   2.229  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                     86469348   12.59  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                       597673555   2.007  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                     24551800   49.38  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                     16948284   71.20  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                   12658956   94.86  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                12711673   95.25  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                    53517573   23.47  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                       46515358   23.34  ns/op         8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                      31373779   46.64  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                      8489103   140.9  ns/op       112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                               9858184   125.0  ns/op        72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                           361391401   3.333  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                            10969280   109.0  ns/op        16  B/op   1  allocs/op

Benchmark_UnsafeBytes/unsafe-12                                 1000000000  0.5338  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                  86281022   13.79  ns/op        16  B/op   1  allocs/op

Benchmark_UnsafeString/unsafe-12                                1000000000  0.3604  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                100000000   11.86  ns/op        16  B/op   1  allocs/op

Benchmark_ByteSize/0-12                                          395636622   3.031  ns/op         0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                           65794736   18.12  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                         58942592   20.66  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                        61452717   19.98  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                        50975253   23.16  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                     59773606   19.77  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                     52359066   22.78  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                  63552447   19.28  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                  53362782   22.54  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                               63670188   19.03  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                               54096292   22.33  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                            64223641   18.77  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                            55165544   22.09  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                         64620789   18.78  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                         54445122   21.99  ns/op        16  B/op   1  allocs/op

Benchmark_FormatUint/small/fiber-12                             1000000000  0.3043  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                            576869341   2.057  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                             126048157   14.77  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/strconv-12                            25518092   58.39  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                               60826864   19.85  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/large/strconv-12                             40973985   40.63  ns/op        24  B/op   1  allocs/op

Benchmark_FormatInt/small_pos/fiber-12                           563756449   2.207  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                         507117070   2.292  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                           532399563   2.090  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                          81482990   14.91  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                           58882579   21.28  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                         53400873   22.44  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                           54462522   22.08  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                         55634112   22.38  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                            41997811   27.73  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                          40634805   29.16  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                            42171849   39.04  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                          35387965   30.66  ns/op        24  B/op   1  allocs/op

Benchmark_FormatUint32/fiber-12                                  126037087   9.696  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint32/strconv-12                                 56969912   21.20  ns/op        16  B/op   1  allocs/op

Benchmark_FormatInt32/fiber-12                                    62566394   19.71  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                  59540052   20.42  ns/op        16  B/op   1  allocs/op

Benchmark_FormatUint16/fiber-12                                  162967428   7.361  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint16/strconv-12                                 74978522   16.09  ns/op         5  B/op   1  allocs/op

Benchmark_FormatInt16/fiber-12                                    75251031   16.00  ns/op         8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                  74987696   16.04  ns/op         8  B/op   1  allocs/op

Benchmark_FormatUint8/fiber-12                                  1000000000  0.3028  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                  81920348   14.83  ns/op         3  B/op   1  allocs/op

Benchmark_FormatInt8/fiber-12                                   1000000000  0.3023  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                   80871388   14.82  ns/op         4  B/op   1  allocs/op

Benchmark_AppendUint/fiber-12                                    134814058   8.961  ns/op         0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                  100000000   10.35  ns/op         0  B/op   0  allocs/op

Benchmark_AppendInt/small_neg/fiber-12                           359453137   3.339  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                         206816691   5.833  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                          123050035   9.767  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                        100000000   10.11  ns/op         0  B/op   0  allocs/op

Benchmark_GetMIME/fiber-12                                        28987636   41.44  ns/op         0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                      17214522   71.37  ns/op         0  B/op   0  allocs/op

Benchmark_ParseVendorSpecificContentType/vendorContentType-12    154834189   7.174  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12   396299224   3.062  ns/op         0  B/op   0  allocs/op

Benchmark_StatusMessage/fiber-12                                1000000000  0.3882  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                               417715003   2.849  ns/op         0  B/op   0  allocs/op

Benchmark_IsIPv4/fiber-12                                         79709500   15.98  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                       51886188   25.12  ns/op         0  B/op   0  allocs/op

Benchmark_IsIPv6/fiber-12                                         20887682   47.90  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                       19282450   62.15  ns/op         0  B/op   0  allocs/op

Benchmark_ParseUint/fiber-12                                     170378602   7.081  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                               161110435   7.483  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                    84852877   14.02  ns/op         0  B/op   0  allocs/op

Benchmark_ParseInt/fiber-12                                      185977370   6.438  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                175932828   6.941  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                     73079564   16.15  ns/op         0  B/op   0  allocs/op

Benchmark_ParseInt32/fiber-12                                    152471434   7.848  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                              140770058   8.489  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                   76423180   16.18  ns/op         0  B/op   0  allocs/op

Benchmark_ParseInt16/fiber-12                                    227571576   5.289  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                              218097106   5.423  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                  100000000   11.22  ns/op         0  B/op   0  allocs/op

Benchmark_ParseInt8/fiber-12                                     271443954   4.481  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                               261060393   4.757  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                   130782296   8.848  ns/op         0  B/op   0  allocs/op

Benchmark_ParseUint32/fiber-12                                   167752328   7.235  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                             154525927   7.583  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                  86833035   14.06  ns/op         0  B/op   0  allocs/op

Benchmark_ParseUint16/fiber-12                                   261109554   4.673  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                             253868272   4.763  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                 130897095   9.147  ns/op         0  B/op   0  allocs/op

Benchmark_ParseUint8/fiber-12                                    351570009   3.477  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                              318383847   3.740  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                  177199309   6.772  ns/op         0  B/op   0  allocs/op

Benchmark_ParseFloat64/fiber-12                                  100000000   10.84  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                            100000000   11.28  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                 48993004   24.95  ns/op         0  B/op   0  allocs/op

Benchmark_ParseFloat32/fiber-12                                  100000000   11.76  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                             99723333   12.20  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                 45525679   26.53  ns/op         0  B/op   0  allocs/op

Benchmark_ToUpper/empty/fiber-12                                 593092693   2.570  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                               627143008   1.881  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/single-lower/fiber-12                          100000000   10.17  ns/op         1  B/op   1  allocs/op
Benchmark_ToUpper/single-lower/default-12                         73451773   16.49  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/single-upper/fiber-12                          610761877   1.980  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/single-upper/default-12                        432002700   2.737  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/numbers/fiber-12                               219887310   5.482  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/numbers/default-12                             100000000   11.58  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                               92281071   13.05  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/default-12                             57422334   21.63  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                        391365019   3.106  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                      229001715   5.235  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-pascal/fiber-12                        90588537   13.34  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get-pascal/default-12                      45610612   26.69  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12              37212859   33.25  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12            10908462   109.6  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/url-camel/fiber-12                              37313818   32.83  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/url-camel/default-12                            10856359   108.4  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/15-char-mixed/fiber-12                          59459043   20.17  ns/op        16  B/op   1  allocs/op
Benchmark_ToUpper/15-char-mixed/default-12                        20162306   59.58  ns/op        16  B/op   1  allocs/op
Benchmark_ToUpper/16-char-mixed/fiber-12                          60514624   20.47  ns/op        16  B/op   1  allocs/op
Benchmark_ToUpper/16-char-mixed/default-12                        19192668   63.74  ns/op        16  B/op   1  allocs/op
Benchmark_ToUpper/24-char-mixed/fiber-12                          44284731   26.10  ns/op        24  B/op   1  allocs/op
Benchmark_ToUpper/24-char-mixed/default-12                        13136126   92.17  ns/op        24  B/op   1  allocs/op
Benchmark_ToUpper/25-char-mixed/fiber-12                          45794679   25.74  ns/op        32  B/op   1  allocs/op
Benchmark_ToUpper/25-char-mixed/default-12                        12740172   94.56  ns/op        32  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                            28999896   41.59  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/default-12                           6709502   179.4  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                            26937060   42.13  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/default-12                           5650413   213.6  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                            43262665   28.13  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                          18352060   66.25  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-first-upper/fiber-12                      27661321   45.10  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-first-upper/default-12                     6342034   181.9  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-last-upper/fiber-12                       26616956   41.79  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-last-upper/default-12                      6641599   180.2  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/very-large-lower/fiber-12                        9267004   130.2  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-lower/default-12                      1643893   721.2  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-mixed/fiber-12                        9065834   131.6  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-mixed/default-12                      1450942   834.5  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-upper/fiber-12                       10477177   113.0  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/very-large-upper/default-12                      4982967   240.7  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/very-large-first-upper/fiber-12                  9188516   131.4  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-first-upper/default-12                1667226   758.6  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-last-upper/fiber-12                   9347838   153.1  ns/op       256  B/op   1  allocs/op
Benchmark_ToUpper/very-large-last-upper/default-12                 1675723   720.9  ns/op       256  B/op   1  allocs/op

Benchmark_ToLower/empty/fiber-12                                 609407731   2.000  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                               632477314   1.901  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/single-lower/fiber-12                          637803123   1.887  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/single-lower/default-12                        484516336   2.448  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/single-upper/fiber-12                          100000000   10.22  ns/op         1  B/op   1  allocs/op
Benchmark_ToLower/single-upper/default-12                         74222977   16.80  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/numbers/fiber-12                               219655238   5.495  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/numbers/default-12                             100000000   11.54  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                              390467918   3.110  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                            278361213   4.361  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                         88689062   13.11  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/default-12                       58492478   20.85  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/http-get-pascal/fiber-12                        89783954   13.10  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-pascal/default-12                      54137679   22.84  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12              35148870   34.78  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12            16135754   73.95  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/url-camel/fiber-12                              36275648   33.21  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/url-camel/default-12                            14417625   83.09  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/15-char-mixed/fiber-12                          62061474   20.85  ns/op        16  B/op   1  allocs/op
Benchmark_ToLower/15-char-mixed/default-12                        19297864   61.86  ns/op        16  B/op   1  allocs/op
Benchmark_ToLower/16-char-mixed/fiber-12                          56775728   21.68  ns/op        16  B/op   1  allocs/op
Benchmark_ToLower/16-char-mixed/default-12                        18371107   64.68  ns/op        16  B/op   1  allocs/op
Benchmark_ToLower/24-char-mixed/fiber-12                          47854680   26.81  ns/op        24  B/op   1  allocs/op
Benchmark_ToLower/24-char-mixed/default-12                        12394261   98.55  ns/op        24  B/op   1  allocs/op
Benchmark_ToLower/25-char-mixed/fiber-12                          43340204   27.89  ns/op        32  B/op   1  allocs/op
Benchmark_ToLower/25-char-mixed/default-12                        12124130   101.2  ns/op        32  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                            43317716   28.15  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                          23931813   50.79  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                            28097258   42.42  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/default-12                           5665676   216.8  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber-12                            26929779   44.31  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/default-12                           6590605   180.7  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-first-upper/fiber-12                      29048421   44.22  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-first-upper/default-12                    10758439   115.3  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-last-upper/fiber-12                       18971841   60.46  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-last-upper/default-12                      9536596   126.3  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/very-large-lower/fiber-12                       10654377   113.1  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/very-large-lower/default-12                      6492882   183.9  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/very-large-mixed/fiber-12                        9478850   130.5  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-mixed/default-12                      1450599   835.5  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-upper/fiber-12                        8886003   135.3  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-upper/default-12                      1689078   734.8  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-first-upper/fiber-12                  9459594   127.4  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-first-upper/default-12                3244156   371.2  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-last-upper/fiber-12                   6312432   177.9  ns/op       256  B/op   1  allocs/op
Benchmark_ToLower/very-large-last-upper/default-12                 2909698   411.4  ns/op       256  B/op   1  allocs/op

Benchmark_AddTrailingSlashString/empty-12                       1000000000  0.3047  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                  1000000000  0.4549  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12               100000000   11.68  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12            1000000000  0.4687  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                 81040472   15.26  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12             1000000000  0.4687  ns/op         0  B/op   0  allocs/op

Benchmark_CalculateTimestamp/fiber-12                           1000000000  0.3204  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                           35641314   33.76  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                     4361775   271.9  ns/op        12  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                   3927289   301.7  ns/op         8  B/op   2  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>
