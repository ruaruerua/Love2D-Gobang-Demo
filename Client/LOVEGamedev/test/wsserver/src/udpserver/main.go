package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"strings"
)

const (
	MsgType_Connect    = "connect"
	MsgType_Connected  = "connected"
	MsgType_Disconnect = "disconnect"
	MsgType_Timeout    = "timeout"
	MsgType_Msg        = "msg"
)
const SERVER_RECV_LEN = 256

func main() {
	address := "localhost:12345"
	println("listen udp on:", address)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, SERVER_RECV_LEN)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ss := bytes.Trim(data, "\x00")
		msg := string(ss)
		fmt.Println("Received:", rAddr, msg)
		if msg == MsgType_Connect {
			fmt.Println("new connect:",rAddr)
			_, err = conn.WriteToUDP([]byte(MsgType_Connected), rAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			upper := strings.ToUpper(msg)
			_, err = conn.WriteToUDP([]byte(upper), rAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Send:", upper)
		}

	}
}
