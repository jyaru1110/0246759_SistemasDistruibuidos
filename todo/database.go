package todo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() *mongo.Collection {
	fmt.Println("Connecting to MongoDB")
	//Set client options
	clientOptions := options.Client().ApplyURI("mongodb://db_mongo:27017").SetAuth(options.Credential{
		Username: "jyaru",
		Password: "12345",
	})

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to MongoDB")

	// Set the database and collection variables
	collection := client.Database("todoapp").Collection("todo")
	return collection
}
