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

Benchmark_ToLowerBytes/empty/fiber-12                           553808877  1.978  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/mut-12                       577751976  2.098  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                         339363715  3.543  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                        100000000  11.85  ns/op         3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/mut-12                    226510630  5.298  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                       63573069  18.38  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                  100000000  11.93  ns/op         3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/mut-12              227504000  5.298  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                 82697782  14.41  ns/op         8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12        48103501  24.73  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/mut-12    92346762  12.94  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12      25210789  47.24  ns/op        48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                      37025133  32.06  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/mut-12                  60762954  19.56  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                    17600511  68.30  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                      38087983  32.14  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/mut-12                  61083208  19.59  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                    16222710  73.69  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                      37883671  32.01  ns/op        64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/mut-12                  63007148  19.76  ns/op         0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                    16186358  74.95  ns/op        64  B/op   1  allocs/op

Benchmark_ToUpperBytes/empty/fiber-12                           616420280  1.944  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/mut-12                       579055560  2.079  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                         365963186  3.296  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                        100000000  11.78  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/mut-12                    225143792  5.328  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                       82009225  14.28  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                  100000000  11.69  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/mut-12              222175198  5.353  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                 62712169  18.94  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12        49937577  24.38  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/mut-12    97075597  12.55  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12      25080837  47.48  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                      37467458  32.27  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/mut-12                  61842121  19.36  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                    15869980  73.70  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                      37423495  32.41  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/mut-12                  60528235  19.50  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                    14926092  81.00  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                      37440182  32.31  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/mut-12                  62480474  19.34  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                    16243902  73.73  ns/op        64  B/op   1  allocs/op

Benchmark_AddTrailingSlashBytes/empty-12                        629800314  1.917  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                   629890261  1.919  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12               100000000  11.97  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12             623784649  1.951  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                 89914021  13.52  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12              626118259  1.920  ns/op         0  B/op   0  allocs/op

Benchmark_EqualFoldBytes/fiber-12                                34780592  34.38  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                              18034400  67.23  ns/op         0  B/op   0  allocs/op

Benchmark_EqualFold/fiber-12                                     35330140  33.87  ns/op         0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                   18055712  66.73  ns/op         0  B/op   0  allocs/op

Benchmark_TrimRight/fiber-12                                    371868604  3.236  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                  371310679  3.242  ns/op         0  B/op   0  allocs/op

Benchmark_TrimRightBytes/fiber-12                               362995164  3.310  ns/op         0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                             367253110  3.250  ns/op         0  B/op   0  allocs/op

Benchmark_TrimLeft/fiber-12                                     403846129  2.974  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                   370522134  3.250  ns/op         0  B/op   0  allocs/op

Benchmark_TrimLeftBytes/fiber-12                                403714339  2.975  ns/op         0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                              401688597  2.997  ns/op         0  B/op   0  allocs/op

Benchmark_Trim/fiber-12                                         254003532  4.710  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default-12                                       240257415  5.011  ns/op         0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                             226416286  5.310  ns/op         0  B/op   0  allocs/op

Benchmark_TrimBytes/fiber-12                                    248018350  4.763  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                  254058271  4.739  ns/op         0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                        213339448  5.610  ns/op         0  B/op   0  allocs/op

Benchmark_TrimSpace/fiber/empty-12                              619603210  1.936  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                            582245630  2.065  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                             332584594  3.618  ns/op    829.27  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/spaces-12                           332178206  3.648  ns/op    822.39  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                         267276278  4.477  ns/op   2010.24  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                       269852575  4.466  ns/op   2015.32  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                 235906568  5.054  ns/op   4748.49  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12               225898935  5.345  ns/op   4490.47  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                  310937485  3.861  ns/op   9843.27  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                254148128  4.717  ns/op   8055.67  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                    270272146  4.438  ns/op   4731.68  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                  272350263  4.429  ns/op   4741.20  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                       268778092  4.463  ns/op   4481.30  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/content-type-12                     271139079  4.423  ns/op   4521.70  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                    336905310  3.617  ns/op   7741.94  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                  288726260  4.151  ns/op   6745.85  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                       269243526  4.478  ns/op   4466.67  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/query-params-12                     270592054  4.427  ns/op   4518.09  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                         224852598  5.334  ns/op  19123.05  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                       202905722  5.903  ns/op  17280.58  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                            613568556  1.956  ns/op   2556.33  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                          451611829  2.665  ns/op   1876.10  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                   238122226  5.033  ns/op   3377.98  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                 224012803  5.391  ns/op   3153.38  MB/s   0  B/op  0  allocs/op

