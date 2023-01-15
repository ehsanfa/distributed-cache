package request_handler

import (
	"fmt"
	"math/rand"
	"time"
	// "runtime"
)

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

func HandleTest() {
	c := new(Cluster)
	c.seeder = Peer{info: PeerInfo{Name: "seeder", Port: Port(7000)}}
	seederConnected := make(chan bool)
	go c.seeder.seederListen(seederConnected)
	if ok := <-seederConnected; !ok {
		panic("seeder out of reach")
	}
	infoReceived := make(chan bool)
	go c.getInfo(infoReceived)
	<-infoReceived
	go func() {
		for {
			<-infoReceived
			fmt.Println("report count", counter)
			counter = 0
		}
	}()
	// fmt.Println(c.nodes, c.info)

	// go c.put()
	for {
		// fmt.Println("goroutine counter", runtime.NumGoroutine())
		value := "2"
		// c.getInfo()
		time.Sleep(1 * time.Millisecond)
		// fmt.Println("handler info", c.info)
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		c.put(fmt.Sprintf("%d", RandStringBytes(4)), fmt.Sprintf("%d", value))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		go c.get(fmt.Sprintf("%d", RandStringBytes(4)))
		// fmt.Println(partition.GetPartition(fmt.Sprintf("%d", key)))
		// fmt.Println("IS HIT?", v == fmt.Sprintf("%d", value))
	}

	// c.serve()
}
