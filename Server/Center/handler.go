package main

import (
	. "Server/Const"
	"Server/DB"
	"Server/Game"
	"Server/Net"
	"context"
	"log"
)

func init(){
	go MatchBegin()
}

func Handler(comm *Net.Comm){
	log.Println("new connect " + comm.RemoteAddr())

	var buff Net.Message

	err := comm.ReadJSON(&buff)
	comm.WriteJSON(&buff)
	if err != nil{
		log.Println(err)
		return
	}
	if buff.CMD != Hello{
		comm.WriteJSON(Net.ErrMsg(ErrProcess))
		return
	}
	var hello Net.HelloMsg
	Net.UnpackMsg(&buff,&hello)

	if err := PlayerOnline(&hello,comm); err != nil{
		comm.WriteJSON(Net.ErrMsg(err))
		return
	}
	dbPlayer := DB.GetPlayer(hello.ID)
	if dbPlayer == nil{
		comm.WriteJSON(Net.ErrMsg(ErrPlayerNotExist))
		return
	}

	player := &Net.Player{
		ID:    dbPlayer.ID,
		Score: dbPlayer.Score,
	}
	ctx,exit := context.WithCancel(context.Background())
	go receiver(player,comm,exit,ctx)
	go sender(player,comm,exit,ctx)
}

func PlayerOnline(hello *Net.HelloMsg,conn *Net.Comm) error{
	connMap[hello.ID] = conn
	log.Println("player online ",hello.ID)
	DB.Set(hello.ID, HallState,0)
	return nil
}

func PlayerOffline(player *Net.Player){
	log.Println("playoff line ", *player)
	if room := Game.FindPlayRoom(player.ID);room != nil{
		ExitRoom(player.ID)
	}
	delete(connMap,player.ID)
	DB.Del(player.ID)
}

//tcp
//func HandlerTcp(c net.Conn){
//	log.Println("new connect " + c.RemoteAddr().String())
//	conn := Net.UpgradeTcpStream(c)
//
//	var buff Net.Message
//
//	err := conn.ReadJSON(&buff)
//	if err != nil{
//		log.Println(err)
//		return
//	}
//	if buff.CMD != Hello{
//		conn.WriteJSON(Net.ErrMsg(ErrProcess))
//		return
//	}
//	var hello Net.HelloMsg
//	Net.UnpackMsg(&buff,&hello)
//
//	comm := Net.NewComm(conn)
//	if err := PlayerOnline(&hello,comm); err != nil{
//		comm.WriteJSON(Net.ErrMsg(err))
//		return
//	}
//	dbPlayer := DB.GetPlayer(hello.ID)
//	if dbPlayer == nil{
//		comm.WriteJSON(Net.ErrMsg(ErrPlayerNotExist))
//		return
//	}
//
//	player := &Net.Player{
//		ID:    dbPlayer.ID,
//		Score: dbPlayer.Score,
//	}
//	ctx,exit := context.WithCancel(context.Background())
//	go func() {
//		for {
//			select {
//			case msg := <-comm.Ch:
//				switch msg.CMD {
//				case Invite:
//					var invitMsg Net.InviteMsg
//					Net.UnpackMsg(msg,&invitMsg)
//					if invitMsg.Action == InviteBack{
//						if room := Game.FindPlayRoom(invitMsg.To);room!=nil{
//							Game.InRoom(room.Id,util.Net2GamePlayer(player))
//						} else {
//							comm.WriteJSON(Net.ErrMsg(ErrRoomNotExist))
//						}
//					} else if invitMsg.Action == InviteCreate{
//						comm.WriteJSON(msg)
//					}
//				case Exit:
//					exit()
//				default:
//					comm.WriteJSON(msg)
//				}
//			case <-ctx.Done():
//				PlayerOffline(player)
//				return
//			}
//		}
//	}()
//	go func(){
//		var buff Net.Message
//		for {
//			var err *Net.Message
//			if err := comm.ReadJSON(&buff);err != nil{
//				exit()
//				return
//			}
//			switch buff.CMD {
//			case Invite:err = InviteEvent(player,&buff)
//			case Room:err = RoomEvent(&buff)
//			case GAME:err = GameEvent(hello.ID,&buff)
//			case Exit:
//				exit()
//				return
//			}
//			if err != nil{
//				comm.WriteJSON(err)
//			}
//		}
//	}()
//}
//
//func PlayerOnline(hello *Net.HelloMsg,conn *Net.Comm) error{
//	connMap[hello.ID] = conn
//	state := DB.GetState(hello.ID)
//
//	if state != OnlineState {
//		return errors.New(ErrPlayIsOnline)
//	}
//	log.Println(DB.Set(hello.ID, HallState))
//	return nil
//}
//
//func PlayerOffline(player *Net.Player){
//	log.Println("playoff line ", *player)
//	delete(connMap,player.ID)
//
//	if room := Game.FindPlayRoom(player.ID);room != nil{
//		Game.ExitRoom(room.Id,player.ID)
//		room.BroadcastMsg(connMap,util.Game2NetRoom(room))
//	}
//	DB.Del(player.ID)
//}
//
//func HandleUDP(comm Net.JsonIO){
//	comm.ReadJSON()
//}

//websocket
//func HandleCenter(w http.ResponseWriter, r *http.Request){
//	log.Println("new connect " + r.RemoteAddr)
//	conn,err := Net.Upgrade(w,r)
//	if err != nil{
//		log.Println(err)
//		return
//	}
//	var buff Net.Message
//
//	err = conn.ReadJSON(&buff)
//	if err != nil{
//		log.Println(err)
//		return
//	}
//	if buff.CMD != Net.Hello{
//		conn.WriteJSON(Net.ErrMsg(Net.ErrProcess))
//		return
//	}
//	hello := Net.UnpackMsg(&buff).(Net.HelloMsg)
//
//	comm := Net.NewComm(conn)
//	if err := PlayerOnline(&hello,comm); err != nil{
//		comm.WriteJSON(Net.ErrMsg(err))
//		return
//	}
//	ctx,exit := context.WithCancel(context.Background())
//	go func() {
//		for {
//			var backData *Net.Message
//			select {
//			case msg := <-comm.Ch:
//				switch msg.CMD {
//				case Net.Invite:backData = msg
//				case Net.Exit:
//					exit()
//				}
//			case <-ctx.Done():
//				PlayerOffline(hello.ID)
//				break
//			}
//			if backData==nil{
//				return
//			}
//			log.Println(backData.String())
//			comm.WriteJSON(backData)
//		}
//	}()
//	go func(){
//		var backData *Net.Message
//		var buff Net.Message
//		for {
//			if err := comm.ReadJSON(&buff);err != nil{
//				exit()
//				return
//			}
//			switch buff.CMD {
//			case Net.Invite:_,_ = InviteEvent(&buff)
//			case Net.Exit:
//				exit()
//				return
//			}
//			comm.WriteJSON(backData)
//		}
//	}()
//}



