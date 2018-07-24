package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/stonelgh/log"
	"github.com/stonelgh/log/error"
)

func main() {
	udp := flag.Bool("u", false, "send a message via UDP")
	dest := flag.String("d", "localhost:15109", "destination of the message")
	head := flag.String("head", "", "head of a message to be sent")
	body := flag.String("body", "", "body of a message to be sent")
	flag.Parse()
	if *udp {
		sendUdpPacket(*dest, *head, *body)
		return
	}

	log.StartUdpListener()
	sample()
	select {}
}

func sendUdpPacket(addr string, head string, body string) {
	uaddr, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		fmt.Println("Failed to resolve UDP addr", addr, e)
		return
	}

	conn, e := net.DialUDP("udp", nil, uaddr)
	if e != nil {
		fmt.Println("Failed to dial UDP addr", addr, e)
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	data := []byte(head + string(0) + body)
	if _, e := conn.Write(data); e != nil {
		fmt.Println("Failed to send UDP packet", e)
		return
	}

	n := 0
	data = make([]byte, log.MAX_UDP_PACKET)
	if n, e = conn.Read(data); e != nil {
		fmt.Println("Failed to recv UDP packet", e)
		return
	}

	data = data[:n]
	resp := log.Response{}
	if e = json.Unmarshal(data, &resp); e != nil {
		fmt.Println("Failed to parse response:", e)
		return
	}
	if resp.Head.Code == error.ErrOK {
		fmt.Println("OK:", resp.Body)
	} else {
		fmt.Println("Error:", resp.Head.Code)
		fmt.Println("Message:", resp.Head.Msg)
	}
}
