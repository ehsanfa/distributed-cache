package request_handler

import (
	"fmt"
	"sort"
	"strconv"
	partition "dbcache/partitioning"
)

/**
 * Refactor this whole sorted partition thing
 */ 
func (c *Cluster) sortPartitions() {
	ps := []int{}
	for p, _ := range c.partitions {
		i, _ := strconv.Atoi(p.Key)
		ps = append(ps, i)
	}
	sort.Slice(ps, func(i,j int) bool {
		return ps[i] < ps[j]
	})
	c.sortedPartitions = []partition.Partition{}
	var key string
	for _, k := range ps {
		key = fmt.Sprintf("%d", k)
		c.sortedPartitions = append(c.sortedPartitions, partition.CreateParition(key))
	}
}

func (c *Cluster) getNearestPartition(p partition.Partition) partition.Partition {
	if _, ok := c.partitions[p]; ok {
		return p
	}
	c.sortPartitions()
	for _, sortedPartition := range c.sortedPartitions {
		if sortedPartition.Compare(p) == 1 {
			return sortedPartition
		}
	}
	return partition.CreateParition("0")
}