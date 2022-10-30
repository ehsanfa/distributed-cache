package request_handler

import (
	"net"
	"fmt"
	"bufio"
	"strings"
)

func (cluster *Cluster) serve() {
	listener, err := net.Listen("tcp", "0.0.0.0:8755")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func(c net.Conn) {
			scanner := bufio.NewScanner(c)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				w := scanner.Text()
				words := strings.Split(w, " ")

				if len(words) == 0 {
					c.Write([]byte("PLEASE SPECIFY AN OPERATION\n"))
					continue
				}

				if len(words) == 1 {
					if words[0] == "exit" {
						break
					}
					c.Write([]byte("INVALID OPERATION\n"))
					continue
				}

				switch strings.ToLower(words[0]) {
				case "get":
					key := words[1]
					c.Write([]byte(cluster.get(key)))
					c.Write([]byte("\n"))
				case "set":
					key := words[1]
					val := words[2]
					cluster.put(key, val)
					c.Write([]byte("OK\n"))
				default:
					c.Write([]byte("INVALID OPERATION\n"))
				}
			}
			c.Write([]byte("DONE!"))
			conn.Close()
		}(conn)
	}

}