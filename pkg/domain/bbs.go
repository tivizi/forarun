package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BBSContext 版块上下文
type BBSContext struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string
}

// BBS 论坛
type BBS struct {
	ID          primitive.ObjectID `bson:"_id"`
	SiteID      primitive.ObjectID
	Name        string
	ParentID    *primitive.ObjectID
	BBSContexts []*BBSContext
	Description string
	CreateTime  time.Time
	Stats       *Stats
}

// Stats 论坛统计
type Stats struct {
	ThreadCount int64
	ReplyCount  int64
}

// Thread 帖子
type Thread struct {
	ID             primitive.ObjectID `bson:"_id"`
	SiteID         primitive.ObjectID
	BBSID          primitive.ObjectID
	BBSContexts    []*BBSContext
	Title          string
	Alias          string
	Author         *Session
	UserAgent      string
	Content        string
	CreateTime     time.Time
	LastActiveTime time.Time
	ViewCount      int64
	GoodCount      int64
	Replies        []*Reply
}

// Reply 帖子回复
type Reply struct {
	ID         primitive.ObjectID `bson:"_id"`
	Content    string
	CreateTime time.Time
	Author     *Session
	UserAgent  string
}

// NewBBS 新论坛
func NewBBS(name, parentID, description string, siteID primitive.ObjectID) error {
	bbs := BBS{
		ID:          primitive.NewObjectID(),
		SiteID:      siteID,
		Name:        name,
		BBSContexts: []*BBSContext{},
		Description: description,
		CreateTime:  time.Now(),
		Stats:       &Stats{},
	}
	if parentID != "0" {
		objID, err := primitive.ObjectIDFromHex(parentID)
		if err != nil {
			return err
		}
		bbs.ParentID = &objID
		parent, err := LoadBBSByID(parentID)
		if err != nil {
			return err
		}
		bbs.BBSContexts = append(bbs.BBSContexts, parent.BBSContexts...)
	}
	bbs.BBSContexts = append(bbs.BBSContexts, &BBSContext{
		ID:   bbs.ID,
		Name: bbs.Name,
	})
	_, err := db.Collection("bbs").InsertOne(context.Background(), bbs)
	return err
}

// LoadSiteBBS 获取所有BBS
func LoadSiteBBS(siteID primitive.ObjectID) ([]*BBS, error) {
	cur, err := db.Collection("bbs").Find(context.Background(), bson.M{"siteid": siteID})
	if err != nil {
		return nil, err
	}
	var bbsList []*BBS
	for cur.Next(context.Background()) {
		var bbs BBS
		err = cur.Decode(&bbs)
		if err != nil {
			log.Println("ERR: Decode BBS. ", err.Error())
			continue
		}
		bbsList = append(bbsList, &bbs)
	}
	return bbsList, nil
}

// LoadBBSByID LoadBBS
func LoadBBSByID(id string) (*BBS, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var bbs BBS
	err = db.Collection("bbs").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&bbs)
	return &bbs, err
}

// Subs 子版块
func (bbs *BBS) Subs() []*BBS {
	cur, err := db.Collection("bbs").Find(context.Background(), bson.M{"parentid": bbs.ID})
	if err != nil {
		return []*BBS{}
	}
	var bbsList []*BBS
	for cur.Next(context.Background()) {
		var b BBS
		cur.Decode(&b)
		bbsList = append(bbsList, &b)
	}
	return bbsList
}

// NewThread 新帖子
func (bbs *BBS) NewThread(title, content, alias, ua string, session *Session) (*Thread, error) {
	if len(alias) == 0 {
		alias = primitive.NewObjectID().Hex()
	}
	if count, err := db.Collection("threads").CountDocuments(context.Background(), bson.M{"$and": bson.A{
		bson.M{"siteid": bbs.SiteID},
		bson.M{"$or": bson.A{
			bson.M{"alias": alias},
			bson.M{"title": title},
		}},
	}}); count > 0 || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("主题名/别名已存在")
	}
	thread := Thread{
		ID:             primitive.NewObjectID(),
		BBSID:          bbs.ID,
		BBSContexts:    bbs.BBSContexts,
		SiteID:         bbs.SiteID,
		Title:          title,
		Alias:          alias,
		Author:         session,
		UserAgent:      ua,
		Content:        content,
		CreateTime:     time.Now(),
		LastActiveTime: time.Now(),
		Replies:        []*Reply{},
	}
	_, err := db.Collection("threads").InsertOne(context.Background(), thread)
	return &thread, err
}

// LoadThreadsByFilterAndOptions 按Filter查询帖子
func LoadThreadsByFilterAndOptions(siteID primitive.ObjectID, size int64, filter *bson.M, options *options.FindOptions) ([]*Thread, error) {
	options.SetLimit(size)
	cur, err := db.Collection("threads").Find(context.Background(), bson.M{"$and": bson.A{
		bson.M{"siteid": siteID},
		filter,
	}}, options)
	if err != nil {
		return nil, err
	}
	var threads []*Thread
	for cur.Next(context.Background()) {
		var thread Thread
		err := cur.Decode(&thread)
		if err != nil {
			log.Println("ERR: decode Thread")
			continue
		}
		threads = append(threads, &thread)
	}
	return threads, nil
}

// LoadThreadsByBBSID LoadThreadsByBBSID
func LoadThreadsByBBSID(siteID primitive.ObjectID, bbsID primitive.ObjectID, page, size int64) ([]*Thread, error) {
	var options options.FindOptions
	options.SetSkip((page - 1) * size).SetSort(bson.M{"createtime": -1})
	return LoadThreadsByFilterAndOptions(siteID, size, &bson.M{"bbsid": bbsID}, &options)
}

// LoadThreadByID 按ID获取帖子
func LoadThreadByID(id string) (*Thread, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var thread Thread
	err = db.Collection("threads").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&thread)
	return &thread, err
}

// LoadThreadByAlias 按别名获取帖子
func LoadThreadByAlias(alias string) (*Thread, error) {
	var thread Thread
	err := db.Collection("threads").FindOne(context.Background(), bson.M{"alias": alias}).Decode(&thread)
	return &thread, err
}

// NewReply 新回复
func (t *Thread) NewReply(content, ua string, session *Session) error {
	_, err := db.Collection("threads").UpdateOne(context.Background(), bson.M{"_id": t.ID}, bson.M{
		"$push": bson.M{"replies": Reply{
			ID:         primitive.NewObjectID(),
			Content:    content,
			CreateTime: time.Now(),
			Author:     session,
			UserAgent:  ua,
		}},
		"$set": bson.M{"lastactivetime": time.Now()},
	})
	return err
}

// IncViewCount 增长浏览数
func (t *Thread) IncViewCount() {
	_, err := db.Collection("threads").UpdateOne(context.Background(), bson.M{"_id": t.ID}, bson.M{"$inc": bson.M{
		"viewcount": 1,
	}})
	if err != nil {
		log.Println(err)
	}
	t.ViewCount++
}

// Good 点赞
func (t *Thread) Good() error {
	_, err := db.Collection("threads").UpdateOne(context.Background(), bson.M{"_id": t.ID}, bson.M{"$inc": bson.M{
		"goodcount": 1,
	}})
	if err != nil {
		return err
	}
	t.GoodCount++
	return nil
}
