package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	v1 := CreateGenClockVersion(0)
	v1 = v1.Increment().(GenClock)
	if v1.Number() != 1 {
		t.Error("versions don't match")
	}

	v2 := CreateGenClockVersion(0)
	v2 = v2.Increment().(GenClock)
	v2 = v2.Increment().(GenClock)
	v2 = v2.Increment().(GenClock)
	if v2.Number() != 3 {
		t.Error("versions after replacing don't match")
	}
}

// func TestMarshal(t *testing.T) {
// 	p := CreateGenClockVersion(0)
// 	p.Increment()
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	if err := enc.Encode(p); err != nil {
// 		t.Error(err)
// 	}
// 	l := GenClock{number: 1}
// 	reader := bytes.NewReader(buf.Bytes())
// 	dec := gob.NewDecoder(reader)
// 	if err := dec.Decode(&l); err != nil {
// 		t.Error(err)
// 	}

// 	if p.Number() != l.Number() {
// 		t.Error("error!!!")
// 	}
// }
