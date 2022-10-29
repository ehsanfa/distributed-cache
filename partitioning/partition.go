package partitioning

import (
	"fmt"
	"math"
	"strconv"
)

const distributionFactor = 13

type Partition struct {
	Key string
}

type Partitioner interface {
	GetName() string
}

func getDegree(bytesTotal int) int{
	return (bytesTotal * distributionFactor) % 360
}

func (p *Partition) getPartitionKey(key string) string{
	num := 0
	for _, b := range []byte(key) {
		num += int(b)
	}
	return fmt.Sprintf("%d", getDegree(num))
}

func (p *Partition) Compare(p1 Partition) int8 {
	key1, _ := strconv.Atoi(p.Key)
	key2, _ := strconv.Atoi(p1.Key)
	if key1 > key2 {
		return 1
	}
	if key1 == key2 {
		return 0
	}
	return -1
}

func GetPartition(key string) Partition {
	p := Partition{}
	k := p.getPartitionKey(key)
	p.Key = k
	return p
}

func CreateParition(key string) Partition {
	return Partition{Key: key}
}

func Initialize(paritionsNum float64) []Partition{
	var partitions []Partition
	var i float64
	deg := float64(360)
	jump := math.Round(deg/paritionsNum*100) / 100
	var key string
	for i = 0; i < deg; i += jump {
		key = fmt.Sprintf("%g", math.Round(i))
		partitions = append(partitions, CreateParition(key))
	}
	return partitions
}
