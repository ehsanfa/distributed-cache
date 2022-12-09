package request_handler

func Handle() {
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
		}
	}()

 	c.serve()
}