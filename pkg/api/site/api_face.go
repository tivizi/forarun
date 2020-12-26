package site

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GET("/api/v1/common/faces", listFaces)
}

func listFaces(c *gin.Context, _ *site.Context) (interface{}, error) {
	return site.GetFaces(), nil
}
