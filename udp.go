package log

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

const (
	MAX_UDP_PACKET = 1500
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
	conn, e := net.ListenUDP("udp", addr)
	if e != nil {
		fmt.Println("Failed to listen UDP addr", addrStr, e)
		return
	}
	defer conn.Close()

	data := make([]byte, MAX_UDP_PACKET)
	for {
		n, raddr, e := conn.ReadFrom(data)
		if e != nil {
			fmt.Println("Failed to read UDP packet:", e)
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
			body = data[idx+1 : n]
		}
		resp, e1 := handleRequest(head, body)
		if e1 != nil {
			resp = NewRespFromError(e1)
		}
		if bresp, e := json.Marshal(resp); e != nil {
			fmt.Println("Failed to encode response:", e)
		} else if _, e = conn.WriteTo(bresp, raddr); e != nil {
			fmt.Println("Failed to send response:", e)
		}
	}
}
