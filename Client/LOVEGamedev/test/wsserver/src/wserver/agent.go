package main

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const ReadChanSize = 1024
const WriteChanSize = 1024
//一个玩家连接对象
type Agent struct {
	conn *websocket.Conn
	readChan chan Message
	writeChan chan Message
	closeSign context.Context
	closeFun context.CancelFunc
	ticker *time.Ticker
}
//消息结构
type Message struct {
	mtype int
	data []byte

}

//new
func NewAgent(conn *websocket.Conn) *Agent {
	//关掉老的连接
	key := conn
	old := GetAgent(key)
	if old!=nil {
		old.Close()
		return nil
	}

	//新建连接
	a := &Agent{
		conn:conn,
		readChan:make(chan Message,ReadChanSize),
		writeChan:make(chan Message,WriteChanSize),
	}
	c,cancel := context.WithCancel(context.Background())
	a.closeSign = c
	a.closeFun = cancel
	a.init()
	//加入管理
	agentMap.Store(conn.RemoteAddr(),a)
	return a
}

//初始化
func (a *Agent) init()  {
	a.ticker = time.NewTicker(time.Second*5)
}

//主生命期
func (a *Agent) run()  {
	go a.readJob()
	go a.writeJob()
	for {
		select {
			case <-a.closeSign.Done():
				//close
				break
			case <-a.ticker.C:
				a.update()
			case msg:=<-a.readChan:
				a.onReadMessage(msg.mtype,msg.data)

		}
	}
	a.onDestroy()
}

func  (a *Agent) readJob()  {
	for {
		mt, message, err := a.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			a.Close()
			break
		}
		a.readChan<-Message{mtype:mt,data:message}
	}
}

func  (a *Agent) writeJob()  {
	w := a.writeChan
	for {
		select {
			case msg := <-w:
				err := a.conn.WriteMessage(msg.mtype,msg.data)
				if err != nil {
					log.Println("write error:", err)
					a.Close()
					break
				}
		}
	}
}

//销毁,在agent主线程执行
func (a *Agent) onDestroy()  {
	a.conn.Close()
	a.ticker.Stop()
	//释放管理
	agentMap.Delete(a.conn)
	log.Println("onDestroy agent:",a)
}


//收到一条消息,在agent主线程执行
func (a *Agent) onReadMessage(t int,data []byte)  {
	log.Println("recv: ", t,data)

	//echo 发回去
	a.Send(t,data)
}

//定时执行,在agent主线程执行
func (a *Agent) update()  {
	log.Println("agent udpate:",a)
	a.SendMessage([]byte("heartbeat"))
}


//发送消息,线程安全
func (a *Agent) Send(messageType int, data []byte)  {
	m := Message{mtype:messageType,data:data}
	a.writeChan<-m
}
//发送文本消息
func (a *Agent) SendMessage(data []byte)  {
	m := Message{mtype:int(websocket.TextMessage),data:data}
	a.writeChan<-m
}
//发起关闭,线程安全
func (a *Agent) Close()  {
	a.closeFun()
	log.Println("requst close:",a)
}


//agent manager
var agentMap = sync.Map{}

func GetAgent(key interface{}) *Agent  {
	old,found := agentMap.Load(key)
	if found {
		return old.(*Agent)
	}
	return nil
}