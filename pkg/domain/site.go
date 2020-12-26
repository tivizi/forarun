package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Site 站点
type Site struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string
	User       *User
	Main       bool
	Hosts      []string
	CreateTime time.Time
	Themes     []*Theme
	Footer     template.HTML
	Header     template.HTML
}

// SiteExtra 站点额外信息
type SiteExtra struct {
	SiteID         primitive.ObjectID
	PkiValidations *map[string]string
}

// Theme 主题
type Theme struct {
	Name      string
	GlobalCSS []byte
	GlobalJS  []byte
	Enabled   bool
}

// Page Site Pages
type Page struct {
	ID         primitive.ObjectID `bson:"_id"`
	SiteID     primitive.ObjectID
	Alias      string
	Name       string
	SysPage    bool
	Body       template.HTML
	Header     template.HTML
	Footer     template.HTML
	CreateTime time.Time
}

// LoadSiteByHost 按Host加载站点
func LoadSiteByHost(host string) (*Site, error) {
	var site Site
	err := db.Collection("sites").FindOne(context.Background(), bson.M{"hosts": bson.M{"$elemMatch": bson.M{"$eq": host}}}).Decode(&site)
	return &site, err
}

// LoadSiteByID 按站点ID加载
func LoadSiteByID(siteID string) (*Site, error) {
	var site Site
	if objID, err := primitive.ObjectIDFromHex(siteID); err == nil {
		err := db.Collection("sites").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&site)
		return &site, err
	}
	return nil, errors.New("Unknown error")
}

// LoadLatestSites 最新的站点
func LoadLatestSites(size int64) ([]*Site, error) {
	var options options.FindOptions
	options.SetLimit(size).SetSort(bson.M{"createtime": -1})
	cur, err := db.Collection("sites").Find(context.Background(), bson.M{}, &options)
	if err != nil {
		return nil, err
	}
	var sites []*Site
	for cur.Next(context.Background()) {
		var site Site
		cur.Decode(&site)
		sites = append(sites, &site)
	}
	return sites, nil
}

// LoadSiteExtra 加载站点额外信息
func LoadSiteExtra(siteid primitive.ObjectID) (*SiteExtra, error) {
	var extra SiteExtra
	err := db.Collection("site-extras").FindOne(context.Background(), bson.M{"siteid": siteid}).Decode(&extra)
	return &extra, err
}

// LoadPages 站点所有页面
func LoadPages(siteID primitive.ObjectID) ([]*Page, error) {
	var options options.FindOptions
	options.SetSort(bson.M{"createtime": -1})
	cur, err := db.Collection("pages").Find(context.Background(), bson.M{"siteid": siteID}, &options)
	if err != nil {
		return nil, err
	}
	var pages []*Page
	for cur.Next(context.Background()) {
		var page Page
		err := cur.Decode(&page)
		if err != nil {
			log.Println("ERR: DecodePage")
			continue
		}
		pages = append(pages, &page)
	}
	return pages, nil
}

// LoadPageByID 按ID加载页面
func LoadPageByID(id string) (*Page, error) {
	var page Page
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = db.Collection("pages").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&page)
	return &page, err
}

// LoadPageByAlias 按别名加载站点页面
func LoadPageByAlias(siteID primitive.ObjectID, alias string) (*Page, error) {
	var page Page
	err := db.Collection("pages").FindOne(context.Background(), bson.M{"siteid": siteID, "alias": alias}).Decode(&page)
	return &page, err
}

// NewSite 新站点
func NewSite(name, host string) (*Site, error) {
	return NewSiteBundleUser(name, host, nil)
}

