package site

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/approot/config"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.POST("/api/v1/site/sites", newSite)
}

// NewSiteReq 新站点
type NewSiteReq struct {
	Name   string `binding:"required"`
	Prefix string `binding:"required"`
}

// @summary 新站点
// @tags 站点功能
// @params req body site.NewSiteReq true "NewSiteReq"
// @router /api/v1/site/sites [post]
func newSite(c *gin.Context, site *site.Context) (interface{}, error) {
	siteDomain := config.GetContext().SiteConfig.Domain
	var newSiteReq NewSiteReq
	if err := c.ShouldBind(&newSiteReq); err != nil {
		return nil, err
	}
	if len(newSiteReq.Prefix) < 2 {
		return nil, errors.New("前缀不少于2位")
	}
	user, err := domain.LoadUserByID(site.Session.UserID)
	if err != nil {
		return nil, errors.New("LoadUser: " + err.Error())
	}
	siteObj, err := domain.NewSiteBundleUser(newSiteReq.Name, strings.Trim(strings.ToLower(newSiteReq.Prefix), " ")+"."+siteDomain, user)
	if err != nil {
		return nil, err
	}
	err = siteObj.Init(site.MainSite)
	return nil, err
}
