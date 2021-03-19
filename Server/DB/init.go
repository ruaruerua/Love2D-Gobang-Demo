package DB

import (
	_ "Server/Config"
	"context"
	"github.com/go-redis/redis/v7"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InitMongo(){
	clientOption := options.Client().ApplyURI("mongodb://" + viper.GetString("Mongo.Addr"))
	log.Println("mongodb://" + viper.GetString("Mongo.Addr"))
	client,err := mongo.Connect(context.Background(), clientOption)
	if err != nil{
		panic(err)
	}

	err = client.Ping(context.Background(),nil)
	if err != nil{
		panic(err)
	}
	mongodb = client.Database(viper.GetString("Mongo.DataBase"))
	log.Println("connect Mongo")

	initCollection()
}

func InitRedis(){
	redisOptions := &redis.Options{}
	redisOptions.Addr = viper.GetString("Redis.Addr")
	redisdb = redis.NewClient(redisOptions)
	if err := redisdb.Ping().Err();err != nil{
		log.Println(err)
	}
}

func initCollection(){
	userCollection = mongodb.Collection(PlayerProp)
}


