package site

import (
	"html/template"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tivizi/forarun/approot/config"
	"github.com/tivizi/forarun/pkg/base"
	"github.com/tivizi/forarun/pkg/domain"
)

var engine *Engine
var siteCache = make(map[string]string)
var siteUniqueCache = make(map[string]*domain.Site)
var mainSite *domain.Site

func init() {
	engine = &Engine{
		engine: gin.Default(),
	}

	// 站点上下文准备
	engine.engine.Use(siteContext())
	// 准入控制 - 认证
	engine.engine.Use(authentication())
	// 准入控制 - 帐号状态
	engine.engine.Use(account())
	// 准入控制 - 鉴权
	engine.engine.Use(authorization())
	// 加载网页模板
	engine.engine.SetFuncMap(funcMapForRender())
	engine.engine.LoadHTMLGlob("templates/*")
	engine.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// 加载主站点信息
	site, err := domain.LoadSiteByHost(config.GetContext().SiteConfig.Domain)
	if err != nil {
		panic("Main Site Load Error: " + err.Error())
	}
	mainSite = site
}

// RemoveSiteCache 清除SiteCache
func RemoveSiteCache(host string) {
	if v, ok := siteCache[host]; ok {
		delete(siteUniqueCache, v)
	}
}

// Engine 站点引擎
type Engine struct {
	engine *gin.Engine
}

// Context 站点上下文
type Context struct {
	Host       string
	RequestURI string
	UserAgent  string
	ClientIP   string
	Site       *domain.Site
	MainSite   *domain.Site
	Session    *domain.Session
}

// GetRequestHandler 站点请求处理器
type GetRequestHandler func(c *gin.Context, siteCtx *Context)

// RequestHandler 站点请求处理器
type RequestHandler func(c *gin.Context, siteCtx *Context) (interface{}, error)

// PageRequestHandler 站点请求处理器
type PageRequestHandler func(c *gin.Context, siteCtx *Context) (string, interface{}, error)

// GET GET
func (r *Engine) GET(relativePath string, handler RequestHandler) {
	r.engine.GET(relativePath, func(c *gin.Context) {
		withContext(c, handler)
	})
}

// GenericGET GenericGET
func (r *Engine) GenericGET(relativePath string, handler GetRequestHandler) {
	r.engine.GET(relativePath, func(c *gin.Context) {
		siteContext, _ := c.Get("siteContext")
		siteCtx := siteContext.(Context)
		handler(c, &siteCtx)
	})
}

// GETPage GETPage
func (r *Engine) GETPage(relativePath string, handler PageRequestHandler) {
	r.engine.GET(relativePath, func(c *gin.Context) {
		siteContext, _ := c.Get("siteContext")
		siteCtx := siteContext.(Context)
		tplPath, model, err := handler(c, &siteCtx)
		if err != nil {
			c.HTML(500, "site_500.html", gin.H{"Message": err.Error()})
			c.Abort()
			return
		}
		if m, ok := model.(gin.H); ok {
			m["SiteContext"] = &siteCtx
		}
		if len(tplPath) != 0 {
			c.HTML(200, tplPath, model)
		}
		onlineState <- &siteCtx
	})
}

// POST POST
func (r *Engine) POST(relativePath string, handler RequestHandler) {
	r.engine.POST(relativePath, func(c *gin.Context) {
		withContext(c, handler)
	})
}

// PUT PUT
func (r *Engine) PUT(relativePath string, handler RequestHandler) {
	r.engine.PUT(relativePath, func(c *gin.Context) {
		withContext(c, handler)
	})
}

// DELETE DELETE
func (r *Engine) DELETE(relativePath string, handler RequestHandler) {
	r.engine.DELETE(relativePath, func(c *gin.Context) {
		withContext(c, handler)
	})
}

// PATCH PATCH
func (r *Engine) PATCH(relativePath string, handler RequestHandler) {
	r.engine.PATCH(relativePath, func(c *gin.Context) {
		withContext(c, handler)
	})
}

// Run 启动引擎
func (r *Engine) Run(addr ...string) (err error) {
	return r.engine.Run(addr...)
}

func withContext(c *gin.Context, handler RequestHandler) {
	siteContext, _ := c.Get("siteContext")
	siteCtx := siteContext.(Context)
	obj, err := handler(c, &siteCtx)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Code": "50000", "Message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"Code":    "20000",
		"Message": "Success",
		"Data":    obj,
	})
	onlineState <- &siteCtx
}

// DefaultEngine 默认站点引擎
func DefaultEngine() *Engine {
	return engine
}

func authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteContext, _ := c.Get("siteContext")
		site := siteContext.(Context)
		var accessKey string
		accessKey = c.Query("accessKey")
		if len(accessKey) == 0 {
			accessKey = c.GetHeader("Authorization")
			if len(accessKey) == 0 {
				ak, err := c.Cookie("accessKey")
				if err != nil {
					onUnauthorized(c, &site)
					return
				}
				accessKey = ak
			}
		}
		session, err := domain.LoadSession(accessKey)
		if err != nil {
			onUnauthorized(c, &site)
			return
		}
		site.Session = session
		c.Set("siteContext", site)
		c.Next()
	}
}

