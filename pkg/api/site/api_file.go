package site

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/tivizi/forarun/pkg/base"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.POST("/api/v1/site/files", uploadFile)
}

func uploadFile(c *gin.Context, site *site.Context) (interface{}, error) {
	f, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	mcli, err := base.MinioCli()
	if err != nil {
		return nil, err
	}
	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	info, err := mcli.PutObject(context.Background(), "forarun-files",
		fmt.Sprintf("upload/%s/%s/%d_%s", site.Site.ID.Hex(), site.Session.Name, time.Now().Unix(), f.Filename),
		file, f.Size, minio.PutObjectOptions{})
	return info, err
}
