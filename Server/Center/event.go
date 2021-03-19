package main

import (
	. "Server/Const"
	"Server/DB"
	"Server/Game"
	"Server/Net"
	"Server/util"
	"log"
	"strconv"
)

func InviteEvent(player *Net.Player,msg *Net.Message)*Net.Message{
	log.Println("InviteEvent")
	if msg.CMD != Invite{
		return Net.ErrMsg(ErrMismatchCMD)
	}
	var inviteMsg Net.InviteMsg
	Net.UnpackMsg(msg,&inviteMsg)
	if inviteMsg.Action == InviteCreate{
		state := DB.GetState(inviteMsg.To)

		if state == HallState || state == RoomState{
			room := Game.FindPlayRoom(player.ID)
			if room == nil{
				room = Game.CreateRoom()
			}
			InRoom(room.Id,util.Net2GamePlayer(player))
			Net.UnpackMsg(msg,&inviteMsg)
			inviteMsg.Ext = room.Id
			connMap[inviteMsg.To].Ch <- Net.PackMsg(&inviteMsg)
			return nil
		}

		return Net.ErrMsg(ErrPlayIsPlaying)
	}
	if inviteMsg.Action == InviteBack{
		connMap[inviteMsg.To].Ch <- msg
		return nil
	}

	return Net.ErrMsg(ErrAction)
}

func RoomEvent(player *Net.Player,msg *Net.Message)*Net.Message{
	if msg.CMD != Room{
		return Net.ErrMsg(ErrMismatchCMD)
	}
	var roomMsg Net.RoomMsg
	Net.UnpackMsg(msg,&roomMsg)
	switch roomMsg.Action {
	case RoomActionIn:return InRoom(roomMsg.RoomID,util.Net2GamePlayer(player))
	case RoomActionExit:return ExitRoom(roomMsg.Sender.ID)
	case RoomActionChangeSite:return ChangeSite(roomMsg.RoomID)
	case RoomActionStart:return ReadyGame(roomMsg.Sender,false)
	case RoomActionAgain:return ReadyGame(roomMsg.Sender,true)
	default:
		log.Println("default room action")
	}

	return nil
}

func ReadyGame(player *Net.Player,changeSite bool) *Net.Message{
	if player == nil{
		return Net.ErrMsg(ErrPlayerOffline)
	}
	if room := Game.FindPlayRoom(player.ID);room != nil{
		if room.ReadyPlay(player.ID) != PlayerReady{
			return Net.ErrMsg(ErrReadyGame)
		}
		room.BroadcastMsg(connMap,Net.PackMsg(Net.NewRoomMsg(RoomActionReady,player,room.Id,"")))
		if room.RoomState == ReadyState{
			StartGame(room.Id,changeSite)
		}
		return nil
	}
	return Net.ErrMsg(ErrRoomNotExist)
}

func GameEvent(playerID string,msg *Net.Message)*Net.Message{
	var gameMsg Net.PlayMsg
	Net.UnpackMsg(msg,&gameMsg)
	if room := Game.FindPlayRoom(playerID);room != nil{
		if room.RoomState != GameState{	return Net.ErrMsg(ErrGame)}

		switch gameMsg.Op {
		case PlayChess:
			gameMsg.Result = room.Play(playerID,Game.Chess(gameMsg.Chess))
			gameMsg.StepNum = room.StepNum
			switch gameMsg.Result {
			case PlayOK:room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
			case PlayEnd:room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
			case PlayHasChess:return Net.ErrMsg(ErrHasChess)
			case PlayOutRange:return Net.ErrMsg(ErrOutRange)
			default:
				return Net.ErrMsg(ErrPlayOp)
			}
		case PlayUndo:
		case PlayUndoBack:
		case PlayDraw:
			site := room.PlayerSite(playerID)
			if site > 1 || room.RoomState != GameState{
				return Net.ErrMsg(ErrAction)
			}
			gameMsg.Ext = playerID
			room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
		case PlayDrawBack:
			site := room.PlayerSite(playerID)
			if site > 1 || room.RoomState != GameState{
				return Net.ErrMsg(ErrAction)
			}
			if gameMsg.Ext == "yes"{
				gameMsg.Result = PlayEnd
				gameMsg.Ext = "draw"
				gameMsg.Op = PlayDrawBack
				room.Draw()
				room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
			} else {
				gameMsg.Result = PlayOK
				gameMsg.Op = PlayDrawBack
				gameMsg.Ext = "noDraw"
				room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
			}
		case PlayRegret:
			if room.Defect(playerID){
				gameMsg.Result = PlayEnd
				gameMsg.Ext = playerID
				room.BroadcastMsg(connMap,Net.PackMsg(&gameMsg))
				return nil
			}
			return Net.ErrMsg(ErrGame)
		default:
			return Net.ErrMsg(ErrAction)
		}
	}
	return nil
}

func GameAgain(msg *Net.RoomMsg)*Net.Message{
	if room := Game.FindPlayRoom(msg.Sender.ID); room != nil{
		if room.Again(msg.Sender.ID){
			room.BroadcastMsg(connMap,util.Game2NetRoom(room))
			return nil
		}
		return Net.ErrMsg(ErrAgain)
	}
	return Net.ErrMsg(ErrRoomNotExist)
}

