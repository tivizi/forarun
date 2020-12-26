package site

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/login.html", login)
	r.GETPage("/logout.html", logout)
	r.GETPage("/register.html", register)
	r.GETPage("/active.html", activePage)
	r.GenericGET("/active/:token", active)
	r.GenericGET("/common/reset-password/:passwd/:token", resetPassword)
}

func login(c *gin.Context, site *site.Context) (string, interface{}, error) {
	return "site_accounts_login.html", gin.H{
		"Site": site.Site,
	}, nil
}

func logout(cc *gin.Context, site *site.Context) (string, interface{}, error) {
	return "site_accounts_logout.html", gin.H{
		"Site": site.Site,
	}, nil
}

func register(c *gin.Context, site *site.Context) (string, interface{}, error) {
	return "site_accounts_register.html", gin.H{
		"Site": site.Site,
	}, nil
}

func activePage(c *gin.Context, site *site.Context) (string, interface{}, error) {
	if site.Session.Active {
		c.Redirect(302, "/")
		return "", nil, nil
	}
	return "site_accounts_active.html", gin.H{}, nil
}

func active(c *gin.Context, site *site.Context) {
	user, err := domain.LoadUserByActiveToken(c.Param("token"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = user.ActiveAccount()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	domain.DeleteUserSession(user.ID.Hex())
	session, err := domain.NewSession(user)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.SetCookie("accessKey", session.SID, 3600*24*7, "", "", false, false)
	c.Redirect(302, "/profile.html")
}
func resetPassword(c *gin.Context, site *site.Context) {
	user, err := domain.LoadUserByActiveToken(c.Param("token"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	err = user.ChangePassword(c.Param("passwd"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Redirect(302, "/login.html")
}
