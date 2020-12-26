package domain

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Resource 资源
type Resource struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string
	ContentType string
	Raw         []byte
	Etag        string
}

// LoadResources 加载所有资源
func LoadResources() ([]*Resource, error) {
	cur, err := db.Collection("resources").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var resources []*Resource
	for cur.Next(context.Background()) {
		var r Resource
		if cur.Decode(&r) == nil {
			resources = append(resources, &r)
		}
	}
	return resources, nil
}

// LoadResource 加载资源
func LoadResource(resourceName string) (*Resource, error) {
	var resource Resource
	err := db.Collection("resources").FindOne(context.Background(), bson.M{"name": resourceName}).Decode(&resource)
	return &resource, err
}

// LoadResourceByID 加载资源
func LoadResourceByID(id string) (*Resource, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var resource Resource
	err = db.Collection("resources").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&resource)
	return &resource, err
}

// NewResource 新资源文件
func NewResource(name, contentType, raw string) error {
	bytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	checksum := md5.Sum(bytes)
	_, err = db.Collection("resources").InsertOne(context.Background(), Resource{
		ID:          primitive.NewObjectID(),
		Name:        name,
		ContentType: contentType,
		Raw:         bytes,
		Etag:        hex.EncodeToString(checksum[:]),
	})

	return err
}

// RawString RawString
func (r *Resource) RawString() string {
	if strings.Index(r.ContentType, "text") != -1 {
		return string(r.Raw)
	}
	return base64.StdEncoding.EncodeToString(r.Raw)
}

// UpdateRaw UpdateRaw
func (r *Resource) UpdateRaw(rawString string) error {
	if strings.Index(r.ContentType, "text") != -1 {
		checksum := md5.Sum([]byte(rawString))
		_, err := db.Collection("resources").UpdateOne(context.Background(), bson.M{"_id": r.ID}, bson.M{"$set": bson.M{
			"raw":  []byte(rawString),
			"etag": hex.EncodeToString(checksum[:]),
		}})
		return err
	}
	bytes, err := base64.StdEncoding.DecodeString(rawString)
	if err != nil {
		return err
	}
	checksum := md5.Sum(bytes)
	_, err = db.Collection("resources").UpdateOne(context.Background(), bson.M{"_id": r.ID}, bson.M{"$set": bson.M{
		"raw":  bytes,
		"etag": hex.EncodeToString(checksum[:]),
	}})
	return err
}