var genricAPI = []string{
	"/common",
	"/active/",
	"/active.html",
	"/api/v1/site/accounts/sessions",
	"/api/v1/site/accounts/active",
	"/logout.html",
}

func account() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range genricAPI {
			if strings.Index(c.Request.RequestURI, v) == 0 {
				c.Next()
				return
			}
		}
		siteContext, _ := c.Get("siteContext")
		site := siteContext.(Context)
		if site.Session != nil && !site.Session.Active {
			c.Redirect(302, "/active.html")
			c.Abort()
		}
		c.Next()
	}
}

func authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteContext, _ := c.Get("siteContext")
		site := siteContext.(Context)
		if strings.Index(c.Request.RequestURI, "/pages/site/admin") == 0 {
			if site.Session == nil || site.Session.UserID != site.Site.User.ID.Hex() {
				c.HTML(403, "site_403.html", gin.H{})
				c.Abort()
				return
			}
			c.Next()
			return
		}
		if strings.Index(c.Request.RequestURI, "/api/v1/site/admin") == 0 {
			if site.Session == nil || site.Session.UserID != site.Site.User.ID.Hex() {
				c.AbortWithStatusJSON(403, gin.H{"Code": 40300, "Message": "Forbidden"})
				return
			}
			c.Next()
			return
		}
		if strings.Index(c.Request.RequestURI, "/pages/admin") == 0 {
			if site.Session == nil || site.Session.Name != "tivizi" {
				c.HTML(403, "site_403.html", gin.H{})
				c.Abort()
				return
			}
			c.Next()
			return
		}
		if strings.Index(c.Request.RequestURI, "/api/v1/admin") == 0 {
			if site.Session == nil || site.Session.Name != "tivizi" {
				c.AbortWithStatusJSON(403, gin.H{"Code": 40300, "Message": "Forbidden"})
				return
			}
			c.Next()
			return
		}
		c.Next()
	}
}

func siteContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Index(c.Request.RequestURI, "/api/v1/admin") == 0 ||
			strings.Index(c.Request.RequestURI, "/common") == 0 {
			ctx := Context{
				Host: c.Request.Host,
				Site: nil,
			}
			c.Set("siteContext", ctx)
			c.Next()
			return
		}
		host := c.Request.Host
		if site := siteFromHost(host); site != nil {
			log.Println("Cache Hit: ", host)
			c.Set("siteContext", newSiteContext(c, site))
			return
		}
		site, err := domain.LoadSiteByHost(c.Request.Host)
		if err != nil {
			c.HTML(404, "site_404.html", gin.H{"Host": c.Request.Host})
			c.Abort()
			return
		}
		siteUniqueCache[site.ID.Hex()] = site
		siteCache[host] = site.ID.Hex()
		c.Set("siteContext", newSiteContext(c, site))
		c.Next()
	}
}

func siteFromHost(host string) *domain.Site {
	if v, ok := siteUniqueCache[siteCache[host]]; ok {
		return v
	}
	return nil
}

func newSiteContext(c *gin.Context, site *domain.Site) Context {
	return Context{
		Host:       c.Request.Host,
		Site:       site,
		MainSite:   mainSite,
		RequestURI: c.Request.RequestURI,
		UserAgent:  c.Request.UserAgent(),
		ClientIP:   c.ClientIP(),
	}
}

var openURI = []string{
	"/common",
	"/pages/public",
	"/login",
	"/register",
	"/threads",
	"/bbs",
	"/api/v1/site/account/sessions",
	"/api/v1/common/",
	"/swagger/",
	"/profile/",
	"/active/",
	"/favicon.ico",
	"/.well-known/pki-validation/",
}

var openPostURI = []string{
	"/api/v1/site/accounts",
}

func onUnauthorized(c *gin.Context, site *Context) {
	if c.Request.RequestURI == "/" {
		return
	}
	for _, v := range openURI {
		if strings.Index(c.Request.RequestURI, v) == 0 {
			return
		}
	}
	if c.Request.Method == "POST" {
		for _, v := range openPostURI {
			if strings.Index(c.Request.RequestURI, v) == 0 {
				return
			}
		}
	}
	if strings.Index(c.GetHeader("Accept"), "html") != -1 {
		c.HTML(401, "site_401.html", site)
		c.Abort()
		return
	}
	c.AbortWithStatusJSON(401, gin.H{
		"Message": "Unauthorized",
		"Code":    "40100",
	})
}

func funcMapForRender() template.FuncMap {
	return template.FuncMap{
		"rwctx":           RenderPageWithContext,
		"intelliTime":     IntelliTime,
		"rwface":          RenderFace,
		"rwubb":           RenderUBB,
		"bbs":             ThreadContent,
		"site":            FromID,
		"region":          base.SimpleRegion,
		"intelliDuration": IntelliDuration,
	}
}
