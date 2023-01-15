package version

import "sync"

type GenClock struct {
	number uint64
	mu     sync.RWMutex
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
