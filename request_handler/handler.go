package request_handler

import (
	"fmt"
	"time"
	"math/rand"
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

func Handle() {
	c := new(Cluster)
	c.seeder = Peer{Name: "seeder", Port: Port(7000)}
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
		}
	}()
	// fmt.Println(c.partitions, c.info)

	// go c.put()
    for {
    	// fmt.Println("goroutine counter", runtime.NumGoroutine())
    	key := "1"
    	value := "2"
    	// c.getInfo()
    	time.Sleep(10 * time.Nanosecond)
    	// fmt.Println("handler info", c.info)
    	c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	// c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	// c.put(fmt.Sprintf("%d", key), fmt.Sprintf("%d", value))
    	c.get(fmt.Sprintf("%d", key))
    	c.get(fmt.Sprintf("%d", key))
    	c.get(fmt.Sprintf("%d", key))
    	c.get(fmt.Sprintf("%d", key))
    	c.get(fmt.Sprintf("%d", key))
    	// fmt.Println(partition.GetPartition(fmt.Sprintf("%d", key)))
    	// fmt.Println("IS HIT?", v == fmt.Sprintf("%d", value))
    }
}