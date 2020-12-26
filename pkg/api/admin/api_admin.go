package admin

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	r := site.DefaultEngine()
	r.POST("/api/v1/admin/sites", newSite)
	r.GET("/api/v1/admin/sites", filterSites)
	r.PATCH("/api/v1/admin/sites/:site/user/:user", bundleSiteUser)
}

// NewSiteReq 新站点请求
type NewSiteReq struct {
	Name   string `binding:"required"`
	Host   string `binding:"required"`
	UserID string `binding:"required"`
}

// SiteFilter 站点筛选
type SiteFilter struct {
	Page int `binding:"required"`
	Size int `binding:"required"`
	Name string
	Host string
}

// @summary 新站点
// @tags 主站管理
// @accept  json
// @produce  json
// @param req body admin.NewSiteReq true "NewSiteReq"
// @router /api/v1/admin/sites [post]
func newSite(c *gin.Context, site *site.Context) (interface{}, error) {
	var newSiteReq NewSiteReq
	if err := c.ShouldBind(&newSiteReq); err != nil {
		return nil, err
	}
	user, err := domain.LoadUserByID(newSiteReq.UserID)
	if err != nil {
		return nil, errors.New("LoadUser: " + err.Error())
	}
	_, err = domain.NewSiteBundleUser(newSiteReq.Name, newSiteReq.Host, user)
	return nil, err
}

// @Summary 列出所有站点
// @tags 主站管理
// @Accept  json
// @Produce  json
// @Router /api/v1/admin/sites [get]
func filterSites(c *gin.Context, site *site.Context) (interface{}, error) {
	var siteFilter SiteFilter
	if err := c.ShouldBindQuery(&siteFilter); err != nil {
		return nil, err
	}
	var filter bson.M
	if len(siteFilter.Host) > 0 {
		filter["hosts"] = bson.M{"$elemMatch": bson.M{"$regex": "/^.*" + siteFilter.Host + ".*$/"}}
	}
	if len(siteFilter.Name) > 0 {
		filter["name"] = bson.M{"$regex": "/^.*" + siteFilter.Name + ".*$/"}
	}
	var options options.FindOptions
	options.
		SetSkip(int64((siteFilter.Page - 1) * siteFilter.Size)).
		SetLimit(int64(siteFilter.Size)).
		SetSort(bson.M{"createtime": -1})
	cur, err := db.Collection("sites").Find(context.Background(), filter, &options)
	if err != nil {
		return nil, err
	}
	var sites []*domain.Site
	for cur.Next(context.Background()) {
		var site domain.Site
		if err := cur.Decode(&site); err == nil {
			sites = append(sites, &site)
			continue
		}
		log.Println("ERR: Site Decoe Err")
	}
	return sites, nil
}

// @Summary 绑定用户和站点
// @tags 主站管理
// @Accept  json
// @Produce  json
// @Param siteID path string true "站点ID"
// @Param userID path string true "用户ID"
// @Router /api/v1/admin/sites/:site/user/:user [patch]
func bundleSiteUser(c *gin.Context, _ *site.Context) (interface{}, error) {
	siteID := c.Param("site")
	userID := c.Param("user")

	site, err := domain.LoadSiteByID(siteID)
	if err != nil {
		return nil, err
	}
	err = site.BundleUser(userID)
	return nil, err
}
