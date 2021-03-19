package Net

import (
	"Server/Const"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Comm struct {
	JsonIO
	Ch chan *Message
}

type UdpJson struct {
	Receiver chan []byte
	addr *net.UDPAddr
	conn *net.UDPConn
}

type UdpMsg struct {
	Data *Message
	Addr string
}

func (u *UdpJson)ReadJSON(v interface{})error{
	if r,ok := <-u.Receiver ;ok{
		err := json.Unmarshal(r,v)
		if v.(*Message).CMD != Const.Heart{
			log.Println("read len:" + strconv.Itoa(len(r)))
			log.Println(v)
		}
		return err
	}
	return errors.New("udp no file")
}

func (u *UdpJson)WriteJSON(v interface{})error{
	m,err := json.Marshal(v)
	if err != nil{
		return err
	}
	if v.(*Message).CMD != Const.Heart {
		log.Println("write len:" + strconv.Itoa(len(m)))
		log.Println(v)
	}
	_,err = u.conn.WriteToUDP(m,u.addr)
	return err
}

func (u *UdpJson)RemoteAddr()string{
	return u.addr.String()
}

type TcpJson struct {
	conn net.Conn
	readBuff []byte
	writeBuff []byte
	*json.Decoder
	*json.Encoder
}

type Message struct {
	CMD int
	Body string
}

type JsonIO interface {
	ReadJSON(interface{})error
	WriteJSON(interface{})error
	RemoteAddr() string
}

func UpgradeTcpStream(c net.Conn)*TcpJson{
	return &TcpJson{
		conn:c,
		Decoder: json.NewDecoder(c),
		Encoder:json.NewEncoder(c),
	}
}

func (m *Message)String() string{
	return "cmd: " + strconv.Itoa(m.CMD) + "   body:"+string(m.Body)
}

type MessageBody interface {
	GetID() int
}

type InviteMsg struct {
	From string
	To string
	Action int
	Ext interface{}
}

type InfoMsg struct {
	Data map[string]interface{}
}

func (*InfoMsg)GetID()int{return Const.Msg}

type Player struct {
	ID string
	Score int
}

type RoomDataMsg struct {
	RoomID string
	RoomState int
	Action int
	Size     int
	FirstPlayer *Player
	LastPlayer   *Player
	End      bool
	Winner   string
	//Snapshot []int
	StepNum int
	Audience []string
}

func (rdm *RoomDataMsg)GetID()int{return Const.RoomData}

type RoomMsg struct {
	Action int
	Sender *Player
	RoomID string
	Ext string
}

type Chess struct {
	X int
	Y int
}

type PlayMsg struct {
	RoomID string
	Chess
	Result int
	StepNum int
	Op int
	Ext string
}

func (*PlayMsg)GetID()int {return Const.GAME}


func (*RoomMsg)GetID()int{ return Const.Room}

func (*InviteMsg)GetID()int{return Const.Invite}


type HelloMsg struct {
	ID string
	Password string
}

func (*HelloMsg)GetID()int{return Const.Hello}

func PackMsg(body MessageBody)*Message{
	bytes,_ := json.Marshal(body)
	return NewMsg(body.GetID(),bytes)
}

func UnpackMsg(msg *Message,v interface{}){
	json.NewDecoder(bytes.NewReader([]byte(msg.Body))).Decode(v)
}

func NewInviteMsg(from,to string,action int,ext interface{})*InviteMsg{
	return &InviteMsg{
		From:   from,
		To:     to,
		Action: action,
		Ext:    ext,
	}
}

func ErrMsg(err interface{})*Message{
	switch err.(type) {
	case error:
		return &Message{
			Const.Err,err.(error).Error(),
		}
	case string:
		return  &Message{
			Const.Err,err.(string),
		}
	case []byte:
		return  &Message{
			Const.Err,string(err.([]byte)),
		}
	default:
		return &Message{
			CMD:  Const.Err,
			Body: "",
		}
	}

}

func NewMsg(cmd int,data []byte)*Message{
	return  &Message{
		cmd,string(data),
	}
}

func NewInfoMsg(args ...interface{})*InfoMsg{
	data := make(map[string]interface{})
	for i := 0; i < len(args) / 2; i++{
		data[args[i * 2].(string)] = args[i * 2 + 1]
	}
	return &InfoMsg{Data: data}
}

func NewHelloMsg(id,password string)*HelloMsg{
	return &HelloMsg{
		ID:       id,
		Password: password,
	}
}

func NewRoomMsg(action int,sender *Player,roomID,ext string)*RoomMsg{
	return &RoomMsg{
		Action: action,
		Sender: sender,
		RoomID: roomID,
		Ext:    ext,
	}
}

func NewPlayMsg(chess Chess,op int,result int)*PlayMsg{
	return &PlayMsg{
		RoomID: "",
		Chess:  chess,
		Result: result,
		Op: op,
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request)(*websocket.Conn, error){
	return upgrade.Upgrade(w,r,nil)
}


func NewComm(conn JsonIO)*Comm{
	return &Comm{
		Ch:make(chan *Message),
		JsonIO:conn,
	}
}

func (comm *Comm)ReadJSON(v interface{})error{
	err := comm.JsonIO.ReadJSON(v)
	return err
}

func (comm *Comm)WriteJSON(v interface{}) error{
	err := comm.JsonIO.WriteJSON(v)
	return err
}

func (comm *Comm)RemoteAddr() string{
	return comm.JsonIO.RemoteAddr()
}

func (t *TcpJson)ReadJSON(v interface{})error{
	err := t.Decoder.Decode(v)
	if v.(*Message).CMD != Const.Heart{
		log.Println("read   ",v)
	}
	return err
}

func (t *TcpJson)WriteJSON(v interface{}) error{
	err := t.Encoder.Encode(v)
	if v.(*Message).CMD != Const.Heart{
		log.Println("write   ", v)
	}
	return err
}


func (t *TcpJson)RemoteAddr()string{
	return t.conn.RemoteAddr().String()
}