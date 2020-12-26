package site

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/approot/config"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GenericGET("/common/static/*id", fileDownload)
}

func fileDownload(c *gin.Context, site *site.Context) {
	prefix := "https"
	if !config.GetContext().MinioConfig.HTTPS {
		prefix = "http"
	}
	c.Header("Cache-Control", "max-age=10000")
	c.Redirect(302, fmt.Sprintf("%s://%s/forarun-files%s", prefix, config.GetContext().MinioConfig.Endpoint, c.Param("id")))
}
