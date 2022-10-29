package request_handler

func (c *Cluster) get(key string) string {
	node := c.pickNode(key)
	// fmt.Println("picked node for get", node.Name)

	// req := CacheRequest{Action: 1, Key: key}
	var resp CacheRequestResponse
	res := node.conn.Go("Node.Get", key, &resp, nil)
	<-res.Done
	// fmt.Println(resp)
	return resp.Value
}

// var putChan chan putReq

func (c *Cluster) put(key, value string) {
	node := c.pickNode(key)
	if node.reqChan == nil {
		node.reqChan = make(chan CacheRequest, 1)
	}
	go func(key, value string) {
		node.reqChan <- CacheRequest{2, key, value}
		// fmt.Println("put request sent")
	}(key, value)
}