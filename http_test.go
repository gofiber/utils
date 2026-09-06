// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"mime"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetMIME(t *testing.T) {
	t.Parallel()
	res := GetMIME(".json")
	require.Equal(t, "application/json", res)

	res = GetMIME(".xml")
	require.Equal(t, "application/xml", res)

	res = GetMIME("xml")
	require.Equal(t, "application/xml", res)

	res = GetMIME("unknown")
	require.Equal(t, MIMEOctetStream, res)

	res = GetMIME(".zst")
	require.Equal(t, "application/zstd", res)

	res = GetMIME("zst")
	require.Equal(t, "application/zstd", res)

	res = GetMIME("cbor")
	require.Equal(t, "application/cbor", res)

	// upper-case extensions match case-insensitively
	res = GetMIME("PNG")
	require.Equal(t, "image/png", res)

	res = GetMIME(".JSON")
	require.Equal(t, "application/json", res)

	// empty case
	res = GetMIME("")
	require.Empty(t, res)

	err := mime.AddExtensionType(".mjs", "application/javascript")
	if err == nil {
		res = GetMIME(".mjs")
		require.Equal(t, "application/javascript", res)
	}
	require.NoError(t, err)

	require.Equal(t, "application/json", GetMIME(".JSON"))
	require.Equal(t, MIMEOctetStream, GetMIME("  .json"))
}

// go test -v -run=^$ -bench=Benchmark_GetMIME -benchmem -count=2
func Benchmark_GetMIME(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = GetMIME(".xml")
			res = GetMIME(".txt")
			res = GetMIME(".png")
			res = GetMIME(".exe")
			res = GetMIME(".json")
		}
		require.Equal(b, "application/json", res)
	})
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res = mime.TypeByExtension(".xml")
			res = mime.TypeByExtension(".txt")
			res = mime.TypeByExtension(".png")
			res = mime.TypeByExtension(".exe")
			res = mime.TypeByExtension(".json")
		}
		require.Equal(b, "application/json", res)
	})
}

func Test_ParseVendorSpecificContentType(t *testing.T) {
	t.Parallel()

	cType := ParseVendorSpecificContentType("application/json")
	require.Equal(t, "application/json", cType)

	// Test with parameters (semicolon IndexByte optimization)
	cType = ParseVendorSpecificContentType("multipart/form-data; boundary=abc123")
	require.Equal(t, "multipart/form-data; boundary=abc123", cType)

	// Test vendor-specific content types (plus IndexByte optimization)
	cType = ParseVendorSpecificContentType("application/vnd.api+json; version=1")
	require.Equal(t, "application/json", cType)
	cType = ParseVendorSpecificContentType("application/vnd.dummy+x-www-form-urlencoded")
	require.Equal(t, "application/x-www-form-urlencoded", cType)

	// Test invalid cases (slash IndexByte optimization)
	cType = ParseVendorSpecificContentType("something invalid")
	require.Equal(t, "something invalid", cType)

	// Additional edge cases for IndexByte optimization
	cType = ParseVendorSpecificContentType("application/vnd.custom+xml; charset=utf-8")
	require.Equal(t, "application/xml", cType)

	cType = ParseVendorSpecificContentType("text/vnd.example+plain")
	require.Equal(t, "text/plain", cType)

	cType = ParseVendorSpecificContentType("aaaaaaaaaaa/vnd.api+json")
	require.Equal(t, "aaaaaaaaaaa/json", cType)

	cType = ParseVendorSpecificContentType("application/vnd.test+json;boundary=test")
	require.Equal(t, "application/json", cType)

	// Edge cases with multiple special characters
	cType = ParseVendorSpecificContentType("application/vnd.api+json+extra; param=value")
	require.Equal(t, "application/json+extra", cType)

	// Semicolon before plus
	cType = ParseVendorSpecificContentType("application/json; charset=utf-8+extra")
	require.Equal(t, "application/json", cType)

	// Empty and single-character inputs
	require.Empty(t, ParseVendorSpecificContentType(""))
	require.Equal(t, "+", ParseVendorSpecificContentType("+"))
	require.Equal(t, ";", ParseVendorSpecificContentType(";"))
	require.Equal(t, "/", ParseVendorSpecificContentType("/"))

	cType = ParseVendorSpecificContentType("multipart/form-data; boundary=dart-http-boundary-ZnVy.ICWq+7HOdsHqWxCFa8g3D.KAhy+Y0sYJ_lBADypu8po3_X")
	require.Equal(t, "multipart/form-data", cType)

	cType = ParseVendorSpecificContentType("multipart/form-data")
	require.Equal(t, "multipart/form-data", cType)

	cType = ParseVendorSpecificContentType("application/vnd.api+json; version=1")
	require.Equal(t, "application/json", cType)

	cType = ParseVendorSpecificContentType("application/vnd.api+json")
	require.Equal(t, "application/json", cType)

	// Uppercase vendor specific type should only be parsed when case-insensitive
	// matching is enabled
	cType = ParseVendorSpecificContentType("APPLICATION/VND.API+JSON")
	require.Equal(t, "APPLICATION/JSON", cType)

	cType = ParseVendorSpecificContentType("APPLICATION/VND.API+JSON", true)
	require.Equal(t, "application/json", cType)

	cType = ParseVendorSpecificContentType("application/vnd.dummy+x-www-form-urlencoded")
	require.Equal(t, "application/x-www-form-urlencoded", cType)

	cType = ParseVendorSpecificContentType("something invalid")
	require.Equal(t, "something invalid", cType)

	cType = ParseVendorSpecificContentType("invalid+withoutSlash")
	require.Equal(t, "invalid+withoutSlash", cType)
}

