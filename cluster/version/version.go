package version

import "encoding"

type Version interface {
	Number() uint64
	Increment()
	ReplaceWith(v2 Version)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// func (v1 *Version) replaceAndIncrement(v2 Version) {
// 	v1.replace(v2)
// 	v1.increment()
// }

// func (v1 *Version) compare(v2 Version) int8 {
// 	if v1.Number > v2.Number {
// 		return 1
// 	}
// 	if v1.Number < v2.Number {
// 		return -1
// 	}
// 	return 0
// }

// func (v *Version) touch() {
// 	v.increment()
// }
