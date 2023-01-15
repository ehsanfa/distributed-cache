package partition

type Partition interface {
	Name() string
	IsEmpty() bool
}

type SimplePartition struct {
	name string
}

func (p SimplePartition) Name() string {
	return p.name
}

func (p SimplePartition) IsEmpty() bool {
	return p == SimplePartition{}
}

func CreateSimplePartition(name string) SimplePartition {
	return SimplePartition{name: name}
}
