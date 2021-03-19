package main

import (
	. "Server/Const"
	"Server/Game"
	"Server/Net"
	"Server/util"
	"context"
)

func receiver(player *Net.Player,comm *Net.Comm,exit context.CancelFunc,ctx context.Context){
	for {
		select {
		case msg := <-comm.Ch:
			switch msg.CMD {
			case Hello:
				msg.Body = "this is server"
				comm.WriteJSON(msg)
			case Invite:
				var inviteMsg Net.InviteMsg
				Net.UnpackMsg(msg,&inviteMsg)
				if inviteMsg.Action == InviteBack{
					if room := Game.FindPlayRoom(inviteMsg.To);room!=nil{
						InRoom(room.Id,util.Net2GamePlayer(player))
					} else {
						comm.WriteJSON(Net.ErrMsg(ErrRoomNotExist))
					}
				} else if inviteMsg.Action == InviteCreate{
					comm.WriteJSON(msg)
				}
			case Exit:
				exit()
			case Heart:
				if msg.Body == "test"{
					comm.WriteJSON(msg)
				} else if msg.Body == "del"{
					exit()
				}
			default:
				comm.WriteJSON(msg)
			}
		case <-ctx.Done():
			PlayerOffline(player)
			return
		}
	}
}

func sender(player *Net.Player,comm *Net.Comm,exit context.CancelFunc,ctx context.Context){
	var buff Net.Message
	for {
		var err *Net.Message
		if err := comm.ReadJSON(&buff);err != nil{
			exit()
			return
		}
		switch buff.CMD {
		case Invite:err = InviteEvent(player,&buff)
		case Room:err = RoomEvent(player,&buff)
		case GAME:err = GameEvent(player.ID,&buff)
		case Match:err = MatchEvent(player)
		case Heart:comm.Ch <- &buff
		case Msg:err = MsgEvent(&buff)
		case Exit:
			exit()
			return
		}
		if err != nil{
			comm.WriteJSON(err)
		}
	}
}