package request_handler

var counter int64 = 0

func (c *Cluster) get(key string) string {
	node := c.pickNode(key)
	// fmt.Println("picked node for get", node.conn)
	if node == nil || node.conn == nil {
		return ""
	}

	// req := CacheRequest{Action: 1, Key: key}
	var resp GetCacheResponse
	res := node.conn.Go("Node.Get", key, &resp, nil)
	<-res.Done
	counter++
	// fmt.Println(resp)
	return resp.Value
}

// var putChan chan putReq

func (c *Cluster) put(key, value string) {
	counter++
	node := c.pickNode(key)
	if node == nil || node.conn == nil {
		return
	}
	if node.reqChan == nil {
		node.reqChan = make(chan CacheRequest, 1)
	}
	go func(key, value string) {
		node.reqChan <- CacheRequest{2, key, value}
		// fmt.Println("put request sent")
	}(key, value)
}
