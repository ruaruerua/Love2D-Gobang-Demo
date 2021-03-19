package main

import (
	"Server/Const"
	"Server/DB"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleLogin(w http.ResponseWriter, r *http.Request){
	log.Println("new connect " + r.RemoteAddr)
	var player DB.Player
	//data,err :=	ioutil.ReadAll(r.Body)
	//log.Println(string(data))
	//if err != nil{
	//	w.Write([]byte("消息读取错误"))
	//	return
	//}
	//json.Unmarshal(data,&player)
	//log.Println(player)
	player.ID = r.RemoteAddr
	dbPlayer := DB.GetPlayer(player.ID)

	if dbPlayer == nil{
		DB.InsertUser(player.ID,player.Password)
	} else if dbPlayer.Password != player.Password{
		w.Write([]byte("密码错误"))
		return
	}

	//if DB.GetState(player.ID) != -1{
	//	w.Write([]byte("用户已登录"))
	//	return
	//}
	DB.Set(player.ID,Const.OnlineState,time.Second *3)
	w.Write([]byte(viper.GetString("Center.Addr") + "|" + 	player.ID + "|" + strconv.Itoa(player.Score)))
}

