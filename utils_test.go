// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"testing"
)

func Test_Utils_ToUpper(t *testing.T) {
	t.Parallel()
	res := ToUpper("/my/name/is/:param/*")
	AssertEqual(t, "/MY/NAME/IS/:PARAM/*", res)
}

func Test_Utils_ToLower(t *testing.T) {
	t.Parallel()
	res := ToLower("/MY/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my/name/is/:param/*", res)
	res = ToLower("/MY1/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my1/name/is/:param/*", res)
	res = ToLower("/MY2/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my2/name/is/:param/*", res)
	res = ToLower("/MY3/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my3/name/is/:param/*", res)
	res = ToLower("/MY4/NAME/IS/:PARAM/*")
	AssertEqual(t, "/my4/name/is/:param/*", res)
}

func Test_Utils_TrimRight(t *testing.T) {
	t.Parallel()
	res := TrimRight("/test//////", '/')
	AssertEqual(t, "/test", res)

	res = TrimRight("/test", '/')
	AssertEqual(t, "/test", res)
}

func Test_Utils_TrimLeft(t *testing.T) {
	t.Parallel()
	res := TrimLeft("////test/", '/')
	AssertEqual(t, "test/", res)

	res = TrimLeft("test/", '/')
	AssertEqual(t, "test/", res)
}
func Test_Utils_Trim(t *testing.T) {
	t.Parallel()
	res := Trim("   test  ", ' ')
	AssertEqual(t, "test", res)

	res = Trim("test", ' ')
	AssertEqual(t, "test", res)
}

// func Test_Utils_AssertEqual(t *testing.T) {
// 	// TODO
// }

// func Test_Utils_setETag(t *testing.T) {
// 	// TODO
// }

func Test_Utils_GetMIME(t *testing.T) {
	t.Parallel()
	res := GetMIME(".json")
	AssertEqual(t, "application/json", res)

	res = GetMIME(".xml")
	AssertEqual(t, "application/xml", res)

	res = GetMIME("xml")
	AssertEqual(t, "application/xml", res)

	res = GetMIME("json")
	AssertEqual(t, "application/json", res)
}

// func Test_Utils_getArgument(t *testing.T) {
// 	// TODO
// }

// func Test_Utils_parseTokenList(t *testing.T) {
// 	// TODO
// }

func Test_Utils_GetString(t *testing.T) {
	t.Parallel()
	res := GetString([]byte("Hello, World!"))
	AssertEqual(t, "Hello, World!", res)
}

func Test_Utils_GetBytes(t *testing.T) {
	t.Parallel()
	res := GetBytes("Hello, World!")
	AssertEqual(t, []byte("Hello, World!"), res)
}

func Test_Utils_extensionMIME(t *testing.T) {
	t.Parallel()
	res := GetMIME(".html")
	AssertEqual(t, "text/html", res)

	res = GetMIME("html")
	AssertEqual(t, "text/html", res)

	res = GetMIME(".msp")
	AssertEqual(t, "application/octet-stream", res)

	res = GetMIME("msp")
	AssertEqual(t, "application/octet-stream", res)
}
