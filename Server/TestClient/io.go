package main

import (
	. "Server/Const"
	"Server/Net"
	"fmt"
	"log"
	"reflect"
)

func (client *Client)Receive(cancel func()){
	var buff Net.Message
	for {
		err:= client.ReadJSON(&buff)
		if err != nil{
			log.Println("err ", err)
			cancel()
			break
		}
		Net.UnpackMsg(&buff,buffPool[buff.CMD])
		buffPool[buff.CMD].Deal(client)
		if buff.CMD != Heart{
			log.Println("deal " + reflect.TypeOf(buffPool[buff.CMD]).String())
		}
	}
}

func (client *Client)Send(){
	var in string
	client.WriteJSON(Net.PackMsg(Net.NewHelloMsg(client.player.ID,"123")))
	for {
		log.Println("input cmd")
		fmt.Scanln(&in)
		switch in {
		case "invite":
			log.Println("input oth")
			fmt.Scanln(&in)
			client.WriteJSON(Net.PackMsg(Net.NewInviteMsg(client.player.ID,in,InviteCreate,nil)))
		case "start":
			client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionStart,client.player,client.room.RoomID,"")))
		case "exit":
			client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionExit,client.player,client.room.RoomID,"")))
		case "match":
			client.WriteJSON(Net.NewMsg(Match,nil))
		case "state":
			log.Println(client.room)
		case "change":
			client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionChangeSite,client.player,client.room.RoomID,"")))
		case "inroom":
			log.Println("input roomID")
			fmt.Scanln(&in)
			client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionIn,client.player,in,"")))
			log.Println("send invite")
		case "again":
			client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionAgain,client.player,client.room.RoomID,"")))
		case "regret":
			client.WriteJSON(Net.PackMsg(Net.NewPlayMsg(Net.Chess{X: -1, Y: -1},PlayRegret,0)))
		case "play":
			var x,y int
			fmt.Scanln(&x)
			fmt.Scanln(&y)
			client.WriteJSON(Net.PackMsg(Net.NewPlayMsg(Net.Chess{X: x, Y: y},0,0)))
		case "draw":
			client.WriteJSON(Net.PackMsg(Net.NewPlayMsg(Net.Chess{},PlayDraw,0)))
		case "data":
			var cmd string
			fmt.Scanln(&cmd)
			apply := make(map[string]interface{})
			switch cmd{
			case CmdMsgRoomData:
				var chche string
				fmt.Scanln(&chche)
				apply["cmd"] = CmdMsgRoomData
				apply["roomID"] = chche
			case CmdMsgAllRoomBase:apply["cmd"] = CmdMsgAllRoomBase
			case CmdMsgPlayerState:
				var chche string
				fmt.Scanln(&chche)
				apply["cmd"] = CmdMsgPlayerState
				apply["playerID"] = chche
			}
			client.WriteJSON(Net.PackMsg(&Net.InfoMsg{Data: apply}))
		default:
			log.Println("err input")
		}
	}
}
