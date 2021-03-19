package main

import (
	. "Server/Const"
	"Server/Net"
	"log"
)

type Client struct {
	player *Net.Player
	host string
	Net.JsonIO
	Snapshot []int
	room *RoomDataMsg
	syncStep int
}

func (client *Client)ReadJSON(v interface{})error{
	err := client.JsonIO.ReadJSON(v)
	return err
}

func (client *Client)WriteJSON(v interface{})error{
	err := client.JsonIO.WriteJSON(v)
	return err
}

func (client *Client)PrintBoard(){
	var tmp []byte
	for _,v := range client.Snapshot{
		tmp = append(tmp,byte(v))
	}
	for i := 0; i < client.room.Size;i++{
		log.Println(string(tmp[i * client.room.Size:client.room.Size*(i + 1)]))
	}
}

var buffPool = make(map[int]DealMsg)

func init(){
	buffPool[Hello] = &HelloMsg{}
	buffPool[Invite] = &InviteMsg{}
	buffPool[Room] = &RoomMsg{}
	buffPool[RoomData] = &RoomDataMsg{}
	buffPool[Err] = &ErrMsg{}
	buffPool[GAME] = &PlayMsg{}
	buffPool[Msg] = &InfoMsg{}
	buffPool[Heart] = &HeartMsg{}
}

type HelloMsg Net.HelloMsg
type ErrMsg Net.Message
type InviteMsg Net.InviteMsg
type RoomMsg Net.RoomMsg
type RoomDataMsg Net.RoomDataMsg
type PlayMsg Net.PlayMsg
type InfoMsg Net.InfoMsg
type HeartMsg Net.Message

type DealMsg interface {
	Deal(client *Client)
}

func (msg *HeartMsg)Deal(client *Client){
	client.WriteJSON(Net.NewMsg(Heart,nil))
}
func (msg *InfoMsg)Deal(client *Client){
	log.Println("deal info")
	log.Println(msg.Data["cmd"])
}
func (msg *HelloMsg)Deal(client *Client){}
func (msg *ErrMsg)Deal(client *Client){}
func (msg *PlayMsg)Deal(client *Client){
	switch msg.Op {
	case PlayChess:
		if client.Snapshot == nil{}
		if client.syncStep % 2 == 1{
			client.Snapshot[msg.X + msg.Y * client.room.Size] = '+'
		} else {
			client.Snapshot[msg.X + msg.Y * client.room.Size] = '-'
		}
		client.syncStep++
		client.PrintBoard()
	case PlayRegret:
		log.Println("play regret" + msg.Ext)
	case PlayDraw:
		if msg.Ext != client.player.ID{
			msg.Ext = "yes"
			msg.Op = PlayDrawBack
			var cov = Net.PlayMsg(*msg)
			client.WriteJSON(Net.PackMsg(&cov))
			log.Println("play PlayDraw " + msg.Ext)
		}
	}
}

func (msg *InviteMsg)Deal(client *Client){
	client.WriteJSON(Net.PackMsg(Net.NewRoomMsg(RoomActionIn,client.player,msg.Ext.(string),"")))
}

func (msg *RoomMsg)Deal(client *Client){
	switch msg.Action {
	case RoomActionIn:
		if msg.Sender.ID == client.player.ID{
			return
		}
		if msg.Ext == "0"{
			client.room.FirstPlayer = msg.Sender
		}
		if msg.Ext == "1"{
			client.room.LastPlayer = msg.Sender
		}
		if msg.Ext == "2"{
			client.room.Audience = append(client.room.Audience,msg.Sender.ID)
		}
	case RoomActionExit:
		if msg.Sender.ID == client.player.ID{
			client.room = nil
			return
		}
		if client.room.FirstPlayer != nil && msg.Sender.ID == client.room.FirstPlayer.ID{
			client.room.FirstPlayer = nil
		}else if client.room.LastPlayer !=nil &&  msg.Sender.ID == client.room.LastPlayer.ID{
			client.room.LastPlayer = nil
		}else {
			for i,v := range client.room.Audience{
				if v == msg.Sender.ID{
					if i < len(client.room.Audience){
						client.room.Audience = append(client.room.Audience[0:i],client.room.Audience[i + 1:]...)
					} else {
						client.room.Audience = client.room.Audience[0:i - 1]
					}
					break
				}
			}
		}
	case RoomActionStart:log.Println("start game")
	default:
		log.Println("default")
	}
}

func (msg *RoomDataMsg)Deal(client *Client){
	client.room = msg
	client.Snapshot = make([]int,msg.Size * msg.Size,msg.Size * msg.Size)
}



