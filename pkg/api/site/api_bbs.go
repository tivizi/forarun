package site

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GET("/api/v1/site/bbs", listBBS)
	r.POST("/api/v1/site/bbs/:bbs/threads", newThread)
	r.POST("/api/v1/site/threads/:tid/reply", newReply)
	r.PATCH("/api/v1/common/threads/:tid/good", newThreadGood)
}

// NewThreadReq 新帖子
type NewThreadReq struct {
	Title   string `binding:"required"`
	Content string `binding:"required"`
	Alias   string
}

// @summary 新帖子
// @tags 站点功能
// @params bbsID path string true "版块ID"
// @params req body true site.NewThreadReq true "NewThreadReq"
// @router /api/v1/site/bbs/:bbs/threads [post]
func newThread(c *gin.Context, site *site.Context) (interface{}, error) {
	bbsID := c.Param("bbs")
	var newThreadReq NewThreadReq
	if err := c.ShouldBind(&newThreadReq); err != nil {
		return nil, err
	}
	bbs, err := domain.LoadBBSByID(bbsID)
	if err != nil {
		return nil, err
	}
	thread, err := bbs.NewThread(newThreadReq.Title, newThreadReq.Content, newThreadReq.Alias, site.UserAgent, site.Session)
	return thread.ID.Hex(), err
}

// @Summary 列出站点所有论坛版块
// @tags 站点功能
// @Accept  json
// @Produce  json
// @Router /api/v1/site/bbs [get]
func listBBS(c *gin.Context, site *site.Context) (interface{}, error) {
	return domain.LoadSiteBBS(site.Site.ID)
}

// NewReplyReq 新回复
type NewReplyReq struct {
	Content string `binding:"required"`
}

// @summary 新帖子回复
// @tags 站点功能
// @params tid path string true "帖子ID"
// @params req body true site.NewReplyReq true "NewReplyReq"
// @router /api/v1/site/threads/:tid/reply [post]
func newReply(c *gin.Context, site *site.Context) (interface{}, error) {
	var newReplyReq NewReplyReq
	if err := c.ShouldBind(&newReplyReq); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	thread, err := domain.LoadThreadByID(c.Param("tid"))
	if err != nil {
		return nil, err
	}
	err = thread.NewReply(newReplyReq.Content, site.UserAgent, site.Session)
	return nil, err
}

// @summary 帖子点赞
// @tags 站点功能
// @params tid path string true "帖子ID"
// @router /api/v1/site/threads/:tid/good [patch]
func newThreadGood(c *gin.Context, site *site.Context) (interface{}, error) {
	thread, err := domain.LoadThreadByID(c.Param("tid"))
	if err != nil {
		return nil, err
	}
	tcKey := "thread-gc-" + thread.ID.Hex()
	if _, err := c.Cookie(tcKey); err == nil {
		return nil, errors.New("已赞")
	}
	session := site.Session
	if session == nil {
		session = &domain.Session{
			UserID: site.ClientIP,
			SID:    site.UserAgent,
		}
	}
	err = thread.Good(session)
	if err == nil {
		c.SetCookie(tcKey, "ok", 3600*24, "/", "", false, false)
	}
	return nil, err
}
