package main

import (
	"fmt"
	"time"
	"math/rand"
	"net/rpc"
	"sort"
	"strconv"
	"runtime"
	partition "dbcache/partitioning"
)

type Port uint16

func (p Port) String() string {
	return fmt.Sprintf("%d", uint32(p))
}

type Peer struct {
	Name   string
	Port   Port
	putChan chan putReq
	conn   *rpc.Client
}

type Cluster struct {
	info   map[Peer]bool
	seeder Peer
	partitions map[partition.Partition]map[*Peer]bool
	sortedPartitions []partition.Partition
}

type ShareInfoResponse struct {
	Info       map[Peer]bool
	Partitions map[partition.Partition]map[Peer]bool
}

type ShareCacheRequest struct{}

type CacheRequest struct {
	Action int8
	Key string
	Value string
	AlreadyAware  map[Peer]bool
}

type CacheRequestResponse struct {
	Ok bool
	Key string
	Value string
}

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

func (c *Cluster) getNodes(p partition.Partition) map[*Peer]bool {
	p = c.getNearestPartition(p)
	peers, ok := c.partitions[p]
	if !ok {
		return map[*Peer]bool{}
	}
	return peers
}

func (c *Cluster) pickNode(key string) *Peer {
	part := partition.GetPartition(key)
	for peer, _ := range c.getNodes(part) {
		return peer
	}
	return &c.seeder
}

func (c *Cluster) get(key string) string {
	node := c.pickNode(key)
	fmt.Println("picked node for get", node.Name)

	req := CacheRequest{Action: 1, Key: key}
	var resp CacheRequestResponse
	res := node.conn.Go("Node.Get", req, &resp, nil)
	<-res.Done
	fmt.Println(resp)
	return resp.Value
}

type putReq struct {
	key string
	val string
}

var putChan chan putReq

func (p *Peer) listen(connectionOpened chan<- bool) {
	if p.putChan == nil {
		p.putChan = make(chan putReq, 20)
	}
	var err error
	p.conn, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
	if err != nil {
		// fmt.Println(err, node)
		fmt.Println(err)
		connectionOpened <- false
		return
	}
	connectionOpened <- true
	defer p.conn.Close()
	for {
		select {
		case v := <-p.putChan:
			req := CacheRequest{Action: 2, Key: v.key, Value: v.val}
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			fmt.Println(resp)
		}
	}
}

func (c *Cluster) put(key, value string) {
	node := c.pickNode(key)
	if node.putChan == nil {
		node.putChan = make(chan putReq, 20)
	}
	go func(key, value string) {
		node.putChan <- putReq{key, value}
	}(key, value)
}

func (c *Cluster) getInfo(infoReceived chan<- bool){
	p := c.seeder
	conn, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(10 * time.Second)
	var connectionOpened chan bool
	for {
		select {
		case <-ticker.C:
			fmt.Println("getting info")
			var resp ShareInfoResponse
			req := ShareCacheRequest{}
			conn.Call("Node.ShareInfo", req, &resp)
			c.info = resp.Info
			fmt.Println("got info ", resp.Info)
			c.partitions = make(map[partition.Partition]map[*Peer]bool)
			var pps map[*Peer]bool
			for part, peers := range resp.Partitions {
				pps = make(map[*Peer]bool)
				for pe, _ := range peers {
					mn := &Peer{Name: pe.Name, Port: pe.Port}
					connectionOpened = make(chan bool)
					go mn.listen(connectionOpened)
					if ok := <-connectionOpened; ok {
						pps[mn] = true
					}
					close(connectionOpened)
				}
				c.partitions[part] = pps
			}
			// c.partitions = resp.Partitions
			c.sortPartitions()
			infoReceived <- true
		}
	}
}

func (c *Cluster) dial(p *Peer) (*rpc.Client, error) {
	fmt.Println(p.Name, p.Port)
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    for i := range b {
        b[i] = letterBytes[r1.Intn(len(letterBytes))]
    }
    return string(b)
}

func main() {
	c := new(Cluster)
	c.seeder = Peer{Name: "seeder", Port: Port(7000)}
	seederConnected := make(chan bool)
	go c.seeder.listen(seederConnected)
	if ok := <-seederConnected; !ok {
		panic("seeder out of reach")
	}
	infoReceived := make(chan bool)
	go c.getInfo(infoReceived)
	<-infoReceived
	go func() {
		for {
			<-infoReceived
		}
	}()
	fmt.Println(c.partitions, c.info)
	// go c.put()
    for {
    	fmt.Println("goroutine counter", runtime.NumGoroutine())
    	key := RandStringBytes(10)
    	value := RandStringBytes(50)
    	// c.getInfo()
    	time.Sleep(10 * time.Nanosecond)
    	// fmt.Println("handler info", c.info)
    	c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	v := c.get(fmt.Sprintf("%d", key))
    	// fmt.Println(partition.GetPartition(fmt.Sprintf("%d", key)))
    	fmt.Println("IS HIT?", v == fmt.Sprintf("%d", value))
    }
}