// NewSiteBundleUser 新站点并绑定用户
func NewSiteBundleUser(name, host string, user *User) (*Site, error) {
	if _, err := LoadSiteByHost(host); err == nil {
		return nil, errors.New("已存在的站点：" + host)
	}

	site := &Site{
		ID:    primitive.NewObjectID(),
		Name:  name,
		User:  user,
		Main:  false,
		Hosts: []string{host},
		Themes: []*Theme{
			{
				Name:      "default",
				GlobalCSS: []byte{},
				GlobalJS:  []byte{},
				Enabled:   true,
			},
		},
		Header:     "",
		Footer:     "<footer>Powered by <a href=\"https://fora.run\">FORARUN</a></footer>",
		CreateTime: time.Now(),
	}

	_, err := db.Collection("sites").InsertOne(context.Background(), site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

// NewPage 新页面
func NewPage(name, header, body, footer, alias string, siteID *primitive.ObjectID) error {
	if len(alias) == 0 {
		alias = primitive.NewObjectID().Hex()
	}
	if strings.Index(alias, "/") != 0 {
		alias = "/" + alias
	}
	if count, err := db.Collection("pages").CountDocuments(context.Background(), bson.M{
		"siteid": siteID,
		"alias":  alias,
	}); count > 0 || err != nil {
		if err != nil {
			return err
		}
		return errors.New("此别名已存在")
	}
	_, err := db.Collection("pages").InsertOne(context.Background(), Page{
		ID:         primitive.NewObjectID(),
		SiteID:     *siteID,
		Name:       name,
		SysPage:    false,
		Alias:      alias,
		Header:     template.HTML(header),
		Body:       template.HTML(body),
		Footer:     template.HTML(footer),
		CreateTime: time.Now(),
	})
	return err
}

// EditPage 修改页面
func EditPage(id, name, header, body, footer string, siteID *primitive.ObjectID) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Collection("pages").UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{
		"name":   name,
		"header": header,
		"body":   body,
		"footer": footer,
	}})
	return err
}

// BundleUser 绑定用户
func (site *Site) BundleUser(userID string) error {
	user, err := LoadUserByID(userID)
	if err != nil {
		return err
	}
	_, err = db.Collection("sites").UpdateOne(context.Background(), bson.M{"_id": site.ID}, bson.M{
		"$set": bson.M{"user": &user},
	})
	return err
}

// Init 初始化站点
func (site *Site) Init(mainSite *Site) error {
	fmt.Println(mainSite)
	threadsNewPage, err := LoadPageByAlias(mainSite.ID, "/threads-new.html")
	if err != nil {
		return errors.New("Template Not Found: " + err.Error())
	}
	threadsNewPage.ID = primitive.NewObjectID()
	threadsNewPage.SiteID = site.ID
	threadsNewPage.CreateTime = time.Now()
	siteMainTemplate, err := LoadPageByAlias(mainSite.ID, "/site-main-template")
	if err != nil {
		return errors.New("Template Not Found: " + err.Error())
	}
	siteMainTemplate.ID = primitive.NewObjectID()
	siteMainTemplate.SiteID = site.ID
	siteMainTemplate.CreateTime = time.Now()
	siteMainTemplate.Alias = "main"
	_, err = db.Collection("pages").InsertMany(context.Background(), []interface{}{threadsNewPage, siteMainTemplate})
	return err
}

// Update 更新
func (site *Site) Update(name, header, footer string) error {
	_, err := db.Collection("sites").UpdateOne(context.Background(), bson.M{
		"_id": site.ID,
	}, bson.M{"$set": bson.M{
		"name":   name,
		"header": header,
		"footer": footer,
	}})
	return err
}

// Editable 是否可编辑
func (site *Site) Editable(session *Session) bool {
	return site.User.ID.Hex() == session.UserID
}

// UpdatePkiValidation 更新PkiValidation
func (site *Site) UpdatePkiValidation(host, fileauth string) error {
	var options options.UpdateOptions
	options.SetUpsert(true)
	_, err := db.Collection("site-extras").UpdateOne(context.Background(), bson.M{"siteid": site.ID},
		bson.M{"$set": bson.M{"pkivalidations." + base64.StdEncoding.EncodeToString([]byte(host)): fileauth}}, &options)
	return err
}

// Delete 删除页面
func (page *Page) Delete() error {
	if page.SysPage {
		return errors.New("禁止删除系统页面")
	}
	_, err := db.Collection("pages").DeleteOne(context.Background(), bson.M{"_id": page.ID})
	return err
}
