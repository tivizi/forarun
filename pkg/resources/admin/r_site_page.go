package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.GETPage("/pages/site/admin/pages.html", pageList)
	r.GETPage("/pages/site/admin/pages/:page/editor.html", pageEditor)

}

func pageList(c *gin.Context, site *site.Context) (string, interface{}, error) {
	pages, err := domain.LoadPages(site.Site.ID)
	return "site_admin_page.html", gin.H{
		"Pages": pages,
	}, err
}

func pageEditor(c *gin.Context, site *site.Context) (string, interface{}, error) {
	pageID := c.Param("page")
	page, err := domain.LoadPageByID(pageID)
	return "site_admin_page_editor.html", gin.H{
		"Site": site.Site,
		"Page": page,
	}, err
}
