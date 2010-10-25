package main

import (
	"fmt"
	"os"
	"net"
	"log"
)

func main() {
	fmt.Printf("Coordinator %s\n", os.Args[0])
	addr, err := net.ResolveTCPAddr(os.Args[1]);
	if err != nil { log.Printf("%s: ", os.Args[0]); log.Exit(err) }
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil { log.Printf("%s: ", os.Args[0]); log.Exit(err) }
	conn.Write([]uint8(os.Args[0]))
	conn.Close()
}
