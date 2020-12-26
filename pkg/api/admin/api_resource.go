package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.POST("/api/v1/admin/resources", newResource)
	r.PATCH("/api/v1/admin/resources/:rid/raw", patchResourceRaw)

}

// NewResourceReq 新资源文件
type NewResourceReq struct {
	ContentType string `binding:"required"`
	Name        string `binding:"required"`
	Raw         string `binding:"required"`
}

// @Summary 新的资源文件
// @tags 主站管理
// @Accept  json
// @Produce  json
// @Param req body admin.NewResourceReq true "NewResourceReq"
// @Router /api/v1/admin/resources [post]
func newResource(c *gin.Context, site *site.Context) (interface{}, error) {
	var newResourceReq NewResourceReq
	if err := c.ShouldBind(&newResourceReq); err != nil {
		return nil, err
	}
	err := domain.NewResource(newResourceReq.Name, newResourceReq.ContentType, newResourceReq.Raw)
	return nil, err
}

// PatchResourceRawReq 更改资源内容
type PatchResourceRawReq struct {
	ResourceString string `binding:"required"`
}

// @summary 修改资源文件内容
// @tags 主站管理
// @accept  json
// @produce  json
// @param rid path string true "资源文件ID"
// @param req body admin.PatchResourceRawReq true "PatchResourceRawReq"
// @router /api/v1/admin/resources/{rid}/raw [patch]
func patchResourceRaw(c *gin.Context, site *site.Context) (interface{}, error) {
	var patchResourceRawReq PatchResourceRawReq
	if err := c.ShouldBind(&patchResourceRawReq); err != nil {
		return nil, err
	}
	resource, err := domain.LoadResourceByID(c.Param("rid"))
	if err != nil {
		return nil, err
	}
	err = resource.UpdateRaw(patchResourceRawReq.ResourceString)
	return nil, err
}
