package main

import (
	"flag"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
)
var addr = flag.String("addr", "localhost:8888", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	a := NewAgent(c)
	if a==nil {
		log.Println("NewAgent error")
		return
	}
	go a.run()
}

func home(w http.ResponseWriter, r *http.Request) {
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/test", echo)
	http.HandleFunc("/", home)
	fmt.Println("server start:",*addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}