package main

import (
	"fmt"
	"log"
	"net/rpc"
	"dbcache/types"
)

func main() {
	client, err := rpc.Dial("tcp", "server:6399")
	if err != nil {
        log.Fatal("dialing:", err)
    }

    
    var resp types.Resp

    err = client.Call("Req.Put", types.Req{2, "hasan", "heshmat"}, &resp)
    if err != nil {
    	fmt.Println(err)
    }
    err = client.Call("Req.Put", types.Req{2, "hooshang", "mammad"}, &resp)
    if err != nil {
    	fmt.Println(err)
    }

    req := types.Req{1, "hasan", ""}
    err = client.Call("Req.Get", req, &resp)
    if err != nil {
    	fmt.Println(err)
    }
    fmt.Println(resp)
    for {}
}