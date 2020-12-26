package site

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/base"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {

	r := site.DefaultEngine()

	r.POST("/api/v1/site/accounts/sessions", newSession)
	r.PATCH("/api/v1/site/accounts/active", sendActiveEmail)
	r.DELETE("/api/v1/site/accounts/sessions", deleteSession)
	r.POST("/api/v1/site/accounts", registerAccount)
	r.POST("/api/v1/site/accounts/passwd/:uid/:passwd/email", sendResetPasswordEmail)
}

// NewSessionReq 新会话
type NewSessionReq struct {
	User     string `binding:"required"`
	Password string `binding:"required"`
}

// RegisterAccountReq 注册帐号
type RegisterAccountReq struct {
	Name     string `binding:"required,min=2"`
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=6"`
}

// @Summary 用户新会话
// @tags 站点功能
// @Accept  json
// @Produce  json
// @Param req body site.NewSessionReq true "NewSessionReq"
// @Router /api/v1/site/accounts/sessions [post]
func newSession(c *gin.Context, _ *site.Context) (interface{}, error) {
	var newSessionReq NewSessionReq
	if err := c.ShouldBind(&newSessionReq); err != nil {
		return nil, err
	}
	user, err := domain.LoadUserByLoginID(newSessionReq.User)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Password != base.PasswordAlgo(newSessionReq.Password) {
		return nil, errors.New("密码错误")
	}
	session, err := domain.NewSession(user)
	if err != nil {
		return nil, err
	}
	c.SetCookie("accessKey", session.SID, 3600*24*7, "", "", false, false)
	return session, nil
}

func deleteSession(c *gin.Context, site *site.Context) (interface{}, error) {
	domain.DeleteSession(site.Session)
	c.SetCookie("accessKey", "", 0, "", "", false, false)
	return nil, nil
}

// @Summary 注册用户
// @tags 站点功能
// @Accept  json
// @Produce  json
// @Param req body site.RegisterAccountReq true "RegisterAccountReq"
// @Router /api/v1/site/accounts [post]
func registerAccount(c *gin.Context, site *site.Context) (interface{}, error) {
	var registerAccountReq RegisterAccountReq
	if err := c.ShouldBind(&registerAccountReq); err != nil {
		return nil, err
	}
	if strings.Index(registerAccountReq.Name, ",") != -1 {
		return nil, errors.New("用户名是由非[,]以外的其他字符组成")
	}
	if strings.Index(registerAccountReq.Name, "  ") != -1 {
		return nil, errors.New("用户名不能有连续的空白符号")
	}
	user, err := domain.NewUser(registerAccountReq.Name, registerAccountReq.Email, registerAccountReq.Password, site.Site.ID)
	go base.SendActiveEmail(site.Site.Name, site.Host, user.ActiveToken, user.Email)
	return nil, err
}

// @Summary 发送激活账户邮件
// @tags 站点功能
// @Accept  json
// @Produce  json
// @Router /api/v1/site/accounts/active [patch]
func sendActiveEmail(c *gin.Context, site *site.Context) (interface{}, error) {
	user, err := domain.LoadUserByID(site.Session.UserID)
	if err != nil {
		return nil, err
	}
	err = base.SendActiveEmail(site.Site.Name, site.Host, user.ActiveToken, user.Email)
	return nil, err
}

// @Summary 发送重置密码邮件
// @tags 站点功能
// @Accept  json
// @Produce  json
// @Router /api/v1/site/accounts/{uid}/passwd/{passwd} [patch]
func sendResetPasswordEmail(c *gin.Context, site *site.Context) (interface{}, error) {
	user, err := domain.LoadUserByEmail(c.Param("uid"))
	if err != nil {
		if user, err = domain.LoadUserByName(c.Param("uid")); err != nil {
			return nil, errors.New("未找到该用户")
		}
	}
	if len(c.Param("passwd")) < 6 {
		return nil, errors.New("密码不能少于6位")
	}
	if len(user.Email) == 0 {
		return nil, errors.New("该用户的邮箱不正确")
	}
	err = user.NewActiveToken()
	if err != nil {
		return nil, err
	}
	err = base.SendResetPasswordEmail(site.Site.Name, site.Host, user.Name, user.Email, user.ActiveToken, c.Param("passwd"))
	return nil, err
}
