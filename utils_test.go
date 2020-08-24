// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"bytes"
	"testing"
)

func Test_Utils_FunctionName(t *testing.T) {
	t.Parallel()
	AssertEqual(t, "github.com/gofiber/utils.Test_Utils_UUID", FunctionName(Test_Utils_UUID))

	AssertEqual(t, "github.com/gofiber/utils.Test_Utils_FunctionName.func1", FunctionName(func() {}))

	var dummyint = 20
	AssertEqual(t, "int", FunctionName(dummyint))
}

func Test_Utils_UUID(t *testing.T) {
	t.Parallel()
	res := UUID()
	AssertEqual(t, 36, len(res))
}

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

func Test_Utils_ToLowerBytes(t *testing.T) {
	t.Parallel()
	res := ToLowerBytes([]byte("/MY/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, bytes.Equal([]byte("/my/name/is/:param/*"), res))
	res = ToLowerBytes([]byte("/MY1/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, bytes.Equal([]byte("/my1/name/is/:param/*"), res))
	res = ToLowerBytes([]byte("/MY2/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, bytes.Equal([]byte("/my2/name/is/:param/*"), res))
	res = ToLowerBytes([]byte("/MY3/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, bytes.Equal([]byte("/my3/name/is/:param/*"), res))
	res = ToLowerBytes([]byte("/MY4/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, bytes.Equal([]byte("/my4/name/is/:param/*"), res))
}

func Test_Utils_EqualsFold(t *testing.T) {
	t.Parallel()
	res := EqualsFold([]byte("/MY/NAME/IS/:PARAM/*"), []byte("/my/name/is/:param/*"))
	AssertEqual(t, true, res)
	res = EqualsFold([]byte("/MY1/NAME/IS/:PARAM/*"), []byte("/MY1/NAME/IS/:PARAM/*"))
	AssertEqual(t, true, res)
	res = EqualsFold([]byte("/my2/name/is/:param/*"), []byte("/my2/name"))
	AssertEqual(t, false, res)
	res = EqualsFold([]byte("/dddddd"), []byte("eeeeee"))
	AssertEqual(t, false, res)
	res = EqualsFold([]byte("/MY3/NAME/IS/:PARAM/*"), []byte("/my3/name/is/:param/*"))
	AssertEqual(t, true, res)
	res = EqualsFold([]byte("/MY4/NAME/IS/:PARAM/*"), []byte("/my4/nAME/IS/:param/*"))
	AssertEqual(t, true, res)
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

	res = Trim(".test", '.')
	AssertEqual(t, "test", res)
}

func Test_Utils_TrimRightBytes(t *testing.T) {
	t.Parallel()
	res := TrimRightBytes([]byte("/test//////"), '/')
	AssertEqual(t, []byte("/test"), res)

	res = TrimRightBytes([]byte("/test"), '/')
	AssertEqual(t, []byte("/test"), res)
}

func Test_Utils_TrimLeftBytes(t *testing.T) {
	t.Parallel()
	res := TrimLeftBytes([]byte("////test/"), '/')
	AssertEqual(t, []byte("test/"), res)

	res = TrimLeftBytes([]byte("test/"), '/')
	AssertEqual(t, []byte("test/"), res)
}
func Test_Utils_TrimBytes(t *testing.T) {
	t.Parallel()
	res := TrimBytes([]byte("   test  "), ' ')
	AssertEqual(t, []byte("test"), res)

	res = TrimBytes([]byte("test"), ' ')
	AssertEqual(t, []byte("test"), res)

	res = TrimBytes([]byte(".test"), '.')
	AssertEqual(t, []byte("test"), res)
}

// func Test_Utils_getArgument(t *testing.T) {
// 	// TODO
// }

// func Test_Utils_parseTokenList(t *testing.T) {
// 	// TODO
// }

func Test_Utils_GetCharPos(t *testing.T) {
	t.Parallel()
	res := GetCharPos("/foo/bar/foobar/test", '/', 3)
	AssertEqual(t, 8, res)
	res = GetCharPos("foo/bar/foobar/test", '/', 3)
	AssertEqual(t, 14, res)
	res = GetCharPos("foo/bar/foobar/test", '/', 1)
	AssertEqual(t, 3, res)
	res = GetCharPos("foo/bar/foobar/test", 'f', 2)
	AssertEqual(t, 8, res)
	res = GetCharPos("foo/bar/foobar/test", 'f', 0)
	AssertEqual(t, 0, res)
}

func Test_Utils_GetTrimmedParam(t *testing.T) {
	t.Parallel()
	res := GetTrimmedParam("*")
	AssertEqual(t, "*", res)
	res = GetTrimmedParam(":param")
	AssertEqual(t, "param", res)
	res = GetTrimmedParam(":param1?")
	AssertEqual(t, "param1", res)
	res = GetTrimmedParam("noParam")
	AssertEqual(t, "noParam", res)
}

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

func Test_Utils_ImmutableString(t *testing.T) {
	t.Parallel()
	res := ImmutableString("Hello, World!")
	AssertEqual(t, "Hello, World!", res)
}

func Test_Utils_GetMIME(t *testing.T) {
	t.Parallel()
	res := GetMIME(".json")
	AssertEqual(t, "application/json", res)

	res = GetMIME(".xml")
	AssertEqual(t, "application/xml", res)

	res = GetMIME("xml")
	AssertEqual(t, "application/xml", res)

	res = GetMIME("unknown")
	AssertEqual(t, MIMEOctetStream, res)
	// empty case
	res = GetMIME("")
	AssertEqual(t, "", res)
}

func Test_Utils_StatusMessage(t *testing.T) {
	t.Parallel()
	res := StatusMessage(204)
	AssertEqual(t, "No Content", res)

	res = StatusMessage(404)
	AssertEqual(t, "Not Found", res)

	res = StatusMessage(426)
	AssertEqual(t, "Upgrade Required", res)

	res = StatusMessage(511)
	AssertEqual(t, "Network Authentication Required", res)

	res = StatusMessage(1337)
	AssertEqual(t, "", res)

	res = StatusMessage(-1)
	AssertEqual(t, "", res)

	res = StatusMessage(0)
	AssertEqual(t, "", res)

	res = StatusMessage(600)
	AssertEqual(t, "", res)
}

func Test_Utils_AssertEqual(t *testing.T) {
	t.Parallel()
	AssertEqual(nil, []string{}, []string{})
	AssertEqual(t, []string{}, []string{})
}
