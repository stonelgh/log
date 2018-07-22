package log

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

var port = 15109 // flog(Flexible Log)

func StartUdpListener() {
	go startUdpListener()
}

func startUdpListener() {
	addrStr := ":" + strconv.Itoa(port)
	addr, e := net.ResolveUDPAddr("udp", addrStr)
	if e != nil {
		fmt.Println("Failed to resolve UDP addr", addrStr, e)
		return
	}
	for {
		conn, e := net.ListenUDP("udp", addr)
		if e != nil {
			fmt.Println("Failed to listen UDP addr", addrStr, e)
			return
		}
		data := []byte{}
		if _, e = conn.Read(data); e != nil {
			fmt.Println("Failed to read UDP packet", e)
			conn.Close()
			continue
		}
		fmt.Println("Recv:", string(data))
		idx := -1
		for i, b := range data {
			if b == 0 {
				idx = i
				break
			}
		}
		head := data
		body := []byte{}
		if idx != -1 {
			head = data[:idx]
			body = data[idx+1:]
		}
		resp, e1 := handleRequest(head, body)
		if e1 != nil {
			resp = NewRespFromError(e1)
		}
		if bresp, e := json.Marshal(resp); e != nil {
			fmt.Println("Failed to encode response")
		} else if _, e = conn.Write(bresp); e != nil {
			fmt.Println("Failed to send response")
		}
		conn.Close()
	}
}
