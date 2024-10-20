package utils

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

func Test_DefaultCBOREncoder(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "Hello World",
		}
		importantString             = `a170696d706f7274616e745f737472696e676b48656c6c6f20576f726c64`
		cborEncoder     CBORMarshal = cbor.Marshal
	)

	raw, err := cborEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(raw), importantString)
}

func Test_DefaultCBORDecoder(t *testing.T) {
	t.Parallel()

	var (
		ss                   sampleStructure
		importantString, err               = hex.DecodeString("a170696d706f7274616e745f737472696e676b48656c6c6f20576f726c64")
		cborDecoder          CBORUnmarshal = cbor.Unmarshal
	)
	if err != nil {
		t.Error("Failed to decode hex string")
	}

	err = cborDecoder(importantString, &ss)
	require.NoError(t, err)
	require.Equal(t, "Hello World", ss.ImportantString)
}

func Test_DefaultCBOREncoderWithEmptyString(t *testing.T) {
	t.Parallel()

	var (
		ss = &sampleStructure{
			ImportantString: "",
		}
		importantString             = `a170696d706f7274616e745f737472696e6760`
		cborEncoder     CBORMarshal = cbor.Marshal
	)

	raw, err := cborEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(raw), importantString)
}

func Test_DefaultCBORDecoderWithEmptyString(t *testing.T) {
	t.Parallel()

	var (
		ss                   sampleStructure
		importantString, err               = hex.DecodeString("a170696d706f7274616e745f737472696e6760")
		cborDecoder          CBORUnmarshal = cbor.Unmarshal
	)
	if err != nil {
		t.Error("Failed to decode hex string")
	}

	err = cborDecoder(importantString, &ss)
	require.NoError(t, err)
	require.Equal(t, "", ss.ImportantString)
}

func Test_DefaultCBOREncoderWithUnitializedStruct(t *testing.T) {
	t.Parallel()

	var (
		ss              sampleStructure
		importantString             = `a170696d706f7274616e745f737472696e6760`
		cborEncoder     CBORMarshal = cbor.Marshal
	)

	raw, err := cborEncoder(ss)
	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(raw), importantString)
}

func Test_DefaultCBORDecoderWithUnitializedStruct(t *testing.T) {
	t.Parallel()

	var (
		ss                   sampleStructure
		emptySs              sampleStructure
		importantString, err               = hex.DecodeString("a170696d706f7274616e745f737472696e6760")
		cborDecoder          CBORUnmarshal = cbor.Unmarshal
	)
	if err != nil {
		t.Error("Failed to decode hex string")
	}

	err = cborDecoder(importantString, &ss)
	require.NoError(t, err)
	require.Equal(t, emptySs, ss)
}
