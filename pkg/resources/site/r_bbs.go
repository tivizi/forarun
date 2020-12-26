package site

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/bbs-:bbsid", bbs)
	r.GETPage("/pages/site/bbs/:bbs/threads-new.html", newThread)
	r.GETPage("/threads-:threadid", threadDetailByID)
	r.GETPage("/threads/:alias", threadDetailByAlias)
}

func bbs(c *gin.Context, site *site.Context) (string, interface{}, error) {
	bbs, err := domain.LoadBBSByID(c.Param("bbsid")[0:strings.Index(c.Param("bbsid"), ".")])
	if err != nil {
		return "", nil, err
	}
	threads, err := domain.LoadThreadsByBBSID(site.Site.ID, bbs.ID, 1, 20)
	if err != nil {
		return "", nil, err
	}
	return "site_bbs.html", gin.H{
		"BBS":     bbs,
		"Threads": threads,
	}, err
}

func newThread(c *gin.Context, site *site.Context) (string, interface{}, error) {
	bbsID := c.Param("bbs")
	bbs, err := domain.LoadBBSByID(bbsID)
	if err != nil {
		return "", nil, err
	}
	return "site_threads_new.html", gin.H{
		"Site": site.Site,
		"BBS":  bbs,
	}, err
}

func threadDetailByID(c *gin.Context, site *site.Context) (string, interface{}, error) {
	threadID := c.Param("threadid")
	thread, err := domain.LoadThreadByID(threadID[0:strings.Index(threadID, ".")])
	if err != nil {
		return "", nil, err
	}
	return threadDetail(thread, c, site)
}

func threadDetailByAlias(c *gin.Context, site *site.Context) (string, interface{}, error) {
	thread, err := domain.LoadThreadByAlias(c.Param("alias"))
	if err != nil {
		return "", nil, err
	}
	return threadDetail(thread, c, site)
}

func threadDetail(thread *domain.Thread, c *gin.Context, site *site.Context) (string, interface{}, error) {
	tcKey := "thread-vc-" + thread.ID.Hex()
	if _, err := c.Cookie(tcKey); err != nil {
		thread.IncViewCount()
		c.SetCookie(tcKey, "ok", 3600*24, "/", "", false, false)
	}
	return "site_threads_detail.html", gin.H{
		"Thread":      thread,
		"BBSContexts": thread.BBSContexts,
		"Site":        site.Site,
	}, nil
}
