package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.GETPage("/pages/site/admin/index.html", adminIndex)
	r.GETPage("/pages/site/admin/site/editor.html", adminSiteEditor)
}

func adminIndex(c *gin.Context, site *site.Context) (string, interface{}, error) {
	onlines, err := domain.LoadOnlines(site.Site.ID)
	if err != nil {
		onlines = []*domain.Online{}
	}
	return "site_admin_index.html", gin.H{
		"Onlines": onlines,
	}, nil
}
func adminSiteEditor(c *gin.Context, site *site.Context) (string, interface{}, error) {
	return "site_admin_site_editor.html", gin.H{}, nil
}
