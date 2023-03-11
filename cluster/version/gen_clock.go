package version

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type GenClock struct {
	number uint64
	mu     sync.RWMutex
}

type marshal struct {
	Number uint64
}

func (v *GenClock) Number() uint64 {
	return v.number
}

func (v *GenClock) Increment() {
	v.mu.Lock()
	v.number++
	v.mu.Unlock()
}

func (v1 *GenClock) ReplaceWith(v2 Version) {
	v1.mu.Lock()
	v1.number = v2.Number()
	v1.mu.Unlock()
}

func CreateGenClockVersion() *GenClock {
	return &GenClock{}
}

func (v *GenClock) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshal{
		Number: v.number,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *GenClock) UnmarshalBinary(data []byte) error {
	m := &marshal{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	v.number = m.Number
	return nil
}
