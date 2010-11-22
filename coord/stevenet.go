package main

import (
	"fmt"
	"os"
	"net"
	"log"
	"strings"
	"time"
)

func dieIfError(err os.Error, msg string) {
	if err != nil { log.Exit("", msg, " in coordinator: ", err) }
}

func setupConnection(addrString string) *net.TCPConn {
	fmt.Printf("Listening with address %s\n", addrString)
	addr, err := net.ResolveTCPAddr(addrString);
	dieIfError(err, "TCP address resolution error")
	listener, err := net.ListenTCP("tcp", addr);
	dieIfError(err, "Listening error")
	
	conn, err := listener.AcceptTCP();
	dieIfError(err, "TCP accept error")
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	
	return conn
}

func receive_from(conn *net.TCPConn) []byte {
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	for err != nil && strings.HasSuffix(err.String(), "temporarily unavailable") {
		time.Sleep(10000)
		size, err = conn.Read(rcvd)
	}
	dieIfError(err, "Receive error")
	return rcvd[0:size]
}
