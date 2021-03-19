package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
)
import "Server/DB"

func TestRedis(test *testing.T)  {
	DB.InitRedis()
	DB.Set("wt",1,time.Second * 3)
	//log.Println(DB.GetState("wt"))
	log.Println(DB.GetAllPlayerState())
	//fmt.Println(DB.Get("wt"))
}

func TestGetPlayer(test *testing.T){
	DB.InitMongo()
	ret := DB.ListUser()
	for _,v := range ret{
		log.Println(v.ID, " ", v.Score)
	}
	if ret == nil {
		test.Fail()
		DB.InsertUser("wt","123")
	}
}



func TestGetAllPlayer(test *testing.T){
	fmt.Println(bson.NewObjectId().Hex())
}