package resources

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.GenericGET("/common/resources/*resources", resourceResponser)
	r.GenericGET("/favicon.ico", favicon)
	r.GenericGET("/.well-known/pki-validation/fileauth.txt", wellKnownPkiValidation)
}

func resourceResponser(c *gin.Context, site *site.Context) {
	resourceName := c.Param("resources")
	resource, err := domain.LoadResource(resourceName)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	etag := c.GetHeader("If-None-Match")
	if len(etag) > 0 && resource.Etag == etag {
		c.Status(304)
		return
	}
	c.Status(200)
	c.Header("Content-Type", resource.ContentType)
	c.Header("Etag", `"`+resource.Etag+`"`)
	c.Writer.Write(resource.Raw)
	c.Writer.Flush()
}

func favicon(c *gin.Context, site *site.Context) {
	f, err := ioutil.ReadFile("assets/favicon.ico")
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Status(200)
	c.Header("Content-Type", "image/icon")
	c.Writer.Write(f)
	c.Writer.Flush()
}

func wellKnownPkiValidation(c *gin.Context, site *site.Context) {
	extra, err := domain.LoadSiteExtra(site.Site.ID)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	pki := *(extra.PkiValidations)
	if fileauth, ok := pki[base64.StdEncoding.EncodeToString([]byte(site.Host))]; ok {
		c.String(200, fileauth)
		return
	}
	c.Status(404)
}
