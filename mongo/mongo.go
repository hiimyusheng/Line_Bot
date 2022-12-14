package mongo

import (
	"context"
	"fmt"
	"line_bot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://localhost:27017/?maxPoolSize=20&w=majority"

func ConnectDB() (mongo.Client, error) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		// if err = client.Disconnect(context.TODO()); err != nil {
		// 	panic(err)
		// }
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return *client, nil
}

func RecieveMessage(message model.Message, client mongo.Client) {
	coll := client.Database("line").Collection("message")
	_, err := coll.InsertOne(context.TODO(), message)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Successfully")
}

func QueryMessage(user string, client mongo.Client) []model.Message {
	var result []model.Message
	filter := bson.D{{"id", user}}
	coll := client.Database("line").Collection("message")
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	for cursor.Next(context.TODO()) {
		var message model.Message
		err := cursor.Decode(&message)
		if err != nil {
			panic(err)
		}
		result = append(result, message)
	}
	return result
}
