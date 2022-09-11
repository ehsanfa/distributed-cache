package stream_formatter

import (
	"fmt"
	"bytes"
	"strconv"
	"encoding/gob"
	"dbcache/types"
)

type Req struct {
	Action int8
	Key string
	Value string
}

type Resp struct {
	Ok bool
	Key string
	Value string
}

func Encode(req interface{}) []byte {
	b := &bytes.Buffer{}
	encoder := gob.NewEncoder(b)
	encoder.Encode(req)
	s := fmt.Sprintf("%d;%s", len(b.Bytes()), string(b.Bytes()))
	return []byte(s)
}

func Decode(str []byte, handleResp func (msg []byte) *types.Resp, ch chan *types.Resp) {
	ls := ""
	i := 0 
	var b byte
	var l int
	for {
		b = str[i]
		if i >= len(str)-1 {
			break
		}
		if string(b) != ";" {
			ls = ls + string(b)
			i++
			continue
		} else {
			l, _ = strconv.Atoi(ls)
			msg := str[i+1:i + l+1]
			ch <- handleResp(msg)
			i = i + l + 1
			ls = ""
		}
	}
}