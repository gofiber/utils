package utils

// XMLMarshal returns v's XML encoding.
type XMLMarshal func(v any) ([]byte, error)

// XMLUnmarshal parses XML data into v, which must be a non-nil pointer.
type XMLUnmarshal func(data []byte, v any) error
