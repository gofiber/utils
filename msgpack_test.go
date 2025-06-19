package utils

import (
	"testing"

	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/require"
)

type sampleStructure2 struct {
	ImportantString string `msgpack:"important_string"`
}

func Test_MsgpackEncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "Hello World",
		}
		msgpackEncoder = msgpack.Marshal
	)

	raw, err := msgpackEncoder(ss)
	require.NoError(t, err)

	var decoded sampleStructure
	err = msgpack.Unmarshal(raw, &decoded)
	require.NoError(t, err)
	require.Equal(t, ss.ImportantString, decoded.ImportantString)
}

func Test_MsgpackDecoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "Hello World",
		}
		msgpackEncoder = msgpack.Marshal
		msgpackDecoder = msgpack.Unmarshal
	)

	raw, err := msgpackEncoder(ss)
	require.NoError(t, err)

	var decoded sampleStructure
	err = msgpackDecoder(raw, &decoded)
	require.NoError(t, err)
	require.Equal(t, "Hello World", decoded.ImportantString)
}
