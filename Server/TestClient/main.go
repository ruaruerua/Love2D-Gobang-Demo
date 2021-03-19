package main

import (
	_ "Server/Config"
	"Server/DB"
	"Server/Net"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func main(){
	var client Client
	//client.host = "-"
	//for client.host[0] != '1'{
	//	log.Println("input id")
	//	fmt.Scanln(&client.id)
	//	client.host = client.Login()
	//}
	client.host = client.Login()
	client.LinkCenterWithUdp()
}


func (client *Client)Login()string{
	body,_ := json.Marshal(DB.Player{ID:"",Password: "123"})
	response,_ := http.Post("http://" + viper.GetString("Login.Addr"),
		"", bytes.NewBuffer(body))
	defer response.Body.Close()

	ret,_ := ioutil.ReadAll(response.Body)
	log.Println(string(ret))
	client.host = strings.Split(string(ret),"|")[0]
	client.player = &Net.Player{
		ID:    "",
		Score: 0,
	}
	client.player.ID = strings.Split(string(ret),"|")[1]
	client.player.Score,_ = strconv.Atoi(strings.Split(string(ret),"|")[2])
	return strings.Split(string(ret),"|")[0]
}

func (client *Client)LinkCenterWithUdp(){
	udpAddr, err := net.ResolveUDPAddr("udp", client.host)
	if err != nil {
		fmt.Println("Resolve udpAddr error", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)

	defer conn.Close()
	if err != nil {
		fmt.Println("connect server error", err)
		return
	}

	ctx,cancel := context.WithCancel(context.Background())
	client.JsonIO = Net.UpgradeTcpStream(conn)
	go client.Send()
	go client.Receive(cancel)
	select {
	case <-ctx.Done():
		return
	}
}

func (client *Client)LinkCenterWithTcp(){
	tcpAddr, err := net.ResolveTCPAddr("tcp", client.host)
	if err != nil {
		fmt.Println("Resolve TCPAddr error", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	defer conn.Close()
	if err != nil {
		fmt.Println("connect server error", err)
		return
	}

	ctx,cancel := context.WithCancel(context.Background())
	client.JsonIO = Net.UpgradeTcpStream(conn)
	go client.Send()
	go client.Receive(cancel)
	select {
	case <-ctx.Done():
		return
	}
}

//func (client *Client)LinkCenter(){
//	u := url.URL{
//		Scheme: "ws",
//		Host:client.host,
//		Path:"/",
//	}
//	ws,_,err := websocket.DefaultDialer.Dial(u.String(),nil)
//
//	if err != nil{
//		log.Println(err)
//		return
//	}
//	ctx,cancel := context.WithCancel(context.Background())
//	client.JsonIO = ws
//	go client.Send()
//	go client.Receive(cancel)
//	select {
//	case <-ctx.Done():
//		return
//	}
//}





//func HeartBeat(ws *websocket.Conn){
//	t := time.Tick(time.Second * 5)
//	ws.WriteJSON(Net.NewMsg(Net.Hello,"wt"))
//	for {
//		select {
//		case <-t:
//
//		}
//	}
//}
