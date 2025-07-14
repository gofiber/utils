package utils

import (
	"testing"

	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/require"
)

type sampleMsgPackStructure struct {
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
		ss = &sampleMsgPackStructure{
			ImportantString: "Hello World",
		}
		msgpackEncoder = msgpack.Marshal
		msgpackDecoder = msgpack.Unmarshal
	)

	raw, err := msgpackEncoder(ss)
	require.NoError(t, err)

	var decoded sampleMsgPackStructure
	err = msgpackDecoder(raw, &decoded)
	require.NoError(t, err)
	require.Equal(t, "Hello World", decoded.ImportantString)
}

func Test_MsgpackDecodeInvalid(t *testing.T) {
	t.Parallel()

	var decoded sampleMsgPackStructure
	err := msgpack.Unmarshal([]byte{0xff, 0xff}, &decoded)
	require.Error(t, err)
}
