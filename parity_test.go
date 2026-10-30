package utils

// Parity index: helper, stdlib oracle, tests.
//
//	{Append,Format,Parse}{Int,Uint}*, ParseNative*   strconv                                parity_ints_test.go, fuzz_test.go
//	ParseFloat32/64, AppendDuration                  strconv, Duration.String               parity_floats_test.go, duration_test.go
//	Append/Format/ParseHTTPDate                      Time.Format, http.ParseTime            parity_net_test.go, edgecase_test.go
//	Append{Query,Path}{Escape,Unescape}              net/url                                edgecase_test.go, parity_net_test.go, fuzz_test.go
//	AppendJSONString, CanonicalHeaderKey             json.Marshal, http.CanonicalHeaderKey  edgecase_test.go
//	SplitHostPort, IsIPv4/6, ParseIPv4/6             net, netip                             parity_net_test.go
//	IncrementIPRange, StatusMessage, ReadFile        netip, http.StatusText, os.ReadFile    parity_net_test.go
//	Unsafe{String,Bytes}, Copy{String,Bytes}         string(b), []byte(s)                   parity_net_test.go
//	GetMIME, Walk, ToString                          mime, filepath.Walk, fmt.Sprint        fuzz_test.go, parity_walk_test.go, parity_convert_test.go
//	Trim*, *CutByte, SplitTrimSeq, To{Lower,Upper}   strings, bytes                         parity_strings_test.go
//	EqualFold, IndexAny2/3, IndexControl*            strings, bytes, unicode                parity_strings_test.go, portable_test.go
//	IndexFold, ContainsFold, Has{Prefix,Suffix}Fold  strings.Index (folded)                 parity_misc_test.go, search_test.go
//	*SecureToken, UUID, FunctionName                 base64, google/uuid, runtime           parity_misc_test.go, common_test.go
//	simd.Memmem, Memchr*, IsASCII                    bytes.Index/IndexByte                  simd/simd_test.go, simd/fuzz_test.go
