package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	r := site.DefaultEngine()
	r.GETPage("/pages/admin/resources.html", resourceEditorIndex)
	r.GETPage("/pages/admin/resources/:rid/editor.html", resourceEditor)
}

func resourceEditor(c *gin.Context, site *site.Context) (string, interface{}, error) {
	resource, err := domain.LoadResourceByID(c.Param("rid"))
	return "admin_resources_editor.html", gin.H{
		"Resource":       resource,
		"ResourceString": resource.RawString(),
	}, err
}

func resourceEditorIndex(c *gin.Context, site *site.Context) (string, interface{}, error) {
	resources, err := domain.LoadResources()
	return "admin_resources.html", gin.H{
		"Resources": resources,
	}, err
}
