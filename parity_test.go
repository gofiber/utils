package utils

// Helpers that stand in for a standard library function are held to it by
// differential tests: the same inputs go to both and the answers, errors
// included, must agree. Where a helper departs from the stdlib on purpose the
// test pins the departure as well, with a check that the stdlib still behaves
// the way the departure assumes, so neither side can drift unnoticed. This is
// the map from helper to test; fuzz targets are named FuzzXxx, and every one of
// them runs its seed corpus as an ordinary test.
//
//	Helper                                   Stdlib                          Held to it by
//	AppendInt, AppendUint, FormatInt*,       strconv.Append*, strconv.Format* parity_ints_test.go (every 8-bit and 16-bit
//	FormatUint*                                                              value), format_test.go (Test_Format_Sweep)
//	ParseInt*, ParseUint*, ParseNative*      strconv.ParseInt, ParseUint,    parity_ints_test.go (every width, error
//	                                         strconv.Atoi                    class and Num), fuzz_test.go (FuzzParse)
//	ParseFloat32, ParseFloat64               strconv.ParseFloat              parsefloat_test.go, parity_floats_test.go
//	AppendDuration                           time.Duration.String            duration_test.go
//	AppendHTTPDate, FormatHTTPDate           time.Time.Format(http.TimeFormat) httpdate_test.go, parity_net_test.go
//	ParseHTTPDate                            http.ParseTime                  parity_net_test.go (values, location and
//	                                                                         *time.ParseError), edgecase_test.go
//	AppendQueryEscape, AppendPathEscape,     url.QueryEscape, url.PathEscape, url_test.go, edgecase_test.go,
//	AppendQueryUnescape, AppendPathUnescape  url.QueryUnescape, url.PathUnescape fuzz_test.go (FuzzURLEscape)
//	AppendJSONString                         json.Marshal of a string        jsonstring_test.go, edgecase_test.go
//	CanonicalHeaderKey                       http.CanonicalHeaderKey         headerkey_test.go, edgecase_test.go
//	SplitHostPort                            net.SplitHostPort               hostport_test.go, parity_net_test.go
//	IsIPv4, IsIPv6                           net.ParseIP                     ips_test.go, parity_net_test.go
//	ParseIPv4, ParseIPv6                     netip.ParseAddr                 ipparse_test.go, parity_net_test.go
//	IncrementIPRange                         netip.Addr.Next, big-int sum    parity_net_test.go
//	StatusMessage                            http.StatusText                 parity_net_test.go
//	GetMIME                                  mime.TypeByExtension            http_test.go, fuzz_test.go (FuzzGetMIME)
//	ReadFile                                 os.ReadFile                     parity_net_test.go
//	Walk                                     filepath.Walk                   parity_walk_test.go
//	ToString                                 fmt.Sprint                      parity_convert_test.go, convert_test.go
//	UnsafeString, UnsafeBytes, CopyString,   string(b), []byte(s)            parity_net_test.go, convert_test.go
//	CopyBytes
//	EqualFold                                strings.EqualFold               parity_strings_test.go, parity_misc_test.go
//	IndexFold, ContainsFold, HasPrefixFold,  strings.Index, HasPrefix,       parity_misc_test.go, search_test.go
//	HasSuffixFold                            HasSuffix on lower-cased copies
//	IndexAny2, IndexAny3                     bytes.IndexByte, strings.IndexAny parity_strings_test.go, portable_test.go
//	Trim, TrimLeft, TrimRight, TrimSpace     strings.Trim*, bytes.Trim*,     parity_strings_test.go, byteseq_test.go
//	                                         strings.TrimSpace
//	CutByte, LastCutByte                     strings.Cut, bytes.Cut          parity_strings_test.go, cut_test.go
//	SplitTrimSeq                             strings.SplitSeq, TrimSpace     parity_strings_test.go, seq_test.go
//	ToLower, ToUpper, Unsafe* variants       strings.ToLower, ToUpper,       parity_strings_test.go, casefold_test.go,
//	(here and in packages bytes, strings)    bytes.ToLower, ToUpper          bytes/case_test.go, strings/case_test.go
//	IndexControl, IndexControlExceptTab      unicode.IsControl               parity_strings_test.go
//	GenerateSecureToken, SecureToken         base64.RawURLEncoding           parity_misc_test.go
//	UUID, UUIDv4                             google/uuid                     common_test.go
//	FunctionName                             runtime.FuncForPC               common_test.go
//	simd.Memmem                              bytes.Index                     simd/simd_test.go, simd/fuzz_test.go (FuzzMemmem)
//	simd.Memchr2, Memchr3, IsASCII, ...      bytes.IndexByte loops           simd/simd_test.go, simd/fuzz_test.go
//
// The portable (non-AVX2) code paths of IsASCII, IndexAny2 and IndexAny3 are
// exercised on AVX2 machines too, by Test_PortableTier in portable_test.go.
