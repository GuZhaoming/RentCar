package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	client, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/rentCar?readPreference=primary&ssl=false&directConnection=true"))

	if err != nil {
		panic(err)
	}
	car := client.Database("rentCar").Collection("account")

	 insertRows(c,car)
	// findRows(c, car)
}

func findRows(c context.Context, car *mongo.Collection) {
	//  res := car.FindOne(c,bson.M{
	//      "open_id":"888",
	//  })
	//  fmt.Printf("%+v\n",res)
	//  var row struct{
	// 	ID primitive.ObjectID `bson:"_id"`
	// 	OpenID string `bson:"open_id"`
	//  }
	//  err := res.Decode(&row)
	//  if err != nil{
	// 	panic(err)
	//  }
	//  fmt.Printf("%+v\n",row)

	cur, err := car.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}

	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}

		err = cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v", row)
	}

}

func insertRows(c context.Context, car *mongo.Collection) {
	res, err := car.InsertMany(c, []interface{}{
		bson.M{
			"open_id": "888",
		},
		bson.M{
			"open_id": "999",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}
