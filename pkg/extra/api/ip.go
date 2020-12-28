package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/base"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()

	r.GET("/open/api/ip/:ip", ip)
}

func ip(c *gin.Context, site *site.Context) (interface{}, error) {
	ip := c.Param("ip")
	if ip == "-" {
		ip = site.ClientIP
	}
	geo, err := base.IPInfo(ip)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"ip":       geo,
		"selfLink": fmt.Sprintf("https://%s/open/api/ip/%s", site.Host, ip),
	}, nil
}