Benchmark_TrimSpaceBytes/fiber/empty-12                         621397833  1.935  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                       507128054  2.360  ns/op         0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                        309226214  3.906  ns/op    768.10  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                      303903385  3.925  ns/op    764.29  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                    253692010  4.736  ns/op   1900.25  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                  253788698  4.713  ns/op   1909.53  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12            226264802  5.313  ns/op   4517.29  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12          214769792  5.591  ns/op   4292.33  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12             284941381  4.185  ns/op   9079.48  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12           238264204  5.031  ns/op   7553.19  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12               251877955  4.779  ns/op   4394.52  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12             253428544  4.738  ns/op   4432.42  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                  251732319  4.785  ns/op   4180.00  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                254985970  4.711  ns/op   4245.65  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12               312281476  3.846  ns/op   7280.24  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12             270232204  4.434  ns/op   6315.42  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                  250553370  4.784  ns/op   4180.31  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                246349831  4.776  ns/op   4187.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                    210534056  5.698  ns/op  17901.04  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                  191590646  6.347  ns/op  16071.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                       596382814  2.008  ns/op   2490.39  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                     388967896  3.026  ns/op   1652.37  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12              221968550  5.371  ns/op   3165.36  MB/s   0  B/op  0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12            211475139  5.659  ns/op   3003.91  MB/s   0  B/op  0  allocs/op

Benchmark_ConvertToBytes/fiber-12                               309964785  3.828  ns/op         0  B/op   0  allocs/op

Benchmark_GenerateSecureToken/16_bytes-12                         4611472  263.8  ns/op        24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                         4086676  283.8  ns/op        48  B/op   1  allocs/op

Benchmark_TokenGenerators/UUIDv4-12                               4165256  290.4  ns/op        64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                          4286131  279.8  ns/op        48  B/op   1  allocs/op

Benchmark_ToString/int-12                                       481441917  2.482  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                      600992764  1.998  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                     463707780  2.569  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                     479085424  2.523  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                     477174490  2.527  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                      553763726  2.168  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                     604430094  1.985  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                    550149152  2.173  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                    562246438  2.141  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                    561967118  2.150  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/string-12                                    604462824  1.975  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                    93515604  12.51  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                      575655411  2.092  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                    24551822  47.81  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                    16832260  72.64  ns/op         4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                  12733485  94.25  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                               12724028  93.37  ns/op        24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                   53674262  22.49  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                      51478485  23.04  ns/op         8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                     32330888  37.29  ns/op        16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                     8759301  137.2  ns/op       112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                              9940117  120.9  ns/op        72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                          392836462  3.023  ns/op         0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                           11243797  107.8  ns/op        16  B/op   1  allocs/op

Benchmark_UnsafeBytes/unsafe-12                                 610311853  1.987  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                 82091510  14.07  ns/op        16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                604025956  1.990  ns/op         0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                               100000000  12.11  ns/op        16  B/op   1  allocs/op

Benchmark_ByteSize/0-12                                         392166794  3.061  ns/op         0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                          66548512  18.53  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                        58134134  20.21  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                       60392170  19.91  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                       54596017  22.59  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                    62401551  19.91  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                    52901695  22.18  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                 63226388  19.36  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                 53183532  22.48  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                              64601944  19.02  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                              54066130  23.61  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                           62276330  18.68  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                           54644703  22.02  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                        63237356  18.96  ns/op        16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                        55583217  21.60  ns/op        16  B/op   1  allocs/op

Benchmark_FormatUint/small/fiber-12                             608818609  1.965  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                           607433466  1.980  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                             64176270  19.30  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/medium/strconv-12                           58552411  20.86  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                              45183770  26.64  ns/op        24  B/op   1  allocs/op
Benchmark_FormatUint/large/strconv-12                            42529689  27.64  ns/op        24  B/op   1  allocs/op

Benchmark_FormatInt/small_pos/fiber-12                          603535461  1.998  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                        602184674  2.013  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                          603278820  1.978  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                         86874946  13.98  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                          61292245  19.44  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                        59327985  21.76  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                          61344859  23.89  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                        59985126  20.41  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                           45635690  26.25  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                         44122918  34.24  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                           46242253  26.05  ns/op        24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                         44057848  28.10  ns/op        24  B/op   1  allocs/op

Benchmark_FormatUint32/fiber-12                                  60739248  19.31  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint32/strconv-12                                58372975  23.66  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                   61600984  19.70  ns/op        16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                 58851537  20.52  ns/op        16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                  79721198  15.13  ns/op         5  B/op   1  allocs/op
Benchmark_FormatUint16/strconv-12                                74082648  19.79  ns/op         5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                   71715822  16.60  ns/op         8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                 77185308  15.68  ns/op         8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                  602110650  1.975  ns/op         0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                 83953016  14.56  ns/op         3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                   609033748  1.969  ns/op         0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                  82871737  14.56  ns/op         4  B/op   1  allocs/op

