package utils

// JSONMarshal returns v's JSON encoding.
type JSONMarshal func(v any) ([]byte, error)

// JSONUnmarshal parses JSON data into v, which must be a non-nil pointer.
type JSONUnmarshal func(data []byte, v any) error
