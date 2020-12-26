package admin

import (
	"github.com/tivizi/forarun/pkg/base"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func init() {
	db = base.DBInstance()
}
