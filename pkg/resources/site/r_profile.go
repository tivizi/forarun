package site

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.GETPage("/profile.html", profile)
	r.GETPage("/profile/:name", profileByName)
}

func profile(c *gin.Context, site *site.Context) (string, interface{}, error) {
	user, err := domain.LoadUserByName(site.Session.Name)
	if err != nil {
		return "", nil, err
	}
	_, err = domain.LoadOnline(site.Site.ID, user.ID)
	return "site_profile.html", gin.H{
		"User":   user,
		"Online": err == nil,
	}, nil
}

func profileByName(c *gin.Context, site *site.Context) (string, interface{}, error) {
	user, err := domain.LoadUserByName(c.Param("name"))
	if err != nil {
		return "", nil, err
	}
	_, err = domain.LoadOnline(site.Site.ID, user.ID)
	return "site_profile.html", gin.H{
		"User":   user,
		"Online": err == nil,
	}, nil
}