func ChangeSite(roomID string) *Net.Message{
	if room := Game.GetRoom(roomID);room != nil && room.ExchangeSite(){
		room.BroadcastMsg(connMap,util.Game2NetRoom(room))
		return nil
	}
	return Net.ErrMsg(ErrSiteErr)
}

func StartGame(roomID string,changeSite bool){
	if room := Game.GetRoom(roomID); room != nil{
		if room.RoomState == ReadyState{
			if changeSite{
				room.ExchangeSite()
			}
			room.StartGame()
			room.BroadcastMsg(connMap,util.Game2NetRoom(room))
		}
	}
}

func MatchBegin(){
	for {
		select {
		case <-queue.ch:
			log.Println("match queue deal start")
			for len(queue.list) >= 2 {
				player1 := queue.Dequeue()
				player2 := queue.Dequeue()
				if DB.GetState(player1.ID) != HallState{
					player1 = nil
				}
				if DB.GetState(player2.ID) != HallState{
					player2 = nil
				}
				msgEnd := Net.NewInfoMsg("cmd",CmdMsgInfo,CmdMsgInfo,"matching")
				if player1 != nil {
					var conn1 = connMap[player1.ID]
					conn1.Ch <- Net.PackMsg(msgEnd)
				}
				if player2 != nil {
					var conn2 = connMap[player2.ID]
					conn2.Ch <- Net.PackMsg(msgEnd)
				}
				if player1 != nil && player2 != nil{
					log.Println("match successful")
					room := Game.CreateRoom()
					Game.InRoom(connMap, room.Id, player1)
					Game.InRoom(connMap, room.Id, player2)
					room.StartGame()
					gameStart := util.Game2NetRoom(room)
					room.BroadcastMsg(connMap, gameStart)
					break
				}
				log.Println("match err")
				if player2 != nil{
					connMap[player2.ID].Ch <- Net.ErrMsg("try match again")
				}
				if player1 != nil{
					connMap[player1.ID].Ch <- Net.ErrMsg("try match again")
				}
			}
			log.Println("match queue deal stop")
		}
		//case <-ctx.Done():
		//	return
		//}
	}
}

func MatchEvent(player *Net.Player) *Net.Message{
	queue.Enqueue(util.Net2GamePlayer(player))
	connMap[player.ID].Ch <- Net.PackMsg(Net.NewInfoMsg("cmd",CmdMsgInfo,CmdMsgInfo,"match send"))
	return nil
}

func MsgEvent(message *Net.Message)*Net.Message{
	var infoMsg Net.InfoMsg
	Net.UnpackMsg(message,&infoMsg)
	cmd,ok := infoMsg.Data["cmd"].(string)
	if !ok{
		return Net.ErrMsg(ErrMismatchCMD)
	}
	switch cmd{
	case CmdMsgAllRoomBase:
		infoMsg.Data[CmdMsgAllRoomBase] = Game.GetAllRoomMsg()
		return Net.PackMsg(&infoMsg)
	case CmdMsgPlayerState:
		infoMsg.Data[CmdMsgPlayerState] = DB.GetPlayerState(infoMsg.Data["playerID"].(string))
		return Net.PackMsg(&infoMsg)
	case CmdMsgRoomData:
		roomID := infoMsg.Data["roomID"].(string)
		if room := Game.GetRoom(roomID);room != nil{
			return util.Game2NetRoom(room)
		} else {
			return Net.ErrMsg(ErrRoomNotExist)
		}
	default:
		return Net.ErrMsg(ErrNotSupportCmd)
	}
	return message
}

func InRoom(roomID string,player *Game.Player) *Net.Message{
	if player == nil{
		return Net.ErrMsg(ErrPlayerOffline)
	}
	oldRoom := Game.FindPlayRoom(player.ID)
	state := Game.InRoom(connMap,roomID,player)

	if state == NotExist{
		return Net.ErrMsg(ErrRoomNotExist)
	}
	if state != AlreadyInRoom{
		room := Game.FindPlayRoom(player.ID)
		connMap[player.ID].Ch <- util.Game2NetRoom(room)

		room.BroadcastMsg(connMap,Net.PackMsg(Net.NewRoomMsg(
			RoomActionIn,util.Game2NetPlayer(player),room.Id,strconv.Itoa(room.PlayerSite(player.ID)))))
	}
	if state == ExitOldRoom{
		oldRoom.BroadcastMsg(connMap,Net.PackMsg(Net.NewRoomMsg(
			RoomActionExit,util.Game2NetPlayer(player),oldRoom.Id,strconv.Itoa(oldRoom.PlayerSite(player.ID)))))
	}
	return nil
}

func ExitRoom(playerID string)*Net.Message{
	if room := Game.FindPlayRoom(playerID);room != nil{
		site := room.PlayerSite(playerID)
		Game.ExitRoom(connMap,room.Id,playerID)
		room.BroadcastMsg(connMap,Net.PackMsg(Net.NewRoomMsg(
			RoomActionExit,&Net.Player{
				ID: playerID,
			},room.Id,strconv.Itoa(site))))
	}
	connMap[playerID].Ch <- Net.NewMsg(Msg,[]byte(MsgExitBack))
	DB.Set(playerID,HallState,0)
	return nil
}
