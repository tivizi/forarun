package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Session 用户会话
type Session struct {
	SID        string
	Name       string
	Avatar     string
	Email      string
	UserID     string
	Active     bool
	CreateTime time.Time
}

// Online 在线用户
type Online struct {
	SiteID         primitive.ObjectID
	User           *Session
	UserAgent      string
	LastActiveTime time.Time
	IP             string
}

var sessionCache = make(map[string]*Session)
var revrseSessionCache = make(map[string][]string)

// NewSession 新会话
func NewSession(user *User) (*Session, error) {
	session := Session{
		SID:        primitive.NewObjectID().Hex(),
		Name:       user.Name,
		Avatar:     user.Avatar,
		Email:      user.Email,
		Active:     user.Active,
		UserID:     user.ID.Hex(),
		CreateTime: time.Now(),
	}
	_, err := db.Collection("sessions").InsertOne(context.Background(), session)
	return &session, err
}

// LoadSession 加载会话
func LoadSession(sessionID string) (*Session, error) {
	if session, ok := sessionCache[sessionID]; ok {
		return session, nil
	}
	var session Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"sid": sessionID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	sessionCache[sessionID] = &session
	revrseSessionCache[session.UserID] = append(revrseSessionCache[session.UserID], sessionID)
	return &session, nil
}

// DeleteSession 删除会话
func DeleteSession(session *Session) {
	db.Collection("sessions").DeleteOne(context.Background(), bson.M{"sid": session.SID})
	db.Collection("onlines").DeleteMany(context.Background(), bson.M{"user.userid": session.UserID})
	delete(sessionCache, session.SID)
}

// DeleteUserSession 删除某用户所有会话
func DeleteUserSession(userid string) {
	db.Collection("sessions").DeleteMany(context.Background(), bson.M{"userid": userid})
	sids := revrseSessionCache[userid]
	if sids != nil {
		for _, v := range sids {
			delete(sessionCache, v)
		}
	}
}

// NewOnlineActive 在线活动
func NewOnlineActive(siteid primitive.ObjectID, ua, clientIP string, session *Session) {
	if session == nil {
		session = &Session{
			UserID: clientIP,
			Name:   "匿名游客",
		}
	}
	var options options.UpdateOptions
	options.SetUpsert(true)
	db.Collection("onlines").UpdateOne(context.Background(), bson.M{
		"siteid":      siteid,
		"user.userid": session.UserID,
	}, bson.M{"$set": bson.M{
		"user":           session,
		"useragent":      ua,
		"ip":             clientIP,
		"lastactivetime": time.Now(),
	}}, &options)
	db.Collection("onlines").DeleteMany(context.Background(), bson.M{"lastactivetime": bson.M{
		"$lt": time.Now().Add(-time.Minute * 10),
	}})
}

// LoadOnlines 在线用户
func LoadOnlines(siteid primitive.ObjectID) ([]*Online, error) {
	var options options.FindOptions
	options.SetSort(bson.M{"lastactivetime": -1})
	cur, err := db.Collection("onlines").Find(context.Background(), bson.M{"siteid": siteid}, &options)
	if err != nil {
		return nil, err
	}
	var onlines []*Online

	for cur.Next(context.Background()) {
		var o Online
		if cur.Decode(&o) != nil {
			continue
		}
		onlines = append(onlines, &o)
	}
	return onlines, nil
}

// LoadOnline 加载在线用户
func LoadOnline(siteid, userid primitive.ObjectID) (*Online, error) {
	var online Online
	err := db.Collection("onlines").FindOne(context.Background(), bson.M{"siteid": siteid, "user.userid": userid.Hex()}).Decode(&online)
	return &online, err
}
