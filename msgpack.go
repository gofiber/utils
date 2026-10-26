package utils

// MsgPackMarshal returns v's MsgPack encoding.
type MsgPackMarshal func(v any) ([]byte, error)

// MsgPackUnmarshal parses MsgPack data into v, which must be a non-nil pointer.
type MsgPackUnmarshal func(data []byte, v any) error
