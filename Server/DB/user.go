package DB

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

func InsertUser(id,pwd string) string{
	result,err := userCollection.InsertOne(context.Background(),Player{
		ID: id,
		Password: pwd,
		Score: 0,
	})
	if err != nil{
		log.Fatal(err)
	}

	return result.InsertedID.(primitive.ObjectID).String()
}

func GetPlayer(id string) *Player{
	var player Player
	if err := userCollection.FindOne(context.Background(),bson.D{{"id",id}}).Decode(&player); err != nil{
		return nil
	}
	return &player
}

func UpdatePlayer(id string,score int){
	filter := bson.D{{"ID", id}}
	update := bson.D{{"$set", bson.D{{"Score", score}}}}
	userCollection.UpdateOne(context.Background(),filter,update)
}

func ListUser()[]*Player{
	ctx,cancel := context.WithTimeout(context.Background(),time.Second * 3)
	defer cancel()
	opts :=options.Find().SetSort(bsonx.Doc{{FieldScore,bsonx.Int32(-1)}})

	cur,err := userCollection.Find(ctx,bson.D{{}},opts)

	if err != nil{
		log.Println(err)
		return nil
	}

	var ret []*Player
	for cur.Next(context.Background()){
		var node Player
		if err = cur.Decode(&node);err != nil{
			log.Println(err)
			return nil
		}

		ret = append(ret, &node)
	}

	return ret
}

func selectPlayerBase()bson.M{
	return bson.M{
		FieldID: 1,
		FieldPWD: 1,
		FieldScore: 1,
	}
}


