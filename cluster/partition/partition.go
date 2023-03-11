package partition

import (
	"encoding"
)

type Partition interface {
	Name() string
	IsEmpty() bool
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
