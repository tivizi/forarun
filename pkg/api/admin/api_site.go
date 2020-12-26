package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.PUT("/api/v1/site/admin/sites", updateSite)
	r.PATCH("/api/v1/site/admin/sites/pki/validation/:host/:fileauth", newPkiFileAuth)
}

// UpdateSiteReq 更新站点
type UpdateSiteReq struct {
	Name   string `binding:"required"`
	Header string
	Footer string
}

// @Summary 更新站点
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Param req body admin.UpdateSiteReq true "UpdateSiteReq"
// @Router /api/v1/site/admin/sites [put]
func updateSite(c *gin.Context, s *site.Context) (interface{}, error) {
	var updateSiteReq UpdateSiteReq
	if err := c.ShouldBind(&updateSiteReq); err != nil {
		return nil, err
	}
	siteObj, err := domain.LoadSiteByHost(s.Host)
	if err != nil {
		return nil, err
	}
	if !siteObj.Editable(s.Session) {
		return nil, errors.New("无权编辑")
	}
	err = siteObj.Update(updateSiteReq.Name, updateSiteReq.Header, updateSiteReq.Footer)
	if err == nil {
		for _, v := range siteObj.Hosts {
			site.RemoveSiteCache(v)
		}
	}
	return nil, err
}

func newPkiFileAuth(c *gin.Context, site *site.Context) (interface{}, error) {
	host := c.Param("host")
	fileauth := c.Param("fileauth")
	err := site.Site.UpdatePkiValidation(host, fileauth)
	return nil, err
}
