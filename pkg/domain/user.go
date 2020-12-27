package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tivizi/forarun/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户
type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	SiteID         primitive.ObjectID
	Name           string
	Password       string
	Avatar         string
	Phone          string
	Email          string
	CreateTime     time.Time
	Coins          *map[string]int64
	OnlineDuration *map[string]OnlineObj
	Active         bool
	ActiveToken    string
}

// OnlineObj 在线时长对象
type OnlineObj struct {
	Duration       int64
	LastActiveTime time.Time
}

// LoadUserByID 按ID加载用户
func LoadUserByID(userID string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var user User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	return &user, err
}

// LoadUserByName 按Name加载用户
func LoadUserByName(name string) (*User, error) {
	var user User
	err := db.Collection("users").FindOne(context.Background(), bson.M{"name": name}).Decode(&user)
	return &user, err
}

// LoadUserByEmail 按邮件加载用户
func LoadUserByEmail(email string) (*User, error) {
	var user User
	err := db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

// LoadUserByActiveToken 按激活码加载用户
func LoadUserByActiveToken(token string) (*User, error) {
	var user User
	err := db.Collection("users").FindOne(context.Background(), bson.M{"activetoken": token}).Decode(&user)
	return &user, err
}

// LoadUserByLoginID 按登录ID加载
func LoadUserByLoginID(user string) (*User, error) {
	var userObj User
	err := db.Collection("users").FindOne(context.Background(), bson.M{"$or": bson.A{
		bson.M{"name": user},
		bson.M{"email": user},
	}}).Decode(&userObj)
	return &userObj, err
}

// NewUser 新用户
func NewUser(name, email, password string, siteID primitive.ObjectID) (*User, error) {
	if count, err := db.Collection("users").CountDocuments(context.Background(), bson.M{"$or": bson.A{
		bson.M{"name": name},
		bson.M{"email": email},
	}}); count > 0 || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("用户名/邮箱已被注册")
	}
	user := User{
		ID:             primitive.NewObjectID(),
		SiteID:         siteID,
		Name:           name,
		Email:          email,
		Coins:          &map[string]int64{},
		OnlineDuration: &map[string]OnlineObj{},
		Password:       base.PasswordAlgo(password),
		CreateTime:     time.Now(),
		ActiveToken:    uuid.New().String(),
	}
	_, err := db.Collection("users").InsertOne(context.Background(), user)
	return &user, err
}

// NewActivedUser 新用户（更多初始化信息）
func NewActivedUser(name, email, phone, password string, coin int64, createTime time.Time, siteID primitive.ObjectID) (*User, error) {
	if count, err := db.Collection("users").CountDocuments(context.Background(), bson.M{"name": name}); count > 0 || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("用户名/邮箱已被注册")
	}
	user := User{
		ID:     primitive.NewObjectID(),
		SiteID: siteID,
		Name:   name,
		Phone:  phone,
		Email:  email,
		Coins: &map[string]int64{
			siteID.Hex(): coin,
		},
		Password:   base.PasswordAlgo(password),
		CreateTime: createTime,
		Active:     true,
	}
	_, err := db.Collection("users").InsertOne(context.Background(), user)
	return &user, err
}

// ActiveAccount 激活
func (u *User) ActiveAccount() error {
	_, err := db.Collection("users").UpdateOne(context.Background(), bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"active":      true,
		"activetoken": "",
	}})
	if err == nil {
		u.Active = true
	}
	return err
}

// NewActiveToken 新验证Token
func (u *User) NewActiveToken() error {
	token := uuid.New().String()
	_, err := db.Collection("users").UpdateOne(context.Background(), bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"activetoken": token,
	}})
	if err == nil {
		u.ActiveToken = token
	}
	return err
}

// ChangePassword 更改密码
func (u *User) ChangePassword(passwd string) error {
	if len(passwd) < 6 {
		return errors.New("密码格式不正确")
	}
	_, err := db.Collection("users").UpdateOne(context.Background(), bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"password":    base.PasswordAlgo(passwd),
		"activetoken": "",
	}})
	return err
}
