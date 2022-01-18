package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)


type MongoDB struct {
	url string
	db string
	colleciton string
}

func InitMongo(url string) (client *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Println("Mongo connect error ", err)
	}
	return client
}

//func GetCollection(db, collection string) (c *mongo.Collection) {
//	return Client.Database(db).Collection(collection)
//}
//
//func Close()  {
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	Client.Disconnect(ctx)
//}