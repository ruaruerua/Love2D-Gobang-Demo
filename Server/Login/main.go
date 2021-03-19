package main

import (
	_ "Server/Config"
	"Server/DB"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/",HandleLogin)
	DB.InitMongo()
	DB.InitRedis()
	if err := http.ListenAndServe(viper.GetString("Login.Addr"),nil); err != nil{
		panic(err)
	}
	log.Println("begin listen")
}


