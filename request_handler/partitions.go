package request_handler

import (
	// "fmt"
	"sort"
	// "strconv"
	partition "dbcache/partitioning"
)

/**
 * Refactor this whole sorted partition thing
 */ 
func (c *Cluster) sortPartitions() {
	ps := []int{}
	for p, _ := range c.nodes {
		i := p.Key
		ps = append(ps, i)
	}
	sort.Slice(ps, func(i,j int) bool {
		return ps[i] < ps[j]
	})
	c.setSortedPartitions([]partition.Partition{})
	for _, k := range ps {
		// key = fmt.Sprintf("%d", k)
		c.setSortedPartitions(append(c.getSortedPartitions(), partition.CreateParition(k)))
	}
}

func (c *Cluster) getNearestPartition(p partition.Partition) partition.Partition {
	if _, ok := c.nodes[p]; ok {
		return p
	}
	c.sortPartitions()
	for _, sortedPartition := range c.getSortedPartitions() {
		if sortedPartition.Compare(p) == 1 {
			return sortedPartition
		}
	}
	return partition.CreateParition(0)
}

func (c *Cluster) setSortedPartitions(ps []partition.Partition) {
	// c.Lock()
	c.sortedPartitions = ps
	// c.Unlock()
}

func (c *Cluster) getSortedPartitions() []partition.Partition {
	// c.RLock()
	// defer c.RUnlock()
	return c.sortedPartitions
}