package TokenMaster

import (
	"context"
	beatrix "github.com/meanOs/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var URI = ""

func Init(mongoUri string) {
	URI = mongoUri
}

func GetToken(id string) (int, Token) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error creating new mongo client", "GETTOKEN")
		return http.StatusInternalServerError, Token{}
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error connecting with new mongo client", "GETTOKEN")
		return http.StatusInternalServerError, Token{}
	}

	var collection = client.Database("Users").Collection("token")

	filter := bson.M{"id": id}

	var t Token

	err = collection.FindOne(ctx, filter).Decode(&t)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error searching for the token?", "GETTOKEN")
		return http.StatusNoContent, Token{}
	}
	return http.StatusOK, t
}

func RemoveToken(id string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error creating new mongo client", "REMOVETOKEN")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error connecting with new mongo client", "REMOVETOKEN")
		return
	}

	filter := bson.M{"id": id}

	var collection = client.Database("Users").Collection("token")
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error deleting token", "REMOVETOKEN")
		return
	}
	return
}

func PutToken(token Token) int {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error creating new mongo client", "PUTTOKEN")
		return http.StatusInternalServerError
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error connecting with new mongo client", "PUTTOKEN")
		return http.StatusInternalServerError
	}

	collection := client.Database("Users").Collection("token")
	_, err = collection.InsertOne(context.Background(), token)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error inserting token", "PUTTOKEN")
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
