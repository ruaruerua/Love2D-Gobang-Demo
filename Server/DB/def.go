package DB

import (
	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FieldID = "json_data.prop.id"
	FieldPWD = "json_data.prop.password"
	FieldScore = "json_data.prop.score"
	FieldHistory = "json_data.prop.history"

	PlayerProp = "PlayerProp"
	GameProp = "GameProp"
)

type Player struct {
	ID string
	Password string
	Score int
	//Histroy string
}
//
//type Histroy struct {
//
//}
//
//type Round struct {
//	BoardX int
//	BoardY int
//}

var userCollection *mongo.Collection
var mongodb *mongo.Database
var redisdb *redis.Client