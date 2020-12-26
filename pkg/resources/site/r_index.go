package site

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/", index)
	r.GETPage("/pages/public/*alias", publicPage)

}

func index(c *gin.Context, site *site.Context) (string, interface{}, error) {
	return renderPublicPageByAlias("main", site)
}

func publicPage(c *gin.Context, site *site.Context) (string, interface{}, error) {
	return renderPublicPageByAlias(c.Param("alias"), site)
}

func renderPublicPageByAlias(alias string, site *site.Context) (string, interface{}, error) {
	page, err := domain.LoadPageByAlias(site.Site.ID, alias)
	if err != nil {
		return "", nil, errors.New("Site Page [" + alias + "] Not Found. Trace: " + err.Error())
	}
	return "site_page.html", gin.H{
		"Page":        page,
		"Site":        site.Site,
		"SiteContext": site,
	}, nil
}
