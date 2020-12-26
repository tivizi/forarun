package site

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/approot/config"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/sites-new.html", newSite)
}

func newSite(v *gin.Context, site *site.Context) (string, interface{}, error) {
	sites, err := domain.LoadLatestSites(10)
	if err != nil {
		sites = []*domain.Site{}
	}
	return "site_sites_new.html", gin.H{
		"Sites":      sites,
		"SiteConfig": config.GetContext().SiteConfig,
	}, nil
}
