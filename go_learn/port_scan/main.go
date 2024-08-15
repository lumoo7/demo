package main

import (
	"fmt"
	"net"
	"time"
)

func isOpen(s string, port int, t time.Duration) {
	address := fmt.Sprintf("%s:%d", s, port)
	_, err := net.DialTimeout("tcp", address, t)
	if err != nil {
		fmt.Printf("\naddress %s timeout", address)
		return
	}
	fmt.Printf("\naddress %s is open", address)
	return
}

func main() {
	for i := 70; i < 81; i++ {
		go isOpen("192.168.33.66", i, time.Microsecond*200)
	}
}
