package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.POST("/api/v1/site/admin/pages", newPage)
	r.PUT("/api/v1/site/admin/pages/:page", editPage)
	r.DELETE("/api/v1/site/admin/pages/:page", deletePage)
}

// NewPageReq 新页面
type NewPageReq struct {
	Name   string `binding:"required"`
	Alias  string
	Header string
	Body   string
	Footer string
}

// @Summary 新的页面
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Param req body admin.NewPageReq true "NewPageReq"
// @Router /api/v1/site/admin/pages [post]
func newPage(c *gin.Context, site *site.Context) (interface{}, error) {
	var newPageReq NewPageReq
	if err := c.ShouldBind(&newPageReq); err != nil {
		return nil, err
	}
	err := domain.NewPage(newPageReq.Name, newPageReq.Header, newPageReq.Body, newPageReq.Footer, newPageReq.Alias, &site.Site.ID)
	return nil, err
}

// @Summary 修改页面
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Param pageID path string true "页面ID"
// @Param req body admin.NewPageReq true "NewPageReq"
// @Router /api/v1/site/admin/pages/{pageID} [post]
func editPage(c *gin.Context, site *site.Context) (interface{}, error) {
	var newPageReq NewPageReq
	if err := c.ShouldBind(&newPageReq); err != nil {
		return nil, err
	}
	pageID := c.Param("page")
	err := domain.EditPage(pageID, newPageReq.Name, newPageReq.Header, newPageReq.Body, newPageReq.Footer, &site.Site.ID)
	return nil, err
}

// @Summary 删除页面
// @tags 站点管理
// @Accept  json
// @Produce  json
// @Param pageID path string true "页面ID"
// @Router /api/v1/site/admin/pages/{pageID} [delete]
func deletePage(c *gin.Context, site *site.Context) (interface{}, error) {
	page, err := domain.LoadPageByID(c.Param("page"))
	if err != nil {
		return nil, err
	}
	err = page.Delete()
	return nil, err
}
