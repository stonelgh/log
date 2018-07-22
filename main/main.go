package main

import (
	"github.com/stonelgh/log"
)

func main() {
	log.StartUdpListener()
	sample()
	select {}
}
