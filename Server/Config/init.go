package Config

import (
	"github.com/spf13/viper"
	"log"
)

func init(){
	setDefault()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil{
		log.Println(err)
	}

}

func setDefault(){
	viper.SetDefault("Login.Addr","127.0.0.1:8888")
	viper.SetDefault("Mongo.Addr","127.0.0.1:27017")
	viper.SetDefault("Mongo.DataBase","Gobang")

	viper.SetDefault("Redis.Addr","127.0.0.1:6379")
}