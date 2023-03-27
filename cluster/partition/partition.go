package partition

type Partition interface {
	Name() string
	IsEmpty() bool
}
