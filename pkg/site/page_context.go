package site

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var ubbreg *regexp.Regexp = regexp.MustCompile("\\[([a-zA-Z0-9]+):(.+)\\]")

// RenderPageWithContext 在上下文渲染页面
func RenderPageWithContext(text template.HTML, site *Context) template.HTML {
	// template render
	tpl, err := template.New("tmp").Funcs(template.FuncMap{
		"intelliTime": IntelliTime,
	}).Parse(string(text))
	if err != nil {
		return text
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, gin.H{
		"Login": site.Session != nil,
		"Admin": site.Session != nil && site.Site.User.ID.Hex() == site.Session.UserID,
		"Site":  site.Site,
	})

	ubbText := buf.String()

	// ubb render
	ubbParams := ubbreg.FindAllStringSubmatch(ubbText, -1)
	for _, params := range ubbParams {
		raw := params[0]
		t := params[1]
		args := params[2]
		render, err := NewRender(t, args, site.Site)
		if err != nil {
			ubbText = strings.ReplaceAll(ubbText, raw, err.Error())
			continue
		}
		ubbText = strings.ReplaceAll(ubbText, raw, render.Render(raw))

	}
	return template.HTML(ubbText)
}

// RenderUBB 渲染UBB
func RenderUBB(text string) string {
	text = RenderAt(text)
	text = RenderFace(text)
	return text
}
