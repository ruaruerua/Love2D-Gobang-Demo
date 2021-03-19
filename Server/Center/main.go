package main

import (
	_ "Server/Config"
	"Server/DB"
	"Server/Net"
	"github.com/spf13/viper"
	"log"
)

func main(){
	DB.InitMongo()
	DB.InitRedis()
	//http.HandleFunc("/",HandleCenter)
	//log.Println("listen " + viper.GetString("Center.Addr"))
	//if err := http.ListenAndServe(viper.GetString("Center.Addr"),nil); err != nil{
	//	panic(err)
	//}
	//if err := Net.ListenAndServe(viper.GetString("Center.Addr"),Handler); err != nil{
	//	log.Println(err)
	//}
	if err := Net.ListenAndServeForUDP(viper.GetString("Center.Addr"),Handler);err != nil{
		log.Println(err)
	}
}
