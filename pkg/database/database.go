package database

import (
	"context"

	"github.com/lonmarsDev/starwars-be/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination mock_database.go -package database github.com/lonmarsDev/starwars-be/pkg/database  DatabaseClientAction
type DatabaseClientAction interface {
	GetAllSavedCharacter(ctx context.Context, col string) (characters []model.Character, err error)
	InsertOne(ctx context.Context, collection string, input interface{}) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, col string, input []interface{}) (string, error)
}

type DatabaseClient struct {
	DatabaseClientAction
}

func NewDatabaseClient(dbConn string) *DatabaseClient {
	return &DatabaseClient{
		newMongoDBclient(dbConn),
	}
}
