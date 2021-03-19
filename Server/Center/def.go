package main

import (
	"Server/Game"
	"Server/Net"
	"sync"
)

var connMap = make(map[string]*Net.Comm)
var queue = SyncQueue{list:make([]*Game.Player,0),ch: make(chan int,1)}
type SyncQueue struct {
	sync.Mutex
	list []*Game.Player
	ch chan int
}

func (q *SyncQueue) Enqueue(player *Game.Player){
	q.Lock()
	defer q.Unlock()
	q.list = append(q.list,player)
	if len(q.list) >= 2{
		q.ch <- len(q.list)
	}
}

func (q *SyncQueue) Dequeue() *Game.Player{
	q.Lock()
	defer q.Unlock()
	if len(q.list) == 0{
		return nil
	}
	ret := q.list[0]

	if len(q.list) == 1{
		q.list = make([]*Game.Player,0)
	} else {
		q.list = q.list[1:]
	}

	return ret
}
