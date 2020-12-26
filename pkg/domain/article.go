package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article 文章
type Article struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      string
	Content    string
	CreateTime time.Time
	Author     *Session
	Replies    []*ArticleReply
	ViewCount  int64
}

// ArticleReply 文章回复
type ArticleReply struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string
	CreateTime time.Time
	UserAgent  string
	Author     *Session
}

// NewArticle 新文章
func NewArticle(title, content string, session *Session) (*Article, error) {
	article := &Article{
		ID:         primitive.NewObjectID(),
		Title:      title,
		Content:    content,
		CreateTime: time.Now(),
		Author:     session,
		ViewCount:  0,
		Replies:    []*ArticleReply{},
	}
	_, err := db.Collection("articles").InsertOne(context.Background(), article)
	return article, err
}

// NewReply 新回复
func (article *Article) NewReply(content, userAgent string, session *Session) error {
	if session == nil {
		session = &Session{
			Name: "佚名",
		}
	}
	_, err := db.Collection("articles").UpdateOne(context.Background(), bson.M{"_id": article.ID}, bson.M{"$push": bson.M{
		"replies": &ArticleReply{
			ID:         primitive.NewObjectID(),
			Content:    content,
			CreateTime: time.Now(),
			UserAgent:  userAgent,
			Author:     session,
		},
	}})
	return err
}
