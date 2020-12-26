package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.POST("/api/v1/site/admin/bbs", newBBS)
	r.GET("/api/v1/site/admin/bbs", listBBS)
	r.PATCH("/api/v1/site/admin/bbs/css", updateBBSCSS)
}

// NewBBSReq 新论坛
type NewBBSReq struct {
	Name        string `binding:"required"`
	ParentID    string `binding:"required"`
	Description string
}

// @Summary 新的论坛版块
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Param req body admin.NewBBSReq true "NewBBSReq"
// @Router /api/v1/site/admin/bbs [post]
func newBBS(c *gin.Context, site *site.Context) (interface{}, error) {
	var newBBSReq NewBBSReq
	if err := c.ShouldBind(&newBBSReq); err != nil {
		return nil, err
	}
	err := domain.NewBBS(newBBSReq.Name, newBBSReq.ParentID, newBBSReq.Description, site.Site.ID)
	return nil, err
}

// @Summary 列出站点所有论坛版块
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Router /api/v1/site/admin/bbs [get]
func listBBS(c *gin.Context, site *site.Context) (interface{}, error) {
	return domain.LoadSiteBBS(site.Site.ID)
}

// UpdateBBSCSSReq 更新论坛CSS
type UpdateBBSCSSReq struct {
	CSS string
}

// @Summary 更新站点论坛CSS
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Router /api/v1/site/admin/bbs/css [patch]
func updateBBSCSS(c *gin.Context, site *site.Context) (interface{}, error) {
	var updateBBSCSSReq UpdateBBSCSSReq
	if err := c.ShouldBind(&updateBBSCSSReq); err != nil {
		return nil, err
	}
	resourceName := fmt.Sprintf("/css/bbs/%s.css", site.Site.ID.Hex())
	resource, err := domain.LoadResource(resourceName)
	if err != nil {
		err = domain.NewResource(resourceName, "text/css", updateBBSCSSReq.CSS)
	} else {
		err = resource.UpdateRaw(updateBBSCSSReq.CSS)
	}
	return nil, err
}
