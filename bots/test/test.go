package main

import (
	"os"
	"../../easynet"
)

func main() {
	listener := easynet.HostWithAddress(os.Args[0])
	defer listener.Close()
	conn := easynet.Accept(listener)
	defer conn.Close()
	
	conn.Write([]uint8("ok"))
}
