package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/pages/site/admin/bbs/index.html", bbsAdmin)
}

func bbsAdmin(c *gin.Context, site *site.Context) (string, interface{}, error) {
	bbsCSS := ""
	if resource, err := domain.LoadResource(fmt.Sprintf("/css/bbs/%s.css", site.Site.ID.Hex())); err == nil {
		bbsCSS = resource.RawString()
	}

	return "site_admin_bbs_index.html", gin.H{
		"Site":   site.Site,
		"BBSCSS": bbsCSS,
	}, nil
}
