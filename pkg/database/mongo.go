package database

import (
	"context"
	"fmt"

	"github.com/lonmarsDev/starwars-be/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBClient struct {
	database *mongo.Database
}

func newMongoDBclient(dbConn string) *MongoDBClient {
	// Set client options
	clientOptions := options.Client().ApplyURI(dbConn)
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	return &MongoDBClient{
		database: client.Database("swdb"),
	}
}

func (c *MongoDBClient) GetAllSavedCharacter(ctx context.Context, col string) ([]model.Character, error) {
	collection := c.database.Collection(col)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all saved characters %+w", err)
	}
	// Iterate over the cursor and decode each document into a Character struct
	var characters []model.Character
	for cursor.Next(ctx) {
		var character model.Character
		err := cursor.Decode(&character)
		if err != nil {
			return nil, fmt.Errorf("error decoding character object %+w", err)
		}
		characters = append(characters, character)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed on cursor iteration %+w", err)
	}
	return characters, nil
}

func (c *MongoDBClient) InsertOne(ctx context.Context, col string, input interface{}) (*mongo.InsertOneResult, error) {
	collection := c.database.Collection(col)
	return collection.InsertOne(context.Background(), input)
}

func (c *MongoDBClient) InsertMany(ctx context.Context, col string, input []interface{}) (string, error) {
	collection := c.database.Collection(col)
	res, err := collection.InsertMany(context.Background(), input)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res.InsertedIDs...), nil
}
