// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"mime"
	"net/http"
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

	// empty case
	res = GetMIME("")
	require.Equal(t, "", res)

	err := mime.AddExtensionType(".mjs", "application/javascript")
	if err == nil {
		res = GetMIME(".mjs")
		require.Equal(t, "application/javascript", res)
	}
	require.NoError(t, err)
}

// go test -v -run=^$ -bench=Benchmark_GetMIME -benchmem -count=2
func Benchmark_GetMIME(b *testing.B) {
	var res string
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = GetMIME(".xml")
			res = GetMIME(".txt")
			res = GetMIME(".png")
			res = GetMIME(".exe")
			res = GetMIME(".json")
		}
		require.Equal(b, "application/json", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
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

	cType = ParseVendorSpecificContentType("multipart/form-data; boundary=dart-http-boundary-ZnVy.ICWq+7HOdsHqWxCFa8g3D.KAhy+Y0sYJ_lBADypu8po3_X")
	require.Equal(t, "multipart/form-data", cType)

	cType = ParseVendorSpecificContentType("multipart/form-data")
	require.Equal(t, "multipart/form-data", cType)

	cType = ParseVendorSpecificContentType("application/vnd.api+json; version=1")
	require.Equal(t, "application/json", cType)

	cType = ParseVendorSpecificContentType("application/vnd.api+json")
	require.Equal(t, "application/json", cType)

	cType = ParseVendorSpecificContentType("application/vnd.dummy+x-www-form-urlencoded")
	require.Equal(t, "application/x-www-form-urlencoded", cType)

	cType = ParseVendorSpecificContentType("something invalid")
	require.Equal(t, "something invalid", cType)

	cType = ParseVendorSpecificContentType("invalid+withoutSlash")
	require.Equal(t, "invalid+withoutSlash", cType)
}

func Benchmark_ParseVendorSpecificContentType(b *testing.B) {
	b.Run("vendorContentType", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			ParseVendorSpecificContentType("application/vnd.api+json; version=1")
		}
	})

	b.Run("defaultContentType", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
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
	require.Equal(t, "", res)

	res = StatusMessage(-1)
	require.Equal(t, "", res)

	res = StatusMessage(0)
	require.Equal(t, "", res)

	res = StatusMessage(600)
	require.Equal(t, "", res)
}

// go test -run=^$ -bench=Benchmark_StatusMessage -benchmem -count=2
func Benchmark_StatusMessage(b *testing.B) {
	b.Run("fiber", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			StatusMessage(http.StatusNotExtended)
		}
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			http.StatusText(http.StatusNotExtended)
		}
	})
}
