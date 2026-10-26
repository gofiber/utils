package utils

// CBORMarshal returns v's CBOR encoding.
type CBORMarshal func(v any) ([]byte, error)

// CBORUnmarshal parses CBOR data into v, which must be a non-nil pointer.
type CBORUnmarshal func(data []byte, v any) error
