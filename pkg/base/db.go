package base

import (
	"context"

	"github.com/tivizi/forarun/approot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func init() {
	mongodb := config.GetContext().MongoConfig
	client, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(mongodb.URI).
			SetAuth(options.Credential{
				Username: mongodb.Username,
				Password: mongodb.Password,
			}))
	if err != nil {
		panic(err)
	}
	db = client.Database("forarun")
}

// DBInstance 数据库实例
func DBInstance() *mongo.Database {
	return db
}