func Test_ParseVendorSpecificContentType_IndexByteOptimization(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{"application/vnd.api+json", "application/json", "plus in middle"},
		{"+json", "+json", "plus at start, no slash"},
		{"application/+json", "application/json", "plus after slash"},
		{"application/json;charset=utf-8", "application/json;charset=utf-8", "semicolon after content type"},
		{";charset=utf-8", ";charset=utf-8", "semicolon at start"},
		{"application/vnd.api+json;version=1", "application/json", "plus before semicolon"},
		{"application/json", "application/json", "normal content type with slash"},
		{"applicationjson", "applicationjson", "no slash in content type"},
		{"app/vnd.test+data/extra", "app/data/extra", "multiple slashes"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.expected, ParseVendorSpecificContentType(tc.input), "input: %s", tc.input)
		})
	}
}

func Benchmark_ParseVendorSpecificContentType(b *testing.B) {
	b.Run("vendorContentType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ParseVendorSpecificContentType("application/vnd.api+json; version=1")
		}
	})

	b.Run("defaultContentType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ParseVendorSpecificContentType("application/json")
		}
	})
}

func Test_StatusMessage(t *testing.T) {
	t.Parallel()
	res := StatusMessage(204)
	require.Equal(t, "No Content", res)

	res = StatusMessage(404)
	require.Equal(t, "Not Found", res)

	res = StatusMessage(426)
	require.Equal(t, "Upgrade Required", res)

	res = StatusMessage(511)
	require.Equal(t, "Network Authentication Required", res)

	res = StatusMessage(1337)
	require.Empty(t, res)

	res = StatusMessage(-1)
	require.Empty(t, res)

	res = StatusMessage(0)
	require.Empty(t, res)

	res = StatusMessage(600)
	require.Empty(t, res)

	res = StatusMessage(100)
	require.Equal(t, "Continue", res)

	res = StatusMessage(http.StatusTeapot)
	require.Equal(t, "I'm a teapot", res)
}

// go test -run=^$ -bench=Benchmark_StatusMessage -benchmem -count=2
func Benchmark_StatusMessage(b *testing.B) {
	b.Run("fiber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StatusMessage(http.StatusNotExtended)
		}
	})
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			http.StatusText(http.StatusNotExtended)
		}
	})
}

// Test_GetMIME_TableCoverage checks every mimeExtensions entry is reachable through the packed-key table.
func Test_GetMIME_TableCoverage(t *testing.T) {
	t.Parallel()
	for ext, want := range mimeExtensions {
		require.NotEmpty(t, ext)
		require.LessOrEqual(t, len(ext), mimeKeyMaxLen, "extension %q does not fit a packed key", ext)
		require.Equal(t, want, GetMIME(ext), "extension %q", ext)
		require.Equal(t, want, GetMIME("."+ext), "extension .%q", ext)
		require.Equal(t, want, GetMIME(strings.ToUpper(ext)), "extension %q upper-cased", ext)
	}
	// Distinct extensions must never share a packed key.
	seen := make(map[uint64]string, len(mimeExtensions))
	for ext := range mimeExtensions {
		key := packExtension(ext)
		prev, dup := seen[key]
		require.False(t, dup, "extensions %q and %q pack to the same key", prev, ext)
		seen[key] = ext
	}
	// Probe sequences must terminate: the table keeps free slots.
	free := 0
	for i := range mimeTable {
		if mimeTable[i].mimeType == "" {
			free++
		}
	}
	require.Positive(t, free)

	// Misses of every length; NUL and non-ASCII bytes are in no mime database, so these are octet-stream everywhere.
	for _, ext := range []string{"\x00", "ht\x00l", "j\xf3on", "1234567\x00", "toolongext\x01", ".\x00\x00\x00\x00\x00\x00\x00\x00"} {
		require.Equal(t, MIMEOctetStream, GetMIME(ext), "extension %q", ext)
	}
}
