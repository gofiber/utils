package unsafeconv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	myString string
	myBytes  []byte
)

func Test_UnsafeString(t *testing.T) {
	t.Parallel()
	b := []byte("hello world")
	s := UnsafeString(b)
	require.Equal(t, "hello world", s)
	require.Empty(t, UnsafeString(nil))
	require.Empty(t, UnsafeString([]byte{}))
}

func Test_UnsafeBytes(t *testing.T) {
	t.Parallel()
	require.Equal(t, []byte("hello world"), UnsafeBytes("hello world"))
	require.Empty(t, UnsafeBytes(""))
}

func Test_Bytes(t *testing.T) {
	t.Parallel()
	require.Equal(t, []byte("from string"), Bytes("from string"))
	require.Equal(t, []byte("from bytes"), Bytes([]byte("from bytes")))
	require.Equal(t, []byte("named"), Bytes(myString("named")))
	require.Equal(t, []byte("named"), Bytes(myBytes("named")))
	require.Empty(t, Bytes(""))
	require.Empty(t, Bytes([]byte(nil)))

	// The []byte view of a []byte input aliases the original storage.
	src := []byte("alias")
	view := Bytes(src)
	require.Same(t, &src[0], &view[0])
}

func Test_Seq(t *testing.T) {
	t.Parallel()
	b := []byte("round trip")
	require.Equal(t, "round trip", Seq[string](b))
	require.Equal(t, []byte("round trip"), Seq[[]byte](b))
	require.Equal(t, myString("round trip"), Seq[myString](b))
	require.Empty(t, Seq[string](nil))
	require.Empty(t, Seq[[]byte](nil))

	// Both instantiations are views over the same storage, not copies.
	viewB := Seq[[]byte](b)
	require.Same(t, &b[0], &viewB[0])

	// Seq is Bytes' inverse for every instantiation.
	require.Equal(t, "inverse", Seq[string](Bytes("inverse")))
	require.Equal(t, []byte("inverse"), Seq[[]byte](Bytes([]byte("inverse"))))
}
