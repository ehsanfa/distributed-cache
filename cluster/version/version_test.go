package version

import "testing"

func TestVersion(t *testing.T) {
	v1 := CreateGenClockVersion()
	v1.Increment()
	if v1.Number() != 1 {
		t.Error("versions don't match")
	}

	v2 := CreateGenClockVersion()
	v2.Increment()
	v2.Increment()
	v2.Increment()
	v1.ReplaceWith(v2)
	if v1.Number() != 3 {
		t.Error("versions after replacing don't match")
	}
}
