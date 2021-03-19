package Net

import (
	"net"
)

func ListenAndServe(address string,handler func(*Comm)) error{
	println("serve on:", address)
	addr, _ := net.ResolveTCPAddr("tcp", address)
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil{
		return err
	}
	defer conn.Close()

	var newConn net.Conn
	for {
		if newConn,err = conn.Accept();err != nil{
			return err
		}
		go handler(NewComm(UpgradeTcpStream(newConn)))
	}
}

func ListenAndServeUdp(address string,handler func(*Comm)) error{
	println("serve on:", address)
	addr, _ := net.ResolveTCPAddr("tcp", address)
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil{
		return err
	}
	defer conn.Close()

	var newConn net.Conn
	for {
		if newConn,err = conn.Accept();err != nil{
			return err
		}
		go handler(NewComm(UpgradeTcpStream(newConn)))
	}
}