Benchmark_AppendUint/fiber-12                                   130507630  9.184  ns/op         0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                 100000000  10.28  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                          356968003  3.330  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                        208894086  5.737  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                         124183686  9.662  ns/op         0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                       120313234  9.983  ns/op         0  B/op   0  allocs/op

Benchmark_GetMIME/fiber-12                                       29476454  41.18  ns/op         0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                     16920612  69.01  ns/op         0  B/op   0  allocs/op

Benchmark_ParseVendorSpecificContentType/vendorContentType-12   129802196  9.246  ns/op         0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12  397311744  3.214  ns/op         0  B/op   0  allocs/op

Benchmark_StatusMessage/fiber-12                                608673082  1.991  ns/op         0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                              444303130  2.741  ns/op         0  B/op   0  allocs/op

Benchmark_IsIPv4/fiber-12                                        79902118  15.07  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv4/default-12                                      35876672  34.15  ns/op        16  B/op   1  allocs/op
Benchmark_IsIPv6/fiber-12                                        26939730  44.64  ns/op         0  B/op   0  allocs/op
Benchmark_IsIPv6/default-12                                      15935161  75.35  ns/op        16  B/op   1  allocs/op

Benchmark_ParseUint/fiber-12                                    163655462  7.337  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                              155114128  7.766  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                   78715622  13.96  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                     187248835  6.394  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                               179853778  6.685  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                    73579004  16.04  ns/op         0  B/op   0  allocs/op

Benchmark_ParseInt32/fiber-12                                   147265508  8.134  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                             142771770  8.404  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                  75315775  15.97  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                   204820798  5.829  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                             196833240  6.097  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                 100000000  11.30  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                    260668231  4.632  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                              255214980  4.693  ns/op         0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                  133742324  8.968  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                  163530393  7.349  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                            155498348  7.685  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                 86542878  13.98  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                  233287539  5.147  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                            220876884  5.449  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                130628672  9.160  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                   327126895  3.656  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                             318704430  3.795  ns/op         0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                 176552492  6.784  ns/op         0  B/op   0  allocs/op

Benchmark_ParseFloat64/fiber-12                                 100000000  11.53  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                           100000000  11.57  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                49127977  24.42  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                 100000000  12.06  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                           100000000  11.89  ns/op         0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                46154216  26.05  ns/op         0  B/op   0  allocs/op

Benchmark_ToUpper/empty/fiber-12                                605515742  1.983  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/mut-12                            362727529  3.316  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                              569960358  2.105  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                              91593726  12.97  ns/op         3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/mut-12                         170938820  7.033  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                            59029579  20.41  ns/op         8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                       361742059  3.348  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/mut-12                   169352232  7.058  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                     224532550  5.362  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12             33448389  32.34  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/mut-12         87597902  13.80  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12           11284638  106.4  ns/op        48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                           30178998  40.48  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/mut-12                       59394913  20.50  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                          6882919  176.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                           42718752  27.67  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/mut-12                       58967574  20.85  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                         18536587  65.66  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                           29335220  42.16  ns/op        64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/mut-12                       57954915  24.87  ns/op         0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                          5507632  218.2  ns/op        64  B/op   1  allocs/op

Benchmark_ToLower/empty/fiber-12                                608674111  1.960  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/mut-12                            363946914  3.296  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                              568967895  2.099  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                             364947778  3.319  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/mut-12                         183344209  6.552  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                           263009422  4.577  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                        88837888  13.15  ns/op         3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/mut-12                   183069525  6.566  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                      59491224  20.62  ns/op         8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12             38666282  31.36  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/mut-12         88002616  13.78  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12           16768167  70.94  ns/op        48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                           43839418  27.52  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/mut-12                       59469973  20.44  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                         23733511  50.75  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                           28452537  42.53  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/mut-12                       57793028  21.80  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                          6648596  177.0  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                           28170004  43.01  ns/op        64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/mut-12                       59480410  20.47  ns/op         0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                          5782359  207.3  ns/op        64  B/op   1  allocs/op

Benchmark_AddTrailingSlashString/empty-12                       618103742  1.956  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                  615537577  1.974  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12               99016713  12.34  ns/op         4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12            609862621  1.959  ns/op         0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                77714750  15.29  ns/op        16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12             612306855  1.959  ns/op         0  B/op   0  allocs/op

Benchmark_CalculateTimestamp/fiber-12                           631784137  1.908  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                          36017288  32.88  ns/op         0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                    4366602  273.8  ns/op         8  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                  3885895  313.8  ns/op         8  B/op   2  allocs/op

Benchmark_GolangXMLEncoder-12                                      521270   2192  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                     532245   2095  ns/op      4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                     331651   3467  ns/op      2854  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>
