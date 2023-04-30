package version

type GenClock struct {
	number uint64
}

func (v GenClock) Number() uint64 {
	return v.number
}

func (v GenClock) Increment() Version {
	return CreateGenClockVersion(v.number + 1)
}

func CreateGenClockVersion(n uint64) GenClock {
	return GenClock{n}
}
