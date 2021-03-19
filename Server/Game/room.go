package Game

import (
	"Server/Const"
	"Server/DB"
	"Server/Net"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"sync"
)

type Room struct{
	Board
	RoomState int
	Audience map[string]bool
	m sync.Mutex
	Id string `bson:"_id" json:"_id,omitempty"`
}

var roomMap sync.Map
var playMap sync.Map

func NewRoom()*Room{
	return &Room{
		Id:bson.NewObjectId().Hex(),
		Board:*NewBoard(),
		RoomState:Const.RoomState,
	}
}

func GetAllRoomMsg()[]string{
	var ret []string
	roomMap.Range(func(k,v interface{})bool{
		ret = append(ret, k.(string) + "|" + strconv.Itoa(v.(*Room).GetRoomPlayerNum()))
		return true
	})
	return ret
}

func (r *Room)GetRoomPlayerNum()int{
	num := 0
	if r.FirstPlayer != nil{
		num++
	}
	if r.LastPlayer != nil{
		num++
	}

	return num + len(r.Audience)
}

func (r *Room)Play(player string,chess Chess)int{
	if r.RoomState != Const.GameState{
		return -1
	}

	state := r.Board.Play(player,chess)
	if state == Const.PlayEnd{
		r.EndGame()
	}
	return state
}

func (r *Room)StartGame(){
	r.changeState(Const.GameState)
	r.Board.StartGame()
}

func (r *Room)EndGame(){
	r.changeState(Const.RoomState)
}

func (r *Room)ExchangeSite()bool{
	return r.Board.ExchangeSite()
}

func (r *Room)Defect(playerID string) bool{
	if r.RoomState == Const.GameState && r.Board.Defect(playerID){
		r.changeState(Const.RoomState)
		return true
	}

	return false
}


func (r *Room)Draw(){
	r.Board.Draw()
	r.changeState(Const.RoomState)
}

func (r *Room)changeState(state int){
	DB.Set(r.FirstPlayer.ID,state,0)
	DB.Set(r.LastPlayer.ID,state,0)
	r.FirstPlayer.State = state
	r.LastPlayer.State = state
	r.RoomState = state
}


/**
TODO 逻辑bug
 */
func CreateRoom()*Room{
	room := NewRoom()
	roomMap.Store(room.Id,room)
	return room
}

func FindPlayRoom(playerID string) *Room{
	var roomID interface{}
	var ok bool
	if roomID,ok=playMap.Load(playerID);!ok{
		return nil
	}
	ret,_ := roomMap.Load(roomID)
	return ret.(*Room)
}

func GetRoom(roomID string)*Room{
	if room,ok:= roomMap.Load(roomID);ok{
		return room.(*Room)
	}
	return nil
}

func (r *Room)ReadyPlay(playerID string)int{
	if !r.End || r.FirstPlayer == nil || r.LastPlayer == nil{
		return Const.PlayerNotEnough
	}
	if r.FirstPlayer.ID == playerID{
		r.FirstPlayer.State = Const.ReadyState
		if r.LastPlayer.State == Const.ReadyState{
			r.RoomState = Const.ReadyState
		}
	} else if r.LastPlayer.ID == playerID{
		r.LastPlayer.State = Const.ReadyState
		if r.FirstPlayer.State == Const.ReadyState{
			r.RoomState = Const.ReadyState
		}
	}else{
		return Const.PlayerNotEnough
	}
	return Const.PlayerReady
}

func (r *Room)GetAllPlayer()[]string{
	var ret []string
	for k,_ := range r.Audience{
		ret = append(ret,k)
	}

	if r.FirstPlayer.ID != ""{
		ret = append(ret, r.FirstPlayer.ID)
	}

	if r.LastPlayer.ID != ""{
		ret = append(ret, r.LastPlayer.ID)
	}

	return ret
}

func (r *Room)Again(playerID string) bool{
	if (r.FirstPlayer != nil && r.FirstPlayer.ID == playerID )||
		(r.LastPlayer != nil && r.LastPlayer.ID == playerID){
		return r.Winner != "" && r.RoomState == Const.RoomState && r.ExchangeSite() && r.End
	}
	return false
}

func InRoom(conn map[string]*Net.Comm,roomID string,player *Player)int{
	state := Const.NotExist
	if id,ok := playMap.Load(player.ID);ok{
		if roomID == id.(string){
			return Const.AlreadyInRoom
		}
		ExitRoom(conn,id.(string),player.ID)
		state = Const.ExitOldRoom
	}

	if r,ok := roomMap.Load(roomID);ok{
		room := r.(*Room)
		room.m.Lock()
		defer room.m.Unlock()
		DB.Set(player.ID,Const.RoomState,0)
		playMap.Store(player.ID,roomID)
	//	if room.RoomState == Const.RoomState{
			if room.FirstPlayer==nil{
				room.FirstPlayer = player
			}else if room.LastPlayer == nil{
				room.LastPlayer = player
			} else {
				if room.Audience == nil{
					room.Audience = make(map[string]bool)
				}
				room.Audience[player.ID] = true
			}
	//	}
		if state != Const.ExitOldRoom{
			state = Const.OnlyInRoom
		}
	}
	return state
}

func ExitRoom(connMap map[string]*Net.Comm,roomID string,playerID string){
	DB.Set(playerID,Const.HallState,0)
	r,ok := roomMap.Load(roomID)
	room := r.(*Room)
	if ok{
		room.m.Lock()
		defer room.m.Unlock()
		if room.Defect(playerID){
			msg := Net.PlayMsg{
				RoomID:  room.Id,
				Chess:   Net.Chess{},
				Result:  Const.PlayEnd,
				StepNum: 0,
				Op:      Const.PlayRegret,
				Ext:     playerID,
			}
			room.BroadcastMsg(connMap,Net.PackMsg(&msg))
		}
		if room.FirstPlayer != nil && room.FirstPlayer.ID == playerID{
			room.FirstPlayer = nil
		}else if room.LastPlayer != nil && room.LastPlayer.ID == playerID {
			room.LastPlayer = nil
		}else if _,ok := room.Audience[playerID];ok{
			delete(room.Audience,playerID)
		}
		playMap.Delete(playerID)
		room.ClearRoom()
	}
}

func (r *Room)BroadcastMsg(connMap map[string]*Net.Comm,msg *Net.Message){
	if r.FirstPlayer != nil{
		connMap[r.FirstPlayer.ID].Ch <- msg
	}
	if r.LastPlayer != nil{
		connMap[r.LastPlayer.ID].Ch <- msg
	}

	for k,_ := range r.Audience{
		connMap[k].Ch <- msg
	}
}

func (r *Room)ClearRoom(){
	if r.FirstPlayer == nil &&
		r.LastPlayer == nil &&
		(r.Audience == nil ||len(r.Audience) == 0){
		roomMap.Delete(r.Id)
		log.Println("close room")
	}
}

func (r *Room)PlayerSite(playerID string) int{
	if r.FirstPlayer != nil && r.FirstPlayer.ID == playerID{
		return Const.SiteFirst
	}
	if r.LastPlayer != nil && r.LastPlayer.ID == playerID{
		return Const.SiteLast
	}
	if r.Audience != nil{
		for k,_ := range r.Audience{
			if k == playerID {
				return Const.SiteViewer
			}
		}
	}
	return Const.SiteNil
}

