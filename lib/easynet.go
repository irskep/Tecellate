package easynet

import (
	// "fmt"
	"json"
	"os"
	"net"
	"log"
	"strings"
	"time"
)

func DieIfError(err os.Error, msg string) {
	if err != nil { log.Exit(msg, ": ", err) }
}

func HostWithAddress(addrString string) *net.TCPListener {
	// fmt.Printf("Listening with address %s\n", addrString)
	addr, err := net.ResolveTCPAddr(addrString);
	DieIfError(err, "TCP address resolution error")
	listener, err := net.ListenTCP("tcp", addr);
	DieIfError(err, "Listening error")
	return listener
}

func Accept(listener *net.TCPListener) *net.TCPConn {
	conn, err := listener.AcceptTCP();
	DieIfError(err, "TCP accept error")
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	conn.SetNoDelay(true)
	return conn
}

func Dial(addrString string) *net.TCPConn {
	addr, err := net.ResolveTCPAddr(addrString);
	DieIfError(err, "TCP address resolution error")
	conn, err := net.DialTCP("tcp", nil, addr)
	for err != nil && strings.HasSuffix(err.String(), "connection refused") {
		time.Sleep(10000)
		conn, err = net.DialTCP("tcp", nil, addr)
	}
	DieIfError(err, "Dial error")
	conn.SetNoDelay(true)
	return conn
}

func ReceiveFromWithError(conn *net.TCPConn) ([]byte, os.Error) {
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	for err != nil && strings.HasSuffix(err.String(), "temporarily unavailable") {
		time.Sleep(10000)
		size, err = conn.Read(rcvd)
	}
	return rcvd[0:size], err
}

func ReceiveFrom(conn *net.TCPConn) []byte {
	msg, err := ReceiveFromWithError(conn)
	// if err != nil {
	// 	panic("receive error!")
	// }
	DieIfError(err, "Receive error")
	return msg
}

func TieConnToChannel(conn *net.TCPConn, c chan []uint8) {
	go func() {
		for {
			rcvd := make([]byte, 4096)
			size, err := conn.Read(rcvd)
			for err != nil && strings.HasSuffix(err.String(), "temporarily unavailable") {
				time.Sleep(10000)
				size, err = conn.Read(rcvd)
			}
			if err != nil {
				return
			} else {
				// fmt.Println(string(rcvd[0:size]))
				c <- rcvd[0:size]
			}
		}
	}()
}

func SendJson(conn *net.TCPConn, obj interface{}) {
	data, err := json.Marshal(obj)
	DieIfError(err, "JSON marshal error")
	conn.Write(data)
}

func ReceiveJson(conn *net.TCPConn, obj interface{}) {
	err := json.Unmarshal(ReceiveFrom(conn), obj)
	DieIfError(err, "JSON unmarshal error")
}
