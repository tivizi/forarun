package site

import (
	"bytes"
	"errors"
	"html/template"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tivizi/forarun/pkg/domain"
)

func init() {
}

// NewRender 新render
func NewRender(t, args string, site *domain.Site) (UBBRender, error) {
	switch t {
	case "threads":
		return &ThreadsRender{BaseRender{t, strings.Split(args, ","), site}}, nil
	default:
		return nil, errors.New("«unsupported resources»")
	}
}

// UBBRender UBB渲染器
type UBBRender interface {
	Render(raw string) string
}

// BaseRender 基础UBB渲染器
type BaseRender struct {
	Type string
	Args []string
	Site *domain.Site
}

func (r *BaseRender) renderWithTemplate(path string, resources interface{}) string {
	// template render
	tpl, err := template.New(path).Funcs(template.FuncMap{
		"intelliTime": IntelliTime,
	}).ParseFiles("templates/" + path)
	if err != nil {
		return "«load template error: " + err.Error() + "»"
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, gin.H{"resources": resources})
	return buf.String()
}

func (r *BaseRender) int64Arg(index int64) (int64, error) {
	if len(r.Args)-1 < int(index) {
		return -1, errors.New("arg" + strconv.FormatInt(index, 10) + " not found")
	}
	return strconv.ParseInt(strings.Trim(r.Args[index], " "), 10, 64)
}

func (r *BaseRender) strArg(index int64) (string, error) {
	if len(r.Args)-1 < int(index) {
		return "", errors.New("arg" + strconv.FormatInt(index, 10) + " not found")
	}
	return strings.Trim(r.Args[index], " "), nil
}
