package main

import (
	"flag"

	"github.com/linde12/gowol"
)

var mac string

func init() {
	flag.StringVar(&mac, "mac", "", "Mac Address of device to wakeup")
	flag.StringVar(&mac, "m", "", "Mac Address of device to wakeup")
}

func main() {
	flag.Parse()
	if packet, err := gowol.NewMagicPacket(mac); err == nil {
		packet.Send("255.255.255.255")
	}
}
