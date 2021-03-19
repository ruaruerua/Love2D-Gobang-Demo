package Net

import (
	"Server/Const"
	"context"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

func ListenAndServeForUDP(address string,handler func(*Comm)) error{
	log.Println("serve on:", address)

	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", addr)
	defer conn.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	go HeartBeat(time.Second * 3,context.Background())
	for {
		data := make([]byte, 1024)
		len, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			continue
		}
		var msg Message
		//str :=string(data[0:len])
		//log.Println(str)
		if err = json.Unmarshal(data[0:len],&msg); err == nil{
			if msg.CMD == Const.Hello{
				addrWithCh := &AddrWithCh{
					make(chan []byte,1),
					true,
				}
				addrWithCh.ch <- data[0:len]
				u := UdpJson{
					Receiver: addrWithCh.ch,
					addr:     rAddr,
					conn:     conn,
				}
				udpMap.Store(rAddr.String(),addrWithCh)
				go handler(NewComm(&u))
			} else if msg.CMD == Const.Heart{
				if v,_ :=udpMap.Load(rAddr.String());v != nil{
					v.(*AddrWithCh).alive = true
					//udpMap.Store(rAddr.String(),v)
				}
			}else if v,_:=udpMap.Load(rAddr.String()) ;v != nil{
				v.(*AddrWithCh).ch <- data[0:len]
			}else{
				conn.WriteToUDP([]byte("error"),rAddr)
			}
		} else {
			log.Println(err)
		}
	}
}


func HeartBeat(d time.Duration,ctx context.Context){
	t := time.Tick(d)
	for {
		select {
		case <-t:
			udpMap.Range(func(k,v interface{})bool{
				if !v.(*AddrWithCh).alive{
					data,_:=json.Marshal(NewMsg(Const.Heart,[]byte("del")))
					v.(*AddrWithCh).ch <- data
					udpMap.Delete(k)
				} else {
					v.(*AddrWithCh).alive = false
					data,_:=json.Marshal(NewMsg(Const.Heart,[]byte("test")))
					v.(*AddrWithCh).ch <- data
				}
				return true
			})
		case <-ctx.Done():
			break
		}
	}
	log.Println("heart exit")
}

type AddrWithCh struct {
	ch chan []byte
	alive bool
}

//var udpMap = make(map[string]chan []byte)

var udpMap sync.Map